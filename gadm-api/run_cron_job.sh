#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

docker compose build gadm-api

docker compose run --rm -e SERVICE_TYPE=cron_job gadm-api
