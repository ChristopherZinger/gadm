package handlers

import (
	"gadm-api/logger"
	"io"
	"net/http"

	geojson "github.com/paulmach/go.geojson"

	db "gadm-api/db"
	rgRepo "gadm-api/db/repo"
	utils "gadm-api/utils"
)

type ReverseGeocodeHandler struct {
	req    *http.Request
	repo   *rgRepo.ReverseGeocodeRepo
	writer http.ResponseWriter
}

func NewReverseGeocodeHandler(
	pgConn *db.PgConn,
	req *http.Request,
	writer http.ResponseWriter,
) *ReverseGeocodeHandler {

	repo := rgRepo.NewReverseGeocodeRepo(pgConn, req.Context())

	return &ReverseGeocodeHandler{
		req:    req,
		writer: writer,
		repo:   repo,
	}
}

type ReverseGeocodeResponse struct {
	Level utils.GadmLevel `json:"level"`
	Id    string          `json:"id"`
}

func (handler *ReverseGeocodeHandler) Handle() {
	if handler.req.Method != http.MethodPost {
		logger.Error("invalid_method method=%s", handler.req.Method)
		http.Error(handler.writer, "method_not_allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(handler.req.Body)
	if err != nil {
		logger.Error("failed_to_read_request_body %v", err)
		http.Error(handler.writer, "internal_server_error", http.StatusInternalServerError)
		return
	}

	geometry, err := geojson.UnmarshalGeometry(body)
	if err != nil {
		logger.Error("failed_to_unmarshal_geometry %v", err)
		http.Error(handler.writer, "invalid_request", http.StatusBadRequest)
		return
	}

	if geometry.Type != "Point" {
		logger.Error("invalid_geometry_type type=%s", geometry.Type)
		http.Error(handler.writer, "invalid_request", http.StatusBadRequest)
		return
	}

	jsonResult, err := handler.repo.GetLocation(
		rgRepo.GetReverseGeocodeParams{
			Lat: geometry.Point[1],
			Lng: geometry.Point[0],
		},
	)
	if err != nil {
		logger.Error("failed_to_get_location %v", err)
		http.Error(handler.writer, "internal_server_error", http.StatusInternalServerError)
		return
	}

	handler.writer.Header().Set("Content-Type", "application/json")
	handler.writer.WriteHeader(http.StatusOK)
	handler.writer.Write(jsonResult)

	return
}
