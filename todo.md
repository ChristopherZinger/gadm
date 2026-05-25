# TODO

[] - refactor geojsona and feature-collection endpoints, access token
[] - run migrations on the server + create population scripts for 
    adm_tree and adm_neighbors tables. Introudce cron job and one-off jobs
[] - hash access tokens


## Feature Ideas 

[] - Tiles builder - user ploads geojson and gets link to map with tiles
[] - support files upload for frontend feature display 

## api

[] open telemetry, graphana, prometheus
[] add middleware for cors
[] create go routine for removing stale or expired tokens from cache
[x] rate limiter or research how to prevent ddos attacks -> Caddy handles this or this lib -> https://github.com/go-chi/httprate


# Interesting Golang Repos For Inspiration

- tiles with pg -> [link](https://medium.com/@lawsontaylor/diy-vector-tile-server-with-postgis-and-fastapi-b8514c95267c)
- https://github.com/tidwall/tile38 // interesting geospatial lib
- https://github.com/tidwall/geojson
- https://pkg.go.dev/github.com/paulmach/go.geojson
- https://pkg.go.dev/github.com/go-spatial/geom/encoding/geojson
- https://github.com/paulmach/go.geojson
