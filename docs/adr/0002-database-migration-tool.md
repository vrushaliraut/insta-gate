# ADR 0002: Choice of Database Migration Tool

- **Status:** Accepted
- **Date:** 2026-08-11

## Context

`insta-gate` relies on PostgreSQL 18.4 for audit logs, tenant configurations, and visitor states. We require a schema migration tool that is version-controlled, explicit, reproducible, and seamlessly integrates into Go and Docker pipelines.

## Decision

We select **`golang-migrate/migrate`** using plain SQL (`.up.sql` and `.down.sql`) files.

## Alternatives Considered & Rejection Reasons

1. **Goose (`pressly/goose`):**
   - _Rejected:_ Goose supports migration files written in Go code. While powerful, embedded Go migrations can introduce runtime complexity, non-reproducible build artifacts, and migration race conditions in containerized deployments.
2. **Atlas (`ariga/atlas`):**
   - _Rejected:_ Atlas relies on declarative schema management and auto-diffing. In an enterprise system requiring immutable append-only audit logs (`audit_logs` table), implicit auto-generated diffs present a risk of destructive schema changes during deployment.
3. **ORM Auto-migrations (GORM / Ent):**
   - _Rejected:_ Auto-migrations are non-deterministic, obscure raw SQL performance (e.g., custom indexes on `passes.status`), and fail to support strict database rollbacks required by our CI/CD safety gates.

## Consequences

- **Positive:**
  - Clean separation of raw SQL DDL files.
  - Direct execution via CLI (`task db:migrate`) and Go binary integration.
  - Immutable history with strict version locking.
- **Negative:**
  - Developers must manually write matching `.down.sql` rollback scripts for every migration.
