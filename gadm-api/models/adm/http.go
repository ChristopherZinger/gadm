package adm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

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

func getLevelIntFromString(level string) (*int, error) {
	if level == "" {
		return nil, nil
	}
	_lv, err := strconv.Atoi(level)
	if err != nil {
		return &_lv, fmt.Errorf("failed_converting_level_to_int %v", err)
	}
	if _lv < 0 || _lv > 5 {
		return &_lv, fmt.Errorf("level_range_error %v", err)
	}
	return &_lv, nil
}

func getBatchSizeIntFromString(batchSize string) (*int, error) {
	if batchSize == "" {
		return nil, nil
	}
	_batchSize, err := strconv.Atoi(batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed_converting_batch_size_to_int %v", err)
	}
	return &_batchSize, nil
}

func (handler *Handler) GetAdmFeatureCollectionHandler(w http.ResponseWriter, r *http.Request, baseUrl url.URL) {
	startAfterId := r.URL.Query().Get("start-after-id")
	startAfterFid := r.URL.Query().Get("start-after-fid")
	batchSize := r.URL.Query().Get("batch-size")
	lvString := r.URL.Query().Get("lv")

	_lv, err := getLevelIntFromString(lvString)
	if err != nil {
		logger.Error("failed_parsing_query_param_lv: %v", err)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	_batchSize, err := getBatchSizeIntFromString(batchSize)
	if err != nil {
		logger.Error("failed_parsing_query_param_batch_size: %v", err)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	optsBuilder := NewAdmQueryOptsBuilder()
	optsBuilder.SetLvAndBatchSize(_lv, _batchSize)
	optsBuilder.SetStartAfterId(startAfterId)
	optsBuilder.SetStartAfterFid(startAfterFid)
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

	lastAdm := result.Features[len(result.Features)-1]
	nextUrl := getAdmsNextUrl(baseUrl, lastAdm, opts)
	if nextUrl != "" {
		w.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"next\"", nextUrl))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getAdmsNextUrl(baseUrl url.URL, lastAdm *geojson.Feature, opts admQueryOpts) string {
	if lastAdm == nil {
		return ""
	}

	lastAdmId := lastAdm.ID
	lastAdmFid := lastAdm.Properties["fid"]
	query := baseUrl.Query()
	if opts.startAfterId != nil {
		query.Set("start-after-id", lastAdmId.(string))
	}
	if opts.startAfterFid != nil {
		query.Set("start-after-fid", lastAdmFid.(string))
	}
	if opts.lv != nil {
		query.Set("lv", fmt.Sprintf("%d", *opts.lv))
	}
	query.Set("batch-size", fmt.Sprintf("%d", opts.batchSize))
	baseUrl.RawQuery = query.Encode()

	return baseUrl.String()
}

func (handler *Handler) AdmGeojsonlHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("geojsonl_handler_called")
	startAfterId := r.URL.Query().Get("start-after-id")
	startAfterFid := r.URL.Query().Get("start-after-fid")
	batchSize := r.URL.Query().Get("batch-size")
	lvString := r.URL.Query().Get("lv")

	_lv, err := getLevelIntFromString(lvString)
	if err != nil {
		logger.Error("failed_parsing_query_param_lv: %v", err)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	_batchSize, err := getBatchSizeIntFromString(batchSize)
	if err != nil {
		logger.Error("failed_parsing_query_param_batch_size: %v", err)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	optsBuilder := NewAdmQueryOptsBuilder()
	optsBuilder.SetLvAndBatchSize(_lv, _batchSize)
	optsBuilder.SetStartAfterId(startAfterId)
	optsBuilder.SetStartAfterFid(startAfterFid)
	optsBuilder.SetIncludeGeometry(true)
	opts, err := optsBuilder.Build()
	if err != nil {
		logger.Error("failed_to_validate_adm_query_params %v", err)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	flusher, err := newFlusher(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ch := make(chan json.RawMessage, 10)
	go func() {
		err := handler.service.getAdmGeojsonlStream(r.Context(), ch, opts)
		if err != nil {
			logger.Error("failed_to_get_adm_geojsonl %v", err)
			return
		}
	}()

	for admJson := range ch {
		err := flusher.flush(admJson)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
