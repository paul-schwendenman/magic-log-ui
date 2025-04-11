package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/paul-schwendenman/magic-log-ui/internal/ingest"
	"github.com/paul-schwendenman/magic-log-ui/internal/logdb"
	"github.com/paul-schwendenman/magic-log-ui/internal/server"
)

var version = "dev"

//go:embed static/*
var staticFiles embed.FS

func main() {
	var (
		dbFile      string
		openBrowser bool
		port        int
		showVersion bool
	)

	flag.StringVar(&dbFile, "db-file", "", "Path to a DuckDB database file. Leave empty for in-memory.")
	flag.BoolVar(&openBrowser, "launch", false, "Automatically open the UI in the default web browser.")
	flag.IntVar(&port, "port", 3000, "Port to serve the web UI on.")
	flag.BoolVar(&showVersion, "version", false, "Print the version and exit.")
	flag.Parse()

	if showVersion {
		fmt.Println("magic-log version:", version)
		return
	}

	ctx := context.Background()

	db := logdb.MustInit(dbFile, ctx)
	logInsert := logdb.MustPrepareInsert(db, ctx)

	// Start the web server
	go server.Start(port, staticFiles, db, ctx)

	// Open browser if flag set
	if openBrowser {
		launchBrowser(port)
	}

	// Start log ingestion from stdin
	go ingest.Start(os.Stdin, logInsert, ctx)

	// Block forever
	select {}
}

func launchBrowser(port int) {
	url := fmt.Sprintf("http://localhost:%d", port)

	var cmd *exec.Cmd
	switch {
	case isCommandAvailable("open"):
		cmd = exec.Command("open", url) // macOS
	case isCommandAvailable("xdg-open"):
		cmd = exec.Command("xdg-open", url) // Linux
	case isCommandAvailable("rundll32"):
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url) // Windows
	default:
		log.Println("⚠️ No supported method to open browser found")
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println("⚠️ Unable to open browser:", err)
	}
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
