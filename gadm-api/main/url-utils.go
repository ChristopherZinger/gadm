package main

import (
	"fmt"
	"gadm-api/logger"
	"net/http"
	"net/url"
	"strconv"
)

type PaginationParams struct {
	Limit      int
	StartAtFid int
}

func expectIntParamInQuery(r *http.Request, paramName string, defaultValue ...int) (int, error) {
	paramStrVal := r.URL.Query().Get(paramName)
	value, err := strconv.Atoi(paramStrVal)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return 0, fmt.Errorf("invalid %s parameter: %w", paramName, err)
	}
	return value, nil
}

func getPaginationParams(r *http.Request) (PaginationParams, error) {
	pageSize, err := expectIntParamInQuery(r, string(PAGE_SIZE_QUERY_KEY), -1)
	if err != nil {
		return PaginationParams{}, fmt.Errorf("failed to parse query parameter 'take': %w", err)
	}

	startAtFid, err := expectIntParamInQuery(r, string(START_AT_QUERY_KEY), -1)
	if err != nil {
		return PaginationParams{}, fmt.Errorf("failed to parse query parameter '%s': %w", START_AT_QUERY_KEY, err)
	}
	return PaginationParams{
		Limit:      pageSize,
		StartAtFid: startAtFid,
	}, nil
}

func getNextPageUrlQueryParams(
	nextStartAtFid int,
	pageSize int,
	filterParams SqlFilterParams) ([]QueryParam, error) {

	nextUrlQueryParams := []QueryParam{
		{
			Key:   string(PAGE_SIZE_QUERY_KEY),
			Value: fmt.Sprintf("%d", pageSize),
		}, {
			Key:   string(START_AT_QUERY_KEY),
			Value: fmt.Sprintf("%d", nextStartAtFid),
		},
	}
	nextFilterQueryParamColName, err := getFilterUrlQueryParameterForFilterableColumnName(filterParams.FilterColName)
	if err == nil {
		nextUrlQueryParams = append(nextUrlQueryParams, QueryParam{
			Key:   nextFilterQueryParamColName,
			Value: filterParams.FilterVal,
		})
	} else {
		return nil, fmt.Errorf(
			"failed_to_get_filter_url_query_parameter_for_filterable_column_name, filter_col_name: %s, %v",
			filterParams.FilterColName, err)
	}

	return nextUrlQueryParams, nil
}

func getFilterUrlQueryParameterForFilterableColumnName(filterableColName string) (string, error) {
	switch filterableColName {
	case "gid_0":
		return "gid-0", nil
	case "gid_1":
		return "gid-1", nil
	case "gid_2":
		return "gid-2", nil
	case "gid_3":
		return "gid-3", nil
	case "gid_4":
		return "gid-4", nil
	case "gid_5":
		return "gid-5", nil
	default:
		return "", fmt.Errorf("unsupported_filterable_column_name %s", filterableColName)
	}
}

func getSqlFilterParamsFromUrl(filterableColNames []string, urlValues url.Values) SqlFilterParams {
	var result SqlFilterParams
	for _, filterColName := range filterableColNames {
		filterUrlParamName, err := getFilterUrlQueryParameterForFilterableColumnName(filterColName)
		if err != nil {
			logger.Error("%v", err)
			continue
		}

		paramStrVal := urlValues.Get(filterUrlParamName)
		if paramStrVal != "" {
			result.FilterColName = filterColName
			result.FilterVal = paramStrVal
			break
		}
	}
	return result
}
