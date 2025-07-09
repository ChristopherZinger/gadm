#!/bin/bash

echo "Running database migration..."

# Set environment variables for goose
export DATABASE_URL="postgres://admin:secret@localhost:5432/gadm?sslmode=disable"

# Run the migration
goose -dir ./gadm-api/migrations postgres "$DATABASE_URL" down

echo "Migration completed!" 