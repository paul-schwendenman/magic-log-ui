package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const maxLimit = 1000
const defaultLimit = 100

func QueryHandler(db *sql.DB, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userQuery := r.URL.Query().Get("q")
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")

		if userQuery == "" {
			http.Error(w, "missing q param", http.StatusBadRequest)
			return
		}

		page, _ := strconv.Atoi(pageStr)
		limit, err := strconv.Atoi(limitStr)

		var timeFilter string

		if fromStr != "" && toStr != "" {
			from, err1 := time.Parse(time.RFC3339, fromStr)
			to, err2 := time.Parse(time.RFC3339, toStr)
			if err1 == nil && err2 == nil {
				timeFilter = fmt.Sprintf(
					"WHERE created_at BETWEEN TIMESTAMP '%s' AND TIMESTAMP '%s'",
					from.Format(time.RFC3339),
					to.Format(time.RFC3339),
				)
			}
		} else if fromStr != "" {
			from, err := time.Parse(time.RFC3339, fromStr)
			if err == nil {
				timeFilter = fmt.Sprintf(
					"WHERE created_at > TIMESTAMP '%s'",
					from.Format(time.RFC3339),
				)
			}
		} else if toStr != "" {
			to, err := time.Parse(time.RFC3339, toStr)
			if err == nil {
				timeFilter = fmt.Sprintf(
					"WHERE created_at < TIMESTAMP '%s'",
					to.Format(time.RFC3339),
				)
			}
		}

		if err != nil || limit <= 0 {
			limit = defaultLimit
		}
		if limit > maxLimit {
			limit = maxLimit
		}
		offset := page * limit

		userQuery = strings.TrimSuffix(strings.TrimSpace(userQuery), ";")

		countQuery := fmt.Sprintf(`WITH q AS (%s) SELECT COUNT(*) FROM q %s`, userQuery, timeFilter)
		var totalRows int
		err = db.QueryRowContext(ctx, countQuery).Scan(&totalRows)
		if err != nil {
			http.Error(w, "count query failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

		safeQuery := fmt.Sprintf(`
			WITH q AS (%s)
			SELECT * FROM q
			%s
			ORDER BY created_at DESC
			LIMIT %d OFFSET %d
		`, userQuery, timeFilter, limit+1, offset)

		rows, err := db.QueryContext(ctx, safeQuery)
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
				"totalPages":      totalPages,
				"page":            page,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
