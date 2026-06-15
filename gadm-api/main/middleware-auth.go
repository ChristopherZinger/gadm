package main

import (
	"net/http"
	"time"

	accessTokenCache "gadm-api/access-token-cache"
	"gadm-api/logger"
	"gadm-api/models/access_token"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAuthMiddleWare(pgPool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := getApiAuthTokenFromRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if token != "" {
				if err := accessTokenCache.TOKEN_CACHE.HandleHitForToken(
					token,
					func(token string) (time.Time, error) {
						accessTokenRepo := access_token.NewAccessTokenRepo(pgPool)
						service := access_token.NewAccessTokenService(accessTokenRepo)
						_token, err := service.GetAccessToken(r.Context(), token)
						if err != nil {
							return time.Time{}, err
						}
						return _token.CreatedAt, nil
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

const (
	NoResultsMsg             = "no_results"
	FailedToQueryDatabaseMsg = "failed_to_query_database"
)
