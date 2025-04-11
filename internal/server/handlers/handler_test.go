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

func TestQueryHandlerSuccess(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	db.ExecContext(ctx, `INSERT INTO logs VALUES (NOW(), 't1', 'info', 'test', '{"msg":"ok"}')`)

	r := httptest.NewRequest("GET", "/query?q="+url.QueryEscape("SELECT * FROM logs"), nil)
	w := httptest.NewRecorder()
	handlers.QueryHandler(db, ctx)(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "test") {
		t.Errorf("Expected log in response, got: %s", w.Body.String())
	}
}

func TestQueryHandlerMissingParam(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	r := httptest.NewRequest("GET", "/query", nil)
	w := httptest.NewRecorder()
	handlers.QueryHandler(db, ctx)(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "missing q param") {
		t.Errorf("Expected error in response, got: %s", w.Body.String())
	}
}
