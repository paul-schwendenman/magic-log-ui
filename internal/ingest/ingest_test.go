package ingest_test

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/paul-schwendenman/magic-log-ui/internal/ingest"
)

func setupTestDB(t *testing.T) (*sql.DB, *sql.Stmt, context.Context) {
	t.Helper()
	db, err := sql.Open("duckdb", "")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = db.Exec(`CREATE TABLE logs (trace_id TEXT, level TEXT, message TEXT, raw JSON, created_at TIMESTAMP)`)
	stmt, err := db.Prepare("INSERT INTO logs VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		t.Fatal(err)
	}
	return db, stmt, context.Background()
}

func queryMessageByTraceID(t *testing.T, db *sql.DB, traceID string) string {
	t.Helper()
	row := db.QueryRow(`SELECT message FROM logs WHERE trace_id = ?`, traceID)
	var msg string
	if err := row.Scan(&msg); err != nil {
		t.Fatalf("Failed to query: %v", err)
	}
	return msg
}

func TestIngest_JSON(t *testing.T) {
	db, stmt, ctx := setupTestDB(t)
	defer db.Close()

	jsonLog := `{"trace_id":"json123","level":"info","message":"json test"}`
	input := strings.NewReader(jsonLog + "\n")

	go ingest.Start(input, stmt, "json", "", ctx)
	time.Sleep(100 * time.Millisecond)

	msg := queryMessageByTraceID(t, db, "json123")
	if msg != "json test" {
		t.Errorf("Expected 'json test', got %s", msg)
	}
}

func TestIngest_TextWithRegex(t *testing.T) {
	db, stmt, ctx := setupTestDB(t)
	defer db.Close()

	textLog := `[INFO] 2025-04-23 10:00:00 service started`
	regex := `\[(?P<level>\w+)] (?P<timestamp>[^ ]+ [^ ]+) (?P<message>.+)`
	input := strings.NewReader(textLog + "\n")

	go ingest.Start(input, stmt, "text", regex, ctx)
	time.Sleep(100 * time.Millisecond)

	row := db.QueryRow(`SELECT message, level FROM logs`)
	var msg, level string
	if err := row.Scan(&msg, &level); err != nil {
		t.Fatalf("Failed to query: %v", err)
	}
	if msg != "service started" {
		t.Errorf("Expected 'service started', got %s", msg)
	}
	if level != "INFO" {
		t.Errorf("Expected level 'INFO', got %s", level)
	}
}

func TestIngest_InvalidJSONFallback(t *testing.T) {
	db, stmt, ctx := setupTestDB(t)
	defer db.Close()

	badJSON := `{ this is not valid json`
	input := strings.NewReader(badJSON + "\n")

	go ingest.Start(input, stmt, "json", "", ctx)
	time.Sleep(100 * time.Millisecond)

	row := db.QueryRow(`SELECT level, message FROM logs`)
	var level, msg string
	if err := row.Scan(&level, &msg); err != nil {
		t.Fatalf("Failed to query: %v", err)
	}
	if level != "raw" {
		t.Errorf("Expected level 'raw', got %s", level)
	}
	if !strings.Contains(msg, "this is not valid") {
		t.Errorf("Expected fallback message, got %s", msg)
	}
}

func TestIngest_RegexNoMatchFallback(t *testing.T) {
	db, stmt, ctx := setupTestDB(t)
	defer db.Close()

	line := `does not match`
	regex := `\[(?P<level>\w+)] (?P<ts>\S+ \S+) (?P<msg>.+)`
	input := bytes.NewReader([]byte(line + "\n"))

	go ingest.Start(input, stmt, "text", regex, ctx)
	time.Sleep(100 * time.Millisecond)

	row := db.QueryRow(`SELECT level, message FROM logs`)
	var level, msg string
	if err := row.Scan(&level, &msg); err != nil {
		t.Fatalf("Failed to query: %v", err)
	}
	if level != "raw" {
		t.Errorf("Expected fallback level 'raw', got %s", level)
	}
	if msg != "does not match" {
		t.Errorf("Expected fallback message 'does not match', got %s", msg)
	}
}
