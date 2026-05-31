# export ndgeojson from postgres
for level in 0 1 2 3 4 5; do
  echo "Exporting level $level"
  ogr2ogr -f GeoJSONSeq geopackage_all_layers/pg_ndgeojson_export/adm_$level.geojsonl \
    PG:"dbname=gadm user=postgres password=secret host=localhost port=5432" \
    -sql "SELECT * FROM gadm.adm_$level" \
    -lco ID_FIELD=fid \
    -lco ID_TYPE=Integer
done

echo "All Done"
