#!/bin/bash
# Script to configure pg_hba.conf for postgres user security

set -e

# Find the pg_hba.conf file location
PG_HBA_FILE=$(find /var/lib/postgresql -name "pg_hba.conf" 2>/dev/null | head -1)

if [ -z "$PG_HBA_FILE" ]; then
    echo "Could not find pg_hba.conf file"
    exit 1
fi

echo "Found pg_hba.conf at: $PG_HBA_FILE"

# Backup original
cp "$PG_HBA_FILE" "$PG_HBA_FILE.backup"

# Create organized HBA configuration using include directives
{
    echo "# Include security rules"
    echo "include /etc/postgresql/hba-security.conf"
    echo ""
    echo "# Default rules from original file"
} > "$PG_HBA_FILE.new"

# Append the original file content
cat "$PG_HBA_FILE" >> "$PG_HBA_FILE.new"

# Copy our organized HBA files to the PostgreSQL config directory
mkdir -p /etc/postgresql
cp /docker-entrypoint-initdb.d/hba-security.conf /etc/postgresql/

# Replace the original file
mv "$PG_HBA_FILE.new" "$PG_HBA_FILE"

echo "pg_hba.conf updated successfully"
echo "Postgres user now restricted to localhost connections only"