/*
Copyright © 2025 Paul Schwendenman
*/
package cmd

import (
	"embed"
	"fmt"
	"log"

	"github.com/paul-schwendenman/magic-log-ui/internal/app"
	"github.com/paul-schwendenman/magic-log-ui/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "dev"

//go:embed all:static/*
var staticFiles embed.FS

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the log server and ingest process",
	Long:  ``,
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
			Version:      version,
		}

		if viper.GetBool("version") {
			fmt.Println("magic-log version:", version)
			return
		}

		if viper.GetBool("list-presets") {
			printPresets()
			return
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
	serverCmd.Flags().Bool("list-presets", false, "List regex/jq presets and exit")
	serverCmd.Flags().Bool("version", false, "Print version and exit")

	viper.BindPFlags(serverCmd.Flags())
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printPresets() {
	_, cfgFile, err := config.ParseArgsAndConfig()

	if err != nil {
		log.Fatalf("❌ %v", err)
	}

	all := app.GetRegexPresets(cfgFile)
	fmt.Println("📜 Available regexp presets:")
	for name := range all {
		fmt.Printf("  - %s\n", name)
	}

	all = app.GetJqPresets(cfgFile)
	fmt.Println("📜 Available jq presets:")
	for name := range all {
		fmt.Printf("  - %s\n", name)
	}
	return
}