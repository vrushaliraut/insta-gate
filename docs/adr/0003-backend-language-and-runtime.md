# ADR 0003: Choice of Backend Language and Runtime

* **Status:** Accepted
* **Date:** 2026-08-11

## Context
The backend engine must verify cryptographic JWT QR passes in sub-50ms, manage WebSocket connections for live occupancy updates, handle background PII scrubbing, and run reliably in lightweight Docker containers.

## Decision
We select **Go (Golang) 1.26.5**.

## Alternatives Considered & Rejection Reasons

1. **Node.js / TypeScript:**
    * *Rejected:* The single-threaded Event Loop in Node.js can suffer latency spikes under heavy CPU-bound cryptographic token verification (RSA/HMAC) and high-throughput concurrent scanning at peak gate traffic hours.
2. **Rust:**
    * *Rejected:* While offering high performance, Rust introduces longer compile times and higher boilerplate overhead for API gateway setup, slowing down execution velocity without delivering meaningful latency gains over Go for this application context.
3. **Python (FastAPI):**
    * *Rejected:* High memory consumption and GIL (Global Interpreter Lock) constraints make Python unsuitable for real-time WebSocket state distribution and high-concurrency microservice workloads.

## Consequences
* **Positive:**
    * Sub-millisecond JWT verification capabilities.
    * Native concurrency model (goroutines) for background workers and WebSockets.
    * Produces small, static single-binary Docker images.
* **Negative:**
    * Requires explicit error handling patterns compared to exception-driven languages.