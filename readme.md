# Magic Log UI

Magic Log UI is a local-first log viewer for streaming structured JSON logs into an in-memory (or persistent) DuckDB database. It provides a real-time browser UI for querying, exploring, and debugging logs with full SQL support and zero setup.

## Features

- Ingests structured JSON logs from stdin
- Query logs live using SQL (DuckDB)
- Real-time log streaming in the browser (WebSocket)
- View past queries
- Optional persistence via `--db-file`
- Fully testable with `go test`
- One-file executable: no setup required


## Build project

This command will build the frontend and then package the app

```
make build
```

## Run tests

```
go test ./...
```

## Manuel Testing

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