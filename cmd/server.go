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
		final, cfgFile, err := config.ParseArgsAndConfig()
		if err != nil {
			log.Fatalf("❌ %v", err)
		}

		if final.ShowVersion {
			fmt.Println("magic-log version:", version)
			return
		}

		if final.ListPresets {
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

		resolvedRegex, err := app.ResolveRegex(final.RegexPreset, final.Regex, cfgFile)
		if err != nil {
			log.Fatalf("❌ %v", err)
		}

		resolvedJqFilter, err := app.ResolveJqFilter(final.JqPreset, final.JqFilter, cfgFile)
		if err != nil {
			log.Fatalf("❌ %v", err)
		}

		app.Run(app.Config{
			DBFile:       final.DBFile,
			Port:         final.Port,
			Launch:       final.Launch,
			Echo:         final.Echo,
			Version:      version,
			LogFormat:    final.LogFormat,
			JqFilter:     resolvedJqFilter,
			CSVFieldsStr: final.CSVFieldsStr,
			HasCSVHeader: final.HasCSVHeader,
			AutoAnalyze:  final.AutoAnalyze,
			ParseRegex:   resolvedRegex,
		}, staticFiles)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
