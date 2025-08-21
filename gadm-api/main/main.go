package main

import (
	"context"
	"log"
	"net/http"
	"os"

	db "gadm-api/db"
	handlers "gadm-api/handlers"
	"gadm-api/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DATABASE_URL_ENV_VAR = "DATABASE_URL"

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Warning("could_not_load_env_file %v", err)
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

	pgConn := db.CreatePgConnector(dbPool)
	mux := http.NewServeMux()

	geojsonlHandlers, err := handlers.CreateGeojsonlHandlers(pgConn)
	for _, handlerInfo := range geojsonlHandlers {
		mux.HandleFunc(handlerInfo.Url, handlerInfo.Handler)
	}

	featureCollectionHandlers, err := handlers.CreateFeatureCollectionHandlers(pgConn)
	for _, handlerInfo := range featureCollectionHandlers {
		mux.HandleFunc(handlerInfo.Url, handlerInfo.Handler)
	}

	createAccessTokenHandlerInfo := handlers.GetAccessTokenCreationHandlerInfo(pgConn)
	mux.HandleFunc(createAccessTokenHandlerInfo.Url, createAccessTokenHandlerInfo.Handler)

	reverseGeocodeHandlerInfo := handlers.GetReverseGeocodeHandlerInfo(pgConn)
	mux.HandleFunc(reverseGeocodeHandlerInfo.Url, reverseGeocodeHandlerInfo.Handler)

	handler := GetAuthMiddleWare(pgConn)(LoggingMiddleware(mux))

	logger.Info("server_starting_on_port_8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
