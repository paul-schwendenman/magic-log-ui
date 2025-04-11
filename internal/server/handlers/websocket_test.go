package handlers_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"github.com/gorilla/websocket"
	_ "github.com/marcboeker/go-duckdb"
	"github.com/paul-schwendenman/magic-log-ui/internal/server/handlers"
	"github.com/paul-schwendenman/magic-log-ui/internal/shared"
)

func TestWebSocketBroadcast(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(handlers.WebSocketHandler(db, ctx)))
	defer srv.Close()

	u := "ws" + strings.TrimPrefix(srv.URL, "http")
	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer c.Close()

	stmt, _ := db.Prepare("INSERT INTO logs VALUES (?, ?, ?, ?, ?)")
	entry := shared.LogEntry{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"trace_id":  "t2",
		"level":     "info",
		"message":   "websocket test",
		"raw":       `{"ws":"ok"}`,
	}
	raw := string(shared.MustJson(entry))
	stmt.ExecContext(ctx, entry["timestamp"], entry["trace_id"], entry["level"], entry["message"], raw)
	handlers.Broadcast(entry)

	time.Sleep(100 * time.Millisecond)

	_, msg, err := c.ReadMessage()
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if !strings.Contains(string(msg), "websocket test") {
		t.Errorf("Expected broadcast message, got: %s", msg)
	}
}

func TestWebSocketMalformedMessage(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(handlers.WebSocketHandler(db, ctx)))
	defer srv.Close()

	u := "ws" + strings.TrimPrefix(srv.URL, "http")
	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte("{invalid json"))
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
}

func TestWebSocketDisconnect(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(handlers.WebSocketHandler(db, ctx)))
	defer srv.Close()

	u := "ws" + strings.TrimPrefix(srv.URL, "http")
	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	c.Close()
	time.Sleep(200 * time.Millisecond)
}
