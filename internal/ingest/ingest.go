package ingest

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"regexp"
	"time"

	"github.com/paul-schwendenman/magic-log-ui/internal/server/handlers"
	"github.com/paul-schwendenman/magic-log-ui/internal/shared"
)

func Start(input io.Reader, stmt *sql.Stmt, logFormat string, parseRegexStr string, ctx context.Context) {
	scanner := bufio.NewScanner(input)
	parseLogLine := makeLogParser(logFormat, parseRegexStr)

	for scanner.Scan() {
		line := scanner.Text()
		entry := parseLogLine(line)
		storeLogEntry(entry, stmt, ctx)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("⚠️ Error while scanning input: %v", err)
	} else {
		log.Println("📭 STDIN closed — no longer receiving logs")
	}
}

func makeLogParser(format string, regexStr string) func(string) shared.LogEntry {
	if format == "text" && regexStr != "" {
		re, err := regexp.Compile(regexStr)
		if err != nil {
			log.Fatalf("Invalid --parse-regex: %v", err)
		}
		return func(line string) shared.LogEntry {
			return parseTextLogWithRegex(line, re)
		}
	}

	// default to JSON
	return func(line string) shared.LogEntry {
		return parseJSONLog(line)
	}
}

func parseJSONLog(line string) shared.LogEntry {
	var entry shared.LogEntry
	if err := json.Unmarshal([]byte(line), &entry); err == nil {
		return entry
	}
	return shared.LogEntry{
		"message": line,
		"level":   "raw",
	}
}

func parseTextLogWithRegex(line string, re *regexp.Regexp) shared.LogEntry {
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return shared.LogEntry{
			"message": line,
			"level":   "raw",
		}
	}

	entry := shared.LogEntry{}
	for i, name := range re.SubexpNames() {
		if i > 0 && name != "" {
			entry[name] = matches[i]
		}
	}

	if _, ok := entry["level"]; !ok {
		entry["level"] = "info"
	}
	if _, ok := entry["message"]; !ok {
		entry["message"] = line
	}

	return entry
}

func storeLogEntry(entry shared.LogEntry, stmt *sql.Stmt, ctx context.Context) {
	raw := string(shared.MustJson(entry))
	createdAt := time.Now().UTC()

	_, err := stmt.ExecContext(
		ctx,
		entry["trace_id"],
		entry["level"],
		entry["message"],
		raw,
		createdAt,
	)

	if err != nil {
		log.Printf("❌ Failed to insert log: %v", err)
		return
	}

	entry["created_at"] = createdAt
	handlers.Broadcast(entry)
}
