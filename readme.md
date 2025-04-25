# Magic Log UI

Magic Log UI is a local-first log viewer for streaming structured JSON logs into an in-memory (or persistent) DuckDB database. It provides a real-time browser UI for querying, exploring, and debugging logs with full SQL support and zero setup.

## Features

- Ingests structured JSON or plain text logs from stdin
- Regex-based parsing for text logs, with support for custom or preset patterns (e.g. apache, nginx, sveltekit)
- Configurable via CLI or .magiclogrc in your home directory (TOML)
- Query logs in real-time using SQL (DuckDB in-memory or persistent)
- Real-time browser UI with WebSocket streaming
- View and re-run past queries
- One-file executable — no external setup, just run and go
- Fully testable Go codebase with coverage support
- Optional log persistence with --db-file
- Homebrew install available: brew install paul-schwendenman/magic-log

## Installation

```
brew tap paul-schwendenman/magic-log-ui
brew install magic-log
```

## Build project

This command will build the frontend and then package the app

```
make build
```

## Run tests

```
go test ./...
```

## Configuration

`magic-log` reads configuration options from `$HOME/.magiclogrc`

For example:

```
[defaults]
db_file = "logs.db"
port = 4000
launch = true
log_format = "text"
parse_preset = "sveltekit"
parse_regex = ""
```

## Manual Testing

### Publish a few messages

```
yes '{"timestamp": "2025-04-10T12:00:00Z", "trace_id": "abc123", "level": "info", "message": "ping"}' | head -n 100 | go run ./cmd/main.go
```

### Publish messages continuously

```
while true; do
  echo "{\"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\", \"level\": \"info\", \"message\": \"live test log\", \"trace_id\": \"abc123\"}"
  sleep 1
done | go run ./cmd/main.go
```

### Random logs

JSON format:

```
./generate_logs.sh | go run ./cmd/main.go
```

String with regex:

```
./generate_logs.sh --format text | \
  go run ./cmd/main.go \
    --log-format=text \
    --parse-regex="\\[(?P<level>\\w+)] (?P<time>\\S+) \\[(?P<trace_id>\\w+)] (?P<msg>.+)"
```

### Pipe logs from CSV

```
./csv_echoer.py ~/Downloads/extract.csv --column Message | go run ./cmd/main.go
```
