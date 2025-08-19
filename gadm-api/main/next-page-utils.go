package main

import (
	"context"
	"fmt"

	db "gadm-api/db"
	"gadm-api/logger"
)

func getNextFid(
	ctx context.Context,
	pgConn *PgConn,
	startAt int,
	pageSize int,
	filterParams db.SqlFilterParams) (int, error) {

	sql, args, err := db.GetNextFidSqlQuery(startAt, pageSize, filterParams)

	var nextFid int
	err = pgConn.db.QueryRow(ctx, sql, args...).Scan(&nextFid)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		return 0, fmt.Errorf("failed_to_get_next_fid %v", err)
	}
	return nextFid, nil
}
