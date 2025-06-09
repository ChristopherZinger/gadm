ogr2ogr \
    -f PostgreSQL PG:"dbname=gadm user=chris password=secret host=127.0.0.1 port=5432" \
    ../source/geopackage_all_layers/gadm_410-levels.gpkg