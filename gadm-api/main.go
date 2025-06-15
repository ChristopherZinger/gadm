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
	"github.com/spf13/viper"
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

type QueryParams struct {
	Take   int
	Offset int
}

type GadmLvPaginationOptions struct {
	Limit         int
	StartAfterFid int
}

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	pgUrl := viper.GetString("DATABASE_URL")

	dbPool, err := pgxpool.New(context.Background(), pgUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	http.HandleFunc("/api/v1/geojsonl/lv1", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		take, err := expectIntParamInQuery(r, "take")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		startAfterFid, err := expectIntParamInQuery(r, "startAfter")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var opt GadmLvPaginationOptions
		opt.Limit = take
		opt.StartAfterFid = startAfterFid

		log.Printf("geojsonl/lv1. take: %d, startAfterFid: %d", take, startAfterFid)

		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		err = queryAdmLv1GeoJsonl(ctx, dbPool, w, opt)
		if err != nil {
			log.Printf("Error streaming GeoJSONL: %v", err)
			return
		}
	})

	http.HandleFunc("/api/v1/fc/lv1", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		take, err := expectIntParamInQuery(r, "take")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		startAfterFid, err := expectIntParamInQuery(r, "startAfter")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var opt GadmLvPaginationOptions
		opt.Limit = take
		opt.StartAfterFid = startAfterFid

		log.Printf("fc/lv1. take: %d, startAfterFid: %d", take, startAfterFid)

		featureCollectionRawMsg, err := queryAdmLv0FeatureCollection(ctx, dbPool, opt)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(featureCollectionRawMsg)

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// TODO: improve logging;
func queryAdmLv1GeoJsonl(ctx context.Context, dbPool *pgxpool.Pool, w http.ResponseWriter, opt GadmLvPaginationOptions) error {
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
		WHERE fid > $2
		LIMIT $1;
	`

	const MIN_LIMIT = 1
	const MAX_LIMIT = 20

	rows, err := dbPool.Query(ctx, sqlQuery, clamp(opt.Limit, MIN_LIMIT, MAX_LIMIT), max(opt.StartAfterFid, MIN_FID))
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

func queryAdmLv0FeatureCollection(ctx context.Context, dbPool *pgxpool.Pool, opt GadmLvPaginationOptions) (json.RawMessage, error) {
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
				WHERE fid > $2
				LIMIT $1
			) sub
		`

	const MIN_LIMIT = 1
	const MAX_LIMIT = 3

	var featureCollectionJSON json.RawMessage
	err := dbPool.QueryRow(ctx, sqlQuery, clamp(opt.Limit, MIN_LIMIT, MAX_LIMIT), max(opt.StartAfterFid, MIN_FID)).Scan(&featureCollectionJSON)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return featureCollectionJSON, nil
}

func expectIntParamInQuery(r *http.Request, paramName string) (int, error) {
	paramStrVal := r.URL.Query().Get(paramName)
	value, err := strconv.Atoi(paramStrVal)
	if err != nil {
		return 0, fmt.Errorf("invalid %s parameter: %w", paramName, err)
	}
	return value, nil
}
