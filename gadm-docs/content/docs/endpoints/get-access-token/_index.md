---
weight: 10
title: "get access token"
---

# Get Access Token

---

## Endpoint Info

{{< highlight text "linenos=false" >}}
method:     POST
path:       /api/v1/create-access-token
auth:       not required for this endpoint
{{< /highlight >}}

## Notes

Use this endpoint to create an API access token for an email address.  
The returned token can then be sent as a Bearer token when calling protected endpoints (for example `fc` and `geojsonl`).

Token creation is guarded by a global rate limiter (one request every ~2 seconds across the service instance).

## Request

The email is passed as a query parameter (`email`) on the POST request.

{{< highlight bash "linenos=false" >}}
curl -X POST \
  "{{< param "apiBaseUrl" >}}/api/v1/create-access-token?email=user@example.com"
{{< /highlight >}}

## Success Response

**Status:** `201 Created`

{{< highlight json "linenos=false" >}}
{
  "token": "c5609f65-0d58-4b58-95fb-1a4d9f6479ce",
  "email": "user@example.com",
  "created_at": "2026-04-25T08:40:12.314512Z"
}
{{< /highlight >}}

## Error Responses

{{< highlight text "linenos=false" >}}
400 Bad Request       -> email_not_provided
405 Method Not Allowed-> method_not_allowed
429 Too Many Requests -> rate_limit_exceeded
500 Internal Server Error -> internal_server_error
{{< /highlight >}}

## Using the Token

After creation, send the token in the `Authorization` header:

{{< highlight bash "linenos=false" >}}
curl -H "Authorization: Bearer $TOKEN" \
  "{{< param "apiBaseUrl" >}}{{< param "pathFeatureCollection" >}}lv0"
{{< /highlight >}}
