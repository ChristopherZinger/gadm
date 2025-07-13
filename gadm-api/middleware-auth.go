package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	accessTokenCache "gadm-api/access-token-cache"
	"gadm-api/logger"
)

func GetAuthMiddleWare(pgConn *PgConn) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			logger.Info("auth_header_received header=%s remote_addr=%s path=%s", authHeader, r.RemoteAddr, r.URL.Path)
			var token string

			if authHeader != "" {
				const bearerPrefix = "Bearer "
				if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
					token = authHeader[len(bearerPrefix):]
					logger.Info("token_extracted token=%s", token)
				} else {
					logger.Debug("invalid_bearer_format auth_header=%s", authHeader)
				}
			}

			if token == "" {
				logger.Warning("missing_token remote_addr=%s path=%s", r.RemoteAddr, r.URL.Path)
				http.Error(w, "Unauthorized: Missing access token", http.StatusUnauthorized)
				return
			}

			if err := accessTokenCache.TOKEN_CACHE.HandleHitForToken(
				token,
				func(token string) (time.Time, error) {
					return getTokenCreatedAtFromDb(r.Context(), pgConn, token)
				}); err != nil {

				logger.Error("token_validation_failed %v", err)

				switch err.Error() {
				case accessTokenCache.TokenExpiredMsg:
					http.Error(w, "Unauthorized: Token expired", http.StatusUnauthorized)
					return
				case accessTokenCache.RateLimitExceededMsg:
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
					return
				case FailedToQueryDatabaseMsg:
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				case NoResultsMsg:
					http.Error(w, "Unauthorized: Invalid access token", http.StatusUnauthorized)
					return
				default:
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

const NOT_RESULTS_FOR_QUERY_PG_MSG = "no rows in result set"

func getTokenCreatedAtFromDb(ctx context.Context, pgConn *PgConn, token string) (time.Time, error) {
	sql, args, err := getAccessTokenCreatedAtSqlQuery(token)

	var createdAt time.Time
	err = pgConn.db.QueryRow(ctx, sql, args...).Scan(&createdAt)
	if err != nil {
		logger.Error("%v", err)
		if err.Error() == NOT_RESULTS_FOR_QUERY_PG_MSG {
			return time.Time{}, errors.New(NoResultsMsg)
		}
		return time.Time{}, errors.New(FailedToQueryDatabaseMsg)
	}

	return createdAt, nil
}

const (
	NoResultsMsg             = "no_results"
	FailedToQueryDatabaseMsg = "failed_to_query_database"
)
