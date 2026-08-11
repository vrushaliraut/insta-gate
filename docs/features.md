# Enterprise Smart Visitor Management & Dynamic Identity Platform

This document outlines the complete product blueprint for a commercial-grade **Smart Visitor Management System (SVMS) powered by a Dynamic QR Code Engine**. It merges the cryptographic, time-bound security of a visitor management system with the robust analytics, branding, and dynamic routing of an enterprise QR platform.

The system is built on **Go 1.26.5, TypeScript 7.0, React 19.2, Next.js 16.3, and PostgreSQL 18.4**, deployed via **Fly.io** using **Docker** and **GitHub Actions**, with **Task** orchestrating the development lifecycle.

---

## 1. Core Use Cases & User Journeys

The platform serves four distinct personas. Each persona interacts with a tailored Next.js interface governed by strict Role-Based Access Control (RBAC).

### Use Case A: The Pre-Approved Visitor (Authority Generated)

**Persona:** Host / Employee

1. An employee (Host) schedules a meeting and logs into the Host Portal.
2. They input the visitor's Name, Phone Number, and Purpose of Visit.
3. The platform's QR Engine generates a dynamic, cryptographically signed QR code link and sends it via SMS/Email.
4. The QR code contains conditional routing: if opened on a desktop, it shows a map to the building; if opened on a mobile device, it displays the high-contrast QR pass.
5. The pass is strictly time-bound (e.g., valid only on Tuesday, 9:00 AM - 12:00 PM).

### Use Case B: The Walk-In Visitor (Self-Generated)

**Persona:** Unexpected Visitor / Vendor

1. A delivery driver arrives at the gate and scans a permanent, branded "Walk-In Registration" static QR code.
2. The dynamic routing engine detects their device language and serves a localized Next.js web form.
3. The visitor inputs their details (Name, Phone, ID Photo) and selects the person they wish to meet.
4. The system issues a "Pending" QR pass.
5. The requested Host receives a push notification to approve or deny. Upon approval, the visitor's Pending QR dynamically updates to an "Approved" state via real-time WebSockets, granting a 1-hour entry window.

### Use Case C: Security Verification & Auditing

**Persona:** Security Guard

1. Guards use a highly optimized Next.js Progressive Web App (PWA) on enterprise tablets.
2. The visitor presents their QR pass. The guard scans it.
3. The Go backend cryptographically verifies the token in sub-50ms.
4. The tablet flashes **GREEN** (Valid), showing the visitor's photo, name, and host. (Address and phone number are masked).
5. The system logs the exact millisecond of entry. If the visitor leaves, a second scan logs the exit, updating the live building occupancy count.

### Use Case D: Global Visibility & Facility Management (NEW)

**Persona:** Facility Admin / Super Admin

1. The Facility Admin logs into the Central Admin Dashboard.
2. They view real-time telemetry: "Currently 142 visitors in Building A, 3 pending approvals, 2 overstay alerts."
3. The Admin can drill down into RBAC configurations, adding new Security Guards, provisioning new Departments, or globally revoking a specific visitor's access.
4. They can export immutable audit logs for compliance audits (e.g., "Show me everyone who visited the 4th floor server room last month").

---

## 2. Extensive Feature List

### Role-Based Access Control (RBAC) & Visibility

- **Super Admin:** Manages multiple buildings/tenants. Can access billing, configure global data retention policies, and view cross-facility analytics.
- **Facility Admin:** Scoped to a specific building. Can manage guard accounts, view live occupancy dashboards, force-checkout overstayed visitors, and export daily visitor manifests.
- **Host (Employee):** Scoped to their own visitors. Can generate passes, approve walk-ins, and view their personal visitor history. Cannot see other employees' visitors.
- **Security Guard:** Scoped to the verification edge. Can scan passes, view masked profile data (photo, name, host), and trigger check-in/out state changes. No access to PII (phone/address) or historical logs.

### Dynamic QR & Identity Engine

- **Cryptographic JWT Tokens:** QR codes do not store raw data. They encode a signed JSON Web Token (JWT) that the Go backend decodes. If a user tries to forge a QR code, the signature validation fails instantly.
- **Custom Branding Studio:** Facility Admins can customize the appearance of the QR passes (colors, company logos embedded in the QR matrix, custom error correction levels).
- **State Machine Engine:** Passes exist in strict, immutable states: `DRAFT` -> `PENDING` -> `APPROVED` -> `SCANNED_IN` -> `SCANNED_OUT` -> `EXPIRED` / `REVOKED`.
- **Instant Revocation:** A Host or Admin can click "Revoke" on the dashboard. The Go backend invalidates the token in Redis, rendering the visitor's QR code completely useless within milliseconds.

### Advanced Admin Dashboard & Analytics

- **Live Occupancy Heatmaps:** Visual representations of building load. Tracks peak hours to help optimize security guard shift scheduling.
- **Overstay Alerting System:** If a pass was issued for a 2-hour window and the visitor has not been scanned out by hour 3, the dashboard flags the visitor in red and alerts the Facility Admin and Host.
- **Granular Audit Trails:** Every action is logged using Go structured logging. Records include: `Timestamp`, `Actor ID`, `Action (e.g., APPROVED_PASS)`, `Target ID`, and `IP Address`.

### Security & Data Privacy Compliance

- **Dynamic Data Masking:** The Next.js frontend strictly limits what the Guard persona can see, ensuring compliance with privacy laws by hiding home addresses and personal contact info at the gate.
- **Automated PII Purging (GDPR/DPDP):** A Go background worker runs nightly, anonymizing visitor records older than 30 days (deleting Name, Phone, and Photo, but keeping the "1 visitor entered" statistic).

---

## 3. System Architecture & Workflows

The platform leverages a highly concurrent backend to ensure the security line never halts, coupled with modular React frontends for different personas.

```text
===================================================================================
                       UNIFIED ENTERPRISE SVMS ARCHITECTURE
===================================================================================

      [ WALKIN DEVICE ]       [ HOST PORTAL ]      [ GUARD TABLET ]    [ ADMIN PANEL ]
      (Next.js 16.3 Web)      (Next.js Web)        (Next.js PWA)       (Next.js Web)
      * Self-Registration     * Invite Gen         * Fast Verify       * Live Dashboards
      * Dynamic Routing       * Approvals          * Check In/Out      * RBAC & Logs
             |                      |                    |                   |
             +----------------------+--------------------+-------------------+
                                            |
                                  [ FLY.IO EDGE NETWORK ]
                                    (Global Load Balancer)
                                            |
===================================================================================
                             GO 1.26.5 BACKEND (DOCKERIZED)
===================================================================================
                                            |
                                  [ API GATEWAY / ROUTER ]
                                            |
        +-----------------------------------+-----------------------------------+
        |                                   |                                   |
 [ IDENTITY & QR ENGINE ]        [ RBAC & STATE MANAGER ]            [ AUDIT & TELEMETRY ]
 * JWT Signing/Validation        * Role Permissions                  * Slog (JSON Logging)
 * Image/Vector Generation       * Pass TTL & Validity               * Occupancy Aggregation
 * Conditional Routing           * WebSockets (Live Updates)         * PII Purging Worker
        |                                   |                                   |
        +-----------------------------------+-----------------------------------+
                                            |
===================================================================================
                                   DATA PERSISTENCE
===================================================================================
                                            |
        +-----------------------------------+-----------------------------------+
        |                                                                       |
 [ REDIS (In-Memory) ]                                               [ POSTGRESQL 18.4 ]
 * Session State                                                     * Relational User Data
 * Revocation Blocklists                                             * Immutable Audit Logs
 * Sub-millisecond Lookups                                           * Historical Analytics
===================================================================================

```

---

## 4. Engineering & CI/CD Operations

The engineering lifecycle is optimized for reliability, utilizing a task runner, extensive testing, and automated deployments.

### Development Orchestration (`Taskfile`)

Instead of bash scripts or Makefiles, development is orchestrated via `task`. The `Taskfile.yml` handles:

- `task dev`: Concurrently boots the Go backend, Next.js frontend, and local Postgres/Redis via Docker Compose.
- `task db:migrate`: Runs strict up/down database schema migrations.
- `task lint`: Enforces Go `golangci-lint` and TS `eslint` standards.

### Quality Assurance & Testing

- **Unit Tests:** Go's native testing framework verifies cryptographic token generation, RBAC permission matrices, and state machine transitions.
- **End-to-End (E2E) Tests:** Playwright is used to simulate the full multi-persona journey: A simulated "Host" creates a pass -> A simulated "Guard" scans it -> A simulated "Admin" verifies the occupancy count updated.
- **Structured Logging:** The Go backend utilizes `slog` to emit JSON-formatted logs. This ensures that errors, state changes, and security events are highly searchable and easily parsed by external observability tools.

### Deployment Pipeline (GitHub Actions -> Fly.io)

1. **Pull Request Stage:** Pushing code to GitHub triggers a GitHub Actions workflow. The workflow spins up isolated containers to run the Go unit tests and the Playwright E2E suite. Code cannot be merged unless coverage and E2E checks pass.
2. **Build Stage:** Upon merging to `main`, GitHub Actions utilizes multi-stage `Dockerfiles` to compile the Go binary and build the Next.js static assets, creating optimized, lightweight production images.
3. **Release Stage:** The Docker images are pushed to a registry, and the Fly.io CLI is triggered. Fly.io orchestrates a zero-downtime rolling deployment, spinning up the new Go 1.26.5 containers at the edge while gracefully draining connections from the old ones. Postgres 18.4 is managed as a highly available cluster within the Fly.io private network.
