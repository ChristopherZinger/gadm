package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	accessTokenCache "gadm-api/access-token-cache"
	db "gadm-api/db"
	"gadm-api/logger"
)

func GetAuthTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	logger.Debug("auth_header_received header=%s remote_addr=%s path=%s", authHeader, r.RemoteAddr, r.URL.Path)
	var token string

	if authHeader != "" {
		const bearerPrefix = "Bearer "
		if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
			token = authHeader[len(bearerPrefix):]
			logger.Debug("token_extracted token=%s", token)
			if token == "" {
				logger.Debug("missing_token remote_addr=%s path=%s", r.RemoteAddr, r.URL.Path)
				return "", errors.New("empty_token")
			}
			return token, nil
		} else {
			logger.Debug("invalid_bearer_format auth_header=%s", authHeader)
			return "", errors.New("invalid_bearer_format")
		}
	}
	logger.Debug("missing_token remote_addr=%s path=%s", r.RemoteAddr, r.URL.Path)
	return "", errors.New("missing_token")
}

func GetAuthMiddleWare(pgConn *PgConn) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := GetAuthTokenFromRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
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
					http.Error(w, "token_expired", http.StatusUnauthorized)
					return
				case accessTokenCache.RateLimitExceededMsg:
					http.Error(w, "rate_limit_exceeded", http.StatusTooManyRequests)
					return
				case FailedToQueryDatabaseMsg:
					http.Error(w, "internal_server_error", http.StatusInternalServerError)
					return
				case NoResultsMsg:
					http.Error(w, "invalid_access_token", http.StatusUnauthorized)
					return
				default:
					http.Error(w, "internal_server_error", http.StatusInternalServerError)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

const NOT_RESULTS_FOR_QUERY_PG_MSG = "no rows in result set"

func getTokenCreatedAtFromDb(ctx context.Context, pgConn *PgConn, token string) (time.Time, error) {
	sql, args, err := db.GetAccessTokenCreatedAtSqlQuery(token)

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
