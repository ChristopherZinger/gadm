package main

import (
	"log"
	"net/http"
	"os"

	db "gadm-api/db"
	gameloop "gadm-api/game-loop"
	handlers "gadm-api/handlers"
	"gadm-api/infra/pg"
	"gadm-api/jobs"
	"gadm-api/logger"
	"gadm-api/models/adm"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Warning("could_not_load_env_file %v", err)
	}

	switch os.Getenv("SERVICE_TYPE") {
	case "rest_api":
		startRestApi()
	case "cron_job":
		jobName := os.Getenv("CRON_JOB_NAME")
		if jobName == "" {
			logger.Fatal("missing_cron_job_name")
		}
		switch jobName {
		case "populate_adm_tree":
			jobs.PopulateAdmTreeJob()
		case "populate_adm_neighbors":
			jobs.PopulateAdmNeighborsJob()
		default:
			logger.Fatal("unknown_cron_job_name %s", jobName)
		}
	case "game_loop":
		gameloop.GameLoop()
	default:
		logger.Fatal("invalid_api_type %s", os.Getenv("API_TYPE"))
	}
}

const MAX_PG_CONNS = int32(45)

func startRestApi() {
	dbPool := pg.InitPgPool(MAX_PG_CONNS)
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
