package ingest_test

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"net"
	"testing"
	"time"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/paul-schwendenman/magic-log-ui/internal/shared"
)

func TestStdinIngestion(t *testing.T) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.Exec(`CREATE TABLE logs (timestamp TIMESTAMP, trace_id TEXT, level TEXT, message TEXT, raw JSON)`)

	pr, pw := net.Pipe()
	defer pr.Close()
	defer pw.Close()

	stmt, _ := db.Prepare("INSERT INTO logs VALUES (?, ?, ?, ?, ?)")
	ctx := context.Background()

	go func() {
		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			var entry shared.LogEntry
			line := scanner.Text()
			json.Unmarshal([]byte(line), &entry)
			raw := string(shared.MustJson(entry))
			stmt.ExecContext(ctx, entry["timestamp"], entry["trace_id"], entry["level"], entry["message"], raw)
		}
	}()

	entry := `{"timestamp":"2025-04-11T10:00:00Z","trace_id":"stdin123","level":"info","message":"stdin test"}`
	pw.Write([]byte(entry + "\n"))
	time.Sleep(100 * time.Millisecond)

	row := db.QueryRow(`SELECT message FROM logs WHERE trace_id = 'stdin123'`)
	var msg string
	if err := row.Scan(&msg); err != nil {
		t.Fatalf("Failed to query: %v", err)
	}
	if msg != "stdin test" {
		t.Errorf("Expected 'stdin test', got %s", msg)
	}
}
