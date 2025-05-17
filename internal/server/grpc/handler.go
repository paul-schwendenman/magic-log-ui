package grpc

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"

	"github.com/paul-schwendenman/magic-log-ui/internal/ingest"
	"github.com/paul-schwendenman/magic-log-ui/internal/rpc"
	"github.com/paul-schwendenman/magic-log-ui/internal/server/handlers"
	"github.com/paul-schwendenman/magic-log-ui/internal/shared"
)

type Server struct {
	DBStmt *sql.Stmt
	rpc.UnimplementedLogIngestorServer
}

func (s *Server) SendLogs(stream rpc.LogIngestor_SendLogsServer) error {
	for {
		batch, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&rpc.Ack{Message: "OK"})
		}
		if err != nil {
			log.Printf("❌ gRPC receive failed: %v", err)
			return err
		}

		for _, jsonStr := range batch.JsonLogs {
			var entry shared.LogEntry
			if err := json.Unmarshal([]byte(jsonStr), &entry); err != nil {
				log.Printf("⚠️ Failed to parse log JSON: %v", err)
				continue
			}

			entry["agent_name"] = batch.AgentName

			if err := ingest.InsertTransformed(s.DBStmt, jsonStr, entry); err != nil {
				log.Printf("❌ DB insert failed: %v", err)
			}

			handlers.Broadcast(entry)
		}
	}
}
