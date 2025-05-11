/*
Copyright © 2025 Paul Schwendenman
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and exit",
	Long:  "Print the current version of magic-log and exit.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("magic-log version:", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
