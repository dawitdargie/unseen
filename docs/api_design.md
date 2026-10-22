# Unseen V1 — API Contract

Base path: /api/v1
Format: JSON
Public endpoints are read-only and require no login.
Only published discoveries are returned.

Design baseline: 2026-10-22. Implementation order: [implementation_plan.md](implementation_plan.md).

## 1. List and Search Discoveries

GET /discoveries

Optional query parameters:
- domain: earth or universe
- category: exact category slug; requires domain. Unknown/mismatched categories return 400.
- q: search text, maximum 200 Unicode characters after trimming. Blank means no search.
- page: positive integer, default 1
- limit: integer from 1 to 50, default 12

Ordering:
- Without q: published_at DESC, then id DESC.
- With q: English full-text rank DESC, then published_at DESC, then id DESC.
- Title has higher search weight than summary, which has higher weight than explanation.
- A nonblank query containing only stop words/no searchable terms returns no matches.
- Malformed/repeated recognized query parameters or numeric overflow return 400; unrecognized parameters are ignored.
- Counts/results in a request share a consistent snapshot; separate pages are not a frozen catalog.

Example request:

GET /api/v1/discoveries?domain=earth&category=ecosystems&page=1&limit=12

Example response — 200 OK:

```json
{
  "data": [
    {
      "id": 1,
      "slug": "socotra",
      "title": "Socotra",
      "summary": "An island known for its unusual endemic life.",
      "domain": "earth",
      "category": {
        "slug": "ecosystems",
        "name": "Ecosystems"
      },
      "image": {
        "url": "https://example.org/socotra.jpg",
        "alt": "Satellite view of Socotra",
        "kind": "photograph",
        "credit": "Example image credit",
        "license": "Public domain",
        "license_url": null,
        "source_url": "https://example.org/image-source"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 12,
    "total": 1
  }
}
```

No matches return 200 with an empty data array and total 0.
A page beyond the results also returns an empty array.

## 2. Get Discovery Details

GET /discoveries/{slug}

Example request:

GET /api/v1/discoveries/socotra

Example response — 200 OK:

```json
{
  "data": {
    "id": 1,
    "slug": "socotra",
    "title": "Socotra",
    "summary": "An island known for its unusual endemic life.",
    "explanation": "Socotra supports distinctive ecosystems with many species found nowhere else.",
    "domain": "earth",
    "category": {
      "slug": "ecosystems",
      "name": "Ecosystems"
    },
    "location": "Yemen",
    "facts": {
      "notable_feature": "High endemism"
    },
    "image": {
      "url": "https://example.org/socotra.jpg",
      "alt": "Satellite view of Socotra",
      "kind": "photograph",
      "credit": "Example image credit",
      "license": "Public domain",
      "license_url": null,
      "source_url": "https://example.org/image-source"
    },
    "source": {
      "provider": "wikipedia",
      "url": "https://en.wikipedia.org/wiki/Socotra",
      "credit": "Wikipedia contributors",
      "license": "CC BY-SA 4.0",
      "license_url": "https://creativecommons.org/licenses/by-sa/4.0/",
      "revision_id": null,
      "revision_url": null,
      "adaptation_note": null
    },
    "published_at": "2026-10-21T08:00:00Z"
  }
}
```

Unknown or unpublished discoveries return 404.

## 3. Get Daily Discovery

GET /discoveries/daily

Example request:

GET /api/v1/discoveries/daily

Response — 200 OK:
- data: the same discovery object used by the detail endpoint.
- date: today's UTC date, formatted YYYY-MM-DD.

Example response structure:

```json
{
  "date": "2026-10-21",
  "data": {}
}
```

The empty object above represents the full detail object.

Use today's still-published selection. Otherwise return the most recent earlier daily selection whose discovery remains published; otherwise the newest published discovery (published_at DESC, id DESC).
If no published discovery exists, return 404. GET performs no selection writes or repairs.
The date field is today's UTC date, including fallback responses.
Daily command selection is deterministic/idempotent and excludes yesterday when alternatives exist; see [database_design.md](database_design.md).

Reserve the slug daily. Register the literal /discoveries/daily route with correct precedence over the dynamic slug route.

## 4. Get Related Discoveries

GET /discoveries/{slug}/related

Optional parameter:
- limit: integer from 1 to 12, default 4

Example request:

GET /api/v1/discoveries/socotra/related?limit=4

Response — 200 OK:
- data: an array of discovery cards, as in the listing endpoint.
- Match the same category and exclude the current discovery.
- Return only published rows, ordered published_at DESC, id DESC.
- Return an empty array when there are no matches.
- Return 404 if the requested discovery is unknown or unpublished.

```json
{
  "data": []
}
```

## 5. List Categories

GET /categories

Optional parameter:
- domain: earth or universe

Return all configured categories, including empty ones, ordered domain ASC, name ASC, slug ASC. Initial Earth slugs: landscapes, caves, ecosystems; Universe: planets-moons, nebulae, galaxies.

Example request:

GET /api/v1/categories?domain=earth

Example response — 200 OK:

```json
{
  "data": [
    {
      "slug": "caves",
      "name": "Caves",
      "domain": "earth"
    },
    {
      "slug": "ecosystems",
      "name": "Ecosystems",
      "domain": "earth"
    },
    {
      "slug": "landscapes",
      "name": "Landscapes",
      "domain": "earth"
    }
  ]
}
```

## 6. Errors

Use the same error structure throughout the content API; health endpoints have the small status payloads specified below.

Example — 400 Bad Request:

```json
{
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "limit must be between 1 and 50",
    "field": "limit"
  }
}
```

Example — 404 Not Found:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Discovery not found"
  }
}
```

Status codes:
- 200: successful request.
- 400: invalid query parameter or category/domain combination.
- 404: discovery unavailable or unknown content API route.
- 405: unsupported method on a known endpoint; include the Allow header.
- 500: unexpected internal error.
- 503: database temporarily unavailable.

Never expose database errors or credentials in responses.

## 7. Shared Rules

- Dates and timestamps use UTC.
- Optional scalar fields return null; empty facts return {}.
- Explanations, credits, and adaptation notes are plain text; normalize provider markup before storing/rendering.
- Malformed/repeated recognized query parameters return 400 on all content endpoints; unknown parameters are ignored.
- Image attribution is included in cards and details.
- Example image URLs and credits are placeholders, not approved assets.
- Return stored content even when external providers are unavailable.
- Importing and publishing use internal Go commands, not public endpoints.

## 8. Health Endpoints

Outside the versioned content API:

- GET /health/live: 200 with {"status":"ok"} when the HTTP process can respond; no provider/database check.
- GET /health/ready: 200 with {"status":"ready"} after a database ping with a short timeout; 503 with {"status":"unavailable"} otherwise.
- Do not expose credentials, raw database failures, or infrastructure details.

Kubernetes liveness uses live; readiness uses ready. A provider outage does not make the stored-content API unready.
