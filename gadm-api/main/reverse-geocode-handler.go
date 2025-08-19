package main

import (
	"context"
	"gadm-api/logger"
	"io"
	"net/http"

	geojson "github.com/paulmach/go.geojson"

	db "gadm-api/db"
	utils "gadm-api/utils"
)

type ReverseGeocodeHandler struct {
	pgConn *db.PgConn
	req    *http.Request
	writer http.ResponseWriter
	ctx    context.Context
}

func NewReverseGeocodeHandler(
	ctx context.Context,
	req *http.Request,
	writer http.ResponseWriter,
	pgConn *db.PgConn) *ReverseGeocodeHandler {
	return &ReverseGeocodeHandler{pgConn: pgConn, req: req, writer: writer, ctx: ctx}
}

type ReverseGeocodeResponse struct {
	Level utils.GadmLevel `json:"level"`
	Id    string          `json:"id"`
}

func (handler *ReverseGeocodeHandler) handle() {
	if handler.req.Method != http.MethodPost {
		logger.Error("invalid_method method=%s", handler.req.Method)
		http.Error(handler.writer, "method_not_allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(handler.req.Body)
	if err != nil {
		logger.Error("failed_to_read_request_body %v", err)
		http.Error(handler.writer, "internal_server_error", http.StatusInternalServerError)
		return
	}

	geometry, err := geojson.UnmarshalGeometry(body)
	if err != nil {
		logger.Error("failed_to_unmarshal_geometry %v", err)
		http.Error(handler.writer, "invalid_request", http.StatusBadRequest)
		return
	}

	if geometry.Type != "Point" {
		logger.Error("invalid_geometry_type type=%s", geometry.Type)
		http.Error(handler.writer, "invalid_request", http.StatusBadRequest)
		return
	}

	lng := geometry.Point[0]
	lat := geometry.Point[1]
	sql, args, err := db.GetReverseGeocodeSqlQuery(db.Point{Lat: lat, Lng: lng})
	if err != nil {
		logger.Error("failed_to_build_reverse_geocode_sql %v", err)
		http.Error(handler.writer, "internal_server_error", http.StatusInternalServerError)
		return
	}

	var jsonResult []byte
	err = handler.pgConn.Db.QueryRow(handler.ctx, sql, args...).Scan(&jsonResult)
	if err != nil {
		logger.Error("failed_to_query_database %v", err)
		http.Error(handler.writer, "no_result_found", http.StatusNotFound)
		return
	}

	handler.writer.Header().Set("Content-Type", "application/json")
	handler.writer.WriteHeader(http.StatusOK)
	handler.writer.Write(jsonResult)

	return
}
