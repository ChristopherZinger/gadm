package main

import (
	"context"
	"fmt"
	"gadm-api/logger"
)

func (s *Server) getNextFid(ctx context.Context, tableName string,
	orderByColumnName string, startAt int, pageSize int,
	filterParams SqlFilterParams) (int, error) {

	sql, args, err := getNextFidSqlQuery(tableName, orderByColumnName,
		startAt, pageSize, filterParams)

	var nextFid int
	err = s.db.QueryRow(ctx, sql, args...).Scan(&nextFid)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		return 0, fmt.Errorf("failed_to_get_next_fid %v", err)
	}
	return nextFid, nil
}
