// cmd/main.go
package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	duckdb "github.com/marcboeker/go-duckdb"
	"log"
	"net/http"
	"os"
	"sync"
)

//go:embed static/*
var staticFiles embed.FS

type LogEntry map[string]any

var (
	db        *duckdb.Conn
	logInsert *duckdb.Stmt
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

func main() {
	// Init DB
	var err error
	db, err = duckdb.Open("") // in-memory
	if err != nil {
		log.Fatal(err)
	}
	db.Exec(`CREATE TABLE logs (
		timestamp TIMESTAMP,
		trace_id TEXT,
		level TEXT,
		message TEXT,
		raw JSON
	);`)

	logInsert, _ = db.Prepare("INSERT INTO logs VALUES (?, ?, ?, ?, ?)")

	// Start HTTP server
	http.HandleFunc("/query", handleQuery)
	http.HandleFunc("/ws", handleWS)
	http.HandleFunc("/", serveStatic)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))
	go http.ListenAndServe(":3000", nil)
	fmt.Println("🌐 Serving UI at http://localhost:3000")

	// Start reading stdin
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
		logInsert.Exec(timestamp, traceID, level, msg, string(mustJson(entry)))
		broadcast(entry)
	}
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "missing q param", 400)
		return
	}
	rows, err := db.Query(q)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	cols := rows.Columns()
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
	f, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Write(f)
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
