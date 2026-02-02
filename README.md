# Contracts

A self-hosted contract management application for tracking subscriptions, renewals, and recurring expenses. Built with a Go backend and React frontend.

## Features

- **Contract tracking** — Store contract details including dates, pricing, notice periods, and renewal terms
- **Category organization** — Group contracts by category (insurance, telecom, etc.)
- **Renewal monitoring** — Dashboard showing upcoming renewals with color-coded urgency indicators
- **Email reminders** — Configurable SMTP-based reminder emails for approaching renewals
- **Batch import** — Import contracts from JSON via file upload or paste
- **Multi-user** — JWT authentication with per-user data isolation
- **Observability** — Prometheus metrics, structured logging, health/readiness probes

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | React 19, TypeScript, Vite, TanStack Router/Query, shadcn/ui, Tailwind CSS |
| Backend | Go (stdlib `net/http`), BadgerDB, JWT, slog |
| Runtime | Docker, Alpine Linux |

## Quick Start

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Bun](https://bun.sh)
- [Task](https://taskfile.dev) — `go install github.com/go-task/task/v3/cmd/task@latest`
- [Air](https://github.com/air-verse/air) — `go install github.com/air-verse/air@latest`

### Development

```sh
task dev
```

This starts the Go backend (port 8080) with live reload and the Vite dev server (port 5173) with API proxying. Open http://localhost:5173.

### Docker

```sh
task build   # build image
task up       # start with docker compose
task down     # stop
```

## Project Structure

```
frontend/          React SPA (Vite, TanStack, shadcn/ui)
backend/           Go API server (net/http, BadgerDB)
Dockerfile         Multi-stage build (Bun → Go → Alpine)
docker-compose.yml Single-service deployment with named volume
```

See `frontend/CLAUDE.md` and `backend/CLAUDE.md` for per-project details.

## API

All endpoints under `/api/v1/`. Auth endpoints are public; everything else requires a JWT bearer token.

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/register` | Register user |
| POST | `/auth/login` | Login (returns JWT) |
| GET/POST | `/categories` | List / create categories |
| GET/PUT/DELETE | `/categories/{id}` | Category CRUD |
| GET/POST | `/categories/{id}/contracts` | Contracts in category |
| GET | `/contracts` | List all contracts |
| GET/PUT/DELETE | `/contracts/{id}` | Contract CRUD |
| GET | `/contracts/upcoming-renewals` | Renewals by date |
| POST | `/contracts/import` | Batch JSON import |
| GET/PUT | `/settings` | Renewal preferences |
| PUT | `/settings/password` | Change password |
| GET | `/summary` | Dashboard stats |

Health (`/healthz`), readiness (`/readyz`), and Prometheus metrics (`/metrics`) are available at the root.

## License

[MIT](LICENSE)
