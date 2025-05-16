package agent

import (
	"context"
	"log"
	"time"

	"github.com/paul-schwendenman/magic-log-ui/internal/rpc"
	"google.golang.org/grpc"
)

type Streamer struct {
	Client     rpc.LogIngestorClient
	AgentName  string
	ServerAddr string

	FlushInterval time.Duration
	RetryDelay    time.Duration
}

func NewStreamer(conn *grpc.ClientConn, agentName string) *Streamer {
	return &Streamer{
		Client:        rpc.NewLogIngestorClient(conn),
		AgentName:     agentName,
		FlushInterval: 2 * time.Second,
		RetryDelay:    2 * time.Second,
	}
}

func (s *Streamer) Run(ctx context.Context, logChan <-chan string) error {
	var buffer []string
	ticker := time.NewTicker(s.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case line, ok := <-logChan:
			if !ok {
				return s.flushWithRetry(ctx, buffer)
			}
			buffer = append(buffer, line)

		case <-ticker.C:
			if len(buffer) > 0 {
				if err := s.flushWithRetry(ctx, buffer); err != nil {
					log.Printf("⚠️  Final flush error: %v", err)
				}
				buffer = nil
			}
		}
	}
}

func (s *Streamer) flushWithRetry(ctx context.Context, logs []string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		stream, err := s.Client.SendLogs(ctx)
		if err != nil {
			log.Printf("❌ Failed to open stream: %v", err)
			time.Sleep(s.RetryDelay)
			continue
		}

		err = stream.Send(&rpc.LogBatch{
			AgentName: s.AgentName,
			JsonLogs:  logs,
		})
		if err != nil {
			log.Printf("⚠️  Failed to send logs: %v", err)
			time.Sleep(s.RetryDelay)
			continue
		}

		return nil
	}
}
