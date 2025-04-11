package server_test

import (
	"context"
	"database/sql"
	"net/http"
	"testing"

	"embed"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/paul-schwendenman/magic-log-ui/internal/server"
)

//go:embed fake_static/*
var fakeStatic embed.FS

func TestServer_StartBasicRoutes(t *testing.T) {
	db, _ := sql.Open("duckdb", "")
	ctx := context.Background()

	// Dry-run to ensure nothing panics
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Server panicked: %v", r)
			}
		}()
		server.Start(34567, fakeStatic, db, ctx)
	}()

	// Give server time to boot
	resp, err := http.Get("http://localhost:34567/")
	if err == nil {
		defer resp.Body.Close()
		t.Logf("Startup request succeeded with code %d", resp.StatusCode)
	}
}
