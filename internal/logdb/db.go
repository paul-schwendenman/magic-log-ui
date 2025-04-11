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

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS logs (
		timestamp TIMESTAMP,
		trace_id TEXT,
		level TEXT,
		message TEXT,
		raw JSON
	);`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func MustPrepareInsert(db *sql.DB, ctx context.Context) *sql.Stmt {
	stmt, err := db.PrepareContext(ctx, "INSERT INTO logs VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}
