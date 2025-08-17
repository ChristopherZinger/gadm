---
weight: 10
title: "geojsonl"
---

# Geojsonl

---

## Endpoint Info

```
method:     GET
path:       {{< param "pathGeojsonl">}}<LEVEL>
LEVEL:      'lv0' | 'lv1' | 'lv2' | 'lv3' | 'lv4' | 'lv5'
```

## Notes

This endpoint stream buffers of [new line delimited](https://www.interline.io/blog/geojsonl-extracts/)
[GeoJSON Features](https://datatracker.ietf.org/doc/html/rfc7946#section-3.2).

The GADM dataset defines up to 6 levels of administrative boundaries, but availability varies by country. For example, France includes all 6 levels, while smaller countries like Monaco only provide level 1.

When making a request, you must specify the desired level directly in the URL path. This tells the API which administrative layer to return for your query.

## Example

```
curl -H "Authorization: Bearer $TOKEN" \
    "{{<param "apiBaseUrl">}}{{< param "pathGeojsonl">}}lv0"
```

## Pagination

The API uses cursor-based pagination. When more results are available,
the response includes a `Link`
[header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Link)
with the URL for the next page:

**Response Header:**

```
Link: {{< param "apiBaseUrl" >}}{{< param "pathGeojsonl" >}}lv0?{{< param "queryParamStartAt">}}=1234; rel="next"
```

**Example with curl:**

```
# First request
curl -I -H "Authorization: Bearer $TOKEN" \
    "{{< param "apiBaseUrl" >}}{{< param "pathGeojsonl">}}lv1"

# Response headers will include:
# Link: {{< param "apiBaseUrl" >}}{{<param "pathGeojsonl">}}lv1?start-after=ABC123; rel="next"

# Use the Link header URL for the next page
curl -H "Authorization: Bearer $TOKEN" \
    "{{<param "apiBaseUrl">}}{{< param "pathFeatureCollection">}}lv1?{{< param "queryParamStartAt">}}=ABC123"
```

## Query Parameters

### Page Size

You can control page size with `{{<param "queryParamPageSize">}}` parameter

There are hard limits on maximum number of features you can retrieve in a single page.

**Max number of features per level**

- level 0, 1 - 20 items
- level 2, 3 - 50 items
- level 4, 5 - 100 items

**Example**

```
{{< param "pathGeojsonl">}}<LEVEL>?{{< param "queryParamPageSize">}}=5
```

### Page number

All features include original GADM `fid` field. This field is used to order
and paginate the results. You can specify `fid` value you want to start with,
however it's better to leverage the cursor included in the `Link` header.

**Example**

```
{{< param "pathGeojsonl">}}<LEVEL>?{{< param "queryParamStartAt">}}=14039
```

### Filtering

You can retrieve boundaries that belong to the same parent administration unit
by providing parent's `GID` value in a query parameter.

**Usage**

```
// PARENT_GID_PARAM: 'gid-0' | 'gid-1' | 'gid-2' | 'gid-3' | 'gid-4' | 'gid-5'
// CHILD_LEVEL: 'lv0' |'lv1' |'lv2' |'lv3' |'lv4' |'lv5'
{{< param "pathGeojsonl">}}<CHILD_LEVEL>?<PARENT_GID_PARAM>=FRA
```

**Example**

```
// get all administrative units at level 3 for France
{{< param "pathGeojsonl">}}lv3?{{< param "queryParamGid">}}0=FRA
```
