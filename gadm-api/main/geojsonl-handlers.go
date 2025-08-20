package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "gadm-api/db"
	dbutils "gadm-api/db/utils"
	utils "gadm-api/utils"

	"gadm-api/logger"

	"github.com/jackc/pgx/v5"
)

const MIN_FID = 1

type GeojsonlHandlerQueryConfig struct {
	QueryLimitConfig
	FilterableColumns []string
}

type GadmGeojsonlHandler struct {
	pgConn    *db.PgConn
	req       *http.Request
	writer    http.ResponseWriter
	ctx       context.Context
	gadmLevel utils.GadmLevel
	config    GeojsonlHandlerQueryConfig
}

type QueryLimitConfig struct {
	minLimit int
	maxLimit int
}

var supportedGadmLevelsForGeojsonl = []utils.GadmLevel{utils.GadmLevel0, utils.GadmLevel1, utils.GadmLevel2, utils.GadmLevel3, utils.GadmLevel4, utils.GadmLevel5}

var geojsonHandlerQueryConfig = map[utils.GadmLevel]GeojsonlHandlerQueryConfig{
	utils.GadmLevel0: {
		QueryLimitConfig: QueryLimitConfig{minLimit: 1, maxLimit: 20},
	},
	utils.GadmLevel1: {
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
		FilterableColumns: arrayToStrings(db.GidName0),
	},
	utils.GadmLevel2: {
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 50},
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1),
	},
	utils.GadmLevel3: {
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 50},
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1, db.GidName2),
	},
	utils.GadmLevel4: {
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 100},
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1, db.GidName2, db.GidName3),
	},
	utils.GadmLevel5: {
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 100},
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1, db.GidName2, db.GidName3, db.GidName4),
	},
}

func (handler *GadmGeojsonlHandler) setGeojsonlStreamingResponseHeaders(nextUrl string) {
	handler.writer.Header().Set("Content-Type", "application/x-ndjson")
	handler.writer.Header().Set("Cache-Control", "no-cache")
	handler.writer.Header().Set("Connection", "keep-alive")
	if nextUrl != "" {
		handler.writer.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"next\"", nextUrl))
	}
}

func CreateGeojsonlHandlers(pgConn *db.PgConn) ([]HandlerInfo, error) {
	handlerInfos := []HandlerInfo{}
	for _, gadmLevel := range supportedGadmLevelsForGeojsonl {
		url := getGeojsonlUrl(gadmLevel)
		handler := func(w http.ResponseWriter, r *http.Request) {
			handler := newGadmGeojsonlHandler(pgConn, r, w, gadmLevel)
			handler.handle()
		}
		handlerInfos = append(handlerInfos, HandlerInfo{url: url, handler: handler})
	}
	return handlerInfos, nil
}

func (handler *GadmGeojsonlHandler) handle() {
	paginationParams, err := getPaginationParams(handler.req)
	if err != nil {
		logger.Error("failed_to_get_pagination_params %v", err)
		http.Error(handler.writer, "invalid_query_parameters", http.StatusBadRequest)
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
	} else {
		logger.Error("failed_to_get_next_url %v", err)
	}
	logger.Debug("_nexUrl %s", nextUrl)

	handler.setGeojsonlStreamingResponseHeaders(nextUrl)

	err = handler.queryAdmGeoJsonl(db.SqlQueryParams{
		LimitValue:      pageSize,
		StartAtValue:    startAtFid,
		SqlFilterParams: filterParams,
	})
	if err != nil {
		logger.Error("failed_to_stream_geojsonl %v", err)
		http.Error(
			handler.writer,
			"failed_to_stream_geojsonl",
			http.StatusInternalServerError)
		return
	}
}

func (handler *GadmGeojsonlHandler) getNextPageUrl(
	startAtFid int,
	pageSize int,
	filterParams db.SqlFilterParams) (string, error) {

	nextStartAtFid, err := dbutils.GetNextPageFid(dbutils.NextPageParams{
		Context:       handler.ctx,
		PgConn:        handler.pgConn,
		StartAt:       startAtFid,
		PageSize:      pageSize,
		FilterColName: filterParams.FilterColName,
		FilterVal:     filterParams.FilterVal,
	})
	if err != nil {
		return "", fmt.Errorf("failed_to_get_next_start_at_fid %v", err)
	}

	nextUrlParams, err := getNextPageUrlQueryParams(nextStartAtFid, pageSize, filterParams)
	if err != nil {
		return "", fmt.Errorf("failed_to_get_next_page_url_query_params %v", err)
	}
	return getGeojsonlUrl(handler.gadmLevel, nextUrlParams...), nil
}

func (handler *GadmGeojsonlHandler) queryAdmGeoJsonl(queryParams db.SqlQueryParams) error {
	sql, args, err := db.BuildGeojsonFeatureSqlQuery(
		handler.gadmLevel,
		queryParams.FilterVal,
		queryParams.FilterColName,
		queryParams.StartAtValue,
		queryParams.LimitValue,
	)
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		return fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := handler.pgConn.Db.Query(handler.ctx, sql, args...)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		return fmt.Errorf("failed to query database: %w", err)
	}
	defer rows.Close()

	err = handler.streamRows(rows)
	if err != nil {
		logger.Error("failed_to_stream_rows %v", err)
		return fmt.Errorf("failed to stream rows: %w", err)
	}

	return nil
}

func (handler *GadmGeojsonlHandler) streamRows(rows pgx.Rows) error {
	flusher, err := NewFlusher(handler.writer, handler.ctx)
	if err != nil {
		logger.Error("failed_to_create_flusher %v", err)
		return fmt.Errorf("failed to create flusher: %w", err)
	}

	for rows.Next() {
		var featureJSON json.RawMessage
		if err := rows.Scan(&featureJSON); err != nil {
			logger.Error("failed_to_scan_row %v", err)
			return fmt.Errorf("failed to scan row: %w", err)
		}

		if err := flusher.Flush(featureJSON); err != nil {
			logger.Error("failed_to_flush_feature %v", err)
			return fmt.Errorf("failed to flush feature: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed_to_iterate_rows %v", err)
		return fmt.Errorf("row iteration error: %w", err)
	}

	return nil
}

func newGadmGeojsonlHandler(
	pgConn *db.PgConn,
	r *http.Request,
	w http.ResponseWriter,
	gadmLevel utils.GadmLevel,
) *GadmGeojsonlHandler {

	return &GadmGeojsonlHandler{
		pgConn:    pgConn,
		req:       r,
		writer:    w,
		ctx:       r.Context(),
		config:    geojsonHandlerQueryConfig[gadmLevel],
		gadmLevel: gadmLevel,
	}
}
