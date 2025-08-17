TOKEN="d2776cdd-c12e-4f3d-a433-cf8d93437487"
# page 1
curl -I -H "Authorization: Bearer $TOKEN" "localhost:8080/api/v1/fc/lv3?page-size=10&start-at=0&gid-1=POL.10_1" | tee geojsonl.json 

# # page 2
# curl -I -H "Authorization: Bearer $TOKEN" "localhost:8080/api/v1/geojsonl/lv3?gid-0=POL&page-size=10&start-at=108435" | tee geojsonl_2.json 
# curl -I -H "Authorization: Bearer $TOKEN" "localhost:8080/api/v1/geojsonl/lv3?gid-0=POL&page-size=10&start-at=108445" | tee geojsonl_2.json 

# # page 3
# curl  "localhost:8080/api/v1/fc/lv3?gid-0=POL&page-size=10&start-at=108445" | tee geojsonl_3.json 