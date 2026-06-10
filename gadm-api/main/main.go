package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	db "gadm-api/db"
	gameloop "gadm-api/game-loop"
	"gadm-api/handlers"
	"gadm-api/infra/pg"
	"gadm-api/jobs"
	"gadm-api/logger"
	"gadm-api/models/adm"

	"github.com/jackc/pgx/v5/pgxpool"
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

	baseApiPath := "/api/v1"
	mux.Handle(baseApiPath+"/", http.StripPrefix(
		baseApiPath,
		getApiHandlers(pgConn, dbPool, baseApiPath),
	))

	mux.Handle("/ws", http.HandlerFunc(getWebsocketHandler))

	handler := LoggingMiddleware(mux)

	logger.Info("server_starting_on_port_8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func getApiHandlers(pgConn *db.PgConn, dbPool *pgxpool.Pool, baseApiPath string) http.Handler {
	mux := http.NewServeMux()

	createAccessTokenHandlerInfo := handlers.GetAccessTokenCreationHandlerInfo(pgConn)
	mux.HandleFunc(createAccessTokenHandlerInfo.Url, createAccessTokenHandlerInfo.Handler)

	admRepo := adm.NewAdmRepo(dbPool)
	admService := adm.NewAdmService(admRepo)
	admHandler := adm.NewAdmNeighborsHandler(admService)
	mux.HandleFunc("/adm-neighbors", admHandler.AdmNeighborsHandler)
	mux.HandleFunc("/reverse-geocode", admHandler.AdmForLatLngHandler)
	mux.HandleFunc("/geojsonl", admHandler.AdmGeojsonlHandler)

	fcPath := "/fc"
	fcBaseUrl := url.URL{Path: path.Join(baseApiPath, fcPath)}
	mux.HandleFunc(
		fcPath,
		func(w http.ResponseWriter, r *http.Request) {
			admHandler.GetAdmFeatureCollectionHandler(w, r, fcBaseUrl)
		},
	)

	handler := GetAuthMiddleWare(pgConn)(mux)
	return handler
}
