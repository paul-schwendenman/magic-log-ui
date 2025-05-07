/*
Copyright © 2025 Paul Schwendenman

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/paul-schwendenman/magic-log-ui/internal/config"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long: ``,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		val, err := config.GetConfigValue(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "❌", err)
			os.Exit(1)
		}
		if val == "" {
			os.Exit(1)
		}
		fmt.Println(val)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.SetConfigValue(args[0], args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "❌", err)
			os.Exit(1)
		}
	},
}

var configUnsetCmd = &cobra.Command{
	Use:   "unset <key>",
	Short: "Unset a config value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.UnsetConfigValue(args[0]); err != nil {
			fmt.Fprintln(os.Stderr, "❌", err)
			os.Exit(1)
		}
	},
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, "❌ Failed to load config:", err)
			os.Exit(1)
		}
		if errs := cfg.Validate(); len(errs) > 0 {
			fmt.Fprintln(os.Stderr, "❌ Config is invalid:")
			for _, e := range errs {
				fmt.Fprintln(os.Stderr, "   -", e)
			}
			os.Exit(1)
		}
		fmt.Println("✅ Config is valid")
	},
}


func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
	configCmd.AddCommand(configValidateCmd)
	rootCmd.AddCommand(configCmd)
}
