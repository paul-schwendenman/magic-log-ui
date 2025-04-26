package logdb

import (
	"context"
	"database/sql"
	"log"

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
			regex_pattern TEXT,
			jq_filter TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_created_at ON logs(created_at);
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func MustPrepareInsert(db *sql.DB, ctx context.Context) *sql.Stmt {
	stmt, err := db.PrepareContext(ctx, `
	INSERT INTO logs (
	  trace_id,
	  level,
	  message,
	  raw_log,
	  parsed_log,
	  log,
	  created_at,
	  timestamp,
	  regex_pattern,
	  jq_filter
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
  `)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}
