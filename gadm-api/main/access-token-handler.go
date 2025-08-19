package main

import (
	"context"
	"encoding/json"
	"gadm-api/logger"
	"net/http"
	"time"

	accessTokenCache "gadm-api/access-token-cache"
	db "gadm-api/db"
)

type AccessTokenCreationHandler struct {
	pgConn *PgConn
	req    *http.Request
	writer http.ResponseWriter
	ctx    context.Context
}

func NewAccessTokenCreationHandler(
	pgConn *PgConn,
	req *http.Request,
	writer http.ResponseWriter,
	ctx context.Context) *AccessTokenCreationHandler {
	return &AccessTokenCreationHandler{pgConn: pgConn, req: req, writer: writer, ctx: ctx}
}

type AccessTokenCreationResponse struct {
	Token     string    `json:"token"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (handler *AccessTokenCreationHandler) handle() {
	if handler.req.Method != http.MethodPost {
		logger.Error("invalid_method method=%s", handler.req.Method)
		http.Error(handler.writer, "method_not_allowed", http.StatusMethodNotAllowed)
		return
	}

	token, err := GetAuthTokenFromRequest(handler.req)
	if err != nil {
		logger.Error("failed_to_get_auth_token %v", err)
		http.Error(handler.writer, "unauthorized "+err.Error(), http.StatusUnauthorized)
		return
	}

	email := handler.req.URL.Query().Get("email")
	if email == "" {
		logger.Error("missing_email_parameter")
		http.Error(handler.writer, "email_not_provided", http.StatusBadRequest)
		return
	}

	sql, args, err := db.GetAccessTokenSqlQuery(token)
	if err != nil {
		logger.Error("failed_to_get_access_token_sql_query %v", err)
		http.Error(handler.writer, "internal_server_error", http.StatusInternalServerError)
		return
	}

	var createdAt time.Time
	var canGenerateAccessTokens bool
	err = handler.pgConn.db.QueryRow(handler.ctx, sql, args...).Scan(&createdAt, &canGenerateAccessTokens)
	if err != nil {
		if err.Error() == NOT_RESULTS_FOR_QUERY_PG_MSG {
			logger.Error("token_not_found token=%s", token)
			http.Error(handler.writer, "invalid_access_token_", http.StatusUnauthorized)
			return
		}
		logger.Error("failed_to_query_access_token %v", err)
		http.Error(handler.writer, "internal_server_error", http.StatusInternalServerError)
		return
	}

	if accessTokenCache.IsTokenExpired(createdAt) {
		logger.Error("token_expired token=%s", token)
		http.Error(handler.writer, "token_expired", http.StatusUnauthorized)
		return
	}

	if !canGenerateAccessTokens {
		logger.Error("insufficient_permissions token=%s", token)
		http.Error(handler.writer, "insufficient_permissions", http.StatusForbidden)
		return
	}

	insertSql, insertArgs, err := db.GetInsertAccessTokenWithReturningSqlQuery(email)
	if err != nil {
		logger.Error("failed_to_get_insert_access_token_sql_query %v", err)
		http.Error(handler.writer, "internal_server_error", http.StatusInternalServerError)
		return
	}

	var createdToken string
	var tokenCreatedAt time.Time
	err = handler.pgConn.db.QueryRow(handler.ctx, insertSql, insertArgs...).Scan(&createdToken, &tokenCreatedAt)
	if err != nil {
		logger.Error("failed_to_insert_access_token %v", err)
		http.Error(handler.writer, "internal_server_error", http.StatusInternalServerError)
		return
	}

	handler.writer.Header().Set("Content-Type", "application/json")
	handler.writer.WriteHeader(http.StatusCreated)

	response := AccessTokenCreationResponse{
		Token:     createdToken,
		Email:     email,
		CreatedAt: tokenCreatedAt,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		logger.Error("failed_to_marshal_response %v", err)
		http.Error(handler.writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	handler.writer.Write(responseJSON)
}
