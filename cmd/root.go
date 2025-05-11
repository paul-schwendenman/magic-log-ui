/*
Copyright © 2025 Paul Schwendenman

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "magic-log",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverCmd.Run(cmd, args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.magiclogrc)")
	rootCmd.Flags().String("db-file", "", "Path to a DuckDB database file")
	rootCmd.Flags().Int("port", 3000, "Port to serve the web UI on")
	rootCmd.Flags().Bool("launch", false, "Open the UI in a browser")
	rootCmd.Flags().Bool("echo", false, "Echo parsed stdin input to stdout")
	rootCmd.Flags().Bool("no-auto-analyze", false, "Disable automatic ANALYZE of logs table")
	rootCmd.Flags().String("log-format", "json", "Log format: json, csv or plain text")
	rootCmd.Flags().String("regex", "", "Custom regex to parse logs (use with text format)")
	rootCmd.Flags().String("regex-preset", "", "Regex preset to use")
	rootCmd.Flags().String("jq", "", "A jq expression to apply to parsed logs")
	rootCmd.Flags().String("jq-preset", "", "jq preset to use")
	rootCmd.Flags().String("csv-fields", "", "Comma-separated field names for CSV logs")
	rootCmd.Flags().Bool("has-csv-header", true, "Whether CSV logs include a header row")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".magiclogrc")
	}

	viper.BindPFlags(rootCmd.Flags())
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Config not loaded: %v\n", err)
	}
}
