package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	pgUrl := os.Getenv("DATABASE_URL")
	ctx := context.Background()

	dbPool, err := pgxpool.New(context.Background(), pgUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	http.HandleFunc("/api/v1/lv1", func(w http.ResponseWriter, r *http.Request) {
		take := 10
		takeStr := r.URL.Query().Get("take")
		_take, err := strconv.Atoi(takeStr)
		if err != nil {
			panic(err)
		}
		take = _take

		offsetStr := r.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			panic(err)
		}

		var opt SearchQueryOptions
		opt.Limit = take
		opt.Offset = offset

		countries, err := queryAdm0(dbPool, ctx, opt)
		if err != nil {
			panic(err)
		}

		jsonResponse, err := json.Marshal(countries)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type SearchQueryOptions struct {
	Limit  int
	Offset int
}

type Lv1QueryResult struct {
	Fid     int    `json:"fid"`
	Gid0    string `json:"gid_0"`
	Country string `json:"country"`
}

func queryAdm0(dbPool *pgxpool.Pool, ctx context.Context, opt SearchQueryOptions) ([]Lv1QueryResult, error) {
	sqlQuery := `select fid as Fid, gid_0 as Gid0, country as Country
		from adm_0 
		limit $1
		offset $2;`

	rows, _ := dbPool.Query(ctx, sqlQuery, opt.Limit, opt.Offset)

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[Lv1QueryResult])
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return result, nil
}
