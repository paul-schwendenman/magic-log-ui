package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func QueryHandler(db *sql.DB, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userQuery := r.URL.Query().Get("q")
		if userQuery == "" {
			http.Error(w, "missing q param", http.StatusBadRequest)
			return
		}

		userQuery = strings.TrimSuffix(strings.TrimSpace(userQuery), ";")

		safeQuery := fmt.Sprintf(`
			WITH q AS (%s)
			SELECT * FROM q
			LIMIT 1000
		`, userQuery)

		rows, err := db.QueryContext(ctx, safeQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
}
