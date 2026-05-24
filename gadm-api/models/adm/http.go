package adm

import (
	"encoding/json"
	"io"
	"net/http"

	"gadm-api/logger"

	geojson "github.com/paulmach/go.geojson"
)

type Handler struct {
	service *Service
}

func NewAdmNeighborsHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) AdmNeighborsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handler.getAdmNeighborsHandler(w, r)
		return
	}
	http.Error(w, "method_not_allowed", http.StatusMethodNotAllowed)
}

func (handler *Handler) getAdmNeighborsHandler(w http.ResponseWriter, r *http.Request) {
	admId := r.URL.Query().Get("adm-id")
	if admId == "" {
		http.Error(w, "missing_adm_id", http.StatusBadRequest)
		return
	}

	result, err := handler.service.GetAdmNeighbors(r.Context(), admId)
	if err != nil {
		logger.Error("failed_to_get_adm_neighbors %v", err)
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (handler *Handler) AdmForLatLngHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handler.getAdmForLatLngHandler(w, r)
		return
	}
	http.Error(w, "method_not_allowed", http.StatusMethodNotAllowed)
}

func (handler *Handler) getAdmForLatLngHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("failed_to_read_request_body %v", err)
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
		return
	}

	geometry, err := geojson.UnmarshalGeometry(body)
	if err != nil {
		logger.Error("failed_to_unmarshal_geometry %v", err)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	if geometry.Type != "Point" {
		logger.Error("invalid_geometry_type type=%s", geometry.Type)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	lng := geometry.Point[0]
	lat := geometry.Point[1]

	result, err := handler.service.GetAdmForLatLng(r.Context(), lat, lng)
	if err != nil {
		logger.Error("failed_to_get_adm_for_lat_lng %v", err)
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
