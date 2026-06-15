package access_token

import (
	"context"
	"fmt"
	"gadm-api/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type accessToken struct {
	Id                      int       `db:"id" json:"id"`
	Token                   string    `db:"token" json:"token"`
	Email                   string    `db:"email" json:"email"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time `db:"updated_at" json:"updated_at"`
	CanGenerateAccessTokens bool      `db:"can_generate_access_tokens" json:"can_generate_access_tokens"`
}

type accessTokenRepo struct {
	db *pgxpool.Pool
}

func NewAccessTokenRepo(db *pgxpool.Pool) *accessTokenRepo {
	return &accessTokenRepo{db}
}


func (repo *accessTokenRepo) getAccessToken(ctx context.Context, token string) (*accessToken, error) {
	sql, args, err := getAccessTokenSqlQuery(token)
	if err != nil {
		return nil, fmt.Errorf("failed_to_get_access_token_sql_query %v", err)
	}

	var row pgx.Rows
	row, err = repo.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed_to_get_access_token %v", err)
	}

	var _accessToken accessToken
	_accessToken, err = pgx.CollectOneRow(row, pgx.RowToStructByNameLax[accessToken])
	if err != nil {
		return nil, fmt.Errorf("failed_to_get_access_token %v", err)
	}

	return &_accessToken, nil
}
