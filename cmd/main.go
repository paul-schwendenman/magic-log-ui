// cmd/main.go
package main

import (
	"bufio"
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/marcboeker/go-duckdb"
)

//go:embed static/*
var staticFiles embed.FS

type LogEntry map[string]any

var (
	db        *sql.DB
	logInsert *sql.Stmt
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	ctx       = context.Background()
)

func main() {
	var err error
	db, err = sql.Open("duckdb", "") // in-memory
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.ExecContext(ctx, `CREATE TABLE logs (
		timestamp TIMESTAMP,
		trace_id TEXT,
		level TEXT,
		message TEXT,
		raw JSON
	);`)
	if err != nil {
		log.Fatal(err)
	}

	logInsert, err = db.PrepareContext(ctx, "INSERT INTO logs VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/query", handleQuery)
	http.HandleFunc("/ws", handleWS)
	http.HandleFunc("/", serveStatic)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))
	go http.ListenAndServe(":3000", nil)
	fmt.Println("🌐 Serving UI at http://localhost:3000")

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			var entry LogEntry
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				entry = LogEntry{"message": line, "level": "raw"}
			}
			timestamp := entry["timestamp"]
			traceID := entry["trace_id"]
			level := entry["level"]
			msg := entry["message"]
			raw := string(mustJson(entry))
			logInsert.ExecContext(ctx, timestamp, traceID, level, msg, raw)
			broadcast(entry)
		}
	}()

	// Prevent exit when stdin closes
	select {}
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "missing q param", 400)
		return
	}
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var results []map[string]any
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range ptrs {
			ptrs[i] = &vals[i]
		}
		rows.Scan(ptrs...)
		row := map[string]any{}
		for i, col := range cols {
			row[col] = vals[i]
		}
		results = append(results, row)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	path := "static" + r.URL.Path
	if r.URL.Path == "/" {
		path = "static/index.html"
	}

	data, err := staticFiles.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Simple content type handling
	if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
	} else if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(path, ".html") {
		w.Header().Set("Content-Type", "text/html")
	}

	w.Write(data)
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()
}

func broadcast(entry LogEntry) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		if err := conn.WriteJSON(entry); err != nil {
			conn.Close()
			delete(clients, conn)
		}
	}
}

func mustJson(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
