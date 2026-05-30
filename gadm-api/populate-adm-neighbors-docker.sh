#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}/.."

exec docker compose run --rm --build \
  -e SERVICE_TYPE=cron_job \
  -e CRON_JOB_NAME=populate_adm_neighbors \
  gadm-api 
