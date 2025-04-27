# Magic Log UI

Magic Log UI is a local-first log viewer for streaming structured JSON logs into an in-memory (or persistent) DuckDB database. It provides a real-time browser UI for querying, exploring, and debugging logs with full SQL support and zero setup.

## Features

- Ingests structured JSON **or** plain text logs from stdin
- Regex-based parsing for text logs, with support for custom or preset patterns (e.g., apache, nginx, sveltekit)
- JQ-style transformations to reshape or extract fields during ingestion
- Configurable via CLI flags **or** a `.magiclogrc` config file (TOML), with optional `MAGIC_LOG_CONFIG` override
- Query logs live in real-time using SQL (DuckDB — in-memory or persistent database modes)
- Real-time browser UI with dynamic WebSocket streaming
- View, save, and re-run past queries from the browser
- Auto-analyze feature keeps query performance fast (optional --no-auto-analyze flag)
- Environment variable override support for flexible deployment
- One-file executable — no external database or server setup needed
- Fully testable Go codebase with coverage reporting
- Optional log persistence with `--db-file`
- Homebrew install available: `brew install paul-schwendenman/magic-log`

## Installation

```
brew tap paul-schwendenman/magic-log-ui
brew install magic-log
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

It is possible to manage the config using the CLI, like so:

```
magic-log config set log_format text
magic-log config set port 4001
magic-log config set regex_presets.myapp 'regex...'
```

You can also fetch values:

```
magic-log config get log_format
magic-log config get port
magic-log config get regex_presets.myapp
```

And remove them:

```
magic-log config unset log_format
magic-log config unset port
magic-log config unset regex_presets.myapp
```

## Development

### Build project

This command will build the frontend and then package the app

```
make build
```

### Run tests

```
go test ./...
```

## Manual Testing

The following sections contain some snippets for testing the app manually by producing various logs

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

### JQ Filter Examples

You can reshape incoming logs during ingestion using `--jq-filter`, based on JQ syntax.

#### Rename fields and keep the rest

```
./generate_logs.sh | go run ./cmd/main.go
--jq-filter='{message: .msg} + .'
```

#### Add new static field

```
./generate_logs.sh | go run ./cmd/main.go
--jq-filter='{app: "magic-log", trace_id: .trace_id, message: .msg}'
```

#### Drop fields

```
./generate_logs.sh | go run ./cmd/main.go
--jq-filter='del(.time, .msg)'
```
