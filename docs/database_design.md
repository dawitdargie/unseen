# Unseen V1 — Database Design

Database: PostgreSQL. Design baseline: 2026-10-22.

Create tables incrementally through versioned SQL migrations, at the milestones in [implementation_plan.md](implementation_plan.md). This model stores the current publication and reviewable candidates; a complete editorial CMS/history is deferred.

## 1. Shared Types

- Entity IDs: `bigint GENERATED ALWAYS AS IDENTITY` primary keys unless specified otherwise.
- Timestamps: `timestamptz`, serialized in UTC. Daily keys are SQL `date` calculated in UTC.
- Text: `text`; require nonblank values where specified. Slugs are lowercase URL-safe values.
- Domain/status/provider/image-kind values use CHECK constraints.
- Foreign keys restrict deletion of referenced records. Unpublish discoveries instead.
- Category determines domain; do not duplicate a discovery-domain column.

## 2. Categories

- `id`: primary key.
- `domain`: required `earth` or `universe`.
- `name`, `slug`: required.
- Unique `(domain, slug)`.

Initial rows: Earth `landscapes`, `caves`, `ecosystems`; Universe `planets-moons`, `nebulae`, `galaxies`. Frontend decorative labels are not categories.

## 3. Discoveries — Current Approved Snapshot

- `id`: primary key; `category_id`: required foreign key.
- `slug`: globally unique, required; stable after publication. Reserve `daily`.
- `title`, `summary`, `explanation`: required nonblank content.
- `location`: optional label; not always geographic for Universe content.
- `facts`: JSONB object, default `{}`; CHECK that its type is object.
- `source_provider`: required `wikipedia` or `nasa`.
- `source_record_id`: stable provider ID stored as text.
- `source_url`, `text_credit`, `text_license`: required when published.
- `text_license_url`: optional.
- `source_revision_id`, `source_revision_url`: optional available article provenance.
- `text_adaptation_note`: optional description of adaptations when applicable.
- `source_content_hash`: optional for manual seeds; set when applying an imported candidate.
- `snapshot_version`: positive bigint, default 1; increment for every content or publication-status change.
- `image_url`, `image_source_url`, `image_alt`, `image_credit`, `image_license`: required when published.
- `image_kind`: `photograph`, `composite`, `illustration`, or `visualization`; required when published.
- `image_license_url`: optional.
- `status`: `draft` or `published`, default `draft`.
- `created_at`, `updated_at`: required timestamps.
- `last_fetched_at`: optional until fetched from a provider.
- `published_at`: optional until first publication; preserve on refresh.

Rules:

- Unique `(source_provider, source_record_id)` prevents repeat record imports, not several NASA assets of one subject. Editorial selection handles subject duplicates.
- CHECK constraints enforce required public fields/allowed values and nonnull `published_at` for published rows. Drafts can have incomplete attribution but still need valid identity/category/core text.
- Application validation checks URLs and editorial rules that SQL cannot prove. License fields describe the applicable license/reuse terms; never invent a universal license for a provider.
- Publish only after content/image relevance and reuse review. Public queries restrict to published rows.
- Refresh updates content/provenance/`updated_at`/`last_fetched_at`, preserving slug and original `published_at`.
- Manual edits to hashed content must recompute the canonical `source_content_hash` or clear it; never compare an old imported hash to edited content. Status-only changes can retain the hash while still incrementing `snapshot_version`.
- Reads never mutate publication/daily-selection state.

## 4. Import Targets — Curated Identities

- `id`: primary key; `category_id`: required foreign key.
- `provider`: `wikipedia` or `nasa`.
- `target`: stable English Wikipedia page ID or NASA record ID as text.
- `enabled`: boolean, default true.
- `last_attempted_at`, `last_succeeded_at`: optional timestamps.
- `last_error`: optional bounded diagnostic text, without secrets/full upstream bodies.
- Unique `(provider, target)`.

Maintain the approved target catalog explicitly. Search manually to select IDs; scheduled search-query ingestion is deferred. One source identity has one category assignment. Catalog synchronization must not silently reclassify published content.

Daily candidate import selects enabled targets not yet published; weekly refresh selects enabled published targets. An unchanged successful fetch updates target metadata without creating a new publication date or candidate.

Classify permanent failures (confirmed missing/invalid targets or unsupported metadata) by setting enabled=false and last_error with a reviewable reason. Scheduled retries skip disabled targets; manual review explicitly re-enables them. Transient failures keep the target enabled and use bounded retries. Pending editorial review is not itself a provider failure.

## 5. Import Candidates — Proposed Snapshots

- `id`: primary key.
- `import_target_id`: required foreign key.
- `discovery_id`: optional foreign key; null for new content until application.
- `source_provider`, `source_record_id`: required identity.
- `category_id`: required foreign key, matching the approved target when staged.
- `payload`: JSONB object containing `schema_version` plus normalized proposed content/image/provenance using discovery field names.
- `content_hash`: canonical editorial hash; excludes volatile fetch timestamps.
- `base_snapshot_version`: nonnegative bigint; current discovery snapshot_version when staged, or 0 if no discovery exists.
- `status`: `pending`, `applied`, or `rejected`, default `pending`.
- `fetched_at`, `created_at`: required timestamps.
- `reviewed_at`: required once applied/rejected; optional while pending.
- `review_note`: optional review decision/concern.
- Unique `(source_provider, source_record_id, base_snapshot_version, content_hash)`.

Relational identity/category/state stay outside JSONB. A versioned validator checks the payload and that its duplicated identity/category match the relational values. Read the current discovery hash/version together when staging. If the proposed editorial hash already matches, record fetch success without a new candidate; otherwise capture that version as the base. A concurrent publication after this read makes the candidate stale at application, which is safe. Structurally invalid responses update target error state/logs without entering this table. Structurally valid but questionable content stays pending.

Apply in one transaction:

1. Acquire a transaction-level advisory lock for the source identity, including when no discovery row exists yet.
2. Lock/check the candidate. Applied is a no-op; rejected cannot be applied implicitly. Revalidate approved category/source/payload.
3. Lock/read the current discovery when present and require its snapshot_version to equal base_snapshot_version (absence is version 0). A mismatch is stale: reimport and review against the current version. Preserve provider revision checks where meaningful.
4. Insert the new snapshot at version 1 or update it at current version + 1, preserving established slug and first-publication time. Manual edits and publish/unpublish commands must increment this version too.
5. Mark candidate applied with discovery/review metadata and commit together.

No network fetch occurs inside this transaction. Unique source identity remains the final duplicate guard. Identical fetches against the same base do not reopen reviewed candidates; reconsidering rejection requires an explicit reviewed transition back to pending. After an applied A → B change, a real upstream reversion to A is a new candidate against B's newer base version, even though A's hash appeared before. No historical applied candidate is replayed automatically.

Per-provider import-run coordination uses a dedicated connection/session advisory lock released on that same connection. Do not return a locked connection to the pool or keep a database transaction open across external HTTP requests. This prevents manual/scheduled overlap without treating a CronJob setting as a global lock.

## 6. Daily Discoveries

- `discovery_date`: UTC SQL `date`, primary key.
- `discovery_id`: required foreign key.
- `selected_at`: required timestamp.

Select at `00:05 UTC`, independently of importing; seed manually during setup. A still-published selection makes rerun a no-op. Otherwise choose deterministically from ascending eligible IDs using UTC days since 1970-01-01 modulo the eligible list length, excluding yesterday before choosing when alternatives exist. One eligible record may repeat; none means no new row.

Transactions, unique date, appropriate locks, and publication rechecks protect concurrent insert/replacement. Unpublishing keeps historical rows. GET uses today's published selection, otherwise the most recent earlier daily selection still published, otherwise newest published discovery. GET never repairs storage; the command may replace today's invalid selection.

## 7. Relationships and Indexes

One category has many discoveries/targets; a target has many candidate snapshots; candidates can reference their current discovery; a discovery can be selected on multiple dates.

Add indexes with the relevant query:

- Unique category domain/slug, discovery slug/source identity, target identity, candidate identity/base-version/hash.
- Partial published feed `(published_at DESC, id DESC)`.
- Partial published category feed `(category_id, published_at DESC, id DESC)`.
- Weighted English title/summary/explanation search vector with GIN; query/index expressions must agree.
- Foreign-key lookup indexes and candidate review `(status, created_at, id)` where useful.

Inspect plans before speculative indexes. Offset pagination is sufficient. Counts/results within one request must share a consistent snapshot (one statement or a read-only repeatable-read transaction); separate page requests do not freeze the catalog.

## 8. Keep Outside the Database

- Image binaries: approved external URLs initially.
- Credentials/schedules: runtime secrets/configuration.
- Search/filter UI state: URL parameters.
- Related discoveries: calculate from published category membership.
- Detailed logs/full provider response archives: logging and sanitized fixtures.

Test constraints/migrations, publication visibility, SQL behavior, duplicate staging, atomic/stale apply, concurrent daily selection, and failure preservation with the milestones that introduce them.
