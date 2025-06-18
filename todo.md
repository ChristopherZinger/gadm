# TODO

## api

[] api endpoint tests

[] add cursor (next url) value in return header

[] add middleware for cors and logging

[] Add more structured error handling

[] Add request validation

[] check if geometry can be null according to pg definitions and if not remove geom is not null check from sql queries

[] rate limiter or research how to prevent ddos attacks

[] implement endpoint access token

[] figure out how to log

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
