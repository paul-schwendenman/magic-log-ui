/*
Copyright © 2025 Paul Schwendenman
*/
package cmd

import (
	"embed"
	"log"

	"github.com/paul-schwendenman/magic-log-ui/internal/app"
	"github.com/paul-schwendenman/magic-log-ui/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed all:static/*
var staticFiles embed.FS

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the local web UI and begin ingesting logs",
	Long: `Starts the Magic Log web interface and begins ingesting logs from stdin.

The logs are parsed using the selected format (json, csv, or regex),
optionally filtered using jq expressions, and stored in a DuckDB database
(either in-memory or on-disk).

You can also configure presets, query past logs, and auto-analyze your data.

Examples:
  pnpm dev | magic-log server --port 5000 --log-format json
  cat logs.txt | magic-log server --regex-preset apache --log-format text`,
	Run: func(cmd *cobra.Command, args []string) {
		fileCfg, err := config.Load()
		if err != nil {
			log.Fatalf("❌ Failed to load config: %v", err)
		}

		resolvedRegex, err := app.ResolveRegex(
			viper.GetString("regex_preset"),
			viper.GetString("regex"),
			fileCfg,
		)
		if err != nil {
			log.Fatalf("❌ %v", err)
		}

		resolvedJq, err := app.ResolveJqFilter(
			viper.GetString("jq_preset"),
			viper.GetString("jq"),
			fileCfg,
		)
		if err != nil {
			log.Fatalf("❌ %v", err)
		}

		cfg := app.Config{
			DBFile:       viper.GetString("db-file"),
			Port:         viper.GetInt("port"),
			Launch:       viper.GetBool("launch"),
			Echo:         viper.GetBool("echo"),
			LogFormat:    viper.GetString("log-format"),
			ParseRegex:   resolvedRegex,
			JqFilter:     resolvedJq,
			CSVFieldsStr: viper.GetString("csv-fields"),
			HasCSVHeader: viper.GetBool("has-csv-header"),
			AutoAnalyze:  !viper.GetBool("no-auto-analyze"),
			Version:      Version,
		}

		app.Run(cfg, staticFiles)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().String("db-file", "", "Path to a DuckDB database file")
	serverCmd.Flags().Int("port", 3000, "Port to serve the web UI on")
	serverCmd.Flags().Bool("launch", false, "Open the UI in a browser")
	serverCmd.Flags().Bool("echo", false, "Echo parsed stdin input to stdout")
	serverCmd.Flags().Bool("no-auto-analyze", false, "Disable automatic ANALYZE of logs table")
	serverCmd.Flags().String("log-format", "json", "Log format: json, csv or plain text")
	serverCmd.Flags().String("regex", "", "Custom regex to parse logs (use with text format)")
	serverCmd.Flags().String("regex-preset", "", "Regex preset to use")
	serverCmd.Flags().String("jq", "", "A jq expression to apply to parsed logs")
	serverCmd.Flags().String("jq-preset", "", "jq preset to use")
	serverCmd.Flags().String("csv-fields", "", "Comma-separated field names for CSV logs")
	serverCmd.Flags().Bool("has-csv-header", true, "Whether CSV logs include a header row")

	viper.BindPFlags(serverCmd.Flags())
}
