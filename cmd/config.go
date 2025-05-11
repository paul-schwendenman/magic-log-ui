/*
Copyright © 2025 Paul Schwendenman
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/paul-schwendenman/magic-log-ui/internal/app"
	"github.com/paul-schwendenman/magic-log-ui/internal/config"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long: `View, edit, or validate your magic-log configuration file.

This command allows you to inspect and modify config keys, such as setting presets,
changing log formats, or specifying database settings.

Examples:
  magic-log config get port
  magic-log config set jq_preset simple '{message: .msg}'
  magic-log config unset regex
  magic-log config validate`,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Long: `Retrieves the value of a configuration key from the config file.

Keys can refer to top-level fields (like 'port') or nested presets (like 'regex_presets.apache').

Examples:
  magic-log config get port
  magic-log config get jq_preset
  magic-log config get regex_presets.apache`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		val, err := app.GetConfigValue(args[0])
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
	Long: `Sets a configuration key to the given value.

Keys can refer to top-level settings like 'port', 'launch', or 'log_format',
or to nested preset keys like 'regex_presets.apache'.

This command will validate the input and update your config file accordingly.

Examples:
  magic-log config set port 4000
  magic-log config set launch true
  magic-log config set regex_presets.apache '(?P<ip>\\S+) ...'
  magic-log config set jq_preset simple '{message: .msg}'`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.SetConfigValue(args[0], args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "❌", err)
			os.Exit(1)
		}
	},
}

var configUnsetCmd = &cobra.Command{
	Use:   "unset <key>",
	Short: "Unset a config value",
	Long: `Removes a configuration key from the config file.

This can be used to delete top-level fields or remove individual preset values.

Examples:
  magic-log config unset port
  magic-log config unset jq_preset
  magic-log config unset regex_presets.apache`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.UnsetConfigValue(args[0]); err != nil {
			fmt.Fprintln(os.Stderr, "❌", err)
			os.Exit(1)
		}
	},
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration",
	Long: `Checks the configuration file for syntax and logic errors.

This command loads your config file and runs validation on top-level settings
as well as any defined regex or jq presets.

Examples:
  magic-log config validate`,
	Args: cobra.NoArgs,
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
	configGetCmd.ValidArgsFunction = app.CompleteConfigKeys
	configSetCmd.ValidArgsFunction = app.CompleteKnownConfigKeys
	configUnsetCmd.ValidArgsFunction = app.CompleteConfigUnsetKeys

	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
	configCmd.AddCommand(configValidateCmd)
	rootCmd.AddCommand(configCmd)
}
