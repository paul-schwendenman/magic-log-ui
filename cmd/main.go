// cmd/main.go
package main

import (
	"bufio"
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/marcboeker/go-duckdb"
)

var version = "dev"

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
	var dbFile string
	var openBrowser bool
	var port int
	var showVersion bool
	flag.StringVar(&dbFile, "db-file", "", "Path to a DuckDB database file. Leave empty for in-memory.")
	flag.BoolVar(&openBrowser, "launch", false, "Automatically open the UI in the default web browser.")
	flag.IntVar(&port, "port", 3000, "Port to serve the web UI on.")
	flag.BoolVar(&showVersion, "version", false, "Print the version and exit.")
	flag.Parse()

	if showVersion {
		fmt.Println("magic-log version:", version)
		return
	}

	if dbFile == "" {
		log.Println("🧠 Using in-memory DuckDB")
	} else {
		log.Println("💾 Using persistent DuckDB:", dbFile)
	}

	var err error
	db, err = sql.Open("duckdb", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS logs (
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
	go http.ListenAndServe(":"+strconv.Itoa(port), nil)
	log.Printf("🌐 Serving UI at http://localhost:%d\n", port)

	// Conditionally open browser (macOS/Linux/Windows)
	if openBrowser {
		go func() {
			url := fmt.Sprintf("http://localhost:%d", port)
			var cmd *exec.Cmd
			if _, err := exec.LookPath("open"); err == nil {
				cmd = exec.Command("open", url) // macOS
			} else if _, err := exec.LookPath("xdg-open"); err == nil {
				cmd = exec.Command("xdg-open", url) // Linux
			} else if _, err := exec.LookPath("rundll32"); err == nil {
				cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url) // Windows
			} else {
				log.Println("⚠️ No supported method to open browser found")
				return
			}
			if err := cmd.Start(); err != nil {
				log.Println("⚠️ Unable to open browser:", err)
			}
		}()
	}

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
		log.Println("📭 STDIN closed — no longer receiving logs")
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
	log.Println("📡 Incoming WebSocket connection...")
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WS upgrade failed:", err)
		return
	}
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()
	log.Println("✅ WebSocket connected")

	// Send backlog
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

	// Detect disconnect
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

func broadcast(entry LogEntry) {
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

func mustJson(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
