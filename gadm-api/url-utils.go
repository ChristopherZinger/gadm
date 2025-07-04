package main

import (
	"fmt"
	"gadm-api/logger"
	"net/http"
	"net/url"
	"strconv"
)

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
