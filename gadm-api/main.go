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

		var opt SearchQueryOptions
		opt.Limit = take
		opt.Offset = offset

		featureCollectionRawMsg, err := queryAdmLv0FeatureCollection(ctx, dbPool, opt)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(featureCollectionRawMsg)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type SearchQueryOptions struct {
	Limit  int
	Offset int
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
