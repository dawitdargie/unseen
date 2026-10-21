# Unseen — Data Sources and Feasibility

Status: Initial samples verified; full feasibility testing remains.
Verification date: 2026-10-21

## 1. Content Direction

Unseen helps visitors discover extraordinary places on Earth and
fascinating objects and phenomena in the universe.

The experience should make visitors ask:
“How can this exist?”

Content must be visually compelling, scientifically grounded, and
supported by clear source attribution.

## 2. Earth Content Scope
Unseen explores places and phenomena that make visitors ask:
“How can this exist on Earth?”

Earth discoveries include beautiful places, unbelievable landscapes,
extreme environments, and strange natural phenomena.

### Content Categories

- Alien-looking landscapes.
- Unusual geological formations.
- Extreme environments.
- Underground worlds and caves.
- Isolated ecosystems and unusual adaptations.
- Unusual lakes and natural chemical phenomena.
- Volcanoes and geothermal landscapes.
- Rare natural events and optical phenomena.
- Remarkable scenic places with a distinctive scientific story.

### Candidate Topics Available on Wikipedia

| Topic | Discovery angle | Wikipedia article |
|---|---|---|
| Socotra | An isolated island with unusual endemic life | https://en.wikipedia.org/wiki/Socotra |
| Danakil Depression | Extreme landscapes, salt formations, and geothermal activity | https://en.wikipedia.org/wiki/Danakil_Depression |
| Movile Cave | An unusual underground ecosystem supported by chemosynthesis | https://en.wikipedia.org/wiki/Movile_Cave |
| Lake Natron | An unusual alkaline lake and its ecosystem | https://en.wikipedia.org/wiki/Lake_Natron |

These are candidate topics, not automatically approved discovery records.
Each requires suitable text, a relevant image, and verified reuse conditions.

### Selection Rule

Every Earth discovery must answer:
1. What makes this place or phenomenon unusual?
2. How does it exist or occur?
3. Where is it found?
4. What evidence supports the explanation?

Do not use fabricated mysteries, exaggerated claims, or misleading images.
Preserve scientific uncertainty where explanations remain incomplete.

## 3. Earth Sources

### Wikipedia — Explanations and Topic Information

API endpoint:
https://en.wikipedia.org/w/api.php

Use Wikipedia to retrieve:
- Article title and stable page identifier.
- Introductory explanation.
- Original article reference.
- Available article image information.
- Relevant topic information where supported.

Wikipedia contains extraordinary places and natural phenomena.
However, it does not provide a ready-made, consistently classified
“unbelievable places” feed.

### Wikimedia Commons — Images and Image Metadata

API endpoint:
https://commons.wikimedia.org/w/api.php

Use Commons to retrieve:
- Relevant image files.
- Image dimensions.
- File description-page URL.
- Creator and credit information when available.
- License and attribution requirements.

Do not assume every Wikipedia article image is suitable for publication.
Do not assume every image is hosted on Commons or has the same license.

### Verified Earth Sample

Topic: Socotra

Observed results:
- Wikipedia API returned HTTP 200.
- Response included a page ID, explanatory text, and an image URL.
- The tested image returned HTTP 200 as JPEG.
- Commons returned public-domain metadata for that particular image.

Additional article availability checked:
- Danakil Depression.
- Movile Cave.
- Lake Natron.

Article availability does not verify the complete import pipeline,
image suitability, or every factual claim.

### Earth Selection Approach

For the first version:
1. Maintain a curated list of extraordinary topics.
2. Fetch their content automatically.
3. Validate explanations, images, and reuse conditions.
4. Publish approved discoveries.
5. Refresh stored records periodically.

Later:
- Use focused searches or relevant categories to find candidates.
- Review candidates before publishing.
- Expand the approved topic list.

Do not populate the feed with random Wikipedia pages.
Automatic fetching and editorial selection are separate responsibilities.

## 4. Universe Content Scope

Include:
- Planets and moons.
- Stars and stellar evolution.
- Nebulae.
- Galaxies.
- Other cosmic objects and phenomena supported by suitable records.

Each discovery should explain what the visitor is seeing and why it matters.

Distinguish:
- Observational images.
- Processed or composite images.
- Scientific visualizations.
- Artist impressions.

Do not present an illustration as a direct photograph.

## 5. Universe Source

### NASA Image and Video Library

API endpoint:
https://images-api.nasa.gov

Verified request:
https://images-api.nasa.gov/search?q=nebula&media_type=image&page_size=1

Use suitable records for:
- Titles.
- Astronomy images.
- Scientific descriptions.
- Source identifiers.
- Credit information.
- Available keywords and dates.

### Verified Universe Sample

Record ID: PIA14417
Title: Weighing in on the Dumbbell Nebula

Observed results:
- Search returned HTTP 200.
- Record included a description, credits, and image URLs.
- The tested image returned HTTP 200 as JPEG.

### Limitations

- The library is a media archive, not a complete astronomical database.
- Scientific facts are not provided in a uniform object schema.
- Search results can contain irrelevant material or duplicate subjects.
- Image dates do not necessarily represent discovery dates.
- Older descriptions may need checking against newer scientific knowledge.
- Image processing and illustration labels must be preserved.
- Missing scientific facts must not be invented.

### Universe Selection Approach

1. Search using focused astronomy topics and image filtering.
2. Validate relevance and image-description consistency.
3. Preserve scientific qualifications and uncertainty.
4. Check attribution and reuse conditions.
5. Publish approved discoveries.
6. Refresh and expand the collection periodically.

Use additional authoritative sources if a detail page needs facts
that the selected NASA record does not provide.

## 6. How Unseen Becomes Dynamic

Dynamic content does not require newly discovered places every day.

Unseen will:
1. Fetch approved topics and search for candidate records on a schedule.
2. Validate and normalize external responses.
3. Save accepted records in PostgreSQL.
4. Update existing records without creating duplicates.
5. Serve feeds, search, filters, and detail pages from PostgreSQL.
6. Select a daily discovery from eligible published records.
7. Keep existing content available when external services fail.

Fetching new content and rotating the daily discovery are separate jobs.

The frontend should not depend on live external API requests
for every visitor interaction.

## 7. Publication Requirements

Publish a discovery only when it has:
- A stable source identifier.
- An original source URL.
- An accurate title.
- A useful explanation.
- A relevant, accessible image.
- Known reuse conditions.
- Required attribution information.
- An Earth or Universe classification.
- An appropriate category.
- A clear reason it belongs in Unseen.

A successful HTTP response alone does not make a record publishable.

Hold incomplete or questionable records for review.

## 8. Licensing and Attribution

### Wikipedia Text

- Follow the applicable attribution and share-alike requirements.
- Preserve the source article reference.
- Record the retrieved revision where available.
- Identify adaptations where required.

### Wikimedia Images

- Check each file's individual license.
- Preserve creator, credit, license, and file-page information.
- Exclude files whose reuse conditions cannot be established.
- A working thumbnail URL is not proof of reuse permission.

### NASA Media

- Follow NASA's media-use guidelines.
- Preserve supplied credits.
- Inspect third-party rights notices.
- Do not assume every NASA-hosted asset is unrestricted.
- Do not imply NASA endorses Unseen.

Text and image licensing must be handled separately.

## 9. Initial Refresh Policy

These are proposed Unseen settings, not provider update guarantees:

- Candidate import: once daily.
- Daily discovery selection: once daily.
- Daily selection timezone: UTC.
- Existing record refresh: weekly initially.
- Failed requests: bounded retries with backoff.
- Invalid updates: retain the last valid published version.
- Missing images or unclear rights: hold the candidate for review.

Respect current provider request policies and rate limits.
The application schedule must adapt to those limits.

## 10. Failure Handling

The ingestion process must handle:
- Request timeouts.
- Rate limiting.
- Provider outages.
- Missing fields.
- Redirected or missing articles.
- Unavailable image URLs.
- Invalid responses.
- Duplicate imports.
- Changed attribution or licensing metadata.

Record import outcomes and errors for troubleshooting.

Do not overwrite valid published content with an incomplete response.
Do not repeatedly retry permanent failures without review.

## 11. Excluded for Now

### NASA Astronomy Picture of the Day — APOD

The response tested on 2026-10-21 contained a nebula explanation
but supplied the NASA logo as its image.

Do not integrate APOD until image and metadata consistency pass
additional testing.

This observation concerns the tested response, not every APOD record.

## 12. Remaining Feasibility Checks

Before closing Step 03:
- Test multiple records across every planned category.
- Verify extraordinary Earth topics have suitable publishable images.
- Check that explanations support useful detail pages.
- Verify image relevance, resolution, credits, and licenses.
- Confirm current rate limits and request policies.
- Test pagination and duplicate detection.
- Test redirects and missing fields.
- Test timeouts, throttling, and unavailable images.
- Check how source changes are detected.
- Save representative sample responses.
- Record observed failure behavior.
- Confirm attribution can be displayed correctly.

## 13. Decision

Use these as the first sources to evaluate and integrate:

- Earth explanations: Wikipedia.
- Earth images: Wikimedia Commons.
- Universe images and descriptions: NASA Image and Video Library.

Use curated selection plus automatic fetching and validation.

These sources have passed initial sample checks.
Step 03 remains open until the broader feasibility checks pass.

No external source is guaranteed to remain available or error-free.
Unseen must preserve useful content independently of temporary outages.

## 14. Official References

NASA Image Library API documentation:
https://images.nasa.gov/docs/images.nasa.gov_api_docs.pdf

NASA media guidelines:
https://www.nasa.gov/nasa-brand-center/images-and-media/

MediaWiki Action API:
https://www.mediawiki.org/wiki/API:Main_page

MediaWiki image metadata:
https://www.mediawiki.org/wiki/API:Imageinfo

MediaWiki request guidelines:
https://www.mediawiki.org/wiki/API:Etiquette

Wikimedia terms:
https://foundation.wikimedia.org/wiki/Policy:Terms_of_Use