package access_token

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

type accessTokenService struct {
	repo *accessTokenRepo
}

func NewAccessTokenService(repo *accessTokenRepo) *accessTokenService {
	return &accessTokenService{repo}
}

func (service *accessTokenService) createAccessToken(ctx context.Context, email string) (string, error) {
	token := generateAccessToken()
	hashedToken := hashAccessToken(token)
	err := service.repo.createAccessToken(ctx, email, hashedToken)
	if err != nil {
		return "", fmt.Errorf("failed_to_create_access_token %v", err)
	}
	return token, nil
}

func generateAccessToken() string {
	return uuid.New().String()
}

func hashAccessToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (service *accessTokenService) GetAccessToken(ctx context.Context, token string) (*accessToken, error) {
	hashedToken := hashAccessToken(token)
	_accessToken, err := service.repo.getAccessToken(ctx, hashedToken)
	if err != nil {
		return nil, fmt.Errorf("failed_to_get_access_token_created_at %v", err)
	}
	return _accessToken, nil
}
