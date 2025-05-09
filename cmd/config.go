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
	Long:  ``,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Args:  cobra.ExactArgs(1),
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
	Args:  cobra.ExactArgs(2),
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
	Args:  cobra.ExactArgs(1),
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
	configGetCmd.ValidArgsFunction = completeConfigKeys
	configSetCmd.ValidArgsFunction = completeConfigKeyValues
	configUnsetCmd.ValidArgsFunction = completeConfigKeys

	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
	configCmd.AddCommand(configValidateCmd)
	rootCmd.AddCommand(configCmd)
}

func completeConfigKeys(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg, _, err := app.LoadConfigMap()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var keys []string
	for k, v := range cfg {
		switch section := v.(type) {
		case map[string]any:
			for subk := range section {
				keys = append(keys, fmt.Sprintf("%s.%s", k, subk))
			}
		default:
			keys = append(keys, k)
		}
	}

	return keys, cobra.ShellCompDirectiveNoFileComp
}

func completeConfigKeyValues(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		// completing first arg: the key
		return getAllConfigKeys(), cobra.ShellCompDirectiveNoFileComp
	}

	if len(args) == 1 {
		key := args[0]
		return suggestValuesForKey(key), cobra.ShellCompDirectiveNoFileComp
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func getAllConfigKeys() []string {
	cfg, _, err := app.LoadConfigMap()
	if err != nil {
		return nil
	}

	var keys []string
	for k, v := range cfg {
		if section, ok := v.(map[string]any); ok {
			for subk := range section {
				keys = append(keys, fmt.Sprintf("%s.%s", k, subk))
			}
		} else {
			keys = append(keys, k)
		}
	}
	return keys
}

func suggestValuesForKey(key string) []string {
	switch key {
	case "log_format":
		return []string{"json", "text"}
	case "launch", "has_csv_header":
		return []string{"true", "false"}
	case "regex_preset":
		return getKeysFromSection("regex_presets")
	case "jq_preset":
		return getKeysFromSection("jq_presets")
	default:
		return nil
	}
}

func getKeysFromSection(section string) []string {
	cfg, _, err := app.LoadConfigMap()
	if err != nil {
		return nil
	}
	sec, ok := cfg[section].(map[string]any)
	if !ok {
		return nil
	}
	var keys []string
	for k := range sec {
		keys = append(keys, k)
	}
	return keys
}
