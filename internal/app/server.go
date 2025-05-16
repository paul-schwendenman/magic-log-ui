package app

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/paul-schwendenman/magic-log-ui/internal/ingest"
	"github.com/paul-schwendenman/magic-log-ui/internal/logdb"
	"github.com/paul-schwendenman/magic-log-ui/internal/server"
	"github.com/spf13/viper"
)

type Config struct {
	DBFile       string
	Port         int
	Launch       bool
	Echo         bool
	LogFormat    string
	ParseRegex   string
	JqFilter     string
	CSVFieldsStr string
	HasCSVHeader bool
	AutoAnalyze  bool
	Version      string
}

func Run(config Config, staticFiles embed.FS) {
	log.Println("⚙️  Using config file:", viper.ConfigFileUsed())
	ctx := context.Background()

	db := logdb.MustInit(config.DBFile, ctx)
	logInsert := logdb.MustPrepareInsert(db, ctx)

	if config.AutoAnalyze {
		logdb.StartAutoAnalyze(db, ctx)
	}

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

	go ingest.Start(os.Stdin, logInsert, config.LogFormat, config.ParseRegex, config.JqFilter, config.CSVFieldsStr, config.HasCSVHeader, config.Echo, ctx)

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
