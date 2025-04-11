package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/paul-schwendenman/magic-log-ui/internal/shared"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

func WebSocketHandler(db *sql.DB, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("📡 Incoming WebSocket connection...")

		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("❌ WS upgrade failed:", err)
			return
		}

		clientsMu.Lock()
		clients[conn] = true
		clientsMu.Unlock()
		log.Println("✅ WebSocket connected")

		// Send last log (for now)
		go func(c *websocket.Conn) {
			rows, err := db.QueryContext(ctx, `SELECT * FROM logs ORDER BY timestamp DESC LIMIT 1`)
			if err == nil {
				defer rows.Close()
				cols, _ := rows.Columns()
				for rows.Next() {
					vals := make([]any, len(cols))
					ptrs := make([]any, len(cols))
					for i := range ptrs {
						ptrs[i] = &vals[i]
					}
					rows.Scan(ptrs...)
					entry := map[string]any{}
					for i, col := range cols {
						entry[col] = vals[i]
					}
					c.WriteJSON(entry)
				}
			}
		}(conn)

		// Disconnect listener
		go func(c *websocket.Conn) {
			for {
				if _, _, err := c.NextReader(); err != nil {
					log.Println("👋 WebSocket disconnected")
					clientsMu.Lock()
					delete(clients, c)
					clientsMu.Unlock()
					c.Close()
					return
				}
			}
		}(conn)
	}
}

func Broadcast(entry shared.LogEntry) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		if err := conn.WriteJSON(entry); err != nil {
			conn.Close()
			delete(clients, conn)
			log.Println("❌ Failed to write to WebSocket, removed client")
		}
	}
}
