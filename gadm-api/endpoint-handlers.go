package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gadm-api/logger"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const MIN_FID = 0

type GeoJSONFeature struct {
	Type       string                 `json:"type"`
	Geometry   json.RawMessage        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type GeoJSONFeatureCollection struct {
	Type     string           `json:"type"`
	Features []GeoJSONFeature `json:"features"`
}

type PaginationParams struct {
	Limit         int
	StartAfterFid int
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

func setFeatureCollectionResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func (s *Server) handleGeoJsonlLv0(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	paginationParams, err := getPaginationParams(r)
	if err != nil {
		logger.Error("failed_to_get_pagination_params %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	const MIN_LIMIT = 1
	const MAX_LIMIT = 20

	paginationParams.Limit = clamp(paginationParams.Limit, MIN_LIMIT, MAX_LIMIT)

	setGeojsonlStreamingResponseHeaders(w)

	err = s.queryAdmLv0GeoJsonl(ctx, w, paginationParams)
	if err != nil {
		logger.Error("failed_to_stream_geojsonl %v", err)
		return
	}
}

func (s *Server) handleGeoJsonlLv1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	paginationParams, err := getPaginationParams(r)
	if err != nil {
		logger.Error("failed_to_get_pagination_params %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	const MIN_LIMIT = 1
	const MAX_LIMIT = 20
	paginationParams.Limit = clamp(paginationParams.Limit, MIN_LIMIT, MAX_LIMIT)

	setGeojsonlStreamingResponseHeaders(w)

	err = s.queryAdmLv1GeoJsonl(ctx, w, paginationParams)
	if err != nil {
		logger.Error("failed_to_stream_geojsonl %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) queryAdmLv1GeoJsonl(ctx context.Context, w http.ResponseWriter, paginationParams PaginationParams) error {
	featurePropertiesSqlExpr := buildGeojsonFeaturePropertiesSqlExpression(
		Adm1.FID, Adm1.GID0, Adm1.Country, Adm1.GID1, Adm1.Name1, Adm1.Varname1,
		Adm1.NlName1, Adm1.Type1, Adm1.Engtype1, Adm1.Cc1, Adm1.Hasc1,
		Adm1.Iso1,
	)

	jsonBuildObjectFeature := fmt.Sprintf(
		`json_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(%[1]s)::json,
			'properties', %[2]s
		)`, Adm1.Geometry, featurePropertiesSqlExpr,
	)

	query := squirrel.Select(jsonBuildObjectFeature).
		From(ADM_1_TABLE).
		Where(squirrel.Expr(fmt.Sprintf("%s > $1", Adm1.FID), max(paginationParams.StartAfterFid, MIN_FID))).
		OrderBy(fmt.Sprintf("%s ASC", Adm1.FID)).
		Limit(uint64(paginationParams.Limit))

	sql, args, err := query.ToSql()
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

func (s *Server) handleFeatureCollectionLv0(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	paginationParams, err := getPaginationParams(r)
	if err != nil {
		logger.Error("failed_to_get_pagination_params %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	const MIN_LIMIT = 1
	const MAX_LIMIT = 3 // TODO: create configuration for limits per endpoint per gadm level;
	paginationParams.Limit = clamp(paginationParams.Limit, MIN_LIMIT, MAX_LIMIT)

	featureCollectionRawMsg, err := s.queryAdmLv0FeatureCollection(ctx, paginationParams)
	if err != nil {
		logger.Error("failed_to_query_feature_collection %v", err)
		panic(err)
	}

	setFeatureCollectionResponseHeaders(w)
	w.Write(featureCollectionRawMsg)

}

// TODO: improve logging;
func (s *Server) queryAdmLv0GeoJsonl(ctx context.Context, w http.ResponseWriter, paginationParams PaginationParams) error {
	jsonBuildObjectProperties := buildGeojsonFeaturePropertiesSqlExpression(Adm0.FID, Adm0.GID0, Adm0.Country)
	jsonBuildObject := fmt.Sprintf(
		`json_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(%[1]s)::json,
			'properties', %[2]s
		)`,
		Adm0.Geometry, jsonBuildObjectProperties,
	)

	query := squirrel.Select(jsonBuildObject).
		From(ADM_0_TABLE).
		Where(squirrel.Expr(fmt.Sprintf("%s > $1", Adm0.FID), max(paginationParams.StartAfterFid, MIN_FID))).
		OrderBy(fmt.Sprintf("%s ASC", Adm0.FID)).
		Limit(uint64(paginationParams.Limit))

	sql, args, err := query.ToSql()
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

func (s *Server) queryAdmLv0FeatureCollection(ctx context.Context, paginationParams PaginationParams) (json.RawMessage, error) {
	jsonBuildObject := fmt.Sprintf(
		`json_build_object(
				'type', 'FeatureCollection',
				'features', json_agg(
					json_build_object(
						'type', 'Feature',
						'geometry', ST_AsGeoJSON(%[1]s)::json,
						'properties', json_build_object(
							'%[2]s', %[2]s,
							'%[3]s', %[3]s,
							'%[4]s', %[4]s
						)
					)
				)
			)`,
		Adm0.Geometry, Adm0.FID, Adm0.GID0, Adm0.Country,
	)

	subQuery := squirrel.Select(fmt.Sprintf("%s, %s, %s, %s", Adm0.FID, Adm0.GID0, Adm0.Country, Adm0.Geometry)).
		From(ADM_0_TABLE).
		Where(squirrel.Expr(fmt.Sprintf("%s > $1", Adm0.FID), max(paginationParams.StartAfterFid, MIN_FID))).
		OrderBy(fmt.Sprintf("%s ASC", Adm0.FID)).
		Limit(uint64(paginationParams.Limit))

	mainQuery := squirrel.Select(jsonBuildObject).FromSelect(subQuery, "sub")

	sql, args, err := mainQuery.ToSql()
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	var featureCollectionJSON json.RawMessage
	err = s.db.QueryRow(ctx, sql, args...).Scan(&featureCollectionJSON)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	return featureCollectionJSON, nil
}

func buildGeojsonFeaturePropertiesSqlExpression(columns ...string) string {
	v := "json_build_object("
	len := len(columns)
	for i, column := range columns {
		v += fmt.Sprintf("'%s', %s", column, column)
		if i < len-1 {
			v += ", "
		}
	}
	v += ")"
	return v
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
