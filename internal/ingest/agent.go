package ingest

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

type ParserOptions struct {
	LogFormat    string
	ParseRegex   string
	JqFilter     string
	CSVFieldsStr string
	HasCSVHeader bool
	Echo         bool
}

func StartForwarding(input io.Reader, out chan<- string, opts ParserOptions) {
	scanner := bufio.NewScanner(input)
	p := buildParsers(opts.LogFormat, opts.ParseRegex, opts.JqFilter, opts.CSVFieldsStr, opts.HasCSVHeader)
	headerExtracted := false

	for scanner.Scan() {
		rawLine := scanner.Text()

		if opts.LogFormat == "csv" && opts.HasCSVHeader && !headerExtracted {
			header, err := readCSV(rawLine)
			if err != nil {
				log.Printf("❌ Failed to read CSV header: %v", err)
				continue
			}
			if p.csvFields == nil {
				p.csvFields = header
			}
			headerExtracted = true
			continue
		}

		parsed := extract(rawLine, p)
		transformed := transform(parsed, p)

		data, err := json.Marshal(transformed)
		if err != nil {
			log.Printf("❌ Failed to marshal transformed log: %v", err)
			continue
		}

		if opts.Echo {
			log.Printf("➡️  %s", data)
		}

		out <- string(data)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("⚠️ Error while scanning: %v", err)
	}
}
