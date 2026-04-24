# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o file-server .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/file-server .

# `public` is gitignored; runtime may mount tiles here. Empty dir keeps FileServer valid.
RUN mkdir -p public

EXPOSE 8090

CMD ["./file-server"]
