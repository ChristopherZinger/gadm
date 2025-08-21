package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	accessTokenCache "gadm-api/access-token-cache"
	db "gadm-api/db"
	"gadm-api/handlers"
	"gadm-api/logger"
)

func GetAuthMiddleWare(pgConn *db.PgConn) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := handlers.GetAuthTokenFromRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if token != "" {
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
			}
			next.ServeHTTP(w, r)
		})
	}
}

func getTokenCreatedAtFromDb(ctx context.Context, pgConn *db.PgConn, token string) (time.Time, error) {
	sql, args, err := db.GetAccessTokenCreatedAtSqlQuery(token)

	var createdAt time.Time
	err = pgConn.Db.QueryRow(ctx, sql, args...).Scan(&createdAt)
	if err != nil {
		logger.Error("%v", err)
		if err.Error() == handlers.NOT_RESULTS_FOR_QUERY_PG_MSG {
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
