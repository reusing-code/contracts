# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Prerequisites

- **Go 1.22+** — https://go.dev/dl/
- **Bun** — JavaScript runtime/package manager: https://bun.sh
- **Task** — Task runner (taskfile.dev): `go install github.com/go-task/task/v3/cmd/task@latest` or `pacman -S go-task`
- **Air** — Live reload for Go: `go install github.com/air-verse/air@latest`
- **Docker + Docker Compose** (optional) — For containerized deployment

## Commands (from repo root)

- `task dev` — Start backend (Air) + frontend (Vite) concurrently
- `task build` — Build Docker image
- `task up` / `task down` — Docker Compose up/down

See `frontend/CLAUDE.md` and `backend/CLAUDE.md` for per-project commands.

## Project structure

- `frontend/` — React 19 + TypeScript SPA (Vite, TanStack Router/Query, shadcn/ui)
- `backend/` — Go API server (stdlib net/http, BadgerDB, slog, Prometheus)
- `Dockerfile` — Multi-stage build: frontend (Bun) → backend (Go) → Alpine runtime
- `docker-compose.yml` — Single-service deployment with named volume for DB

## Dev workflow

`task dev` from root starts both services. Vite dev server (port 5173) proxies `/api/*`, `/healthz`, `/readyz`, `/metrics` to the Go backend (port 8080). Open http://localhost:5173 in the browser.
