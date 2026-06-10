package handlers

import (
	"net/http"
	"net/url"
)

type HandlerInfo struct {
	Url     string
	Handler func(w http.ResponseWriter, r *http.Request)
}

func getBaseApiUrl() *url.URL {
	u := &url.URL{
		Path: "/api/v1/",
	}
	return u
}
