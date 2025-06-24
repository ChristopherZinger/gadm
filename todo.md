# TODO

## MVP

[x] - feature collection endpoints
[x] - geojsonl endpoints
[] - security
[x] - data in production
[] - rate limiting
[] - access tokens
[] - docs - svelte app or swagger

## api

[] api endpoint tests

[] add cursor (next url) value in return header

[] use constants for query parameter names - take, startAfter

[] add tracing id to logger

[] add middleware for cors

[] structured json logging

[] write tests for exiting endpoints and utility functions

[] Add more structured error handling

[] Add request validation

[] check if geometry can be null according to pg definitions and if not remove geom is not null check from sql queries

[] implement endpoint access token

[] find monitoring platform for analytics and observability

[x] rate limiter or research how to prevent ddos attacks -> Caddy handles this

## tile server ?

[] https://developmentseed.org/titiler/#packages
[] research postgis raster and tiling capabilities

## frontend

[] svelte app with leaflet map

## db

[] introduce indexes on gis\_ field

# Interesting Golang Repos For Inspiration

- https://github.com/tidwall/tile38 // interesting geospatial lib
- https://github.com/tidwall/geojson
- https://pkg.go.dev/github.com/paulmach/go.geojson
- https://pkg.go.dev/github.com/go-spatial/geom/encoding/geojson
- https://github.com/paulmach/go.geojson
