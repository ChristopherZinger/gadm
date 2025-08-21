package db

import (
	"context"
	"encoding/json"
	"fmt"
	"gadm-api/db"
	gadmUtils "gadm-api/utils"
)

type FeatureCollectionRepo struct {
	pgConn *db.PgConn
	ctx    context.Context
}

func NewFeatureCollectionRepo(
	pgConn *db.PgConn,
	ctx context.Context,
) *FeatureCollectionRepo {
	return &FeatureCollectionRepo{pgConn: pgConn, ctx: ctx}
}

type GetFeatureCollectionParams struct {
	GadmLevel     gadmUtils.GadmLevel
	FilterValue   string
	FilterColName string
	StartAtFid    int
	PageSize      int
}

func (repo *FeatureCollectionRepo) GetFeatureCollection(
	params GetFeatureCollectionParams,
) (json.RawMessage, error) {
	sql, args, err := db.BuildGadmFeatureCollectionSelectBuilder(db.GadmFeatureCollectionSelectBuilderParams{
		GadmLevel:     params.GadmLevel,
		FilterValue:   params.FilterValue,
		FilterColName: params.FilterColName,
		StartAtFid:    params.StartAtFid,
		PageSize:      params.PageSize,
	}).ToSql()

	var featureCollectionJSON json.RawMessage
	err = repo.pgConn.Db.QueryRow(repo.ctx, sql, args...).
		Scan(&featureCollectionJSON)
	if err != nil {
		return nil, fmt.Errorf("failed_to_query_database %v", err)
	}
	return featureCollectionJSON, nil
}
