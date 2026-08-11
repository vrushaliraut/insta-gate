# ADR 0008: Choice of Structured Logging Framework

- **Status:** Accepted
- **Date:** 2026-08-11

## Context

Security compliance and operational auditing require structured, machine-parsable JSON logs across all backend services,
with correlation tracing (`request_id`).

## Decision

We select Go's native standard library **`log/slog`**.

## Alternatives Considered & Rejection Reasons

1. **Uber Zap (`go.uber.org/zap`):**
   - _Rejected:_ Adds an external dependency. While Zap is extremely fast, `log/slog` delivers equivalent JSON
     structured logging performance without external maintenance risks.
2. **Zerolog (`rs/zerolog`):**
   - _Rejected:_ Requires non-standard logging patterns and extra allocation tuning. `log/slog` is integrated into the
     Go standard library, ensuring long-term compatibility.

## Consequences

- **Positive:**
  - Zero third-party dependency footprint.
  - Native support for JSON formatting, log levels, and contextual key-value attributes.
- **Negative:**
  - Requires custom wrapper middleware to inject contextual `request_id` values from HTTP headers into the `slog`
    context.
