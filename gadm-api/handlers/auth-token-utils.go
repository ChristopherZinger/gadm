package handlers

import (
	"errors"
	"gadm-api/logger"
	"net/http"
)

const NOT_RESULTS_FOR_QUERY_PG_MSG = "no rows in result set"

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
