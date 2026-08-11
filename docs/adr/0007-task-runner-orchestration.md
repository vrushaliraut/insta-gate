# ADR 0007: Choice of Development Task Orchestrator

* **Status:** Accepted
* **Date:** 2026-08-11

## Context

A monorepo containing Go, Next.js, Docker, and database tools requires a cross-platform task runner to unify developer
commands (`dev`, `lint`, `db:migrate`, `build`).

## Decision

We select **`Task` (`Taskfile.yml`)**.

## Alternatives Considered & Rejection Reasons

1. **GNU Make / Makefile:**
    * *Rejected:* Make suffers from cross-platform syntax incompatibilities between macOS (BSD Make) and Linux/Windows (
      GNU Make), relies on tab-formatting gotchas, and handles environment variable passing poorly across multi-folder
      monorepos.
2. **npm Scripts:**
    * *Rejected:* Binds repository-wide orchestrator logic exclusively to the Node.js ecosystem, forcing backend-only
      workflows to depend on `package.json` configurations.

## Consequences

* **Positive:**
    * Cross-platform consistency across Windows, macOS, and Linux.
    * Readable YAML syntax with built-in environment variable handling and task dependency chaining.
* **Negative:**
    * Developers must install the `task` CLI binary on their local workstations.