# TODO

## MVP 1

[x] - feature collection endpoints
[x] - geojsonl endpoints
[] - security
[x] - data in production
[x] - rate limiting
[x] - access tokens
[] - docs - svelte app or swagger

## MVP 2

[] - reverse geo-location
[] - countries endpoint (flags but no geometry?)
[] - info without geometries endpoint

## MVP 3

[] - Tiles (raster or vector)

## api

[] write tests for exiting endpoints and utility functions
[] - test pagination
[] - test filtering (with pagination)
[] - test min max limits
[] - test sql injections
[] add tracing id to logger
[] add middleware for cors
[] structured json logging
[] Add more structured error handling
[] Add request validation
[] check if geometry can be null according to pg definitions and if not remove geom is not null check from sql queries
[] implement endpoint access token - for creation and removal
[] create go routine for removing stale or expired tokens from cache
[] find monitoring platform for analytics and observability
[x] rate limiter or research how to prevent ddos attacks -> Caddy handles this or this lib -> https://github.com/go-chi/httprate

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
