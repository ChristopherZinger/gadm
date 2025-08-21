package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gadm-api/logger"
	"net/http"
	"time"

	db "gadm-api/db"
)

type AccessTokenCreationHandler struct {
	pgConn *db.PgConn
	req    *http.Request
	writer http.ResponseWriter
	ctx    context.Context
}

func NewAccessTokenCreationHandler(
	pgConn *db.PgConn,
	req *http.Request,
	writer http.ResponseWriter,
	ctx context.Context) *AccessTokenCreationHandler {
	return &AccessTokenCreationHandler{pgConn: pgConn, req: req, writer: writer, ctx: ctx}
}

func GetAccessTokenCreationHandlerInfo(pgConn *db.PgConn) HandlerInfo {
	return HandlerInfo{
		Url: fmt.Sprintf("%screate-access-token", getBaseApiUrl().Path),
		Handler: func(w http.ResponseWriter, r *http.Request) {
			NewAccessTokenCreationHandler(pgConn, r, w, r.Context()).handle()
		},
	}
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

	if !tokenCreationRateLimiter.IsAllowed() {
		logger.Warning("global_rate_limit_exceeded remote_addr=%s", handler.req.RemoteAddr)
		errorMsg := fmt.Sprintf("rate_limit_exceeded")
		http.Error(handler.writer, errorMsg, http.StatusTooManyRequests)
		return
	}

	email := handler.req.URL.Query().Get("email")
	if email == "" {
		logger.Error("missing_email_parameter")
		http.Error(handler.writer, "email_not_provided", http.StatusBadRequest)
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
	err = handler.pgConn.Db.QueryRow(handler.ctx, insertSql, insertArgs...).Scan(&createdToken, &tokenCreatedAt)
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
