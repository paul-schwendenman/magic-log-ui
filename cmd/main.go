package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/paul-schwendenman/magic-log-ui/internal/config"
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
	Echo       bool
	LogFormat  string
	ParseRegex string
	JqFilter   string
	Version    string
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "config" {
		handleConfigCommand(os.Args[2:])
		return
	}

	final, cfgFile, err := config.ParseArgsAndConfig()
	if err != nil {
		log.Fatalf("❌ %v", err)
	}

	if final.ShowVersion {
		fmt.Println("magic-log version:", version)
		return
	}

	if final.ListPresets {
		all := getAllPresets(cfgFile)
		fmt.Println("📜 Available parse presets:")
		for name := range all {
			fmt.Printf("  - %s\n", name)
		}
		return
	}

	resolvedRegex, err := resolveRegex(final.ParsePreset, final.ParseRegex, cfgFile)
	if err != nil {
		log.Fatalf("❌ %v", err)
	}

	Run(Config{
		DBFile:     final.DBFile,
		Port:       final.Port,
		Launch:     final.Launch,
		Echo:       final.Echo,
		Version:    version,
		LogFormat:  final.LogFormat,
		ParseRegex: resolvedRegex,
	})

}

func Run(config Config) {
	ctx := context.Background()

	db := logdb.MustInit(config.DBFile, ctx)
	logInsert := logdb.MustPrepareInsert(db, ctx)

	if config.DBFile == "" {
		log.Println("🧠 Connected to in-memory DuckDB database")
	} else {
		absPath, err := filepath.Abs(config.DBFile)
		if err != nil {
			absPath = config.DBFile // fallback
		}
		log.Printf("💾 Connected to DuckDB file: %s\n", absPath)
	}

	go server.Start(config.Port, staticFiles, db, ctx)

	if config.Launch {
		launchBrowser(config.Port)
	}

	go ingest.Start(os.Stdin, logInsert, config.LogFormat, config.ParseRegex, config.JqFilter, config.Echo, ctx)

	select {}
}

func resolveRegex(preset, raw string, cfg *config.Config) (string, error) {
	if raw != "" && preset == "" {
		return raw, nil
	}

	if preset != "" {
		all := getAllPresets(cfg)
		if regex, ok := all[preset]; ok {
			return regex, nil
		}
		return "", fmt.Errorf("unknown preset: %s", preset)
	}

	return "^(?P<message>.*)$", nil
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

func getAllPresets(cfg *config.Config) map[string]string {
	presets := map[string]string{
		"apache": `(?P<ip>\S+) (?P<ident>\S+) (?P<user>\S+) \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<path>\S+) (?P<protocol>\S+)" (?P<status>\d{3}) (?P<size>\d+|-)`,
	}

	// Merge in user-defined (override if names match)
	for k, v := range cfg.Presets {
		presets[k] = v
	}

	return presets
}
