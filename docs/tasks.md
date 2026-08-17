# Insta-Gate - Jira/Trello Master Implementation Plan

## Epic 1: Foundation, Infrastructure & CI/CD

**Description:** Setup the monorepo, local development environment, CI/CD pipelines with automated testing gates, and base cloud infrastructure on Fly.io.

- [ x ] **Story 1.1: Monorepo & Taskfile Initialization**

  - [ x ] Task 1.1.1: Initialize Git repository with `/backend`, `/frontend`, `/database`, `/deploy` directories.
  - [ x ] Task 1.1.2: Create `Taskfile.yml` with commands: `dev`, `db:migrate`, `db:rollback`, `lint`, `test:unit`, `test:e2e`, `build`.
  - [ x ] Task 1.1.3: Configure Docker Compose for local Postgres 18.4 and Redis with persistent volume mounting.
  - [ x ] Task 1.1.4: Setup `pre-commit` hooks for formatting (gofmt, prettier) and linting.
  - **Test Coverage:** Add a CI step that verifies `task dev` boots successfully without crash.
  - **AC:** Running `task dev` starts backend, frontend, Postgres, and Redis concurrently. `task lint` passes with zero errors.

- [ ] **Story 1.2: Go 1.26.5 Backend Scaffolding**

  - [ x ] Task 1.2.1: Initialize Go module, setup standard layout (`/cmd`, `/internal`, `/pkg`).
  - [ x ] Task 1.2.2: Implement `log/slog` (JSON formatter) with request ID injection middleware.
  - [ x ] Task 1.2.3: Implement environment variable loading (`godotenv`/`viper`) with validation (fail fast on missing vars).
  - [ x ] Task 1.2.4: Implement standard JSON error response struct (`{"error": {"code": "", "message": ""}}`).
  - **Test Coverage:** Write unit tests for environment loading (missing vars, valid vars) and JSON error formatting.
  - **AC:** Go server boots, connects to local DB/Redis, logs structured JSON, and fails fast if env vars are missing.

- [ ] **Story 1.3: Next.js 16.3 Frontend Scaffolding**

  - [ ] Task 1.3.1: Initialize Next.js 16.3 App Router with TypeScript 7.0 strict mode.
  - [ ] Task 1.3.2: Setup Tailwind CSS, ESLint, Prettier, and absolute path imports (`@/components`, etc.).
  - [ ] Task 1.3.3: Configure Next.js standalone output mode for Docker optimization.
  - **Test Coverage:** Setup Jest + React Testing Library. Add a dummy component test to ensure pipeline runs frontend tests.
  - **AC:** Next.js dev server runs, `npm run test` passes dummy test, `npm run build` succeeds with standalone output.

- [ ] **Story 1.4: GitHub Actions CI/CD & Fly.io Deployment**
  - [ ] Task 1.4.1: Create GitHub Actions PR workflow: lint, unit test, build. Block merge on failure.
  - [ ] Task 1.4.2: Create GitHub Actions `main` workflow: build Docker images, push to GitHub Container Registry.
  - [ ] Task 1.4.3: Create `fly.toml` for backend and frontend with health checks and rolling deployment strategy.
  - [ ] Task 1.4.4: Provision Fly.io Postgres 18.4 cluster and Redis. Configure secrets via Fly CLI.
  - **Test Coverage:** Add a manual dispatch step in GitHub Actions to run a smoke test against the deployed Fly.io `/healthz` endpoint.
  - **AC:** Push to `main` triggers zero-downtime deployment to Fly.io; health checks pass.

## Epic 2: Database Schema, Redis & Persistence Layer

**Description:** Design and implement the PostgreSQL schema, migrations, and Redis connection pooling.

- [ ] **Story 2.1: PostgreSQL Migrations & Connection Pool**

  - [ ] Task 2.1.1: Integrate `golang-migrate` and setup migration file structure.
  - [ ] Task 2.1.2: Implement `pgxpool` connection wrapper with context timeouts.
  - **Test Coverage:** Write integration tests using a Dockerized Postgres test container to verify migrations apply cleanly up and down.
  - **AC:** `task db:migrate` applies schema; `task db:rollback` reverts it without manual DB deletion.

- [ ] **Story 2.2: Identity & RBAC Schema**

  - [ ] Task 2.2.1: Create `tenants` (Super Admin scope) and `buildings` tables.
  - [ ] Task 2.2.2: Create `departments` and `users` tables (Roles: SUPER_ADMIN, FACILITY_ADMIN, HOST, GUARD). Add foreign keys.
  - [ ] Task 2.2.3: Add DB constraints: e.g., unique email per tenant, role validation check.
  - **Test Coverage:** Write DB integration tests attempting to insert invalid roles or missing foreign keys (must fail).
  - **AC:** Schema accepts all defined roles; DB rejects invalid foreign keys or duplicate emails.

- [ ] **Story 2.3: Visitor & Pass Schema**

  - [ ] Task 2.3.1: Create `visitors` table (Name, Phone, Address, ID Photo URL, `is_anonymized` boolean).
  - [ ] Task 2.3.2: Create `passes` table with `status` enum, `valid_from`, `valid_to`, timestamps.
  - [ ] Task 2.3.3: Add indexes on `passes.status` and `passes.building_id`.
  - **Test Coverage:** Write integration tests verifying the `status` enum restricts invalid data.
  - **AC:** Pass state enum restricts invalid data; time fields use UTC `timestamptz`; indexes exist.

- [ ] **Story 2.4: Audit & Telemetry Schema**

  - [ ] Task 2.4.1: Create `audit_logs` table (`timestamp`, `actor_id`, `action`, `target_id`, `ip_address`, `metadata` JSONB).
  - [ ] Task 2.4.2: Implement DB trigger to prevent `UPDATE`/`DELETE` on `audit_logs` table (append-only).
  - **Test Coverage:** Write DB tests attempting to `UPDATE` an audit log; verify SQL exception is thrown.
  - **AC:** Manual or programmatic attempt to delete/update an audit log fails at the database level.

- [ ] **Story 2.5: Redis & Caching Layer**
  - [ ] Task 2.5.1: Implement Go Redis client wrapper.
  - [ ] Task 2.5.2: Implement session storage logic (UserID -> Session Token).
  - [ ] Task 2.5.3: Implement Pass Revocation Blocklist logic (PassID -> TTL).
  - **Test Coverage:** Write unit/integration tests for Redis connection drops, verifying fallback or error handling.
  - **AC:** Backend can read/write sessions and check revocation blocklist in sub-millisecond latency.

## Epic 3: Core Backend Services & API Gateway

**Description:** Build the API routing, middleware, authentication, and strict RBAC enforcement.

- [ ] **Story 3.1: API Gateway & Middleware Stack**

  - [ ] Task 3.1.1: Setup `chi` router (or Go 1.22+ ServeMux). Group routes by `/api/v1/...`.
  - [ ] Task 3.1.2: Implement middleware: Request ID, JSON Logging, Panic Recovery, CORS, Rate Limiting.
  - **Test Coverage:** Unit test middleware chain. Assert Rate Limiter blocks >100 req/sec. Assert Panic Recovery returns 500 JSON.
  - **AC:** All requests are logged with Request ID; rate limiting prevents >100 req/sec per IP; panics are caught.

- [ ] **Story 3.2: Authentication Service (JWT & Sessions)**

  - [ ] Task 3.2.1: Implement Host/Admin login endpoint generating Access (15m) and Refresh (7d) JWTs.
  - [ ] Task 3.2.2: Implement Guard PWA device-login endpoint (PIN/Device ID) generating a scoped JWT.
  - [ ] Task 3.2.3: Store refresh tokens in Redis.
  - **Test Coverage:** Unit tests for JWT generation/validation. Integration tests for login endpoint (valid creds, invalid creds, expired token).
  - **AC:** Tokens are generated, stored in Redis, and invalid credentials return 401. Expired access tokens trigger refresh flow.

- [ ] **Story 3.3: RBAC & Scoping Middleware**
  - [ ] Task 3.3.1: Create `RequireRole` middleware.
  - [ ] Task 3.3.2: Create `ScopeBuilding` middleware injecting `building_id` from JWT into request context.
  - [ ] Task 3.3.3: Create `ScopeHost` middleware ensuring Hosts can only query `passes` where `host_id == their_id`.
  - **Test Coverage:** Write a comprehensive table-driven unit test for RBAC matrix (Super Admin can access X, Host gets 403 on X, etc.).
  - **AC:** Host trying to query another Host's visitors receives 403 Forbidden. Guard trying to access Admin panel gets 403.

## Epic 4: Dynamic QR & Identity Engine

**Description:** The cryptographic core. Generate, sign, route, and validate dynamic QR passes.

- [ ] **Story 4.1: JWT QR Token Generation & Validation**

  - [ ] Task 4.1.1: Implement Go service to generate JWT for QR codes (Payload: `pass_id`, `building_id`, `state`, `exp`).
  - [ ] Task 4.1.2: Implement RSA or HMAC signing logic using secure keys from env vars.
  - [ ] Task 4.1.3: Implement validation endpoint for Guard scans. Check Redis blocklist -> Check DB state -> Return payload.
  - **Test Coverage:** Unit tests for tampered tokens, expired tokens, and tokens with invalid signatures.
  - **AC:** Forged JWTs fail signature validation; expired JWTs return 401. Valid JWTs return pass payload.

- [ ] **Story 4.2: QR Image & Vector Generation**

  - [ ] Task 4.2.1: Integrate Go QR generation library (e.g., `go-qrcode`).
  - [ ] Task 4.2.2: Implement high-contrast SVG vector generation.
  - [ ] Task 4.2.3: Implement Custom Branding logic: overlay transparent PNG logo in center, adjust Error Correction Level to `High`.
  - **Test Coverage:** Unit tests to verify SVG output contains logo data and correct error correction level.
  - **AC:** Generated SVGs scan reliably even with a logo overlay; branding changes update immediately.

- [ ] **Story 4.3: Conditional Routing Engine**

  - [ ] Task 4.3.1: Create `/r/:pass_token` endpoint that reads `User-Agent`.
  - [ ] Task 4.3.2: If Mobile -> redirect to `/pass/:id`. If Desktop -> redirect to `/map/:id`.
  - **Test Coverage:** Unit tests using mocked HTTP headers (iPhone vs Windows Chrome) to verify correct 302 redirect target.
  - **AC:** Scanning a QR code on a phone opens the pass; opening on a laptop opens the map.

- [ ] **Story 4.4: Pass State Machine & Instant Revocation**
  - [ ] Task 4.4.1: Implement Go state machine logic. Prevent invalid transitions (e.g., `SCANNED_OUT` -> `APPROVED`).
  - [ ] Task 4.4.2: Implement `/api/v1/passes/:id/revoke` endpoint.
  - [ ] Task 4.4.3: On revoke, update DB status to `REVOKED` and set PassID in Redis blocklist with TTL.
  - **Test Coverage:** Table-driven unit tests for all valid/invalid state transitions. Integration test for Redis blocklist population on revoke.
  - **AC:** Revoked passes fail validation within milliseconds; state machine rejects illegal transitions.

## Epic 5: Real-Time & Notification Services

**Description:** WebSocket infrastructure for live updates, and integrations for SMS/Email/Push.

- [ ] **Story 5.1: WebSocket Hub Implementation**

  - [ ] Task 5.1.1: Implement Go WebSocket hub managing client connections and rooms (e.g., `building_123`).
  - [ ] Task 5.1.2: Emit events to rooms on state change (`PASS_APPROVED`, `PASS_SCANNED`).
  - [ ] Task 5.1.3: Implement frontend React Hook `useWebSocket` to subscribe to rooms.
  - **Test Coverage:** Backend unit tests for hub room management (joining, leaving, broadcasting). Frontend RTL tests for WebSocket context rendering events.
  - **AC:** When a guard scans a pass, the Admin dashboard updates occupancy instantly without refresh.

- [ ] **Story 5.2: SMS & Email Dispatch Service**

  - [ ] Task 5.2.1: Integrate Twilio (SMS) and SendGrid (Email) SDKs.
  - [ ] Task 5.2.2: Create Go worker pool for async message dispatch to prevent API blocking.
  - [ ] Task 5.2.3: Implement dynamic link sending: `https://svms.app/r/:pass_token`.
  - **Test Coverage:** Mock Twilio/SendGrid APIs. Unit test the worker pool to ensure it processes queued jobs and handles API failures (retries).
  - **AC:** Host creates invite -> Visitor receives SMS with link in < 5 seconds. API doesn't block during dispatch.

- [ ] **Story 5.3: Push Notifications (FCM)**
  - [ ] Task 5.3.1: Integrate Firebase Admin SDK in Go backend.
  - [ ] Task 5.3.2: Create endpoint for Next.js Host app to register FCM device tokens.
  - [ ] Task 5.3.3: Trigger push to Host when a Walk-in visitor selects them.
  - **Test Coverage:** Mock FCM SDK in Go. Unit test that a walk-in pass creation triggers the FCM send function with the correct payload.
  - **AC:** Walk-in submission triggers a push notification on the Host's browser/device.

## Epic 6: Frontend - Host Portal (Next.js)

**Description:** The employee dashboard for pre-approving visitors and managing walk-ins.

- [ ] **Story 6.1: Host Dashboard & Visitor History**

  - [ ] Task 6.1.1: Build Host layout with sidebar (Schedule, Pending, History).
  - [ ] Task 6.1.2: Implement paginated, server-side rendered Visitor History table.
  - [ ] Task 6.1.3: Implement "Revoke" button action on active passes.
  - **Test Coverage:** Jest/RTL tests for table rendering with mock data. Test that clicking "Revoke" triggers the API call and updates UI state.
  - **AC:** Host sees only their invited visitors; revocation updates UI instantly via WebSocket.

- [ ] **Story 6.2: Pre-Approval Flow**

  - [ ] Task 6.2.1: Build "Schedule Visit" form (Name, Phone, Purpose, Time Window, Building selection).
  - [ ] Task 6.2.2: Implement form validation (Zod / Yup) and submission to backend Pass generation endpoint.
  - [ ] Task 6.2.3: Show success state with preview of the branded QR code.
  - **Test Coverage:** RTL form validation tests (required fields, phone format). RTL test for successful submission rendering QR preview.
  - **AC:** Host can schedule a visit; backend triggers SMS/Email to visitor; form validates input correctly.

- [ ] **Story 6.3: Walk-In Approval View**
  - [ ] Task 6.3.1: Build "Pending Approvals" UI component listening to Host's WebSocket channel.
  - [ ] Task 6.3.2: Implement Approve/Deny buttons with 1-click action.
  - **Test Coverage:** Mock WebSocket provider in Jest. Test that receiving a `PENDING_PASS` event renders the card. Test Approve button calls correct endpoint.
  - **AC:** Walk-in request appears instantly; approving updates the visitor's pending phone screen.

## Epic 7: Frontend - Walk-In Self-Registration (Next.js)

**Description:** The mobile-optimized, localized web form for unexpected visitors.

- [ ] **Story 7.1: Dynamic Landing & Localization**

  - [ ] Task 7.1.1: Build `/walk-in/:building_id` route. Implement `Accept-Language` header parsing.
  - [ ] Task 7.1.2: Setup `next-intl` for localized form labels based on parsed language.
  - **Test Coverage:** Unit tests for language parsing logic. RTL tests to verify Spanish vs English labels render based on mock headers.
  - **AC:** Scanning the static building QR routes to the form; language matches device settings.

- [ ] **Story 7.2: Walk-In Form & Host Selection**

  - [ ] Task 7.2.1: Build mobile-first form: Name, Phone, ID Photo upload (camera capture).
  - [ ] Task 7.2.2: Implement Host selection dropdown (searchable, no PII shown to visitor).
  - [ ] Task 7.2.3: Submit form to backend, transition pass to `PENDING`.
  - **Test Coverage:** RTL tests for photo upload simulation. Test host dropdown filtering logic.
  - **AC:** Visitor photo captured, form submitted, backend triggers Host push notification.

- [ ] **Story 7.3: Pending & Approved States**
  - [ ] Task 7.3.1: Build "Waiting for Approval" UI with loading spinner.
  - [ ] Task 7.3.2: Listen to WebSockets for `PASS_APPROVED` or `PASS_DENIED` events.
  - [ ] Task 7.3.3: On approval, render high-contrast dynamic QR code with 1-hour countdown timer.
  - **Test Coverage:** RTL tests mocking WebSocket events. Verify waiting screen transitions to QR display on `PASS_APPROVED`. Test countdown timer logic.
  - **AC:** Screen automatically transitions from waiting to displaying the QR code without manual refresh.

## Epic 8: Frontend - Guard Tablet PWA (Next.js)

**Description:** The highly optimized edge-app for security guards.

- [ ] **Story 8.1: PWA Setup & Scanner Integration**

  - [ ] Task 8.1.1: Configure Next.js PWA manifest, Service Worker for offline shell caching.
  - [ ] Task 8.1.2: Integrate `html5-qrcode` or ZXing JS for continuous camera scanning.
  - [ ] Task 8.1.3: Implement Kiosk Mode / Wake Lock API to prevent tablet screen from sleeping.
  - **Test Coverage:** E2E Playwright test verifying PWA manifest is reachable. Unit test Wake Lock API fallback logic.
  - **AC:** Tablet can be installed to home screen; camera scans QR codes continuously; screen stays awake.

- [ ] **Story 8.2: Sub-50ms Verification Flow**

  - [ ] Task 8.2.1: On scan, fire optimized `fetch` to `/api/v1/verify` with keep-alive.
  - [ ] Task 8.2.2: Implement full-screen UI flash: GREEN checkmark (Valid), RED X (Invalid/Expired/Revoked).
  - [ ] Task 8.2.3: Play distinct audio beeps for success/failure.
  - **Test Coverage:** RTL tests for UI state changes based on API responses (200 OK vs 401 Unauthorized). Test audio playback trigger.
  - **AC:** Scan to visual flash feedback takes < 100ms total round trip.

- [ ] **Story 8.3: Dynamic Data Masking & Check-In/Out**
  - [ ] Task 8.3.1: Render verified visitor data: Photo, Name, Host Name. Explicitly omit Phone/Address in frontend types.
  - [ ] Task 8.3.2: Auto-detect pass state. If `APPROVED`, trigger Check-In. If `SCANNED_IN`, trigger Check-Out.
  - **Test Coverage:** TypeScript strict type checking to ensure PII fields cannot be accessed. RTL test verifying Check-In vs Check-Out API payloads.
  - **AC:** Guard cannot see PII; scanning an already scanned-in pass successfully checks them out.

## Epic 9: Frontend - Admin Dashboard & Analytics (Next.js)

**Description:** Facility and Super Admin global visibility, telemetry, and configuration.

- [ ] **Story 9.1: Live Telemetry & Occupancy Heatmaps**

  - [ ] Task 9.1.1: Build main Admin dashboard querying aggregated DB metrics via WebSocket/API.
  - [ ] Task 9.1.2: Implement real-time counters: Current Visitors, Pending Approvals, Overstay Alerts.
  - [ ] Task 9.1.3: Integrate `Recharts` to render occupancy heatmaps (hour-of-day vs day-of-week).
  - **Test Coverage:** Jest tests for data transformation functions (raw API data -> Recharts format). RTL test for counter rendering.
  - **AC:** Dashboard updates in real-time; heatmap correctly aggregates historical scan data.

- [ ] **Story 9.2: Facility Management & RBAC**

  - [ ] Task 9.2.1: Build CRUD UI for Guard accounts, Departments, and Buildings.
  - [ ] Task 9.2.2: Build UI for Super Admin to manage multiple buildings and global retention policies.
  - [ ] Task 9.2.3: Implement global "Force Checkout" button for overstayed visitors.
  - **Test Coverage:** RTL tests for form validation on creating new guards. Test Force Checkout confirmation modal.
  - **AC:** Admins can provision new guards; Super Admin can switch building context globally.

- [ ] **Story 9.3: Custom Branding Studio**

  - [ ] Task 9.3.1: Build UI to upload company logo (PNG/SVG), select primary/secondary colors.
  - [ ] Task 9.3.2: Implement live preview of the QR code using the Go backend's generation endpoint.
  - **Test Coverage:** RTL test for file upload validation (rejecting non-images). Test that color picker updates state.
  - **AC:** Admin uploads logo; preview QR updates instantly with logo embedded and error correction adjusted.

- [ ] **Story 9.4: Immutable Audit Log Explorer**
  - [ ] Task 9.4.1: Build highly filterable data table (Date Range, Actor, Action, Target).
  - [ ] Task 9.4.2: Implement CSV export functionality using stream generation for large datasets.
  - **Test Coverage:** Unit tests for CSV generation logic (escaping commas, handling nulls). RTL tests for date range filtering.
  - **AC:** Admin can filter for "4th floor server room entries last month" and export to CSV safely.

## Epic 10: Security, Privacy & Background Workers

**Description:** Compliance automation (GDPR/DPDP), overstay alerting, and data masking enforcement.

- [ ] **Story 10.1: Overstay Detection Worker**

  - [ ] Task 10.1.1: Implement Go background ticker running every 5 minutes.
  - [ ] Task 10.1.2: Query DB for `SCANNED_IN` passes where `current_time > valid_to`.
  - [ ] Task 10.1.3: Trigger WebSocket alert to Admin Dashboard and Push Notification to Host.
  - **Test Coverage:** Inject mock time into Go worker. Test that a pass 1 minute overstay triggers the alert, and a pass 1 minute under does not.
  - **AC:** Visitor staying past their window triggers an alert on the dashboard within 5 minutes.

- [ ] **Story 10.2: Automated PII Purging Worker (GDPR/DPDP)**

  - [ ] Task 10.2.1: Implement Go nightly cron job (e.g., 2 AM).
  - [ ] Task 10.2.2: Find `visitors` where `created_at < NOW() - 30 days`.
  - [ ] Task 10.2.3: Nullify Name/Phone, delete ID Photo from storage, set `is_anonymized = true`. Keep `passes` record for stats.
  - **Test Coverage:** DB integration test with mock data inserted 31 days ago and 29 days ago. Verify 31-day data is anonymized, 29-day data is untouched.
  - **AC:** 31-day-old visitor records have PII stripped, but occupancy statistics remain accurate.

- [ ] **Story 10.3: Strict API Data Masking**
  - [ ] Task 10.3.1: Implement Go DTOs (Data Transfer Objects) for Guard endpoints. Exclude Phone, Address, Email.
  - [ ] Task 10.3.2: Add unit tests verifying JSON response payloads do not contain PII keys for Guard roles.
  - **Test Coverage:** Contract tests! Write a Go test that serializes the Guard response struct to JSON bytes and asserts the string does not contain `"phone"` or `"address"`.
  - **AC:** Network packet sniffing on Guard endpoint reveals no PII fields. Contract test passes in CI.

## Epic 11: QA, E2E Testing & Launch Readiness

**Description:** End-to-end multi-persona testing, performance validation, and production hardening.

- [ ] **Story 11.1: Backend Unit & Integration Testing Suite**

  - [ ] Task 11.1.1: Write Go unit tests for JWT generation, state machine transitions, and RBAC matrix.
  - [ ] Task 11.1.2: Write integration tests using a test Postgres container for pass revocation and Redis blocklist.
  - [ ] Task 11.1.3: Setup Go coverage reporting (`go test -cover`) in GitHub Actions. Gate PRs at >80% coverage.
  - **Test Coverage:** Meta-task to ensure all previous test tasks are unified and passing.
  - **AC:** Code coverage > 80% for `/internal` packages. CI blocks PR on failure.

- [ ] **Story 11.2: Playwright E2E Multi-Persona Suites**

  - [ ] Task 11.2.1: Setup Playwright with multiple contexts (Host browser, Walk-in mobile browser, Guard tablet).
  - [ ] Task 11.2.2: E2E Test - Use Case A: Host pre-approves -> Guard scans in -> Admin sees occupancy +1.
  - [ ] Task 11.2.3: E2E Test - Use Case B: Walk-in registers -> Host approves via portal -> Walk-in screen shows QR -> Guard scans.
  - [ ] Task 11.2.4: E2E Test - Revocation: Host revokes -> Guard scans revoked pass -> Tablet flashes RED.
  - **Test Coverage:** This story _is_ the test coverage for the entire user journey.
  - **AC:** All E2E suites pass in GitHub Actions CI on every PR.

- [ ] **Story 11.3: Load Testing & Performance Optimization**

  - [ ] Task 11.3.1: Use `k6` or `vegeta` to simulate 200 concurrent Guard scans to `/api/v1/verify`.
  - [ ] Task 11.3.2: Optimize DB queries / add missing indexes if p99 latency > 50ms.
  - [ ] Task 11.3.3: Verify Redis blocklist scales correctly under concurrent read load.
  - **Test Coverage:** K6 script committed to repo. Run K6 against staging environment.
  - **AC:** 95th percentile latency for QR verification remains under 50ms under heavy load.

- [ ] **Story 11.4: Observability & Structured Logging Validation**
  - [ ] Task 11.4.1: Ensure all `slog` logs output valid JSON with `request_id` traceability.
  - [ ] Task 11.4.2: Verify error logs capture stack traces and user context without leaking PII.
  - [ ] Task 11.4.3: Setup basic health check endpoints (`/healthz`, `/readyz`) for Fly.io load balancer.
  - **Test Coverage:** Automated script that greps local docker logs for `request_id` and validates JSON structure.
  - **AC:** Logs are parsable by standard JSON log aggregators; Fly.io correctly routes traffic only to healthy instances. Health endpoints return 200 OK.
