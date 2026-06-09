#!/usr/bin/env bash
# Smoke-test every API endpoint. Exits non-zero on the first non-200 response.
#
# Override defaults via env vars, e.g.:
#   BASE_URL=http://localhost:8080 TOKEN=... ADM_ID=... bash call-api.sh

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
TOKEN="${TOKEN:-3b2b552e-33e5-41c3-8547-9670ab57dca1}"
ADM_ID="${ADM_ID:-b9b06fa6-4697-4fa4-ad1a-1a2068fe97e7}"
DEFAULT_POINT_BODY='{"type":"Point","coordinates":[26.1,41.7]}'

API_BASE="${BASE_URL}/api/v1"

PASS=0
FAIL=0
FAILED=()

# call METHOD URL [EXPECTED_STATUS] [BODY]
call() {
  local method="$1"
  local url="$2"
  local expected="${3:-200}"
  local body="${4:-}"

  local args=(-sS -o /dev/null -w '%{http_code}'
              -X "$method"
              -H "Authorization: Bearer ${TOKEN}")
  if [[ -n "$body" ]]; then
    args+=(-H 'Content-Type: application/json' --data "$body")
  fi

  local status
  status=$(curl "${args[@]}" "$url" || echo "000")

  if [[ "$status" == "$expected" ]]; then
    printf '  PASS  %-4s %s -> %s\n' "$method" "$url" "$status"
    PASS=$((PASS + 1))
  else
    printf '  FAIL  %-4s %s -> %s (expected %s)\n' \
      "$method" "$url" "$status" "$expected"
    FAIL=$((FAIL + 1))
    FAILED+=("$method $url")
  fi
}

echo "--- adm-neighbors ---"
call GET  "${API_BASE}/adm-neighbors?adm-id=${ADM_ID}"
call POST "${API_BASE}/adm-neighbors"   200 "$DEFAULT_POINT_BODY"

echo "--- reverse-geocode ---"
call POST "${API_BASE}/reverse-geocode" 200 "$DEFAULT_POINT_BODY"

sleep 1

echo "--- feature collection ---"
for lv in 0 1 2 3 4 5; do
  call GET "${API_BASE}/fc?lv=${lv}&batch-size=2"
done

sleep 1

echo "--- geojsonl ---"
for lv in 0 1 2 3 4 5; do
  call GET "${API_BASE}/geojsonl/lv${lv}?page-size=2"
done

# create-access-token is disabled by default because it writes to the DB and
# is globally rate limited. Enable with: WITH_TOKEN_CREATE=1 bash call-api.sh
if [[ "${WITH_TOKEN_CREATE:-0}" == "1" ]]; then
  echo "--- create-access-token ---"
  call POST "${API_BASE}/create-access-token?email=smoketest@example.com" 201
fi

echo
echo "Passed: ${PASS}    Failed: ${FAIL}"

if [[ "$FAIL" -gt 0 ]]; then
  echo "Failed endpoints:"
  printf '  - %s\n' "${FAILED[@]}"
  exit 1
fi


# curl -N "${API_BASE}/geojsonl?lv=0" \
#   -H "Authorization: Bearer ${TOKEN}"








