package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("failed to open duckdb: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE logs (
		timestamp TIMESTAMP,
		trace_id TEXT,
		level TEXT,
		message TEXT,
		raw JSON
	)`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	return db
}

func TestMustJson(t *testing.T) {
	input := map[string]any{"hello": "world"}
	out := mustJson(input)

	var decoded map[string]any
	if err := json.Unmarshal(out, &decoded); err != nil {
		t.Fatalf("Expected valid JSON: %v", err)
	}
	if decoded["hello"] != "world" {
		t.Errorf("Expected 'world', got %v", decoded["hello"])
	}
}

func TestQueryLogs(t *testing.T) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE logs (timestamp TIMESTAMP, trace_id TEXT, level TEXT, message TEXT, raw JSON)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO logs VALUES (NOW(), 'abc123', 'info', 'test message', '{"test":true}')`)
	if err != nil {
		t.Fatal(err)
	}

	row := db.QueryRow(`SELECT level FROM logs WHERE trace_id = 'abc123'`)
	var level string
	if err := row.Scan(&level); err != nil {
		t.Fatal(err)
	}
	if level != "info" {
		t.Errorf("Expected level=info, got %v", level)
	}
}

func TestHandleQuerySuccess(t *testing.T) {
	db = setupTestDB(t)
	ctx = context.Background()
	_, err := db.ExecContext(ctx, `INSERT INTO logs VALUES (NOW(), 't1', 'info', 'test', '{"msg":"ok"}')`)
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	r := httptest.NewRequest("GET", "/query?q="+url.QueryEscape("SELECT * FROM logs"), nil)
	w := httptest.NewRecorder()
	handleQuery(w, r)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", res.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	if !strings.Contains(buf.String(), "test") {
		t.Errorf("Response does not contain expected content: %v", buf.String())
	}
}

func TestHandleQueryMissingParam(t *testing.T) {
	r := httptest.NewRequest("GET", "/query", nil)
	w := httptest.NewRecorder()
	handleQuery(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "missing q param") {
		t.Errorf("Unexpected error message: %v", w.Body.String())
	}
}
