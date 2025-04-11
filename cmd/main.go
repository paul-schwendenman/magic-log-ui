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
	var dbFile string
	var openBrowser bool
	var port int
	var showVersion bool

	flag.StringVar(&dbFile, "db-file", "", "Path to DuckDB file. Leave empty for in-memory.")
	flag.BoolVar(&openBrowser, "launch", false, "Open the UI in the browser.")
	flag.IntVar(&port, "port", 3000, "Port to serve the web UI on.")
	flag.BoolVar(&showVersion, "version", false, "Print version and exit.")
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
		url := fmt.Sprintf("http://localhost:%d", port)
		open(url)
	}

	// Start log ingestion from stdin
	go ingest.Start(os.Stdin, logInsert, ctx)

	// Block forever
	select {}
}

func open(url string) {
	var cmd *exec.Cmd
	if _, err := exec.LookPath("open"); err == nil {
		cmd = exec.Command("open", url)
	} else if _, err := exec.LookPath("xdg-open"); err == nil {
		cmd = exec.Command("xdg-open", url)
	} else if _, err := exec.LookPath("rundll32"); err == nil {
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	} else {
		log.Println("⚠️ No supported method to open browser found")
		return
	}
	if err := cmd.Start(); err != nil {
		log.Println("⚠️ Unable to open browser:", err)
	}
}
