package handlers_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/paul-schwendenman/magic-log-ui/internal/server/handlers"
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

func TestQueryHandler_Success(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	db.ExecContext(ctx, `INSERT INTO logs VALUES (NOW(), 't1', 'info', 'hello test', '{"msg":"ok"}')`)

	req := httptest.NewRequest("GET", "/query?q="+url.QueryEscape("SELECT * FROM logs"), nil)
	w := httptest.NewRecorder()
	handler := handlers.QueryHandler(db, ctx)
	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
	}
	body := w.Body.String()
	if !strings.Contains(body, "hello test") {
		t.Errorf("Expected query to return data, got: %s", body)
	}
}
