---
weight: 1
bookFlatSection: true
title: "Authorization"
---

# Authorization

Each request must include an Authorization header with a valid Bearer token.
This token acts as your key to the service, allowing the API to verify your
identity and permissions. Without it, requests will be denied with an
HTTP 401 Unauthorized error.

{{< highlight bash "linenos=false" >}}
curl -H "Authorization: Bearer <TOKEN>" \
    "{{< param "apiBaseUrl" >}}/api/v1/geojsonl/lv1"
{{< /highlight >}}

You can obtain access token [programmatically](/docs/endpoints/get-access-token)
