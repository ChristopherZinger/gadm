package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"gadm-api/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type PgConn struct {
	db *pgxpool.Pool
}

var DATABASE_URL_ENV_VAR = "DATABASE_URL"

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Warning("could_not_load_env_file %v", err)
		return
	}

	pgUrl := os.Getenv(DATABASE_URL_ENV_VAR)
	if pgUrl == "" {
		logger.Fatal("missing_db_url_env_variable %s", DATABASE_URL_ENV_VAR) // TODO: expect var from config file util
	}

	dbPool, err := pgxpool.New(context.Background(), pgUrl)
	if err != nil {
		logger.Error("failed_to_connect_to_database %v", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	pgConn := getPgConnector(dbPool)
	mux := http.NewServeMux()

	geojsonlHandlers, err := CreateGeojsonlHandlers(pgConn)
	for _, handlerInfo := range geojsonlHandlers {
		mux.HandleFunc(handlerInfo.url, handlerInfo.handler)
	}

	featureCollectionHandlers, err := CreateFeatureCollectionHandlers(pgConn)
	for _, handlerInfo := range featureCollectionHandlers {
		mux.HandleFunc(handlerInfo.url, handlerInfo.handler)
	}

	createAccessTokenUrl := fmt.Sprintf("%screate-access-token", getBaseApiUrl().Path)
	mux.HandleFunc(createAccessTokenUrl, func(w http.ResponseWriter, r *http.Request) {
		NewAccessTokenCreationHandler(pgConn, r, w, r.Context()).handle()
	})

	handler := GetAuthMiddleWare(pgConn)(LoggingMiddleware(mux))

	logger.Info("server_starting_on_port_8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func getPgConnector(db *pgxpool.Pool) *PgConn {
	return &PgConn{db: db}
}
