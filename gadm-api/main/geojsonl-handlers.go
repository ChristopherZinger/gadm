package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "gadm-api/db"
	repo "gadm-api/db/repo"
	dbutils "gadm-api/db/utils"
	gadmUtils "gadm-api/utils"
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
	req          *http.Request
	writer       http.ResponseWriter
	ctx          context.Context
	gadmLevel    gadmUtils.GadmLevel
	config       GeojsonlHandlerQueryConfig
	nextPageRepo *dbutils.NextPageRepo
	geojsonlRepo *repo.GeojsonlRepo
}

type QueryLimitConfig struct {
	minLimit int
	maxLimit int
}

var supportedGadmLevelsForGeojsonl = []gadmUtils.GadmLevel{gadmUtils.GadmLevel0, gadmUtils.GadmLevel1, gadmUtils.GadmLevel2, gadmUtils.GadmLevel3, gadmUtils.GadmLevel4, gadmUtils.GadmLevel5}

var geojsonHandlerQueryConfig = map[gadmUtils.GadmLevel]GeojsonlHandlerQueryConfig{
	gadmUtils.GadmLevel0: {
		QueryLimitConfig: QueryLimitConfig{minLimit: 1, maxLimit: 20},
	},
	gadmUtils.GadmLevel1: {
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
		FilterableColumns: utils.ArrayToStrings(db.GidName0),
	},
	gadmUtils.GadmLevel2: {
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 50},
		FilterableColumns: utils.ArrayToStrings(db.GidName0, db.GidName1),
	},
	gadmUtils.GadmLevel3: {
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 50},
		FilterableColumns: utils.ArrayToStrings(db.GidName0, db.GidName1, db.GidName2),
	},
	gadmUtils.GadmLevel4: {
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 100},
		FilterableColumns: utils.ArrayToStrings(db.GidName0, db.GidName1, db.GidName2, db.GidName3),
	},
	gadmUtils.GadmLevel5: {
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 100},
		FilterableColumns: utils.ArrayToStrings(db.GidName0, db.GidName1, db.GidName2, db.GidName3, db.GidName4),
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

	filterParams := getSqlGidLevelFilterParamsFromUrl(
		handler.config.FilterableColumns,
		handler.req.URL.Query())
	pageSize := utils.Clamp(paginationParams.Limit,
		handler.config.QueryLimitConfig.minLimit,
		handler.config.QueryLimitConfig.maxLimit)
	startAtFid := utils.Max(paginationParams.StartAtFid, MIN_FID)

	nextUrl, err := handler.getNextPageUrl(startAtFid, pageSize, filterParams)
	if err != nil {
	} else {
		logger.Error("failed_to_get_next_url %v", err)
	}
	logger.Debug("_nexUrl %s", nextUrl)

	handler.setGeojsonlStreamingResponseHeaders(nextUrl)

	err = handler.queryAdmGeoJsonl(db.SqlAdmQueryParams{
		PageSize:             pageSize,
		StartAtFid:           startAtFid,
		GidLevelFilterParams: filterParams,
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
	filterParams db.GidLevelFilterParams,
) (string, error) {
	nextStartAtFid, err := handler.nextPageRepo.GetNextPageFid(dbutils.NextPageParams{
		StartAt:       startAtFid,
		PageSize:      pageSize,
		FilterColName: filterParams.GidLevelName,
		FilterVal:     filterParams.GidValue,
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

func (handler *GadmGeojsonlHandler) queryAdmGeoJsonl(queryParams db.SqlAdmQueryParams) error {
	rows, err := handler.geojsonlRepo.GetGeojsonl(repo.GetGeojsonlParams{
		GadmLevel:     handler.gadmLevel,
		FilterVal:     queryParams.GidValue,
		FilterColName: queryParams.GidLevelName,
		StartAtValue:  queryParams.StartAtFid,
		LimitValue:    queryParams.PageSize,
	})
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		return fmt.Errorf("failed to build sql query: %w", err)
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
	gadmLevel gadmUtils.GadmLevel,
) *GadmGeojsonlHandler {

	ctx := r.Context()
	nextPageRepo := dbutils.NewNextPageRepo(pgConn, ctx)
	geojsonlRepo := repo.NewGeojsonlRepo(pgConn, ctx)

	return &GadmGeojsonlHandler{
		req:          r,
		writer:       w,
		ctx:          ctx,
		config:       geojsonHandlerQueryConfig[gadmLevel],
		gadmLevel:    gadmLevel,
		nextPageRepo: nextPageRepo,
		geojsonlRepo: geojsonlRepo,
	}
}
