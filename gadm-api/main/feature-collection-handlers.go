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
	QueryLimitConfig
	// filterable columns need to be ordered by importance level
	// in case of gid filtering it means from gid_0 down to gid_5
	// if user passes multiple filtering params in the url
	// the first matching filterable column will be
	// the only one used for filtering
	FilterableColumns []string
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
		QueryLimitConfig: QueryLimitConfig{minLimit: 1, maxLimit: 5},
	},
	GadmLevel1: {
		FilterableColumns: arrayToStrings(db.GidName0),
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 5},
	},
	GadmLevel2: {
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1),
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
	},
	GadmLevel3: {
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1, db.GidName2),
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
	},
	GadmLevel4: {
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1, db.GidName2, db.GidName3),
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 50},
	},
	GadmLevel5: {
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1, db.GidName2, db.GidName3, db.GidName4),
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

	sql, args, err := buildGadmFeatureCollectionSelectBuilder(
		handler.gadmLevel, filterParams.FilterVal, filterParams.FilterColName, startAtFid, pageSize).ToSql()

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
