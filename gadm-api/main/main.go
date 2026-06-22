package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	gameloop "gadm-api/game-loop"
	"gadm-api/infra/pg"
	"gadm-api/jobs"
	"gadm-api/logger"
	"gadm-api/models/access_token"
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
	mux := http.NewServeMux()

	baseApiPath := "/api/v1"
	mux.Handle(baseApiPath+"/", http.StripPrefix(
		baseApiPath,
		getApiHandlers(dbPool, baseApiPath),
	))

	mux.Handle("/ws", http.HandlerFunc(getWebsocketHandler))

	handler := LoggingMiddleware(mux)

	logger.Info("server_starting_on_port_8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func getApiHandlers(dbPool *pgxpool.Pool, baseApiPath string) http.Handler {
	mux := http.NewServeMux()

	accessTokenRepo := access_token.NewAccessTokenRepo(dbPool)
	accessTokenService := access_token.NewAccessTokenService(accessTokenRepo)
	accessTokenHandler := access_token.NewAccessTokenHandler(accessTokenService)
	tokenCreationRateLimiter := access_token.NewAccessTokenCreationRateLimiter()
	mux.HandleFunc("/create-access-token", func(w http.ResponseWriter, r *http.Request) {
		accessTokenHandler.CreateAccessTokenHandler(w, r, tokenCreationRateLimiter)
	})

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

	handler := GetAuthMiddleWare(dbPool)(mux)
	return handler
}
