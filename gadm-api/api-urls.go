package main

import (
	"fmt"
	"net/url"
)

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

func getApiUrl(endpointType EndpointType, gadmLevel GadmLevel, queryParams ...QueryParam) string {
	u := &url.URL{
		Path: fmt.Sprintf("/api/v1/%s/lv%d", endpointType, gadmLevel),
	}

	q := u.Query()
	for _, param := range queryParams {
		q.Set(param.Key, param.Value)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

func getFeatureCollectionUrl(gadmLevel GadmLevel, queryParams ...QueryParam) string {
	return getApiUrl(featureCollectionEndpoint, gadmLevel, queryParams...)
}

func getGeojsonlUrl(gadmLevel GadmLevel, queryParams ...QueryParam) string {
	return getApiUrl(geojsonlEndpoint, gadmLevel, queryParams...)
}
