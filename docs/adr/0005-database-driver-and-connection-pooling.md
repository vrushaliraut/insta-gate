# ADR 0005: Choice of PostgreSQL Driver and Connection Pooling

* **Status:** Accepted
* **Date:** 2026-08-11

## Context

`insta-gate` requires microsecond-level database interaction with PostgreSQL 18.4 to support sub-50ms QR validation at
enterprise gates. The database schema heavily utilizes PostgreSQL-native data types (JSONB for audit metadata, UUIDs,
UTC `timestamptz`, and Enums). We must choose a data access layer that maximizes throughput, minimizes memory
allocations, and provides strict control over query execution plans.

## Decision

We select **`jackc/pgx/v5` (`pgxpool`)** for database connectivity and pooling, bypassing the standard `database/sql`
library. For developer velocity, this is expected to be paired with explicit SQL builders or generators (like `sqlc`)
rather than an Object-Relational Mapper (ORM).

## Alternatives Considered & Rejection Reasons

1. **GORM (`gorm.io/gorm`):**
    * *Rejected:* While GORM provides excellent developer ergonomics and rapid prototyping speed, it relies heavily on
      Go runtime reflection and the generic `database/sql` interface. This introduces measurable CPU overhead and memory
      allocations. Furthermore, GORM's abstraction obscures the underlying query planner, making it difficult to
      optimize complex indexing for the `passes` table or prevent N+1 query regressions during telemetry aggregation.
2. **Standard `database/sql` with `lib/pq`:**
    * *Rejected:* `lib/pq` is officially in maintenance mode. The standard library abstraction prevents utilizing
      PostgreSQL-specific wire protocol optimizations, such as direct binary encoding/decoding of `JSONB` and array
      types.
3. **Ent (`entgo.io/ent`):**
    * *Rejected:* Though more type-safe and performant than GORM due to its code-generation approach, Ent's graph-based
      architecture introduces unnecessary complexity for a system that fundamentally requires straightforward,
      high-speed relational transactions and immutable audit appending.

## Consequences

* **Positive:**
    * Highest possible throughput via the native PostgreSQL wire protocol.
    * Direct, zero-allocation support for complex data types (`JSONB`, `UUID`, `timestamptz`).
    * `pgxpool` provides superior connection lifecycle management and pipeline querying capabilities compared to
      standard library pooling.
    * Explicit SQL usage ensures predictable query execution plans, which is vital for database performance tuning.
* **Negative:**
    * Loss of auto-migration and rapid CRUD helper functions provided by ORMs.
    * Increased verbosity for basic database operations (mitigated by adopting `sqlc` or generic repository patterns).