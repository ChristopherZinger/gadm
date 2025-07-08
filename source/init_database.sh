
echo "RESET GADM DB"
docker cp ./reset_db_and_roles.sql gadm-db:/reset.sql
docker exec -it gadm-db psql -U chris -d postgres -f /reset_db_and_roles.sql

sleep 2

echo "INITIALIZE GADM DB"
docker cp ./init_db_and_roles.sql gadm-db:/init_db_and_roles.sql
docker exec -it gadm-db psql -U chris -d gadm -f /init_db_and_roles.sql

sleep 2

echo "POPULATE GADM DB as ADMIN user"
echo "Ingesting GADM data into PostgreSQL"
ogr2ogr \
    -f PostgreSQL PG:"dbname=gadm user=admin password=secret host=127.0.0.1 port=5432" \
    -lco SCHEMA=gadm \
    -progress \
   ./geopackage_all_layers/gadm_410-levels.gpkg  

echo "Done"