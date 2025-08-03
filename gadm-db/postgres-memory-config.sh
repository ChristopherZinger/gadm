#!/bin/bash
# PostgreSQL memory configuration for large imports

set -e

# Create custom postgresql.conf settings for memory optimization
cat >> "$PGDATA/postgresql.conf" <<EOF

# Memory tuning for 1GB VM GADM imports
shared_buffers = '${POSTGRES_SHARED_BUFFERS:-64MB}'
work_mem = '${POSTGRES_WORK_MEM:-1MB}'
maintenance_work_mem = '${POSTGRES_MAINTENANCE_WORK_MEM:-16MB}'
effective_cache_size = '${POSTGRES_EFFECTIVE_CACHE_SIZE:-128MB}'

# Additional settings for 1GB VM bulk imports
checkpoint_completion_target = 0.9
wal_buffers = 4MB
checkpoint_timeout = 30min
max_wal_size = 512MB
min_wal_size = 80MB

# Logging for debugging memory issues
log_temp_files = 10MB
log_checkpoints = on
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '

EOF

echo "PostgreSQL memory configuration applied"