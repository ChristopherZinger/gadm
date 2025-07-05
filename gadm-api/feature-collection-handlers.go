package main

import (
	"encoding/json"
	"fmt"
	"gadm-api/logger"
	"net/http"
)

type FeatureCollectionHandlerQueryConfig struct {
	TableName string
	GeoJSONFeatureConfig
	QueryLimitConfig

	// filterable columns need to be ordered by importance level
	// in case of gid filtering it means from gid_0 down to gid_5
	// if user passes multiple filtering params in the url
	// the first matching filterable column will be
	// the only one used for filtering
	FilterableColumns []string
	OrderByColumnName string
}

var featureCollectionHandlerQueryConfig = map[GadmLevel]FeatureCollectionHandlerQueryConfig{
	GadmLevel0: {
		TableName: ADM_0_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm0.FID, Adm0.GID0, Adm0.Country},
			GeometryColumnName:     Adm0.Geometry,
		},
		OrderByColumnName: Adm0.FID,
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 5},
	},
	GadmLevel1: {
		TableName: ADM_1_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm1.FID, Adm1.GID0, Adm1.Country,
				Adm1.GID1, Adm1.Name1, Adm1.Varname1, Adm1.NlName1, Adm1.Type1,
				Adm1.Engtype1, Adm1.Cc1, Adm1.Hasc1, Adm1.Iso1},
			GeometryColumnName: Adm1.Geometry,
		},
		OrderByColumnName: Adm1.FID,
		FilterableColumns: []string{Adm1.GID0},
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 5},
	},
	GadmLevel2: {
		TableName: ADM_2_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm2.FID, Adm2.GID0, Adm2.Country,
				Adm2.GID1, Adm2.Name1, Adm2.NlName1, Adm2.GID2, Adm2.Name2,
				Adm2.Varname2, Adm2.NlName2, Adm2.Type2, Adm2.Engtype2, Adm2.Cc2,
				Adm2.Hasc2,
			},
			GeometryColumnName: Adm2.Geometry,
		},
		OrderByColumnName: Adm2.FID,
		FilterableColumns: []string{Adm2.GID0, Adm2.GID1},
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
	},
	GadmLevel3: {
		TableName: ADM_3_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm3.FID, Adm3.GID0, Adm3.Country,
				Adm3.GID1, Adm3.Name1, Adm3.NlName1, Adm3.GID2, Adm3.Name2,
				Adm3.NlName2, Adm3.GID3, Adm3.Name3, Adm3.Varname3,
				Adm3.NlName3, Adm3.Type3, Adm3.Engtype3, Adm3.Cc3,
				Adm3.Hasc3,
			},
			GeometryColumnName: Adm3.Geometry,
		},
		OrderByColumnName: Adm3.FID,
		FilterableColumns: []string{Adm3.GID0, Adm3.GID1, Adm3.GID2},
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 20},
	},
	GadmLevel4: {
		TableName: ADM_4_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm4.FID, Adm4.GID0, Adm4.Country,
				Adm4.GID1, Adm4.Name1, Adm4.GID2, Adm4.Name2, Adm4.GID3,
				Adm4.Name3, Adm4.GID4, Adm4.Name4, Adm4.Varname4, Adm4.Type4,
				Adm4.Engtype4, Adm4.Cc4,
			},
			GeometryColumnName: Adm4.Geometry,
		},
		OrderByColumnName: Adm4.FID,
		FilterableColumns: []string{Adm4.GID0, Adm4.GID1, Adm4.GID2, Adm4.GID3},
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 50},
	},
	GadmLevel5: {
		TableName: ADM_5_TABLE,
		GeoJSONFeatureConfig: GeoJSONFeatureConfig{
			FeaturePropertiesNames: []string{Adm5.FID, Adm5.GID0, Adm5.Country,
				Adm5.GID1, Adm5.Name1, Adm5.GID2, Adm5.Name2, Adm5.GID3,
				Adm5.Name3, Adm5.GID4, Adm5.Name4, Adm5.GID5, Adm5.Name5,
				Adm5.Type5, Adm5.Engtype5, Adm5.Cc5,
			},
			GeometryColumnName: Adm5.Geometry,
		},
		OrderByColumnName: Adm5.FID,
		FilterableColumns: []string{Adm5.GID0, Adm5.GID1, Adm5.GID2, Adm5.GID3, Adm5.GID4},
		QueryLimitConfig:  QueryLimitConfig{minLimit: 1, maxLimit: 50},
	},
}

var supportedGadmLevelsForFeatureCollection = []GadmLevel{GadmLevel0, GadmLevel1, GadmLevel2, GadmLevel3, GadmLevel4, GadmLevel5}

func setFeatureCollectionResponseHeaders(w http.ResponseWriter, nextUrl string) {
	w.Header().Set("Content-Type", "application/json")
	if nextUrl != "" {
		w.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"next\"", nextUrl))
	}
}

func CreateFeatureCollectionHandlers(s *Server) ([]HandlerInfo, error) {
	handlerInfos := []HandlerInfo{}
	for _, gadmLevel := range supportedGadmLevelsForFeatureCollection {
		url := getFeatureCollectionUrl(gadmLevel)
		handler := func(w http.ResponseWriter, r *http.Request) {
			s.featureCollectionEndpointHandler(w, r, gadmLevel)
		}
		handlerInfos = append(handlerInfos, HandlerInfo{url: url, handler: handler})
	}
	return handlerInfos, nil
}

func (s *Server) featureCollectionEndpointHandler(w http.ResponseWriter, r *http.Request, gadmLevel GadmLevel) {
	ctx := r.Context()

	paginationParams, err := getPaginationParams(r)
	if err != nil {
		logger.Error("failed_to_get_pagination_params %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handlerConfig := featureCollectionHandlerQueryConfig[gadmLevel]

	filterParams := getSqlFilterParamsFromUrl(handlerConfig.FilterableColumns, r.URL.Query())
	pageSize := clamp(paginationParams.Limit,
		handlerConfig.QueryLimitConfig.minLimit,
		handlerConfig.QueryLimitConfig.maxLimit)
	startAtFid := max(paginationParams.StartAtFid, MIN_FID)

	nextPageUrlParams, err := s.getNextPageUrlQueryParams(ctx, gadmLevel, startAtFid, pageSize, filterParams)
	var nextUrl string
	if err != nil {
		logger.Error("failed_to_get_next_fid %v", err)
	} else {
		nextUrl = getFeatureCollectionUrl(gadmLevel, nextPageUrlParams...)
	}

	sql, args, err := buildFeatureCollectionSqlQuery(gadmLevel, SqlQueryParams{
		StartAtValue:    startAtFid,
		LimitValue:      pageSize,
		SqlFilterParams: filterParams,
	})
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var featureCollectionJSON json.RawMessage
	err = s.db.QueryRow(ctx, sql, args...).Scan(&featureCollectionJSON)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setFeatureCollectionResponseHeaders(w, nextUrl)
	w.Write(featureCollectionJSON)
}
