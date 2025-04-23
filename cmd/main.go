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

type Config struct {
	DBFile     string
	Port       int
	Launch     bool
	LogFormat  string
	ParseRegex string
	Version    string
}

func main() {
	var (
		dbFile       string
		openBrowser  bool
		port         int
		logFormat    string
		parseRegex   string
		parsePreset  string
		showVersion  bool
	)

	flag.StringVar(&dbFile, "db-file", "", "Path to a DuckDB database file. Leave empty for in-memory.")
	flag.BoolVar(&openBrowser, "launch", false, "Automatically open the UI in the default web browser.")
	flag.IntVar(&port, "port", 3000, "Port to serve the web UI on.")
	flag.StringVar(&logFormat, "log-format", "json", "Log format to parse: json or text.")
	flag.StringVar(&parseRegex, "parse-regex", "", "Regex to parse text logs (only used if --log-format=text).")
	flag.StringVar(&parsePreset, "parse-preset", "", "Preset regex name for text logs (e.g. 'apache'). Overrides --parse-regex.")
	flag.BoolVar(&showVersion, "version", false, "Print the version and exit.")
	flag.Parse()

	if showVersion {
		fmt.Println("magic-log version:", version)
		return
	}

	resolvedRegex, err := resolveRegex(parsePreset, parseRegex)
	if err != nil {
		log.Fatalf("❌ %v", err)
	}

	Run(Config{
		DBFile:     dbFile,
		Port:       port,
		Launch:     openBrowser,
		LogFormat:  logFormat,
		ParseRegex: resolvedRegex,
		Version:    version,
	})
}

func Run(config Config) {
	ctx := context.Background()

	db := logdb.MustInit(config.DBFile, ctx)
	logInsert := logdb.MustPrepareInsert(db, ctx)

	go server.Start(config.Port, staticFiles, db, ctx)

	if config.Launch {
		launchBrowser(config.Port)
	}

	go ingest.Start(os.Stdin, logInsert, config.LogFormat, config.ParseRegex, ctx)

	select {}
}

func resolveRegex(preset, raw string) (string, error) {
	if raw != "" && preset == "" {
		return raw, nil
	}
	if preset != "" {
		switch preset {
		case "apache":
			return `(?P<ip>\S+) (?P<ident>\S+) (?P<user>\S+) \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<path>\S+) (?P<protocol>\S+)" (?P<status>\d{3}) (?P<size>\d+|-)`, nil
		default:
			return "", fmt.Errorf("unknown preset: %s", preset)
		}
	}
	return "", nil
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
