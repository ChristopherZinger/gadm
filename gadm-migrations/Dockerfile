# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the migration CLI
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gadm-migrate ./cmd

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/gadm-migrate .

# Copy migrations directory from gadm-db
COPY ./migrations/ ./migrations/

# Create a wrapper script for common operations
RUN echo '#!/bin/sh' > /app/migrate.sh && \
    echo 'exec ./gadm-migrate --db-url="$DATABASE_URL" --migrations-dir="./migrations" "$@"' >> /app/migrate.sh && \
    chmod +x /app/migrate.sh

# Default command
ENTRYPOINT ["/app/migrate.sh"]
CMD ["up"]