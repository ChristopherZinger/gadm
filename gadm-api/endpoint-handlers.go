package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
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

func (s *Server) handleGeoJsonlLv1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	paginationParams, err := getPaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	const MIN_LIMIT = 1
	const MAX_LIMIT = 20

	paginationParams.Limit = clamp(paginationParams.Limit, MIN_LIMIT, MAX_LIMIT)

	setGeojsonlStreamingResponseHeaders(w)

	err = s.queryAdmLv1GeoJsonl(ctx, w, paginationParams)
	if err != nil {
		log.Printf("Error streaming GeoJSONL: %v", err)
		return
	}
}

func (s *Server) handleFeatureCollectionLv1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	paginationParams, err := getPaginationParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	const MIN_LIMIT = 1
	const MAX_LIMIT = 3 // TODO: create configuration for limits per endpoint per gadm level;
	paginationParams.Limit = clamp(paginationParams.Limit, MIN_LIMIT, MAX_LIMIT)

	featureCollectionRawMsg, err := s.queryAdmLv0FeatureCollection(ctx, paginationParams)
	if err != nil {
		panic(err)
	}

	setFeatureCollectionResponseHeaders(w)
	w.Write(featureCollectionRawMsg)

}

// TODO: improve logging;
func (s *Server) queryAdmLv1GeoJsonl(ctx context.Context, w http.ResponseWriter, paginationParams PaginationParams) error {
	jsonBuildObject := fmt.Sprintf(
		`json_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(%[1]s)::json,
			'properties', json_build_object(
				'%[2]s', %[2]s,
				'%[3]s', %[3]s,
				'%[4]s', %[4]s
			)
		)`,
		Adm0.Geometry, Adm0.FID, Adm0.GID0, Adm0.Country,
	)

	query := squirrel.Select(jsonBuildObject).
		From(ADM_0_TABLE).
		Where(squirrel.Expr(fmt.Sprintf("%s > $1", Adm0.FID), max(paginationParams.StartAfterFid, MIN_FID))).
		OrderBy(fmt.Sprintf("%s ASC", Adm0.FID)).
		Limit(uint64(paginationParams.Limit))

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to query database: %w", err)
	}
	defer rows.Close()

	flusher, err := NewFlusher(w, ctx)
	if err != nil {
		return fmt.Errorf("failed to create flusher: %w", err)
	}

	for rows.Next() {
		var featureJSON json.RawMessage
		if err := rows.Scan(&featureJSON); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		if err := flusher.Flush(featureJSON); err != nil {
			return fmt.Errorf("failed to flush feature: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("row iteration error: %w", err)
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
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	var featureCollectionJSON json.RawMessage
	err = s.db.QueryRow(ctx, sql, args...).Scan(&featureCollectionJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	return featureCollectionJSON, nil
}
