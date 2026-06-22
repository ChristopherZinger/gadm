package access_token

import (
	"encoding/json"
	"fmt"
	"gadm-api/logger"
	"net/http"
	"time"
)

type accessTokenHandler struct {
	service *accessTokenService
}

func NewAccessTokenHandler(service *accessTokenService) *accessTokenHandler {
	return &accessTokenHandler{service: service}
}

type limiter interface {
	IsAllowed() bool
}

func (handler *accessTokenHandler) CreateAccessTokenHandler(w http.ResponseWriter, req *http.Request, limiter limiter) {
	if req.Method != http.MethodPost {
		logger.Error("invalid_method method=%s", req.Method)
		http.Error(w, "method_not_allowed", http.StatusMethodNotAllowed)
		return
	}

	if !limiter.IsAllowed() {
		logger.Warning("global_rate_limit_exceeded remote_addr=%s", req.RemoteAddr)
		errorMsg := fmt.Sprintf("rate_limit_exceeded")
		http.Error(w, errorMsg, http.StatusTooManyRequests)
		return
	}

	email := req.URL.Query().Get("email")
	if email == "" {
		logger.Error("missing_email_parameter")
		http.Error(w, "email_not_provided", http.StatusBadRequest)
		return
	}

	token, err := handler.service.createAccessToken(req.Context(), email)
	if err != nil {
		logger.Error("failed_to_create_access_token %v", err)
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	m := map[string]string{
		"token":      token,
		"email":      email,
		"created_at": time.Now().Format(time.RFC3339),
	}

	responseJSON, err := json.Marshal(m)
	if err != nil {
		logger.Error("failed_to_marshal_response %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(responseJSON)
}
