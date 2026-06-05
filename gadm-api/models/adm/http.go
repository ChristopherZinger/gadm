package adm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gadm-api/logger"
	"gadm-api/utils"

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
	if r.Method == http.MethodPost {
		handler.postAdmNeighborsForPointHandler(w, r)
		return
	}
	logger.Error("method_not_allowed %s", r.Method)
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
		handler.getAdmForPointHandler(w, r)
		return
	}

	logger.Error("method_not_allowed %s", r.Method)
	http.Error(w, "method_not_allowed", http.StatusMethodNotAllowed)
}

func (handler *Handler) postAdmNeighborsForPointHandler(w http.ResponseWriter, r *http.Request) {
	point, err := getPointFromRequestBody(r)
	if err != nil {
		logger.Error("failed_to_get_point_from_request_body %v", err)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	result, err := handler.service.GetAdmNeighborsForPoint(r.Context(), point)
	if err != nil {
		logger.Error("failed_to_get_adm_neighbors_for_point %v", err)
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (handler *Handler) getAdmForPointHandler(w http.ResponseWriter, r *http.Request) {
	point, err := getPointFromRequestBody(r)
	if err != nil {
		logger.Error("failed_to_get_lat_lng_from_request_body %v", err)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	result, err := handler.service.GetAdmForPoint(r.Context(), point)
	if err != nil {
		logger.Error("failed_to_get_adm_for_lat_lng %v", err)
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getPointFromRequestBody(r *http.Request) (utils.Point, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.Point{}, fmt.Errorf("failed_to_read_request_body: %w", err)
	}

	geometry, err := geojson.UnmarshalGeometry(body)
	if err != nil {
		return utils.Point{}, fmt.Errorf("failed_to_unmarshal_geometry: %w", err)
	}

	if geometry.Type != "Point" {
		return utils.Point{}, fmt.Errorf("invalid_geometry_type: type %s", geometry.Type)
	}

	return utils.NewPointLngLat(geometry.Point[0], geometry.Point[1]), nil
}

func (handler *Handler) GetAdmFeatureCollectionHandler(w http.ResponseWriter, r *http.Request) {
	startAfterId := r.URL.Query().Get("start-after-id")
	startAfterFid := r.URL.Query().Get("start-after-fid")
	batchSize := r.URL.Query().Get("batch-size")
	lv := r.URL.Query().Get("lv")

	optsBuilder := NewAdmQueryOptsBuilder()

	if startAfterId != "" {
		optsBuilder.SetStartAfterId(startAfterId)
	}
	if startAfterFid != "" {
		optsBuilder.SetStartAfterFid(startAfterFid)
	}
	if batchSize != "" {
		_batchSize, err := strconv.Atoi(batchSize)
		if err != nil {
			logger.Error("failed_to_convert_lv_to_int %v", err)
			http.Error(w, "invalid_request", http.StatusBadRequest)
			return
		}
		optsBuilder.SetBatchSize(_batchSize)
	}
	if lv != "" {
		_lv, err := strconv.Atoi(lv)
		if err != nil {
			logger.Error("failed_to_convert_lv_to_int %v", err)
			http.Error(w, "invalid_request", http.StatusBadRequest)
			return
		}
		optsBuilder.SetLv(_lv)
	}
	optsBuilder.SetIncludeGeometry(true)
	opts, err := optsBuilder.Build()
	if err != nil {
		logger.Error("failed_to_build_adm_query_opts %v", err)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	result, err := handler.service.GetAdmsFc(r.Context(), opts)
	if err != nil {
		logger.Error("failed_to_get_adm_feature_collection %v", err)
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
