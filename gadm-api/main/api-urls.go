package main

import (
	"fmt"
	"net/http"
	"net/url"

	utils "gadm-api/utils"
)

type HandlerInfo struct {
	url     string
	handler func(w http.ResponseWriter, r *http.Request)
}

type EndpointType string

const (
	featureCollectionEndpoint EndpointType = "fc"
	geojsonlEndpoint          EndpointType = "geojsonl"
)

type PaginationParamKeys string

const (
	PAGE_SIZE_QUERY_KEY PaginationParamKeys = "page-size"
	START_AT_QUERY_KEY  PaginationParamKeys = "start-at"
)

type QueryParam struct {
	Key   string
	Value string
}

func getBaseApiUrl() *url.URL {
	u := &url.URL{
		Path: "/api/v1/",
	}
	return u
}

func getApiUrl(endpointType EndpointType, gadmLevel utils.GadmLevel, queryParams ...QueryParam) string {
	u := &url.URL{
		Path: fmt.Sprintf("%s%s/lv%d", getBaseApiUrl().Path, endpointType, gadmLevel),
	}

	q := u.Query()
	for _, param := range queryParams {
		q.Set(param.Key, param.Value)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

func getFeatureCollectionUrl(gadmLevel utils.GadmLevel, queryParams ...QueryParam) string {
	return getApiUrl(featureCollectionEndpoint, gadmLevel, queryParams...)
}

func getGeojsonlUrl(gadmLevel utils.GadmLevel, queryParams ...QueryParam) string {
	return getApiUrl(geojsonlEndpoint, gadmLevel, queryParams...)
}
