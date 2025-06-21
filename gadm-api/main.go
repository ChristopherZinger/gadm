package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Server struct {
	db *pgxpool.Pool
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
		log.Printf("Make sure .env file exists in the project root directory")
		return
	}

	pgUrl := os.Getenv("DATABASE_URL")
	if pgUrl == "" {
		log.Fatal("DATABASE_URL not found in .env file or environment variables")
	}

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
