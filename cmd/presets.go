/*
Copyright © 2025 Paul Schwendenman
*/
package cmd

import (
	"fmt"

	"github.com/paul-schwendenman/magic-log-ui/internal/config"
	"github.com/spf13/cobra"
)

var presetsCmd = &cobra.Command{
	Use:   "presets",
	Short: "List available regex and jq presets",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Println("❌ Failed to load config:", err)
			return
		}

		fmt.Println("📜 Available regexp presets:")
		for name := range cfg.RegexPresets {
			fmt.Printf("  - %s\n", name)
		}

		fmt.Println("📜 Available jq presets:")
		for name := range cfg.JQPresets {
			fmt.Printf("  - %s\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(presetsCmd)
}
