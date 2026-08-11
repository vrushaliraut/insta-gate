# ADR 0009: Choice of Deployment and Hosting Platform

- **Status:** Accepted
- **Date:** 2026-08-11

## Context

The platform requires global edge hosting with low-latency execution for sub-50ms QR scans, persistent PostgreSQL
storage, Redis caching, and zero-downtime container updates.

## Decision

We select **Fly.io with Docker containers and GitHub Actions CI/CD**.

## Alternatives Considered & Rejection Reasons

1. **AWS (ECS / EKS + RDS):**
   - _Rejected:_ Excessive operational overhead, higher baseline cost, and complex IAM/VPC networking configurations
     for the required scope compared to Fly.io's developer-focused edge networking.
2. **Vercel + External Managed Database (Railway/Render):**
   - _Rejected:_ Vercel cannot host persistent Go WebSocket connections or background cron workers (e.g., nightly PII
     anonymization and overstay tickers) due to serverless execution limits.

## Consequences

- **Positive:**
  - Low-latency routing at edge locations near enterprise gates.
  - Unified container deployment for Go backend, Next.js frontend, and persistence layers.
  - Native GitHub Actions integration for zero-downtime rolling releases.
- **Negative:**
  - Requires strict monitoring of multi-region database replication latency if scaled globally in future phases.
