# Unseen V1 — Components and Responsibilities

## 1. Frontend — Next.js + TypeScript
- Display the feed, daily discovery, Earth/Universe pages, and details.
- Provide search, filters, pagination, and related discoveries.
- Show attribution and responsive loading, empty, and error states.

## 2. Backend API — Go
- Serve published discoveries from PostgreSQL.
- Handle search, filters, pagination, details, and daily discovery.
- Validate requests and provide consistent errors and health endpoints.

## 3. Database — PostgreSQL
- Store discoveries, categories, sources, and publication status.
- Store image URLs, credits, licenses, and daily selections.
- Prevent duplicate imports and record import outcomes.

## 4. Content Importer — Go
- Fetch content from Wikipedia, Wikimedia Commons, and NASA.
- Normalize and validate records, images, and attribution.
- Save valid content, prevent duplicates, and handle source failures.
- Hold questionable content for review.

## 5. Scheduler
- Trigger imports and content refreshes.
- Run daily discovery selection independently of external fetching.

## 6. External Sources
- Wikipedia: Earth explanations.
- Wikimedia Commons: Earth images and license metadata.
- NASA Image Library: Universe images, descriptions, and credits.

## Shared Responsibilities
- Environment configuration and secure secrets.
- Structured logs and important behavior tests.
- Database migrations, backups, and recovery instructions.

## Structure
One repository. API and importer share relevant Go code.
Docker packages the application; Kubernetes and AWS follow during learning.