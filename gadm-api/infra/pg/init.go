package pg

import (
	"context"
	logger "gadm-api/logger"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DATABASE_URL_ENV_VAR = "DATABASE_URL"

func InitPgPool(poolSize int32) *pgxpool.Pool {
	databaseUrl := os.Getenv(DATABASE_URL_ENV_VAR)
	if databaseUrl == "" {
		logger.Fatal("missing_db_url_env_variable %s", databaseUrl) // TODO: expect var from config file util
	}

	cfg, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		logger.Fatal("failed_to_parse_db_url %v", err)
	}
	cfg.MaxConns = poolSize

	dbPool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		logger.Fatal("failed_to_connect_to_database %v", err)
	}

	return dbPool
}
