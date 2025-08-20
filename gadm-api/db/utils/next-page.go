package db

import (
	"context"
	"fmt"

	db "gadm-api/db"
	"gadm-api/logger"
)

type NextPageRepo struct {
	pgConn *db.PgConn
	ctx    context.Context
}

func NewNextPageRepo(pgConn *db.PgConn, ctx context.Context) *NextPageRepo {
	return &NextPageRepo{pgConn: pgConn, ctx: ctx}
}

type NextPageParams struct {
	StartAt       int
	PageSize      int
	FilterColName string
	FilterVal     string
}

func (repo *NextPageRepo) GetNextPageFid(
	params NextPageParams,
) (int, error) {

	sql, args, err := db.GetNextFidSqlQuery(db.GetNextFidSqlQueryParams{
		StartAtFid:    params.StartAt,
		PageSize:      params.PageSize,
		FilterColName: params.FilterColName,
		FilterVal:     params.FilterVal,
	})

	var nextFid int
	err = repo.pgConn.Db.QueryRow(repo.ctx, sql, args...).Scan(&nextFid)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		return 0, fmt.Errorf("failed_to_get_next_fid %v", err)
	}
	return nextFid, nil
}
