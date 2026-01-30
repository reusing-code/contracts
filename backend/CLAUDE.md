# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Prerequisites

- **Go 1.22+** — https://go.dev/dl/
- **Task** — Task runner (taskfile.dev): `go install github.com/go-task/task/v3/cmd/task@latest` or `pacman -S go-task`
- **Air** — Live reload for Go: `go install github.com/air-verse/air@latest`
- **golangci-lint** (optional) — Linter: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` or `pacman -S golangci-lint`

## Commands

- `task dev` — Start with Air (live reload), sets `CORS_ORIGIN` for frontend dev server
- `task build` — Build the binary
- `task test` — Run tests
- `task lint` — Run golangci-lint

## Architecture

Go 1.22+ stdlib `net/http` with method+pattern routing. Single binary that optionally serves the frontend SPA.

**Storage:** BadgerDB (embedded LSM-tree KV store). Data stored as JSON documents with key prefixes for multi-user namespacing: `u/{userId}/cat/{id}`, `u/{userId}/con/{id}`, `u/{userId}/idx/cat_con/{catId}/{conId}`.

**Config:** Environment variables via `caarlos0/env` struct tags. See `.env.example` for all options.

**Logging:** `log/slog` — JSON handler in production, text in dev.

**Metrics:** Prometheus via `/metrics` endpoint. Tracks `http_requests_total`, `http_request_duration_seconds`, `http_active_requests`.

**Middleware chain (outermost first):** RequestID → Recovery → Metrics → Logging → CORS → handler.

**Auth:** Not yet implemented. All requests use a hardcoded default user ID. The store interface accepts `userID` as a parameter so handlers can pass it from auth context later without store changes.

## Key directories

- `cmd/server/` — Entrypoint
- `internal/config/` — Environment-based configuration
- `internal/model/` — Category and Contract types (JSON tags match frontend Zod schemas)
- `internal/store/` — Store interface + BadgerDB implementation
- `internal/handler/` — HTTP handlers (health, category CRUD, contract CRUD)
- `internal/middleware/` — Request ID, recovery, metrics, logging, CORS
- `internal/server/` — Mux setup, middleware wiring, graceful shutdown, SPA serving

## API

All endpoints under `/api/v1/`. JSON request/response with camelCase field names.

- `GET|POST /api/v1/categories` — List/create categories
- `GET|PUT|DELETE /api/v1/categories/{id}` — Get/update/delete category (delete cascades to contracts)
- `GET|POST /api/v1/categories/{id}/contracts` — List/create contracts in category
- `GET|PUT|DELETE /api/v1/contracts/{id}` — Get/update/delete contract
- `GET /healthz` — Liveness probe
- `GET /readyz` — Readiness probe (checks DB)
- `GET /metrics` — Prometheus metrics
