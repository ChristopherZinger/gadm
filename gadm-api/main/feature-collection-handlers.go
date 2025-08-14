package main

import (
	"context"
	"encoding/json"
	"fmt"
	db "gadm-api/db"
	"gadm-api/logger"
	"net/http"
)

type FeatureCollectionHandlerQueryConfig struct {
	TableName string
	GeoJSONFeatureConfig
	QueryLimitConfig
	// filterable columns need to be ordered by importance level
	// in case of gid filtering it means from gid_0 down to gid_5
	// if user passes multiple filtering params in the url
	// the first matching filterable column will be
	// the only one used for filtering
	FilterableColumns []string
	OrderByColumnName string
}

type GadmFeatureCollectionHandler struct {
	pgConn    *PgConn
	req       *http.Request
	writer    http.ResponseWriter
	ctx       context.Context
	gadmLevel GadmLevel
	config    GeojsonlHandlerQueryConfig
}

var featureCollectionHandlerQueryConfig = map[GadmLevel]FeatureCollectionHandlerQueryConfig{
	GadmLevel0: {
		TableName: db.ADM_0_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{db.Adm0.FID, db.Adm0.GID0, db.Adm0.Country},
			GeometryColumnName:     db.Adm0.Geometry,
		},
		OrderByColumnName: db.Adm0.FID,
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 5},
	},
	GadmLevel1: {
		TableName: db.ADM_1_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{db.Adm1.FID, db.Adm1.GID0, db.Adm1.Country,
				db.Adm1.GID1, db.Adm1.Name1, db.Adm1.Varname1, db.Adm1.NlName1, db.Adm1.Type1,
				db.Adm1.Engtype1, db.Adm1.Cc1, db.Adm1.Hasc1, db.Adm1.Iso1},
			GeometryColumnName: db.Adm1.Geometry,
		},
		OrderByColumnName: db.Adm1.FID,
		FilterableColumns: []string{db.Adm1.GID0},
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 5},
	},
	GadmLevel2: {
		TableName: db.ADM_2_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{db.Adm2.FID, db.Adm2.GID0, db.Adm2.Country,
				db.Adm2.GID1, db.Adm2.Name1, db.Adm2.NlName1, db.Adm2.GID2, db.Adm2.Name2,
				db.Adm2.Varname2, db.Adm2.NlName2, db.Adm2.Type2, db.Adm2.Engtype2, db.Adm2.Cc2,
				db.Adm2.Hasc2,
			},
			GeometryColumnName: db.Adm2.Geometry,
		},
		OrderByColumnName: db.Adm2.FID,
		FilterableColumns: []string{db.Adm2.GID0, db.Adm2.GID1},
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
	},
	GadmLevel3: {
		TableName: db.ADM_3_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{db.Adm3.FID, db.Adm3.GID0, db.Adm3.Country,
				db.Adm3.GID1, db.Adm3.Name1, db.Adm3.NlName1, db.Adm3.GID2, db.Adm3.Name2,
				db.Adm3.NlName2, db.Adm3.GID3, db.Adm3.Name3, db.Adm3.Varname3,
				db.Adm3.NlName3, db.Adm3.Type3, db.Adm3.Engtype3, db.Adm3.Cc3,
				db.Adm3.Hasc3,
			},
			GeometryColumnName: db.Adm3.Geometry,
		},
		OrderByColumnName: db.Adm3.FID,
		FilterableColumns: []string{db.Adm3.GID0, db.Adm3.GID1, db.Adm3.GID2},
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
	},
	GadmLevel4: {
		TableName: db.ADM_4_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{db.Adm4.FID, db.Adm4.GID0, db.Adm4.Country,
				db.Adm4.GID1, db.Adm4.Name1, db.Adm4.GID2, db.Adm4.Name2, db.Adm4.GID3,
				db.Adm4.Name3, db.Adm4.GID4, db.Adm4.Name4, db.Adm4.Varname4, db.Adm4.Type4,
				db.Adm4.Engtype4, db.Adm4.Cc4,
			},
			GeometryColumnName: db.Adm4.Geometry,
		},
		OrderByColumnName: db.Adm4.FID,
		FilterableColumns: []string{db.Adm4.GID0, db.Adm4.GID1, db.Adm4.GID2, db.Adm4.GID3},
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 50},
	},
	GadmLevel5: {
		TableName: db.ADM_5_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{db.Adm5.FID, db.Adm5.GID0, db.Adm5.Country,
				db.Adm5.GID1, db.Adm5.Name1, db.Adm5.GID2, db.Adm5.Name2, db.Adm5.GID3,
				db.Adm5.Name3, db.Adm5.GID4, db.Adm5.Name4, db.Adm5.GID5, db.Adm5.Name5,
				db.Adm5.Type5, db.Adm5.Engtype5, db.Adm5.Cc5,
			},
			GeometryColumnName: db.Adm5.Geometry,
		},
		OrderByColumnName: db.Adm5.FID,
		FilterableColumns: []string{db.Adm5.GID0, db.Adm5.GID1, db.Adm5.GID2, db.Adm5.GID3, db.Adm5.GID4},
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 50},
	},
}

var supportedGadmLevelsForFeatureCollection = []GadmLevel{GadmLevel0, GadmLevel1, GadmLevel2, GadmLevel3, GadmLevel4, GadmLevel5}

func (handler *GadmFeatureCollectionHandler) setFeatureCollectionResponseHeaders(nextUrl string) {
	handler.writer.Header().Set("Content-Type", "application/json")
	if nextUrl != "" {
		handler.writer.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"next\"", nextUrl))
	}
}

func CreateFeatureCollectionHandlers(s *PgConn) ([]HandlerInfo, error) {
	handlerInfos := []HandlerInfo{}
	for _, gadmLevel := range supportedGadmLevelsForFeatureCollection {
		url := getFeatureCollectionUrl(gadmLevel)
		handler := func(w http.ResponseWriter, r *http.Request) {
			newGadmFeatureCollectionHandler(s, r, w, gadmLevel).handle()
		}
		handlerInfos = append(handlerInfos, HandlerInfo{url: url, handler: handler})
	}
	return handlerInfos, nil
}

func (handler *GadmFeatureCollectionHandler) handle() {
	paginationParams, err := getPaginationParams(handler.req)
	if err != nil {
		logger.Error("failed_to_get_pagination_params %v", err)
		http.Error(handler.writer, err.Error(), http.StatusBadRequest)
		return
	}

	filterParams := getSqlFilterParamsFromUrl(
		handler.config.FilterableColumns,
		handler.req.URL.Query())
	pageSize := clamp(paginationParams.Limit,
		handler.config.QueryLimitConfig.minLimit,
		handler.config.QueryLimitConfig.maxLimit)
	startAtFid := max(paginationParams.StartAtFid, MIN_FID)

	nextUrl, err := handler.getNextPageUrl(startAtFid, pageSize, filterParams)
	if err != nil {
		logger.Error("failed_to_get_next_fid %v", err)
	}

	sql, args, err := buildFeatureCollectionSqlQuery(
		handler.gadmLevel,
		SqlQueryParams{ // todo: remove gadm level on for getting config, config should be part of the handler
			StartAtValue:    startAtFid,
			LimitValue:      pageSize,
			SqlFilterParams: filterParams,
		})
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		http.Error(handler.writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var featureCollectionJSON json.RawMessage
	err = handler.pgConn.db.QueryRow(handler.ctx, sql, args...).
		Scan(&featureCollectionJSON)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		http.Error(handler.writer, err.Error(), http.StatusInternalServerError)
		return
	}

	handler.setFeatureCollectionResponseHeaders(nextUrl)
	handler.writer.Write(featureCollectionJSON)
}

func (handler *GadmFeatureCollectionHandler) getNextPageUrl(startAtFid int, pageSize int, filterParams SqlFilterParams) (string, error) {
	nextStartAtFid, err := getNextFid(
		handler.ctx,
		handler.pgConn,
		handler.config.TableName,
		handler.config.OrderByColumnName,
		startAtFid,
		pageSize,
		filterParams,
	)
	if err != nil {
		return "", fmt.Errorf("failed_to_get_next_start_at_fid %v", err)
	}

	nextUrlParams, err := getNextPageUrlQueryParams(
		nextStartAtFid,
		pageSize,
		filterParams)
	if err != nil {
		return "", fmt.Errorf("failed_to_get_next_page_url_query_params %v", err)
	}
	return getGeojsonlUrl(handler.gadmLevel, nextUrlParams...), nil
}

func (handler *GadmFeatureCollectionHandler) getNextFid(
	startAtFid int,
	pageSize int,
	filterParams SqlFilterParams) (int, error) {

	return getNextFid(
		handler.ctx,
		handler.pgConn,
		handler.config.TableName,
		handler.config.OrderByColumnName,
		startAtFid,
		pageSize,
		filterParams,
	)
}

func newGadmFeatureCollectionHandler(
	pgConn *PgConn,
	r *http.Request,
	w http.ResponseWriter,
	gadmLevel GadmLevel,
) *GadmFeatureCollectionHandler {

	return &GadmFeatureCollectionHandler{
		pgConn:    pgConn,
		req:       r,
		writer:    w,
		ctx:       r.Context(),
		config:    geojsonHandlerQueryConfig[gadmLevel],
		gadmLevel: gadmLevel,
	}
}
