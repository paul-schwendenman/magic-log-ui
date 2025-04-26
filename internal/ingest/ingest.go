package ingest

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	"github.com/paul-schwendenman/magic-log-ui/internal/jqfilter"
	"github.com/paul-schwendenman/magic-log-ui/internal/shared"
)

func Start(input io.Reader, stmt *sql.Stmt, logFormat string, parseRegexStr string, jqQuery string, echo bool, ctx context.Context) {
	scanner := bufio.NewScanner(input)

	// Initialize jq filter if any
	jqfilter.Init(jqQuery)

	for scanner.Scan() {
		rawLine := scanner.Text()

		if echo {
			fmt.Println(rawLine)
		}

		var parsed shared.LogEntry
		var parsedLogJson []byte
		var err error
		regexUsed := false

		// Attempt to parse the line
		if logFormat == "json" {
			err = json.Unmarshal([]byte(rawLine), &parsed)
		}

		if err != nil || logFormat == "text" {
			// fallback to regex if given
			parsed, err = parseWithRegex(rawLine, parseRegexStr)
			if err != nil {
				// totally fallback
				parsed = shared.LogEntry{
					"message": rawLine,
					"level":   "raw",
				}
			} else {
				regexUsed = true
			}
		}

		// Save parsed_log
		parsedLogJson = shared.MustJson(parsed)

		// Apply jq transformation
		// transformed := shared.LogEntry(jqfilter.Apply(logEntryToStringMap(parsed)))
		transformed := mapToLogEntry(jqfilter.Apply(logEntryToStringMap(parsed)))

		// Save final log
		finalLogJson := shared.MustJson(transformed)

		// Extract useful fields
		traceID, _ := safeString(transformed, "trace_id")
		level, _ := safeString(transformed, "level")
		message, _ := safeString(transformed, "message")

		if level == "" {
			level = "info"
		}
		if message == "" {
			message = "(no message)"
		}

		// Extract timestamp from log if exists
		timestamp := time.Now().UTC()
		if ts, ok := safeString(transformed, "timestamp"); ok {
			parsedTs, err := time.Parse(time.RFC3339, ts)
			if err == nil {
				timestamp = parsedTs
			}
		}

		regexPattern := ""
		if regexUsed {
			regexPattern = parseRegexStr
		}

		// Insert into DB
		_, err = stmt.ExecContext(
			ctx,
			traceID,
			level,
			message,
			rawLine,
			string(parsedLogJson),
			string(finalLogJson),
			time.Now().UTC(), // created_at
			timestamp,
			nullify(regexPattern),
			nullify(jqQuery),
		)
		if err != nil {
			log.Printf("❌ Failed to insert log: %v", err)
		}

		// Broadcast (optional)
		transformed["created_at"] = timestamp
		// handlers.Broadcast(transformed) // Uncomment if you broadcast after ingest
	}

	if err := scanner.Err(); err != nil {
		log.Printf("⚠️ Error while scanning input: %v", err)
	} else {
		log.Println("📭 STDIN closed — no longer receiving logs")
	}
}

func parseWithRegex(line, pattern string) (shared.LogEntry, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	match := re.FindStringSubmatch(line)
	if match == nil {
		return nil, fmt.Errorf("no match")
	}

	result := shared.LogEntry{}
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result, nil
}

func mapToLogEntry(m map[string]string) shared.LogEntry {
	e := make(shared.LogEntry, len(m))
	for k, v := range m {
		e[k] = v
	}
	return e
}

func logEntryToStringMap(entry shared.LogEntry) map[string]string {
	m := make(map[string]string, len(entry))
	for k, v := range entry {
		m[k] = fmt.Sprintf("%v", v)
	}
	return m
}

func safeString(m map[string]any, key string) (string, bool) {
	v, ok := m[key]
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%v", v), true
}

func nullify(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
