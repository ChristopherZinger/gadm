---
weight: 10
title: "feature collection"
---

# Feature Collection

---

## Endpoint Info

```
method:     GET
path:       {{< param "pathFeatureCollection">}}<LEVEL>
LEVEL:      'lv0' | 'lv1' | 'lv2' | 'lv3' | 'lv4' | 'lv5'
```

## Notes

Feature collection endpoint returns results as
[GeoJSON FeatureCollection](https://datatracker.ietf.org/doc/html/rfc7946#section-3.3)

The GADM dataset defines up to 6 levels of administrative boundaries, but availability varies by country. For example, France includes all 6 levels, while smaller countries like Monaco only provide level 1.

When making a request, you must specify the desired level directly in the URL path. This tells the API which administrative layer to return for your query.

## Example

```
curl -H "Authorization: Bearer $TOKEN" \
    "{{< param "apiBaseUrl" >}}/api/v1/fc/lv1"
```

## Pagination

The API uses cursor-based pagination. When more results are available,
the response includes a `Link`
[header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Link)
with the URL for the next page:

**Response Header:**

```
Link: {{< param "apiBaseUrl" >}}/api/v1/fc/lv1?{{< param "queryParamStartAt">}}=1234; rel="next"
```

**Example with curl:**

```
# First request
curl -I -H "Authorization: Bearer $TOKEN" \
    {{< param "apiBaseUrl" >}}/api/v1/fc/lv1

# Response headers will include:
# Link: {{< param "apiBaseUrl" >}}/api/v1/fc/lv1?start-after=ABC123; rel="next"

# Use the Link header URL for the next page
curl -H "Authorization: Bearer $TOKEN" \
    "{{<param "apiBaseUrl">}}{{< param "pathFeatureCollection">}}lv1?{{< param "queryParamStartAt">}}=ABC123"
```

## Query Parameters

### Page Size

You can control page size with `{{<param "queryParamPageSize">}}` parameter

There are relatively low limits on maximum number of features you can retrieve
at once since the geometries can be rather heavy.

**Max number of features per level**

- level 0, 1 - 5 items
- level 2, 3 - 20 items
- level 4, 5 - 50 items

**Example**

```
{{< param "pathFeatureCollection">}}<LEVEL>?{{< param "queryParamPageSize">}}=5
```

### Page number

All features include original GADM `fid` field. This field is used to order
and paginate the results. You can specify `fid` value you want to start with,
however it's better to leverage the cursor included in the `Link` header.

**Example**

```
{{< param "pathFeatureCollection">}}<LEVEL>?{{< param "queryParamStartAt">}}=14039
```

### Filtering

You can retrieve boundaries that belong to the same parent administration unit
by providing parent's `GID` value in a query parameter.

**Usage**

```
// PARENT_GID_PARAM: 'gid-0' | 'gid-1' | 'gid-2' | 'gid-3' | 'gid-4' | 'gid-5'
// CHILD_LEVEL: 'lv0' |'lv1' |'lv2' |'lv3' |'lv4' |'lv5'
{{< param "pathFeatureCollection">}}<CHILD_LEVEL>?<PARENT_GID_PARAM>=FRA
```

**Example**

```
// get all administrative units at level 3 for France
{{< param "pathFeatureCollection">}}lv3?{{< param "queryParamGid">}}0=FRA
```
