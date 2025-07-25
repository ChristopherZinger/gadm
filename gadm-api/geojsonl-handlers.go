package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gadm-api/logger"

	"github.com/jackc/pgx/v5"
)

const MIN_FID = 1

type GeojsonlHandlerQueryConfig struct {
	GadmLevel GadmLevel
	TableName string
	GeoJSONFeatureConfig
	QueryLimitConfig
	FilterableColumns []string
	OrderByColumnName string
}

type GadmGeojsonlHandler struct {
	pgConn    *PgConn
	req       *http.Request
	writer    http.ResponseWriter
	ctx       context.Context
	gadmLevel GadmLevel
	config    GeojsonlHandlerQueryConfig
}

type QueryLimitConfig struct {
	minLimit int
	maxLimit int
}

type GeoJSONFeatureConfig struct {
	FeaturePropertiesNames []string
	GeometryColumnName     string
}

type GadmLevel int

const (
	GadmLevel0 GadmLevel = iota
	GadmLevel1
	GadmLevel2
	GadmLevel3
	GadmLevel4
	GadmLevel5
)

var supportedGadmLevelsForGeojsonl = []GadmLevel{GadmLevel0, GadmLevel1, GadmLevel2, GadmLevel3, GadmLevel4, GadmLevel5}

var geojsonHandlerQueryConfig = map[GadmLevel]GeojsonlHandlerQueryConfig{
	GadmLevel0: {
		TableName: ADM_0_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm0.FID, Adm0.GID0, Adm0.Country},
			GeometryColumnName:     Adm0.Geometry},
		OrderByColumnName: Adm0.FID,
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
	},
	GadmLevel1: {
		TableName: ADM_1_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm1.FID, Adm1.GID0, Adm1.Country,
				Adm1.GID1, Adm1.Name1, Adm1.Varname1, Adm1.NlName1, Adm1.Type1,
				Adm1.Engtype1, Adm1.Cc1, Adm1.Hasc1, Adm1.Iso1,
			},
			GeometryColumnName: Adm1.Geometry,
		},
		OrderByColumnName: Adm1.FID,
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
		FilterableColumns: []string{Adm1.GID0},
	},
	GadmLevel2: {
		TableName: ADM_2_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm2.FID, Adm2.GID0, Adm2.Country,
				Adm2.GID1, Adm2.Name1, Adm2.NlName1, Adm2.GID2, Adm2.Name2,
				Adm2.Varname2, Adm2.NlName2, Adm2.Type2, Adm2.Engtype2, Adm2.Cc2,
				Adm2.Hasc2,
			},
			GeometryColumnName: Adm2.Geometry,
		},
		OrderByColumnName: Adm2.FID,
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
		FilterableColumns: []string{Adm2.GID0, Adm2.GID1},
	},
	GadmLevel3: {
		TableName: ADM_3_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm3.FID, Adm3.GID0, Adm3.Country,
				Adm3.GID1, Adm3.Name1, Adm3.NlName1, Adm3.GID2, Adm3.Name2,
				Adm3.NlName2, Adm3.GID3, Adm3.Name3, Adm3.Varname3,
				Adm3.NlName3, Adm3.Type3, Adm3.Engtype3, Adm3.Cc3,
				Adm3.Hasc3,
			},
			GeometryColumnName: Adm3.Geometry,
		},
		OrderByColumnName: Adm3.FID,
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
		FilterableColumns: []string{Adm3.GID0, Adm3.GID1, Adm3.GID2},
	},
	GadmLevel4: {
		TableName: ADM_4_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm4.FID, Adm4.GID0, Adm4.Country,
				Adm4.GID1, Adm4.Name1, Adm4.GID2, Adm4.Name2, Adm4.GID3,
				Adm4.Name3, Adm4.GID4, Adm4.Name4, Adm4.Varname4, Adm4.Type4,
				Adm4.Engtype4, Adm4.Cc4,
			},
			GeometryColumnName: Adm4.Geometry,
		},
		OrderByColumnName: Adm4.FID,
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
		FilterableColumns: []string{Adm4.GID0, Adm4.GID1, Adm4.GID2, Adm4.GID3},
	},
	GadmLevel5: {
		TableName: ADM_5_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm5.FID, Adm5.GID0, Adm5.Country,
				Adm5.GID1, Adm5.Name1, Adm5.GID2, Adm5.Name2, Adm5.GID3,
				Adm5.Name3, Adm5.GID4, Adm5.Name4, Adm5.GID5, Adm5.Name5,
				Adm5.Type5, Adm5.Engtype5, Adm5.Cc5,
			},
			GeometryColumnName: Adm5.Geometry,
		},
		OrderByColumnName: Adm5.FID,
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
		FilterableColumns: []string{Adm5.GID0, Adm5.GID1, Adm5.GID2, Adm5.GID3,
			Adm5.GID4},
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

func CreateGeojsonlHandlers(pgConn *PgConn) ([]HandlerInfo, error) {
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

	err = handler.queryAdmGeoJsonl(SqlQueryParams{
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
	filterParams SqlFilterParams) (string, error) {

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

	nextUrlParams, err := getNextPageUrlQueryParams(nextStartAtFid, pageSize, filterParams)
	if err != nil {
		return "", fmt.Errorf("failed_to_get_next_page_url_query_params %v", err)
	}
	return getGeojsonlUrl(handler.gadmLevel, nextUrlParams...), nil
}

func (handler *GadmGeojsonlHandler) queryAdmGeoJsonl(queryParams SqlQueryParams) error {
	sql, args, err := buildGeojsonFeatureSqlQuery(handler.gadmLevel, queryParams)
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		return fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := handler.pgConn.db.Query(handler.ctx, sql, args...)
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
	pgConn *PgConn,
	r *http.Request,
	w http.ResponseWriter,
	gadmLevel GadmLevel,
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
