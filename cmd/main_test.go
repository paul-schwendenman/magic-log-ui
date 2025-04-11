package main

import (
	"database/sql"
	"encoding/json"
	"testing"
)

func TestMustJson(t *testing.T) {
	in := map[string]any{"foo": "bar"}
	out := mustJson(in)

	var decoded map[string]any
	if err := json.Unmarshal(out, &decoded); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}
	if decoded["foo"] != "bar" {
		t.Errorf("Expected foo=bar, got %v", decoded["foo"])
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
