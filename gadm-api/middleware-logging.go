package main

import (
	"net/http"
	"time"

	"gadm-api/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Info("request_started method=%s path=%s remote_addr=%s query_params=%v",
			r.Method, r.URL.Path, r.RemoteAddr, r.URL.Query())

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		logger.Info("request_completed method=%s path=%s duration=%v remote_addr=%s",
			r.Method, r.URL.Path, duration, r.RemoteAddr)
	})
}
