package db

import (
	"context"
	"encoding/json"
	"fmt"
	"gadm-api/db"
	"gadm-api/logger"
)

type ReverseGeocodeRepo struct {
	pgConn *db.PgConn
	ctx    context.Context
}

func NewReverseGeocodeRepo(
	pgConn *db.PgConn,
	ctx context.Context,
) *ReverseGeocodeRepo {
	return &ReverseGeocodeRepo{pgConn: pgConn, ctx: ctx}
}

type GetReverseGeocodeParams struct {
	Lat float64
	Lng float64
}

func (repo *ReverseGeocodeRepo) GetLocation(
	params GetReverseGeocodeParams,
) (json.RawMessage, error) {
	sql, args, err := db.GetReverseGeocodeSqlQuery(
		db.Point{Lat: params.Lat, Lng: params.Lng})

	if err != nil {
		logger.Error(
			"failed_to_build_reverse_geocode_sql %v",
			err,
		)
		return nil, fmt.Errorf("failed_to_build_reverse_geocode_sql %v", err)
	}

	var jsonResult []byte
	err = repo.pgConn.Db.QueryRow(
		repo.ctx,
		sql,
		args...,
	).Scan(&jsonResult)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		return nil, fmt.Errorf("no_result_found")
	}

	return jsonResult, nil
}
