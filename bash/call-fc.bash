TOKEN=
# page 1
 curl -H "Authorization: Bearer $TOKEN" "localhost:8080/api/v1/fc/lv3?page-size=10&start-at=0&gid-0=POL" | tee geojsonl.json 

# page 2
curl -H "Authorization: Bearer $TOKEN" "localhost:8080/api/v1/fc/lv3?gid-0=POL&page-size=10&start-at=108435" | tee geojsonl_2.json 

# page 3
curl  "localhost:8080/api/v1/fc/lv3?gid-0=POL&page-size=10&start-at=108445" | tee geojsonl_3.json 