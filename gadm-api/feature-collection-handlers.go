package main

import (
	"encoding/json"
	"fmt"
	"gadm-api/logger"
	"net/http"
)

type FeatureCollectionHandlerInfo struct {
	queryParams FeatureCollectionQueryParams
	queryLimits QueryLimits
}

// todo: this and geojsonl settings can be abstracted since only limits are different, but is it worth it?
var featureCollectionEndpointInfo = map[GadmLevel]FeatureCollectionHandlerInfo{
	GadmLevel0: {
		queryParams: FeatureCollectionQueryParams{
			TableName:              ADM_0_TABLE,
			FeaturePropertiesNames: []string{Adm0.FID, Adm0.GID0, Adm0.Country},
			GeometryColumnName:     Adm0.Geometry,
			OrderByColumnName:      Adm0.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 5},
	},
	GadmLevel1: {
		queryParams: FeatureCollectionQueryParams{
			TableName: ADM_1_TABLE,
			FeaturePropertiesNames: []string{Adm1.FID, Adm1.GID0, Adm1.Country,
				Adm1.GID1, Adm1.Name1, Adm1.Varname1, Adm1.NlName1, Adm1.Type1,
				Adm1.Engtype1, Adm1.Cc1, Adm1.Hasc1, Adm1.Iso1,
			},
			GeometryColumnName: Adm1.Geometry,
			OrderByColumnName:  Adm1.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 10},
	},
	GadmLevel2: {
		queryParams: FeatureCollectionQueryParams{
			TableName: ADM_2_TABLE,
			FeaturePropertiesNames: []string{Adm2.FID, Adm2.GID0, Adm2.Country,
				Adm2.GID1, Adm2.Name1, Adm2.NlName1, Adm2.GID2, Adm2.Name2,
				Adm2.Varname2, Adm2.NlName2, Adm2.Type2, Adm2.Engtype2, Adm2.Cc2,
				Adm2.Hasc2,
			},
			GeometryColumnName: Adm2.Geometry,
			OrderByColumnName:  Adm2.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 20},
	},
	GadmLevel3: {
		queryParams: FeatureCollectionQueryParams{
			TableName: ADM_3_TABLE,
			FeaturePropertiesNames: []string{Adm3.FID, Adm3.GID0, Adm3.Country,
				Adm3.GID1, Adm3.Name1, Adm3.NlName1, Adm3.GID2, Adm3.Name2,
				Adm3.NlName2, Adm3.GID3, Adm3.Name3, Adm3.Varname3,
				Adm3.NlName3, Adm3.Type3, Adm3.Engtype3, Adm3.Cc3,
				Adm3.Hasc3,
			},
			GeometryColumnName: Adm3.Geometry,
			OrderByColumnName:  Adm3.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 20},
	},
	GadmLevel4: {
		queryParams: FeatureCollectionQueryParams{
			TableName: ADM_4_TABLE,
			FeaturePropertiesNames: []string{Adm4.FID, Adm4.GID0, Adm4.Country,
				Adm4.GID1, Adm4.Name1, Adm4.GID2, Adm4.Name2, Adm4.GID3,
				Adm4.Name3, Adm4.GID4, Adm4.Name4, Adm4.Varname4, Adm4.Type4,
				Adm4.Engtype4, Adm4.Cc4,
			},
			GeometryColumnName: Adm4.Geometry,
			OrderByColumnName:  Adm4.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 50},
	},
	GadmLevel5: {
		queryParams: FeatureCollectionQueryParams{
			TableName: ADM_5_TABLE,
			FeaturePropertiesNames: []string{Adm5.FID, Adm5.GID0, Adm5.Country,
				Adm5.GID1, Adm5.Name1, Adm5.GID2, Adm5.Name2, Adm5.GID3,
				Adm5.Name3, Adm5.GID4, Adm5.Name4, Adm5.GID5, Adm5.Name5,
				Adm5.Type5, Adm5.Engtype5, Adm5.Cc5,
			},
			GeometryColumnName: Adm5.Geometry,
			OrderByColumnName:  Adm5.FID,
		},
		queryLimits: QueryLimits{minLimit: 1, maxLimit: 50},
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

	params := featureCollectionEndpointInfo[gadmLevel].queryParams
	limits := featureCollectionEndpointInfo[gadmLevel].queryLimits

	fcQueryParams := FeatureCollectionQueryParams{
		TableName:              params.TableName,
		FeaturePropertiesNames: params.FeaturePropertiesNames,
		GeometryColumnName:     params.GeometryColumnName,
		OrderByColumnName:      params.OrderByColumnName,
		StartAtValue:           max(paginationParams.StartAtFid, MIN_FID),
		LimitValue:             clamp(paginationParams.Limit, limits.minLimit, limits.maxLimit),
	}

	nextFid, err := s.getNextFid(ctx, params.TableName, params.OrderByColumnName,
		fcQueryParams.StartAtValue, fcQueryParams.LimitValue)
	var nextUrl string
	if err != nil {
		logger.Error("failed_to_get_next_fid %v", err)
	} else {
		nextUrl = getFeatureCollectionUrl(gadmLevel, QueryParam{
			Key:   string(PAGE_SIZE_QUERY_KEY),
			Value: fmt.Sprintf("%d", fcQueryParams.LimitValue),
		}, QueryParam{
			Key:   string(START_AT_QUERY_KEY),
			Value: fmt.Sprintf("%d", nextFid),
		},
		)
	}

	sql, args, err := buildFeatureCollectionSqlQuery(fcQueryParams)
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
