# TODO

## api

[] add query parameters to /api/lv0 endpoint: "?take=10"

[] check if geometry can be null according to pg definitions and if not remove geom is not null check from sql queries

[] rate limiter or research how to prevent ddos attacks

[] implement endpoint access token

[] figure out how to log

[] create endpoint that returns geojson. Checkout libs below:

- https://github.com/tidwall/tile38 // interesting geospatial lib
- https://github.com/tidwall/geojson
- https://pkg.go.dev/github.com/paulmach/go.geojson
- https://pkg.go.dev/github.com/go-spatial/geom/encoding/geojson
- https://github.com/paulmach/go.geojson

## tile server ?

[] https://developmentseed.org/titiler/#packages

## frontend

[] svelte app with leaflet map

## raspberry

[] move pg with gadm to raspberry pi
[] setup cloudflare tunnels

## db

[] introduce indexes on gis\_ field
