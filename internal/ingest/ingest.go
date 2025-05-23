package ingest

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/paul-schwendenman/magic-log-ui/internal/jqfilter"
	"github.com/paul-schwendenman/magic-log-ui/internal/server/handlers"
	"github.com/paul-schwendenman/magic-log-ui/internal/shared"
)

type parsers struct {
	logFormat  string
	parseRegex *regexp.Regexp
	jqEnabled  bool
	jqFilter   string
	csvFields  []string
}

func Start(input io.Reader, stmt *sql.Stmt, logFormat, parseRegexStr, jqQuery, csvFieldsStr string, hasCSVHeader bool, echo bool, ctx context.Context) {
	scanner := attach(input)
	parsers := buildParsers(logFormat, parseRegexStr, jqQuery, csvFieldsStr, hasCSVHeader)
	headerExtracted := false

	for scanner.Scan() {
		rawLine := scanner.Text()

		if logFormat == "csv" && hasCSVHeader && !headerExtracted {
			header, err := readCSV(rawLine)
			if err != nil {
				log.Fatalf("❌ Failed to read CSV header: %v", err)
			}
			if parsers.csvFields == nil {
				parsers.csvFields = header
			}
			headerExtracted = true
			continue // Skip header row
		}

		parsed := extract(rawLine, parsers)
		transformed := transform(parsed, parsers)

		err := load(stmt, rawLine, parsed, transformed, parsers, ctx)
		if err != nil {
			log.Printf("❌ Failed to insert log: %v", err)
		}

		broadcast(transformed, echo)
	}

	handleScannerError(scanner)
}

func attach(input io.Reader) *bufio.Scanner {
	return bufio.NewScanner(input)
}

func buildParsers(logFormat, parseRegexStr, jqQuery, csvFieldsStr string, hasCSVHeader bool) parsers {
	var regex *regexp.Regexp
	if parseRegexStr != "" {
		var err error
		regex, err = regexp.Compile(parseRegexStr)
		if err != nil {
			log.Fatalf("❌ Invalid regex: %v", err)
		}
	}

	jqfilter.Init(jqQuery)

	var csvFields []string
	if logFormat == "csv" && csvFieldsStr != "" {
		csvFields = strings.Split(csvFieldsStr, ",")
	}

	return parsers{
		logFormat:  logFormat,
		parseRegex: regex,
		jqEnabled:  jqQuery != "",
		jqFilter:   jqQuery,
		csvFields:  csvFields,
	}
}

func extract(rawLine string, p parsers) shared.LogEntry {
	var parsed shared.LogEntry
	var err error

	if p.logFormat == "json" {
		err = json.Unmarshal([]byte(rawLine), &parsed)
	}

	if err != nil || p.logFormat == "text" {
		if p.parseRegex != nil {
			parsed, err = parseWithRegex(rawLine, p.parseRegex)
			if err == nil {
				return parsed
			}
		}
		parsed = shared.LogEntry{
			"message": rawLine,
			"level":   "raw",
		}
	}

	if p.logFormat == "csv" {
		parsed, err = parseCSV(rawLine, p.csvFields)
		if err != nil {
			parsed = shared.LogEntry{
				"message": rawLine,
				"level":   "raw",
			}
		}
	}

	return parsed
}

func transform(entry shared.LogEntry, p parsers) shared.LogEntry {
	if p.jqEnabled {
		entry = mapToLogEntry(jqfilter.Apply(logEntryToStringMap(entry)))
	}
	ensureTimestamp(entry)

	return entry
}

func load(stmt *sql.Stmt, rawLine string, parsed, transformed shared.LogEntry, p parsers, ctx context.Context) error {
	traceID, _ := safeString(transformed, "trace_id")
	level, _ := safeString(transformed, "level")
	message, _ := safeString(transformed, "message")

	if level == "" {
		level = "info"
	}
	if message == "" {
		message = "(no message)"
	}

	timestamp := time.Now().UTC()
	if ts, ok := safeString(transformed, "timestamp"); ok {
		parsedTs, err := time.Parse(time.RFC3339, ts)
		if err == nil {
			timestamp = parsedTs
		}
	}

	parsedLogJson := shared.MustJson(parsed)
	finalLogJson := shared.MustJson(transformed)

	regexPattern := ""
	if p.parseRegex != nil && p.logFormat == "text" {
		regexPattern = p.parseRegex.String()
	}

	id := uuid.New().String()

	_, err := stmt.ExecContext(
		ctx,
		id,
		traceID,
		level,
		message,
		rawLine,
		string(parsedLogJson),
		string(finalLogJson),
		time.Now().UTC(),
		timestamp,
		nullify(p.logFormat),
		nullify(regexPattern),
		nullify(p.jqFilter),
		nullify(strings.Join(p.csvFields, ",")),
	)

	return err
}

func broadcast(entry shared.LogEntry, echo bool) {
	handlers.Broadcast(entry)
	if echo {
		out, err := json.Marshal(entry)
		if err != nil {
			log.Printf("⚠️ Failed to encode entry for echo: %v", err)
			return
		}
		fmt.Println(string(out))
	}
}

func handleScannerError(scanner *bufio.Scanner) {
	if err := scanner.Err(); err != nil {
		log.Printf("⚠️ Error while scanning input: %v", err)
	} else {
		log.Println("📬 STDIN closed — no longer receiving logs")
	}
}

func parseWithRegex(line string, re *regexp.Regexp) (shared.LogEntry, error) {
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

func readCSV(rawLine string) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(rawLine))
	return reader.Read()
}

func parseCSV(rawLine string, fields []string) (shared.LogEntry, error) {
	record, err := readCSV(rawLine)
	if err != nil {
		return nil, err
	}

	if len(record) != len(fields) {
		return nil, fmt.Errorf("CSV field count mismatch: expected %d, got %d", len(fields), len(record))
	}

	entry := make(shared.LogEntry)
	for i, field := range fields {
		entry[field] = record[i]
	}

	return entry, nil
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

func ensureTimestamp(entry shared.LogEntry) {
	if entry == nil {
		entry = make(shared.LogEntry)
	}

	now := time.Now().UTC()

	if ts, ok := safeString(entry, "timestamp"); ok {
		parsedTs, err := time.Parse(time.RFC3339, ts)
		if err == nil {
			entry["timestamp"] = parsedTs.Format(time.RFC3339)
			return
		}
	}

	entry["timestamp"] = now.Format(time.RFC3339)
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
