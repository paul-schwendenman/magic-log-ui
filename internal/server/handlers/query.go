package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const maxLimit = 1000
const defaultLimit = 100

func QueryHandler(db *sql.DB, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userQuery := r.URL.Query().Get("q")
		if userQuery == "" {
			http.Error(w, "missing q param", http.StatusBadRequest)
			return
		}

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit <= 0 {
			limit = defaultLimit
		}
		if limit > maxLimit {
			limit = maxLimit
		}
		offset := page * limit

		userQuery = strings.TrimSuffix(strings.TrimSpace(userQuery), ";")

		wrappedQuery := fmt.Sprintf(`
			WITH q AS (%s)
			SELECT * FROM q
			LIMIT %d OFFSET %d
		`, userQuery, limit+1, offset)

		rows, err := db.QueryContext(ctx, wrappedQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()

		cols, _ := rows.Columns()
		results := []map[string]any{}

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

		hasNext := len(results) > limit
		if hasNext {
			results = results[:limit]
		}

		resp := map[string]any{
			"data": results,
			"meta": map[string]any{
				"hasNextPage":     hasNext,
				"hasPreviousPage": page > 0,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
