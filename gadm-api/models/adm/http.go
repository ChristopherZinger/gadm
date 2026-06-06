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

func (handler *Handler) validateAdmFcQueryParams(urlQuery url.Values) (admQueryOpts, error) {
	startAfterId := urlQuery.Get("start-after-id")
	startAfterFid := urlQuery.Get("start-after-fid")
	batchSize := urlQuery.Get("batch-size")
	lvString := urlQuery.Get("lv")

	getLvInt := func(lv string) (int, error) {
		_lv, err := strconv.Atoi(lv)
		if err != nil {
			return 0, fmt.Errorf("failed_int_conversion %v", err)
		}
		if _lv < 0 || _lv > 5 {
			return 0, fmt.Errorf("range_error %v", err)
		}
		return _lv, nil
	}
	var _lv *int
	if lvString != "" {
		__lv, err := getLvInt(lvString)
		if err != nil {
			return admQueryOpts{}, fmt.Errorf("failed_parsing_query_param_lv: %v", err)
		}
		_lv = &__lv
	}

	var _batchSize int
	var err error
	getBatchSizeForLv := func(lv int, batchSize int) int {
		if lv < 2 {
			return utils.Clamp(batchSize, 1, 5)
		}
		if lv < 4 {
			return utils.Clamp(batchSize, 1, 20)
		}
		return utils.Clamp(batchSize, 1, 50)
	}
	if batchSize != "" {
		_batchSize, err = strconv.Atoi(batchSize)
		if err != nil {
			return admQueryOpts{}, fmt.Errorf("failed_int_conversion batch_size %v", err)
		}
		if _lv != nil {
			_batchSize = getBatchSizeForLv(*_lv, _batchSize)
		} else {
			_batchSize = utils.Clamp(_batchSize, 1, 20)
		}
	}

	optsBuilder := NewAdmQueryOptsBuilder()
	optsBuilder.SetBatchSize(_batchSize)
	optsBuilder.SetIncludeGeometry(true)
	if _lv != nil {
		optsBuilder.SetLv(*_lv)
	}
	if startAfterId != "" {
		optsBuilder.SetStartAfterId(startAfterId)
	}
	if startAfterFid != "" {
		optsBuilder.SetStartAfterFid(startAfterFid)
	}

	opts, err := optsBuilder.Build()
	if err != nil {
		return admQueryOpts{}, fmt.Errorf("failed_to_build_adm_query_opts %v", err)
	}
	return opts, nil
}

func (handler *Handler) GetAdmFeatureCollectionHandler(w http.ResponseWriter, r *http.Request, baseUrl url.URL) {
	opts, err := handler.validateAdmFcQueryParams(r.URL.Query())
	if err != nil {
		logger.Error("failed_to_validate_adm_fc_query_params %v", err)
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
