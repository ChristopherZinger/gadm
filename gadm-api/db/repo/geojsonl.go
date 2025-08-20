package db

import (
	"context"
	"fmt"
	"gadm-api/db"
	"gadm-api/logger"
	"gadm-api/utils"

	"github.com/jackc/pgx/v5"
)

type GeojsonlRepo struct {
	pgConn *db.PgConn
	ctx    context.Context
}

func NewGeojsonlRepo(pgConn *db.PgConn, ctx context.Context) *GeojsonlRepo {
	return &GeojsonlRepo{pgConn: pgConn, ctx: ctx}
}

type GetGeojsonlParams struct {
	GadmLevel     utils.GadmLevel
	FilterVal     string
	FilterColName string
	StartAtValue  int
	LimitValue    int
}

func (repo *GeojsonlRepo) GetGeojsonl(
	params GetGeojsonlParams,
) (pgx.Rows, error) {
	sql, args, err := db.BuildGeojsonFeatureSqlQuery(
		params.GadmLevel,
		params.FilterVal,
		params.FilterColName,
		params.StartAtValue,
		params.LimitValue,
	)
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := repo.pgConn.Db.Query(repo.ctx, sql, args...)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	return rows, nil
}
