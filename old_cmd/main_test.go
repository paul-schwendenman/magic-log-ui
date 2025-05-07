package main

import (
	"database/sql"
	"io"
	"os"
	"testing"
	"time"

	_ "github.com/marcboeker/go-duckdb"
)

func TestRunStartsServer(t *testing.T) {
	r, w, _ := os.Pipe()
	defer r.Close()
	defer w.Close()

	// Replace stdin temporarily
	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	// Write a test log line after a small delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		io.WriteString(w, `{"timestamp":"2025-04-11T10:00:00Z","trace_id":"test","level":"info","message":"hello"}`+"\n")
	}()

	go func() {
		Run(Config{
			DBFile: ":memory:",
			Port:   3999, // avoid conflict
			Launch: false,
		})
	}()

	// Give time for server + ingestion to boot
	time.Sleep(300 * time.Millisecond)

	// Verify DB content (manually open DuckDB and check logs)
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// NOTE: We used a separate db here, so this just validates no crash.
}
