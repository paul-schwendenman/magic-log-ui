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

## Command line options

```
$ magic-log --help
Usage:
magic-log [flags]
magic-log config [get|set|unset] <key> [value]

Flags:
  -csv-fields string
        Comma-separated field names for CSV logs (used with --log-format=csv)
  -db-file string
        Path to a DuckDB database file.
  -echo
        Echo parsed stdin input to stdout
  -has-csv-header
        Indicates if CSV logs include a header row (default true)
  -jq string
        A jq expression to apply to parsed logs
  -jq-preset string
        Regex preset to use.
  -launch
        Open the UI in a browser.
  -list-presets
        List available regex and jq presets and exit.
  -log-format string
        Log format: json, csv or plain text. (default "json")
  -no-auto-analyze
        Disable automatic ANALYZE of logs table
  -port int
        Port to serve the web UI on. (default 3000)
  -regex string
        Custom regex to parse logs. Use with "text" format
  -regex-preset string
        Regex preset to use.
  -version
        Print version and exit.

Config:
  The CLI reads config from ~/.magiclogrc by default.
  You can override the config path using the MAGIC_LOG_CONFIG environment variable.

Examples:
  MAGIC_LOG_CONFIG=/path/to/custom.toml magic-log --port 4000
  magic-log config set port 4000
```

## Configuration

`magic-log` reads configuration options from a TOML file `$HOME/.magiclogrc`

For example:

```
db_file = "logs.db"
port = 4000
launch = true
log_format = "text"
regex_preset = "sveltekit"
regex = ""
jq = ""
jq_preset = ""
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

You can also use the UI to manage the config

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

Tip: To run the local build with these snippets use an alias:

```
alias magic-log="go run main.go"
```

And you should see something like:

```
$ magic-log -version
magic-log version: dev
```

Instead of:

```
$ magic-log -version
magic-log version: v0.0.6
```

To unset later use `unalias`

```
$ unalias magic-log
```

### Publish a few messages

```
yes '{"timestamp": "2025-04-10T12:00:00Z", "trace_id": "abc123", "level": "info", "message": "ping"}' | head -n 100 | magic-log
```

### Publish messages continuously

```
while true; do
  echo "{\"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\", \"level\": \"info\", \"message\": \"live test log\", \"trace_id\": \"abc123\"}"
  sleep 1
done | magic-log
```

### Random logs

JSON format:

```
./generate_logs.sh | magic-log
```

String with regex:

```
./generate_logs.sh --format text | \
  magic-log \
    --log-format=text \
    --regex="\\[(?P<level>\\w+)] (?P<time>\\S+) \\[(?P<trace_id>\\w+)] (?P<msg>.+)"
```

### Pipe logs from CSV

```
./generate_logs.sh --format csv | magic-log --log-format=csv --has-csv-header
```

Custom columns names:
```
$ ./generate_logs.sh --format csv | magic-log --log-format=csv --has-csv-header --csv-fields "timestamp,level,message,trace_id"
```

If one column is json that you want to ingest:
```
./csv_echoer.py ~/Downloads/extract.csv --column Message | magic-log

```

### JQ Filter Examples

You can reshape incoming logs during ingestion using `--jq`, based on JQ syntax.

#### Rename fields and keep the rest

```
./generate_logs.sh | magic-log \
--jq='{message: .msg} + .'
```

#### Add new static field

```
./generate_logs.sh | magic-log \
--jq='{app: "magic-log", trace_id: .trace_id, message: .msg}'
```

#### Drop fields

```
./generate_logs.sh | magic-log \
--jq='del(.time, .msg)'
```

### Controlling the tool

#### Launch flag

`magic-log` is capable of opening a new browser tab with the UI using the `--launch` flag

For example:

```
./generate_logs.sh | magic-log --launch
```

You can also use the config to set the default to `true`

```
magic-log config set launch true
```

Once the tool defaults to true you can still pass `--launch` to disable it like so:

```
./generate_logs.sh | magic-log --launch=false
```

#### Setting the db file

`magic-log` can persist logs to disc via a duckdb database or just use an in memory db. The
in-memory is the default.

Persisting logs to disc, using relative path:
```
... | magic-log --db-file="logs.duckdb"
```

Setting a default db via configuration:
```
magic-log config set db_file logs.duckdb
```

Overriding config to use in-memory database:
```
... | magic-log --db-file=""
```
