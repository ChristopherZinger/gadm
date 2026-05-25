package main

import (
	"context"
	"log"
	"net/http"
	"os"

	db "gadm-api/db"
	handlers "gadm-api/handlers"
	"gadm-api/logger"
	"gadm-api/models/adm"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DATABASE_URL_ENV_VAR = "DATABASE_URL"

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Warning("could_not_load_env_file %v", err)
	}

	switch os.Getenv("SERVICE_TYPE") {
	case "rest_api":
		startRestApi()
	case "cron_job":
		runCronJob()
	default:
		logger.Fatal("invalid_api_type %s", os.Getenv("API_TYPE"))
	}
}

const MAX_PG_CONNS = int32(45)

func initPgPool(poolSize int32) *pgxpool.Pool {
	pgUrl := os.Getenv(DATABASE_URL_ENV_VAR)
	if pgUrl == "" {
		logger.Fatal("missing_db_url_env_variable %s", DATABASE_URL_ENV_VAR) // TODO: expect var from config file util
	}

	cfg, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		logger.Fatal("failed_to_parse_db_url %v", err)
	}
	cfg.MaxConns = poolSize

	dbPool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		logger.Fatal("failed_to_connect_to_database %v", err)
	}

	return dbPool
}

func runCronJob() {
	dbPool := initPgPool(MAX_PG_CONNS)
	defer dbPool.Close()

	admRepo := adm.NewAdmRepo(dbPool)
	admService := adm.NewAdmService(admRepo)
	err := admService.PopulateAdmTree(context.Background())
	if err != nil {
		logger.Fatal("failed_to_populate_adm_tree %v", err)
	}
}

func startRestApi() {
	dbPool := initPgPool(MAX_PG_CONNS)
	defer dbPool.Close()
	pgConn := db.CreatePgConnector(dbPool)
	mux := http.NewServeMux()

	geojsonlHandlers, err := handlers.CreateGeojsonlHandlers(pgConn)
	if err != nil {
		logger.Fatal("failed_to_create_geojsonl_handlers %v", err)
	}
	for _, handlerInfo := range geojsonlHandlers {
		mux.HandleFunc(handlerInfo.Url, handlerInfo.Handler)
	}

	featureCollectionHandlers, err := handlers.CreateFeatureCollectionHandlers(pgConn)
	if err != nil {
		logger.Fatal("failed_to_create_feature_collection_handlers %v", err)
	}
	for _, handlerInfo := range featureCollectionHandlers {
		mux.HandleFunc(handlerInfo.Url, handlerInfo.Handler)
	}

	createAccessTokenHandlerInfo := handlers.GetAccessTokenCreationHandlerInfo(pgConn)
	mux.HandleFunc(createAccessTokenHandlerInfo.Url, createAccessTokenHandlerInfo.Handler)

	admRepo := adm.NewAdmRepo(dbPool)
	admService := adm.NewAdmService(admRepo)
	admHandler := adm.NewAdmNeighborsHandler(admService)
	mux.HandleFunc("/api/v1/adm-neighbors", admHandler.AdmNeighborsHandler)
	mux.HandleFunc("/api/v1/reverse-geocode", admHandler.AdmForLatLngHandler)

	handler := GetAuthMiddleWare(pgConn)(LoggingMiddleware(mux))

	logger.Info("server_starting_on_port_8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
