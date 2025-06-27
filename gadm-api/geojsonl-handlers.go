package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gadm-api/logger"

	"github.com/jackc/pgx/v5"
)

const MIN_FID = 0

type PaginationParams struct {
	Limit         int
	StartAfterFid int
}

type GeojsonlHandlerInfo struct {
	queryParams GeoJsonFeatureSqlQueryParams
	queryLimits QueryLimits
}

type QueryLimits struct {
	minLimit int
	maxLimit int
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

var geojsonEndpointInfo = map[GadmLevel]GeojsonlHandlerInfo{
	GadmLevel0: {
		queryParams: GeoJsonFeatureSqlQueryParams{
			TableName:              ADM_0_TABLE,
			FeaturePropertiesNames: []string{Adm0.FID, Adm0.GID0, Adm0.Country},
			GeometryColumnName:     Adm0.Geometry,
			OrderByColumnName:      Adm0.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 20},
	},
	GadmLevel1: {
		queryParams: GeoJsonFeatureSqlQueryParams{
			TableName: ADM_1_TABLE,
			FeaturePropertiesNames: []string{Adm1.FID, Adm1.GID0, Adm1.Country,
				Adm1.GID1, Adm1.Name1, Adm1.Varname1, Adm1.NlName1, Adm1.Type1,
				Adm1.Engtype1, Adm1.Cc1, Adm1.Hasc1, Adm1.Iso1,
			},
			GeometryColumnName: Adm1.Geometry,
			OrderByColumnName:  Adm1.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 20},
	},
	GadmLevel2: {
		queryParams: GeoJsonFeatureSqlQueryParams{
			TableName: ADM_2_TABLE,
			FeaturePropertiesNames: []string{Adm2.FID, Adm2.GID0, Adm2.Country,
				Adm2.GID1, Adm2.Name1, Adm2.NlName1, Adm2.GID2, Adm2.Name2,
				Adm2.Varname2, Adm2.NlName2, Adm2.Type2, Adm2.Engtype2, Adm2.Cc2,
				Adm2.Hasc2,
			},
			GeometryColumnName: Adm2.Geometry,
			OrderByColumnName:  Adm2.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 20},
	},
	GadmLevel3: {
		queryParams: GeoJsonFeatureSqlQueryParams{
			TableName: ADM_3_TABLE,
			FeaturePropertiesNames: []string{Adm3.FID, Adm3.GID0, Adm3.Country,
				Adm3.GID1, Adm3.Name1, Adm3.NlName1, Adm3.GID2, Adm3.Name2,
				Adm3.NlName2, Adm3.GID3, Adm3.Name3, Adm3.Varname3,
				Adm3.NlName3, Adm3.Type3, Adm3.Engtype3, Adm3.Cc3,
				Adm3.Hasc3,
			},
			GeometryColumnName: Adm3.Geometry,
			OrderByColumnName:  Adm3.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 20},
	},
	GadmLevel4: {
		queryParams: GeoJsonFeatureSqlQueryParams{
			TableName: ADM_4_TABLE,
			FeaturePropertiesNames: []string{Adm4.FID, Adm4.GID0, Adm4.Country,
				Adm4.GID1, Adm4.Name1, Adm4.GID2, Adm4.Name2, Adm4.GID3,
				Adm4.Name3, Adm4.GID4, Adm4.Name4, Adm4.Varname4, Adm4.Type4,
				Adm4.Engtype4, Adm4.Cc4,
			},
			GeometryColumnName: Adm4.Geometry,
			OrderByColumnName:  Adm4.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 20},
	},
	GadmLevel5: {
		queryParams: GeoJsonFeatureSqlQueryParams{
			TableName: ADM_5_TABLE,
			FeaturePropertiesNames: []string{Adm5.FID, Adm5.GID0, Adm5.Country,
				Adm5.GID1, Adm5.Name1, Adm5.GID2, Adm5.Name2, Adm5.GID3,
				Adm5.Name3, Adm5.GID4, Adm5.Name4, Adm5.GID5, Adm5.Name5,
				Adm5.Type5, Adm5.Engtype5, Adm5.Cc5,
			},
			GeometryColumnName: Adm5.Geometry,
			OrderByColumnName:  Adm5.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 20},
	},
}

func getPaginationParams(r *http.Request) (PaginationParams, error) {
	take, err := expectIntParamInQuery(r, "take", 10)
	if err != nil {
		return PaginationParams{}, fmt.Errorf("failed to parse query parameter 'take': %w", err)
	}

	startAfterFid, err := expectIntParamInQuery(r, "startAfter", 0)
	if err != nil {
		return PaginationParams{}, fmt.Errorf("failed to parse query parameter 'startAfter': %w", err)
	}
	return PaginationParams{
		Limit:         take,
		StartAfterFid: startAfterFid,
	}, nil
}

func setGeojsonlStreamingResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

type HandlerInfo struct {
	url     string
	handler func(w http.ResponseWriter, r *http.Request)
}

func CreateGeojsonlHandlers(s *Server) ([]HandlerInfo, error) {
	handlerInfos := []HandlerInfo{}
	for _, gadmLevel := range supportedGadmLevelsForGeojsonl {
		url := getGeojsonlUrl(gadmLevel)
		handler := func(w http.ResponseWriter, r *http.Request) {
			s.handleGeoJsonl(w, r, gadmLevel)
		}
		handlerInfos = append(handlerInfos, HandlerInfo{url: url, handler: handler})
	}
	return handlerInfos, nil
}

func (s *Server) handleGeoJsonl(w http.ResponseWriter, r *http.Request, gadmLevel GadmLevel) {
	ctx := r.Context()

	paginationParams, err := getPaginationParams(r)
	if err != nil {
		logger.Error("failed_to_get_pagination_params %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	setGeojsonlStreamingResponseHeaders(w)

	params := geojsonEndpointInfo[gadmLevel].queryParams
	limits := geojsonEndpointInfo[gadmLevel].queryLimits
	err = s.queryAdmGeoJsonl(ctx, w, GeoJsonFeatureSqlQueryParams{
		TableName:              params.TableName,
		FeaturePropertiesNames: params.FeaturePropertiesNames,
		GeometryColumnName:     params.GeometryColumnName,
		OrderByColumnName:      params.OrderByColumnName,
		StartAfterValue:        paginationParams.StartAfterFid,
		LimitValue:             clamp(paginationParams.Limit, limits.minLimit, limits.maxLimit),
	})
	if err != nil {
		logger.Error("failed_to_stream_geojsonl %v", err)
		return
	}
}

func (s *Server) queryAdmGeoJsonl(ctx context.Context, w http.ResponseWriter, queryParams GeoJsonFeatureSqlQueryParams) error {
	sql, args, err := buildGeojsonFeatureSqlQuery(queryParams)
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		return fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		return fmt.Errorf("failed to query database: %w", err)
	}
	defer rows.Close()

	err = streamRows(rows, w, ctx)
	if err != nil {
		logger.Error("failed_to_stream_rows %v", err)
		return fmt.Errorf("failed to stream rows: %w", err)
	}

	return nil
}

func streamRows(rows pgx.Rows, w http.ResponseWriter, ctx context.Context) error {
	flusher, err := NewFlusher(w, ctx)
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
