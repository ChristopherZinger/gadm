package main

import (
	"context"
	"fmt"
	"gadm-api/logger"
	"net/http"

	db "gadm-api/db"
	repo "gadm-api/db/repo"
	dbutils "gadm-api/db/utils"
	utils "gadm-api/utils"
)

type FeatureCollectionHandlerQueryConfig struct {
	QueryLimitConfig
	// filterable columns need to be ordered by importance level
	// in case of gid filtering it means from gid_0 down to gid_5
	// if user passes multiple filtering params in the url
	// the first matching filterable column will be
	// the only one used for filtering
	FilterableColumns []string
}

type GadmFeatureCollectionHandler struct {
	fcRepo    *repo.FeatureCollectionRepo
	pgConn    *db.PgConn
	req       *http.Request
	writer    http.ResponseWriter
	ctx       context.Context
	gadmLevel utils.GadmLevel
	config    GeojsonlHandlerQueryConfig
}

var featureCollectionHandlerQueryConfig = map[utils.GadmLevel]FeatureCollectionHandlerQueryConfig{
	utils.GadmLevel0: {
		QueryLimitConfig: QueryLimitConfig{minLimit: 1, maxLimit: 5},
	},
	utils.GadmLevel1: {
		FilterableColumns: arrayToStrings(db.GidName0),
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 5},
	},
	utils.GadmLevel2: {
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1),
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
	},
	utils.GadmLevel3: {
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1, db.GidName2),
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
	},
	utils.GadmLevel4: {
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1, db.GidName2, db.GidName3),
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 50},
	},
	utils.GadmLevel5: {
		FilterableColumns: arrayToStrings(db.GidName0, db.GidName1, db.GidName2, db.GidName3, db.GidName4),
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 50},
	},
}

var supportedGadmLevelsForFeatureCollection = []utils.GadmLevel{utils.GadmLevel0, utils.GadmLevel1, utils.GadmLevel2, utils.GadmLevel3, utils.GadmLevel4, utils.GadmLevel5}

func (handler *GadmFeatureCollectionHandler) setFeatureCollectionResponseHeaders(nextUrl string) {
	handler.writer.Header().Set("Content-Type", "application/json")
	if nextUrl != "" {
		handler.writer.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"next\"", nextUrl))
	}
}

func CreateFeatureCollectionHandlers(s *db.PgConn) ([]HandlerInfo, error) {
	handlerInfos := []HandlerInfo{}
	for _, gadmLevel := range supportedGadmLevelsForFeatureCollection {
		url := getFeatureCollectionUrl(gadmLevel)
		handler := func(w http.ResponseWriter, r *http.Request) {
			newGadmFeatureCollectionHandler(s, r, w, gadmLevel).handle()
		}
		handlerInfos = append(handlerInfos, HandlerInfo{url: url, handler: handler})
	}
	return handlerInfos, nil
}

func (handler *GadmFeatureCollectionHandler) handle() {
	paginationParams, err := getPaginationParams(handler.req)
	if err != nil {
		logger.Error("failed_to_get_pagination_params %v", err)
		http.Error(handler.writer, err.Error(), http.StatusBadRequest)
		return
	}

	filterParams := getSqlFilterParamsFromUrl(
		handler.config.FilterableColumns,
		handler.req.URL.Query())
	pageSize := clamp(paginationParams.Limit,
		handler.config.QueryLimitConfig.minLimit,
		handler.config.QueryLimitConfig.maxLimit)
	startAtFid := max(paginationParams.StartAtFid, MIN_FID)

	nextUrl, err := handler.getNextPageUrl(startAtFid, pageSize, filterParams)
	if err != nil {
		logger.Error("failed_to_get_next_fid %v", err)
	}

	featureCollectionJSON, err := handler.fcRepo.GetFeatureCollection(repo.GetFeatureCollectionParams{
		GadmLevel:     handler.gadmLevel,
		FilterValue:   filterParams.FilterVal,
		FilterColName: filterParams.FilterColName,
		StartAtFid:    startAtFid,
		PageSize:      pageSize,
	})
	if err != nil {
		logger.Error("failed_to_get_feature_collection %v", err)
		http.Error(handler.writer, err.Error(), http.StatusInternalServerError)
		return
	}

	handler.setFeatureCollectionResponseHeaders(nextUrl)
	handler.writer.Write(featureCollectionJSON)
}

func (handler *GadmFeatureCollectionHandler) getNextPageUrl(startAtFid int, pageSize int, filterParams db.SqlFilterParams) (string, error) {
	nextStartAtFid, err := dbutils.GetNextPageFid(dbutils.NextPageParams{
		Context:       handler.ctx,
		PgConn:        handler.pgConn,
		StartAt:       startAtFid,
		PageSize:      pageSize,
		FilterColName: filterParams.FilterColName,
		FilterVal:     filterParams.FilterVal,
	})

	if err != nil {
		return "", fmt.Errorf("failed_to_get_next_start_at_fid %v", err)
	}

	nextUrlParams, err := getNextPageUrlQueryParams(
		nextStartAtFid,
		pageSize,
		filterParams)
	if err != nil {
		return "", fmt.Errorf("failed_to_get_next_page_url_query_params %v", err)
	}
	return getFeatureCollectionUrl(handler.gadmLevel, nextUrlParams...), nil
}

func newGadmFeatureCollectionHandler(
	pgConn *db.PgConn,
	r *http.Request,
	w http.ResponseWriter,
	gadmLevel utils.GadmLevel,
) *GadmFeatureCollectionHandler {

	ctx := r.Context()

	return &GadmFeatureCollectionHandler{
		pgConn:    pgConn,
		req:       r,
		writer:    w,
		ctx:       ctx,
		config:    geojsonHandlerQueryConfig[gadmLevel],
		gadmLevel: gadmLevel,
		fcRepo:    repo.NewFeatureCollectionRepo(pgConn, ctx),
	}
}
