echo "POPULATE GADM DB as ADMIN user"
echo "Ingesting GADM data into PostgreSQL"
ogr2ogr \
    -f PostgreSQL PG:"dbname=gadm-db user=postgres password=secret host=192.168.178.199 port=5432" \
    -lco SCHEMA=gadm \
    -progress \
   ./geopackage_all_layers/gadm_410-levels.gpkg  


   