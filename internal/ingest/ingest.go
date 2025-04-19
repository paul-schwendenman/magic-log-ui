package ingest

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/paul-schwendenman/magic-log-ui/internal/server/handlers"
	"github.com/paul-schwendenman/magic-log-ui/internal/shared"
)

func Start(input io.Reader, stmt *sql.Stmt, ctx context.Context) {
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()

		var entry shared.LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			entry = shared.LogEntry{
				"message": line,
				"level":   "raw",
			}
		}

		raw := string(shared.MustJson(entry))
		createdAt := time.Now().UTC()

		_, err := stmt.ExecContext(
			ctx,
			entry["trace_id"],
			entry["level"],
			entry["message"],
			raw,
			createdAt,
		)
		if err != nil {
			log.Printf("❌ Failed to insert log: %v", err)
			continue
		}

		entry["created_at"] = createdAt
		handlers.Broadcast(entry)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("⚠️ Error while scanning input: %v", err)
	} else {
		log.Println("📭 STDIN closed — no longer receiving logs")
	}
}
