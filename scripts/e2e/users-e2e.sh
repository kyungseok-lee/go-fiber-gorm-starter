#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
WAIT_SECONDS="${WAIT_SECONDS:-60}"
EMAIL="e2e-$(date +%s)-$$@example.com"

request() {
  local method="$1"
  local path="$2"
  local payload="${3:-}"
  local body_file
  body_file="$(mktemp)"

  local curl_args=(-sS -o "$body_file" -w "%{http_code}" -X "$method")
  if [[ -n "${E2E_API_KEY:-}" ]]; then
    curl_args+=(-H "Authorization: Bearer ${E2E_API_KEY}")
  fi
  if [[ -n "$payload" ]]; then
    curl_args+=(-H "Content-Type: application/json" -d "$payload")
  fi

  local status
  status="$(curl "${curl_args[@]}" "${BASE_URL}${path}")"
  printf '%s\n' "$status"
  cat "$body_file"
  rm -f "$body_file"
}

expect_status() {
  local expected="$1"
  local method="$2"
  local path="$3"
  local payload="${4:-}"
  local output status body

  output="$(request "$method" "$path" "$payload")"
  status="$(printf '%s\n' "$output" | sed -n '1p')"
  body="$(printf '%s\n' "$output" | sed '1d')"

  if [[ "$status" != "$expected" ]]; then
    printf 'Expected %s %s to return %s, got %s\n' "$method" "$path" "$expected" "$status" >&2
    printf 'Response body: %s\n' "$body" >&2
    exit 1
  fi

  printf '%s' "$body"
}

wait_for_ready() {
  local deadline=$((SECONDS + WAIT_SECONDS))
  while (( SECONDS < deadline )); do
    if body="$(expect_status 200 GET /ready 2>/dev/null)"; then
      printf 'Ready: %s\n' "$body"
      return
    fi
    sleep 2
  done

  printf 'Timed out waiting for %s/ready\n' "$BASE_URL" >&2
  exit 1
}

extract_id() {
  sed -n 's/.*"id":\([0-9][0-9]*\).*/\1/p'
}

main() {
  expect_status 200 GET /health >/dev/null
  wait_for_ready

  local create_body user_id
  create_body="$(expect_status 201 POST /v1/users "{\"name\":\"E2E User\",\"email\":\"${EMAIL}\",\"status\":\"active\"}")"
  user_id="$(printf '%s' "$create_body" | extract_id)"
  if [[ -z "$user_id" ]]; then
    printf 'Could not extract created user id from: %s\n' "$create_body" >&2
    exit 1
  fi

  expect_status 400 POST /v1/users "{\"name\":\"Invalid Status\",\"email\":\"invalid-${EMAIL}\",\"status\":\"pending\"}" >/dev/null
  expect_status 409 POST /v1/users "{\"name\":\"E2E User\",\"email\":\"${EMAIL}\",\"status\":\"active\"}" >/dev/null
  expect_status 200 GET "/v1/users/${user_id}" >/dev/null
  expect_status 200 GET "/v1/users?search=${EMAIL}&limit=10" >/dev/null
  expect_status 400 GET "/v1/users?status=pending" >/dev/null
  expect_status 400 PUT "/v1/users/${user_id}" '{"status":"pending"}' >/dev/null
  expect_status 200 PUT "/v1/users/${user_id}" '{"name":"E2E User Updated","status":"inactive"}' >/dev/null
  expect_status 204 DELETE "/v1/users/${user_id}" >/dev/null
  expect_status 404 GET "/v1/users/${user_id}" >/dev/null
  expect_status 404 GET /v1/not-found >/dev/null

  printf 'E2E passed against %s\n' "$BASE_URL"
}

main "$@"
