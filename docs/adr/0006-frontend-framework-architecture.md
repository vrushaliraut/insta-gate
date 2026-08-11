# ADR 0006: Choice of Frontend Framework and Architecture

* **Status:** Accepted
* **Date:** 2026-08-11

## Context

The platform must serve four distinct personas across desktop, web, and tablet environments:

* Host Portal (Desktop web)
* Walk-In Self-Registration (Mobile localized web)
* Guard Tablet (Offline-capable high-speed PWA)
* Admin Dashboard (Analytics and telemetry web)

## Decision

We select **Next.js 16.3 (App Router) with TypeScript 7.0 and React 19.2**.

## Alternatives Considered & Rejection Reasons

1. **Vite + React Single Page Application (SPA):**
    * *Rejected:* Lacks server-side rendering (SSR), dynamic edge route parsing, and built-in optimization for localized
      language loading (`next-intl`) required during dynamic walk-in QR scans.
2. **Remix / React Router v7:**
    * *Rejected:* Next.js provides superior standard PWA caching hooks, server component optimizations, and standalone
      Docker deployment outputs optimized for Fly.io infrastructure.

## Consequences

* **Positive:**
    * Unified codebase for SSR dashboards, localized mobile forms, and PWA tablet interfaces.
    * Standalone Docker build output minimizes container sizes.
    * Server Components reduce JavaScript bundle sizes sent to low-end mobile devices.
* **Negative:**
    * Strict distinction between Server and Client Components requires careful state management design.