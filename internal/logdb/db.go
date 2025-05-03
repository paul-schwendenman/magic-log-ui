package logdb

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/marcboeker/go-duckdb"
)

func MustInit(path string, ctx context.Context) *sql.DB {
	db, err := sql.Open("duckdb", path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS logs (
			id UUID PRIMARY KEY DEFAULT uuid(),
			timestamp TIMESTAMP,
			level TEXT,
			trace_id TEXT,
			message TEXT,
			raw_log TEXT,
			parsed_log JSON,
			log JSON,
			created_at TIMESTAMP DEFAULT current_timestamp,
			log_format TEXT,
			regex_pattern TEXT,
			jq_filter TEXT,
			csv_headers TEXT,
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_created_at ON logs(created_at);
		CREATE INDEX IF NOT EXISTS idx_timestamp ON logs(timestamp);
		CREATE INDEX IF NOT EXISTS idx_trace_id ON logs(trace_id);
		CREATE INDEX IF NOT EXISTS idx_level ON logs(level);
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func MustPrepareInsert(db *sql.DB, ctx context.Context) *sql.Stmt {
	stmt, err := db.PrepareContext(ctx, `
	INSERT INTO logs (
	  id,
	  trace_id,
	  level,
	  message,
	  raw_log,
	  parsed_log,
	  log,
	  created_at,
	  timestamp,
	  log_format,
	  regex_pattern,
	  jq_filter,
	  csv_headers
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
  `)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

func StartAutoAnalyze(db *sql.DB, ctx context.Context) {
	go func() {
		log.Println("🧠 Auto-analyze background job started")

		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, err := db.ExecContext(ctx, `ANALYZE logs`)
				if err != nil {
					log.Printf("⚠️ Failed to analyze logs table: %v", err)
				} else {
					log.Println("🧠 Refreshed statistics on logs table")
				}
			case <-ctx.Done():
				log.Println("🛑 Stopping auto-analyze")
				return
			}
		}
	}()
}
