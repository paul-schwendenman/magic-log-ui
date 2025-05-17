package app

import (
	"context"
	"io"

	"github.com/paul-schwendenman/magic-log-ui/internal/agent"
	"github.com/paul-schwendenman/magic-log-ui/internal/ingest"
	"google.golang.org/grpc"
)

func RunAgent(ctx context.Context, cfg agent.AgentConfig, input io.Reader) error {
	conn, err := grpc.Dial(cfg.ServerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	logChan := make(chan string, 100)

	go func() {
		ingest.StartForwarding(input, logChan, ingest.ParserOptions{
			LogFormat:    cfg.LogFormat,
			ParseRegex:   cfg.ParseRegex,
			JqFilter:     cfg.JqFilter,
			CSVFieldsStr: cfg.CSVFieldsStr,
			HasCSVHeader: cfg.HasCSVHeader,
			Echo:         cfg.Echo,
		})
		close(logChan)
	}()

	streamer := agent.NewStreamer(conn, cfg.Name)
	return streamer.Run(ctx, logChan)
}
