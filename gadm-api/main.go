package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type Server struct {
	db *pgxpool.Pool
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

	server := newServer(dbPool)
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/geojsonl/lv1", server.handleGeoJsonlLv1)
	mux.HandleFunc("/api/v1/geojsonl/lv2", server.handleGeoJsonlLv2)

	mux.HandleFunc("/api/v1/fc/lv1", server.handleFeatureCollectionLv1)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func newServer(db *pgxpool.Pool) *Server {
	return &Server{db: db}
}
