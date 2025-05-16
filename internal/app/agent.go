package app

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/paul-schwendenman/magic-log-ui/internal/agent"
	"github.com/paul-schwendenman/magic-log-ui/internal/ingest"
	"github.com/paul-schwendenman/magic-log-ui/internal/rpc"
	"google.golang.org/grpc"
)

func RunAgent(ctx context.Context, cfg agent.AgentConfig, input io.Reader) error {
	conn, err := grpc.Dial(cfg.ServerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := rpc.NewLogIngestorClient(conn)
	stream, err := client.SendLogs(ctx)
	if err != nil {
		return err
	}

	// Reuse ingest parsing logic
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

	// Forward logs to server in batches
	var buffer []string
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-logChan:
			if !ok {
				return nil // Done
			}
			buffer = append(buffer, line)
		case <-ticker.C:
			if len(buffer) == 0 {
				continue
			}
			err := stream.Send(&rpc.LogBatch{
				AgentName: cfg.Name,
				JsonLogs:  buffer,
			})
			if err != nil {
				log.Printf("⚠️ Failed to send logs: %v", err)
			}
			buffer = nil
		}
	}
}
