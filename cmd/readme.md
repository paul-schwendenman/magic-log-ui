# Magic Log UI



## Testing

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