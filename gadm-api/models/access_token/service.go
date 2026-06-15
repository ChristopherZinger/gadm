package access_token

import (
	"context"
	"fmt"
	"time"
)

type accessTokenService struct {
	repo *accessTokenRepo
}

func NewAccessTokenService(repo *accessTokenRepo) *accessTokenService {
	return &accessTokenService{repo}
}


func (service *accessTokenService) GetAccessToken(ctx context.Context, token string) (*accessToken, error) {
	_accessToken, err := service.repo.getAccessToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed_to_get_access_token_created_at %v", err)
	}
	return _accessToken, nil
}
