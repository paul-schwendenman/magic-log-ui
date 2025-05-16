package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/paul-schwendenman/magic-log-ui/internal/agent"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run in agent mode (read logs from stdin and forward to a server)",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cfg := agent.Config{
			Name:       viper.GetString("agent.name"),
			ServerAddr: viper.GetString("agent.server_addr"),
			Format:     viper.GetString("agent.format"),
			Regex:      viper.GetString("agent.regex"),
			JQ:         viper.GetString("agent.jq"),
			CSVFields:  viper.GetString("agent.csv_fields"),
			HasHeader:  viper.GetBool("agent.csv_has_header"),
			Echo:       viper.GetBool("agent.echo"),
		}

		if cfg.Name == "" {
			cfg.Name = agent.GenerateDefaultName()
		}

		if err := agent.Run(ctx, cfg, os.Stdin); err != nil {
			log.Fatalf("❌ agent failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)

	flags := agentCmd.Flags()
	flags.String("name", "", "Agent name (default: auto-generated)")
	flags.String("server-addr", "", "gRPC address of the server (e.g., localhost:50051)")
	flags.String("format", "json", "Log format (json, csv, text)")
	flags.String("regex", "", "Regex pattern for parsing text logs")
	flags.String("regex-preset", "", "Regex preset to use")
	flags.String("jq", "", "JQ query to transform logs")
	flags.String("jq-preset", "", "JQ preset to use")
	flags.String("csv-fields", "", "Comma-separated CSV field names")
	flags.Bool("csv-has-header", false, "Expect a header row in CSV input")
	flags.Bool("echo", false, "Echo parsed logs to stdout")

	viper.BindPFlags(flags)
	viper.SetDefault("agent.format", "json")
}
