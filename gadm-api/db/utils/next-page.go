package db

import (
	"context"
	"fmt"

	db "gadm-api/db"
	"gadm-api/logger"
)

type NextPageParams struct {
	Context       context.Context
	PgConn        *db.PgConn
	StartAt       int
	PageSize      int
	FilterColName string
	FilterVal     string
}

func GetNextPageFid(
	params NextPageParams,
) (int, error) {

	sql, args, err := db.GetNextFidSqlQuery(db.GetNextFidSqlQueryParams{
		StartAtFid:    params.StartAt,
		PageSize:      params.PageSize,
		FilterColName: params.FilterColName,
		FilterVal:     params.FilterVal,
	})

	var nextFid int
	err = params.PgConn.Db.QueryRow(params.Context, sql, args...).Scan(&nextFid)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		return 0, fmt.Errorf("failed_to_get_next_fid %v", err)
	}
	return nextFid, nil
}
