package ingest

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"

	"github.com/paul-schwendenman/magic-log-ui/internal/server/handlers"
	"github.com/paul-schwendenman/magic-log-ui/internal/shared"
)

func Start(input io.Reader, stmt *sql.Stmt, ctx context.Context) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		var entry shared.LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			entry = shared.LogEntry{"message": line, "level": "raw"}
		}
		raw := string(shared.MustJson(entry))
		stmt.ExecContext(ctx, entry["timestamp"], entry["trace_id"], entry["level"], entry["message"], raw)
		handlers.Broadcast(entry)
	}
	log.Println("📭 STDIN closed — no longer receiving logs")
}
