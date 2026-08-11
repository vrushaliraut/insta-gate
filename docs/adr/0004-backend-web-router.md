# ADR 0004: Choice of Backend Web Framework

- **Status:** Accepted
- **Date:** 2026-08-11

## Context

The Go backend requires a high-performance web framework to handle REST API endpoints, WebSocket upgrades, rapid rate
limiting, and strict RBAC middleware chaining. The system must verify cryptographic JWT QR passes in sub-50ms under
heavy concurrent load at enterprise security gates.

## Decision

We select **Go Fiber v3.4.0**.

## Alternatives Considered & Rejection Reasons

1. **`go-chi/chi` or standard `net/http` ServeMux:**
   - _Rejected:_ While 100% standard library compatible, they require significantly more boilerplate to handle common
     enterprise API requirements (e.g., standardized JSON responses, parameter parsing, and rate limiting). They also
     do not match the zero-allocation, extreme high-throughput performance characteristics of `fasthttp`, which Fiber
     utilizes under the hood.
2. **Gin (`gin-gonic/gin`):**
   - _Rejected:_ Gin uses a custom `gin.Context` pattern but trails behind Fiber in raw benchmark throughput and memory
     efficiency. Fiber's Express.js-like API is more ergonomic for rapid development and includes more robust
     out-of-the-box middleware for modern enterprise needs.

## Consequences

- **Positive:**
  - Extremely high performance and low memory footprint due to the underlying `fasthttp` engine.
  - Expressive, developer-friendly API that drastically reduces boilerplate for request parsing and JSON
    serialization.
  - Comprehensive suite of built-in, officially maintained middleware (Rate Limiting, CORS, WebSockets, Logger).
- **Negative:**
  - Fiber does not use Go's standard `net/http` interfaces. This means third-party packages strictly requiring
    `http.Handler` will need Fiber-specific adapter wrappers (e.g., `fiber/v3/middleware/adaptor`).
  - Developers must be cautious with handler contexts (`ctx`), as Fiber optimizes by reusing them. Passing the context
    to background goroutines requires manual copying to prevent data corruption.
