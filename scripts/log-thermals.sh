#!/usr/bin/env bash
# Logs CPU/SoC thermal zone temperatures to a log file.
# Intended to be invoked from cron.

set -euo pipefail

LOG_FILE="${THERMAL_LOG_FILE:-/var/log/thermals.log}"

mkdir -p "$(dirname "$LOG_FILE")"

timestamp="$(date '+%Y-%m-%d %H:%M:%S%z')"

readings="$(
    paste \
        <(cat /sys/class/thermal/thermal_zone*/type) \
        <(cat /sys/class/thermal/thermal_zone*/temp) \
    | column -s $'\t' -t \
    | sed 's/\(.\)..$/.\1°C/'
)"

{
    echo "===== $timestamp ====="
    echo "$readings"
} >> "$LOG_FILE"
