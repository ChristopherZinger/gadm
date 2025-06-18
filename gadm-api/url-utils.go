package main

import (
	"fmt"
	"net/http"
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
