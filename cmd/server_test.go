package cmd

import (
	"testing"

	"github.com/paul-schwendenman/magic-log-ui/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestServerCmdFlagsDefaults(t *testing.T) {
	if serverCmd.Flags().Lookup("port").DefValue != "3000" {
		t.Errorf("expected default port 3000, got %s", serverCmd.Flags().Lookup("port").DefValue)
	}
	if serverCmd.Flags().Lookup("log-format").DefValue != "json" {
		t.Errorf("expected default log-format json, got %s", serverCmd.Flags().Lookup("log-format").DefValue)
	}
	// Add more default flag checks if needed.
}

func TestServerCmdRun(t *testing.T) {
	// Backup original Run function and restore later.
	origRun := serverCmd.Run
	defer func() { serverCmd.Run = origRun }()

	var capturedCfg app.Config
	// Override Run function to capture the configuration built from viper.
	serverCmd.Run = func(cmd *cobra.Command, args []string) {
		capturedCfg = app.Config{
			DBFile:       viper.GetString("db-file"),
			Port:         viper.GetInt("port"),
			Launch:       viper.GetBool("launch"),
			Echo:         viper.GetBool("echo"),
			LogFormat:    viper.GetString("log-format"),
			ParseRegex:   viper.GetString("regex"),
			JqFilter:     viper.GetString("jq"),
			CSVFieldsStr: viper.GetString("csv-fields"),
			HasCSVHeader: viper.GetBool("has-csv-header"),
			AutoAnalyze:  !viper.GetBool("no-auto-analyze"),
		}
	}

	// Set flag values for testing.
	serverCmd.Flags().Set("db-file", "/tmp/test.db")
	serverCmd.Flags().Set("port", "5000")
	serverCmd.Flags().Set("launch", "true")
	serverCmd.Flags().Set("echo", "true")
	serverCmd.Flags().Set("log-format", "csv")
	serverCmd.Flags().Set("regex", "some-regex")
	serverCmd.Flags().Set("regex-preset", "apache")
	serverCmd.Flags().Set("jq", ".")
	serverCmd.Flags().Set("jq-preset", "simple")
	serverCmd.Flags().Set("csv-fields", "a,b,c")
	serverCmd.Flags().Set("has-csv-header", "false")
	serverCmd.Flags().Set("no-auto-analyze", "true")

	// Execute the command's Run function.
	serverCmd.Run(serverCmd, []string{})

	// Validate the captured configuration.
	if capturedCfg.DBFile != "/tmp/test.db" {
		t.Errorf("expected DBFile /tmp/test.db, got %s", capturedCfg.DBFile)
	}
	if capturedCfg.Port != 5000 {
		t.Errorf("expected Port 5000, got %d", capturedCfg.Port)
	}
	if capturedCfg.Launch != true {
		t.Errorf("expected Launch true, got %v", capturedCfg.Launch)
	}
	if capturedCfg.Echo != true {
		t.Errorf("expected Echo true, got %v", capturedCfg.Echo)
	}
	if capturedCfg.LogFormat != "csv" {
		t.Errorf("expected LogFormat csv, got %s", capturedCfg.LogFormat)
	}
	if capturedCfg.ParseRegex != "some-regex" {
		t.Errorf("expected ParseRegex 'some-regex', got %s", capturedCfg.ParseRegex)
	}
	if capturedCfg.JqFilter != "." {
		t.Errorf("expected JqFilter '.', got %s", capturedCfg.JqFilter)
	}
	if capturedCfg.CSVFieldsStr != "a,b,c" {
		t.Errorf("expected CSVFieldsStr 'a,b,c', got %s", capturedCfg.CSVFieldsStr)
	}
	if capturedCfg.HasCSVHeader != false {
		t.Errorf("expected HasCSVHeader false, got %v", capturedCfg.HasCSVHeader)
	}
	if capturedCfg.AutoAnalyze != false {
		t.Errorf("expected AutoAnalyze false, got %v", capturedCfg.AutoAnalyze)
	}
}
