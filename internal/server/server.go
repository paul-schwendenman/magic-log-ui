package server

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/paul-schwendenman/magic-log-ui/internal/server/handlers"
)

func Start(port int, staticFiles embed.FS, db *sql.DB, ctx context.Context) {
	http.HandleFunc("/query", handlers.QueryHandler(db, ctx))
	http.HandleFunc("/ws", handlers.WebSocketHandler(db, ctx))
	http.HandleFunc("/", handlers.StaticHandler(staticFiles))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("🌐 Serving UI at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
