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
	"github.com/paul-schwendenman/magic-log-ui/internal/server/handlers"
	"github.com/paul-schwendenman/magic-log-ui/internal/shared"
)

func Start(input io.Reader, stmt *sql.Stmt, logFormat string, parseRegexStr string, jqQuery string, echo bool, ctx context.Context) {
	scanner := bufio.NewScanner(input)

	parseLogLine, err := makeLogParser(logFormat, parseRegexStr)
	if err != nil {
		log.Fatalf("❌ Failed to initialize log parser: %v", err)
	}

	jqfilter.Init(jqQuery)

	for scanner.Scan() {
		line := scanner.Text()

		entry, parsed := parseLogLine(line)
		if !parsed {
			log.Printf("❌ Failed to parse log line: %q", line)
		}

		entry = mapToLogEntry(jqfilter.Apply(logEntryToStringMap(entry)))

		raw := string(shared.MustJson(entry))
		createdAt := time.Now().UTC()

		if echo {
			fmt.Println(raw)
		}

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
			continue
		}

		entry["created_at"] = createdAt
		handlers.Broadcast(entry)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("⚠️ Error while scanning input: %v", err)
	} else {
		log.Println("📭 STDIN closed — no longer receiving logs")
	}
}

func makeLogParser(format string, regexStr string) (func(string) (shared.LogEntry, bool), error) {
	if format == "text" && regexStr != "" {
		re, err := regexp.Compile(regexStr)
		if err != nil {
			return nil, err
		}
		return func(line string) (shared.LogEntry, bool) {
			return parseTextLogWithRegex(line, re)
		}, nil
	}

	// default to JSON
	return func(line string) (shared.LogEntry, bool) {
		return parseJSONLog(line)
	}, nil
}

func parseJSONLog(line string) (shared.LogEntry, bool) {
	var entry shared.LogEntry
	if err := json.Unmarshal([]byte(line), &entry); err == nil {
		return entry, true
	}
	return shared.LogEntry{
		"message": line,
		"level":   "raw",
	}, false
}

func parseTextLogWithRegex(line string, re *regexp.Regexp) (shared.LogEntry, bool) {
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return shared.LogEntry{
			"message": line,
			"level":   "raw",
		}, false
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

	return entry, true
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
