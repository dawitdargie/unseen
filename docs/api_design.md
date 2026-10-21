# Unseen V1 — API Contract

Base path: /api/v1
Format: JSON
Public endpoints are read-only and require no login.
Only published discoveries are returned.

## 1. List and Search Discoveries

GET /discoveries

Optional query parameters:
- domain: earth or universe
- category: category slug; requires domain
- q: search text, maximum 200 characters
- page: positive integer, default 1
- limit: integer from 1 to 50, default 12

Ordering:
- Without q: newest published first.
- With q: most relevant first.
- Use discovery ID to break ordering ties.

Example request:

GET /api/v1/discoveries?domain=earth&category=islands&page=1&limit=12

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
        "slug": "islands",
        "name": "Islands"
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
      "slug": "islands",
      "name": "Islands"
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
      "license_url": "https://creativecommons.org/licenses/by-sa/4.0/"
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

If today's selection is unavailable, return an eligible published fallback.
If no published discovery exists, return 404.

Register /discoveries/daily before the dynamic slug route where necessary.

## 4. Get Related Discoveries

GET /discoveries/{slug}/related

Optional parameter:
- limit: integer from 1 to 12, default 4

Example request:

GET /api/v1/discoveries/socotra/related?limit=4

Response — 200 OK:
- data: an array of discovery cards, as in the listing endpoint.
- Match the same category and exclude the current discovery.
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

Example request:

GET /api/v1/categories?domain=earth

Example response — 200 OK:

```json
{
  "data": [
    {
      "slug": "islands",
      "name": "Islands",
      "domain": "earth"
    },
    {
      "slug": "caves",
      "name": "Caves",
      "domain": "earth"
    }
  ]
}
```

## 6. Errors

Use the same error structure everywhere.

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
- 404: discovery unavailable.
- 500: unexpected internal error.
- 503: database temporarily unavailable.

Never expose database errors or credentials in responses.

## 7. Shared Rules

- Dates and timestamps use UTC.
- Optional scalar fields return null; empty facts return {}.
- Explanations are plain text.
- Image attribution is included in cards and details.
- Example image URLs and credits are placeholders, not approved assets.
- Return stored content even when external providers are unavailable.
- Importing and publishing use internal Go commands, not public endpoints.