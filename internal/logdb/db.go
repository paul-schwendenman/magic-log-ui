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
			trace_id TEXT,
			level TEXT,
			message TEXT,
			raw JSON,
			created_at TIMESTAMP DEFAULT current_timestamp
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
		INSERT INTO logs (trace_id, level, message, raw, created_at)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}
