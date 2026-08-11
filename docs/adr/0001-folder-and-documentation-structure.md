# ADR 0001: Monorepo Folder and Documentation Structure

* **Status:** Accepted
* **Date:** 2026-08-11

## Context
The `insta-gate` project is an enterprise system combining a Go backend, Next.js multi-persona frontend, PostgreSQL database, Redis caching layer, and deployment automation. We need a clear, standardized location for product architecture, task tracking, and architectural decisions.

## Decision
We establish a top-level `/docs` directory in the repository containing:
1. `features.md`: Complete product specification and persona workflows.
2. `tasks.md`: Detailed master implementation plan organized by Epics and Stories.
3. `/docs/adr/`: Architecture Decision Records tracking key technological trade-offs.

## Alternatives Considered & Rejection Reasons
* **Scattering docs across service folders (`/backend/README.md`, `/frontend/README.md`):**
    * *Rejected:* Creates fragmented documentation where cross-cutting architectural context (e.g., RBAC or state transitions) is hard to maintain.
* **External Wiki (Notion/Confluence):**
    * *Rejected:* Documentation drifts from code when stored externally. Co-locating docs with code ensures version-controlled history aligned with PRs.

## Consequences
* **Positive:** Centralized source of truth accessible to all developers and CI tools.
* **Negative:** Requires strict PR discipline to update documentation when code changes.