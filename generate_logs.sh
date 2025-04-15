#!/usr/bin/env bash

LEVELS=("debug" "info" "warn" "error")
MESSAGES=(
  "User login"
  "Fetching orders"
  "Database timeout"
  "Cache miss"
  "File uploaded"
  "Payment processed"
  "Webhook received"
  "Trace complete"
  "Invalid credentials"
  "Retrying request"
)
TRACE_ID=$(openssl rand -hex 6)

while true; do
  LEVEL=${LEVELS[$RANDOM % ${#LEVELS[@]}]}
  MESSAGE=${MESSAGES[$RANDOM % ${#MESSAGES[@]}]}
  TRACE=$(openssl rand -hex 6)
  TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

  jq -n -c  --arg ts "$TIME" \
            --arg level "$LEVEL" \
            --arg msg "$MESSAGE" \
            --arg trace_id "$TRACE" \
            '{timestamp: $ts, level: $level, message: $msg, trace_id: $trace_id}'

  sleep 0.2
done
