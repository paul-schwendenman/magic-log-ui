package cmd

import (
	"context"
	"log"
	"os"

	"github.com/paul-schwendenman/magic-log-ui/internal/agent"
	"github.com/paul-schwendenman/magic-log-ui/internal/app"
	"github.com/paul-schwendenman/magic-log-ui/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run in agent mode (read logs from stdin and forward to a server)",
	Run: func(cmd *cobra.Command, args []string) {
		fileCfg, err := config.Load()
		if err != nil {
			log.Fatalf("❌ Failed to load config: %v", err)
		}

		parsing, err := config.ResolveParsingOptions(
			viper.GetString("regex_preset"),
			viper.GetString("regex"),
			viper.GetString("jq_preset"),
			viper.GetString("jq"),
			fileCfg,
		)
		if err != nil {
			log.Fatalf("❌ %v", err)
		}

		cfg := agent.AgentConfig{
			Name:         viper.GetString("name"),
			ServerAddr:   viper.GetString("server_addr"),
			LogFormat:    viper.GetString("log_format"),
			ParseRegex:   parsing.Regex,
			JqFilter:     parsing.JQ,
			CSVFieldsStr: viper.GetString("csv_fields"),
			HasCSVHeader: viper.GetBool("has_csv_header"),
			Echo:         viper.GetBool("echo"),
		}

		if cfg.Name == "" {
			cfg.Name = agent.GenerateAgentName()
		}

		if err := app.RunAgent(context.Background(), cfg, os.Stdin); err != nil {
			log.Fatalf("❌ Agent failed: %v", err)
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
