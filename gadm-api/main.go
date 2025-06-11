package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GeoJSONFeature struct {
	Type       string                 `json:"type"`
	Geometry   json.RawMessage        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type GeoJSONFeatureCollection struct {
	Type     string           `json:"type"`
	Features []GeoJSONFeature `json:"features"`
}

func main() {
	pgUrl := os.Getenv("DATABASE_URL")

	dbPool, err := pgxpool.New(context.Background(), pgUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	http.HandleFunc("/api/v1/lv1", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		take := 10
		takeStr := r.URL.Query().Get("take")
		_take, err := strconv.Atoi(takeStr)
		if err != nil {
			panic(err)
		}
		take = _take

		var offset int
		if r.URL.Query().Has("offset") {
			offsetStr := r.URL.Query().Get("offset")
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				panic(err)
			}
		}

		responseFormat := r.URL.Query().Get("response-format")
		if responseFormat != "" && responseFormat != "feature-collection" && responseFormat != "geojsonl" {
			http.Error(
				w,
				"Invalid response-format parameter. Must be either 'feature-collection' or 'geojsonl'",
				http.StatusBadRequest,
			)
			return
		}

		var opt SearchQueryOptions
		opt.Limit = take
		opt.Offset = offset
		opt.ResponseFormat = responseFormat

		if responseFormat == "geojsonl" {
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			err := queryAdmLv1GeoJsonl(ctx, dbPool, w, opt)
			if err != nil {
				log.Printf("Error streaming GeoJSONL: %v", err)
				return
			}
		} else {
			featureCollectionRawMsg, err := queryAdmLv0FeatureCollection(ctx, dbPool, opt)
			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(featureCollectionRawMsg)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type SearchQueryOptions struct {
	Limit          int
	Offset         int
	ResponseFormat string
}

// TODO: check for max limit and optimize offset; Also fix logging
func queryAdmLv1GeoJsonl(ctx context.Context, dbPool *pgxpool.Pool, w http.ResponseWriter, opt SearchQueryOptions) error {
	sqlQuery := `
		SELECT json_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(geom)::json,
			'properties', json_build_object(
				'fid', fid,
				'gid_0', gid_0,
				'country', country
			)
		)
		FROM adm_0 
		WHERE geom IS NOT NULL 
		LIMIT $1 OFFSET $2;
	`

	rows, err := dbPool.Query(ctx, sqlQuery, opt.Limit, opt.Offset)
	if err != nil {
		return fmt.Errorf("failed to query database: %w", err)
	}
	defer rows.Close()

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("Warning: ResponseWriter doesn't support flushing - data may be buffered")
	}

	for rows.Next() {
		var featureJSON json.RawMessage
		if err := rows.Scan(&featureJSON); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		if _, err := w.Write(featureJSON); err != nil {
			return fmt.Errorf("failed to write feature: %w", err)
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}

		if flusher != nil {
			flusher.Flush()
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("row iteration error: %w", err)
	}

	return nil
}

func queryAdmLv0FeatureCollection(ctx context.Context, dbPool *pgxpool.Pool, opt SearchQueryOptions) (json.RawMessage, error) {
	sqlQuery := `
			SELECT json_build_object(
				'type', 'FeatureCollection',
				'features', json_agg(
					json_build_object(
						'type', 'Feature',
						'geometry', ST_AsGeoJSON(geom)::json,
						'properties', json_build_object(
							'fid', fid,
							'gid_0', gid_0,
							'country', country
						)
					)
				)
			) as feature_collection
			FROM (
				SELECT fid, gid_0, country, geom 
				FROM adm_0 
				WHERE geom IS NOT NULL 
				LIMIT $1 OFFSET $2
			) sub
		`

	var featureCollectionJSON json.RawMessage
	err := dbPool.QueryRow(ctx, sqlQuery, opt.Limit, opt.Offset).Scan(&featureCollectionJSON)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return featureCollectionJSON, nil
}
