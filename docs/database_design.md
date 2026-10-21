# Unseen V1 — Database Design

Database: PostgreSQL

## 1. Categories
- id: primary key
- domain: earth or universe
- name
- slug

Rules:
- Domain and slug together must be unique.

## 2. Discoveries
- id: primary key
- category_id: foreign key
- slug: unique
- title
- summary
- explanation
- location: optional
- facts: optional JSON object
- source_provider
- source_record_id
- source_url
- text_credit
- text_license
- text_license_url: optional
- image_url
- image_source_url
- image_alt
- image_kind: photograph, composite, illustration, or visualization
- image_credit
- image_license
- image_license_url: optional
- status: draft or published
- created_at
- updated_at
- last_fetched_at
- published_at: optional until published

Rules:
- Category determines Earth or Universe.
- Source provider and source record ID together must be unique.
- Publish only after content, image, and reuse conditions are verified.
- Public queries return published discoveries only.

## 3. Import Targets
- id: primary key
- category_id: foreign key
- provider: wikipedia or nasa
- target: approved article or search query
- enabled
- last_attempted_at
- last_succeeded_at
- last_error: optional

Rules:
- Provider and target together must be unique.
- Failed imports must preserve valid existing content.

## 4. Daily Discoveries
- discovery_date: primary key, UTC date
- discovery_id: foreign key
- selected_at

Rules:
- One discovery selection per day.
- Select and display published discoveries only.
- Selection runs independently of importing.

## Relationships
- One category has many discoveries.
- One category has many import targets.
- One discovery can be selected on multiple dates.
- Restrict deletion of referenced records; unpublish instead.

## Initial Indexes
- Published feed: published_at and id.
- Category feed: category_id, published_at, and id.
- Full-text search: title, summary, and explanation.
- Foreign-key columns.

## Keep Outside the Database
- Image files: use approved external URLs initially.
- Credentials and schedules: configuration/secrets.
- Search inputs and filters: request parameters.
- Related discoveries: calculate by category.
- Temporary API responses: memory.
- Detailed operational logs: application logging.

## Implementation
- Create tables and constraints through migrations.
- Seed initial categories.
- Test duplicate imports, publication rules, search, and daily selection.