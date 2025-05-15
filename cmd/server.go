/*
Copyright © 2025 Paul Schwendenman
*/
package cmd

import (
	"embed"

	"github.com/paul-schwendenman/magic-log-ui/internal/app"
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
		cfg := app.Config{
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

		// resolvedRegex, err := app.ResolveRegex(final.RegexPreset, final.Regex, cfgFile)
		// if err != nil {
		// 	log.Fatalf("❌ %v", err)
		// }

		// resolvedJqFilter, err := app.ResolveJqFilter(final.JqPreset, final.JqFilter, cfgFile)
		// if err != nil {
		// 	log.Fatalf("❌ %v", err)
		// }

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
