#!/usr/bin/env bash

set -e

FORMAT="json"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --format)
      FORMAT="$2"
      shift 2
      ;;
    *)
      echo "Usage: $0 [--format json|text]"
      exit 1
      ;;
  esac
done

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

while true; do
  LEVEL=${LEVELS[$RANDOM % ${#LEVELS[@]}]}
  MESSAGE=${MESSAGES[$RANDOM % ${#MESSAGES[@]}]}
  TRACE=$(openssl rand -hex 6)

  if [[ "$FORMAT" == "json" ]]; then
    TIME=$(date -u +%s)
    jq -n -c \
      --arg time "$TIME" \
      --arg level "$LEVEL" \
      --arg msg "$MESSAGE" \
      --arg trace_id "$TRACE" \
      '{time: $time, level: $level, msg: $msg, trace_id: $trace_id}'
  else
    TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    echo "[$LEVEL] $TIME [$TRACE] $MESSAGE"
  fi

  sleep 0.2
done
