# Insta-Gate: Enterprise Smart Visitor Management System

Insta-Gate is a commercial-grade Smart Visitor Management System (SVMS) powered by a Dynamic QR Code Engine. It merges
the cryptographic, time-bound security of a visitor management system with the robust analytics, branding, and dynamic
routing of an enterprise QR platform.

## 🚀 Key Features

- **Dynamic Cryptographic Identity:** QR codes encode cryptographically signed JWTs that are validated in sub-50ms at
  the edge.
- **Strict Role-Based Access Control (RBAC):** Distinct interfaces and data scopes for Super Admins, Facility Admins,
  Hosts (Employees), and Security Guards.
- **Conditional Routing:** A single QR code dynamically routes users to a check-in pass (mobile) or a facility map (
  desktop).
- **Live Occupancy & Telemetry:** WebSocket-powered dashboards provide real-time heatmaps and overstay alerts.
- **Privacy by Design:** Granular data masking at the security gate and automated background workers for GDPR/DPDP
  compliant PII purging.

## 🛠️ Technology Stack

- **Backend:** Go 1.26.5, **Go Fiber v3.4.0** (High-performance web framework), `slog` (structured JSON logging).
- **Frontend:** Next.js 16.3 (App Router), React 19.2, TypeScript 7.0, Tailwind CSS, Playwright (E2E testing).
- **Database & Caching:** PostgreSQL 18.4 (`pgxpool`, `golang-migrate`), Redis (Session state & revocation blocklists).
- **Infrastructure:** Docker, GitHub Actions, Fly.io (Global Edge Network).
- **Orchestration:** `Task` (Taskfile).

## 📂 Repository Structure

This project follows a strict monorepo architecture to ensure synchronized deployments and centralized documentation.

```text
insta-gate/
├── backend/       # Go 1.26.5 API Gateway (Fiber v3.4.0), JWT Engine, WebSockets
├── frontend/      # Next.js 16.3 multi-persona web applications and PWA
├── database/      # PostgreSQL migration files (.up.sql / .down.sql)
├── deploy/        # Dockerfiles, fly.toml, and GitHub Actions workflows
├── docs/          # Project documentation hub
│   ├── adr/       # Architecture Decision Records (ADRs)
│   ├── features.md# Core Use Cases & Feature List
│   └── tasks.md   # Master Implementation Plan (Epics & Stories)
└── Taskfile.yml   # Centralized development task runner
```
