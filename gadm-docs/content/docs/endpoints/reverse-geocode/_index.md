---
weight: 10
title: "reverse geocode"
---

# Reverse Geocode

---

## Endpoint Info

```
method:     POST
path:       {{< param "pathReverseGeocode">}}
body:       GeoJSON Point Geometry
```

## Notes

Reverse geocode endpoint returns lowest available GADM level contains a given
point as a
[GeoJSON Feature](https://datatracker.ietf.org/doc/html/rfc7946#section-3.2)

Request body must include valid
[GeoJSON Point Geometry](https://datatracker.ietf.org/doc/html/rfc7946#section-3.1.2)

## Example

```
curl -X POST -H "Authorization: Bearer $TOKEN"\
    -H "Content-Type: application/json" \
    -d '{"type": "Point", "coordinates": [ -0.017369731082479873, 44.12916742279152]}' \
    "{{< param "apiBaseUrl">}}{{<param "pathReverseGeocode">}}"

```
