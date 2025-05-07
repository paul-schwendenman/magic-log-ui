/*
Copyright © 2025 Paul Schwendenman
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("magic-log version:", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
