package handlers

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

func StaticHandler(staticFiles embed.FS) http.HandlerFunc {
	content, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to mount embedded static files: %v", err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else {
			path = strings.TrimPrefix(path, "/")
		}

		data, err := fs.ReadFile(content, path)

		if err != nil {
			data, err = fs.ReadFile(content, "200.html")
			if err != nil {
				http.Error(w, "200.html not found", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}

		switch {
		case strings.HasSuffix(path, ".js"):
			w.Header().Set("Content-Type", "application/javascript")
		case strings.HasSuffix(path, ".css"):
			w.Header().Set("Content-Type", "text/css")
		case strings.HasSuffix(path, ".html"):
			w.Header().Set("Content-Type", "text/html")
		}

		w.Write(data)
	}
}
