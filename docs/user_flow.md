# Unseen V1 — User Flow

## 1. Open the Website
Visitor opens Home and sees:
- Daily discovery.
- New discoveries.
- Explore Earth and Explore Universe links.
- Search.

## 2. Explore Earth or Universe
Visitor selects Explore Earth or Explore Universe.
→ Opens the discovery listing for that domain.
→ Browses image cards with titles and short descriptions.
→ Optionally selects a category.
→ Opens another page of results or selects a discovery.

## 3. Search
Visitor enters search text and submits.
→ Opens the listing with matching discoveries.
→ Optionally narrows results by domain and category.
→ Selects a discovery.

Rules:
- Home search covers both Earth and Universe.
- Search within a domain keeps that domain selected.
- Changing domain clears the category filter.
- Changing search or filters resets pagination to page 1.
- Search and filters remain in the URL.

## 4. Open a Discovery
Visitor selects a card or the daily discovery.
→ Opens its detail page.
→ Sees the large image, explanation, category, and available location/facts.
→ Sees image credits, licensing, and original content source.
→ Can open the original source in a new tab.

## 5. Continue Exploring
From the detail page, the visitor can:
- Select a related discovery.
- Open the same category's listing.
- Return to the previous results.
- Navigate to Home, Earth, or Universe.

Returning to results preserves search, filters, and page number.

## 6. Interface States
- Loading: show placeholders while fetching.
- No results: explain and offer to clear search or filters.
- No related discoveries: hide that section.
- Missing or unpublished discovery: show a not-found page with explore links.
- Request failure: show a retry action.
- Image failure: show a fallback while keeping the content readable.

## 7. Page Structure
- / — Home
- /explore — Listing and search
- /explore?domain=earth — Earth discoveries
- /explore?domain=universe — Universe discoveries
- /discoveries/{slug} — Discovery detail

No sign-in is required.