# Contracts — Frontend

React SPA for the Contracts application.

## Tech Stack

- **React 19** + **TypeScript** (Vite)
- **Tailwind CSS v4** + **shadcn/ui** for components
- **TanStack Router** for type-safe routing
- **TanStack Query** for server state / data fetching
- **React Hook Form** + **Zod** for forms and validation

## Getting Started

```bash
bun install
bun dev
```

Dev server runs at `http://localhost:5173`.

## Scripts

| Command | Description |
|---------|-------------|
| `bun dev` | Start dev server with HMR |
| `bun run build` | Type-check and build for production |
| `bun run lint` | Run ESLint |
| `bun run preview` | Preview production build |

## Adding UI Components

This project uses [shadcn/ui](https://ui.shadcn.com). Components are not installed as dependencies — they're copied into `src/components/ui/` and can be freely customized.

```bash
npx shadcn@latest add button
npx shadcn@latest add table
```

Browse available components at [ui.shadcn.com](https://ui.shadcn.com).

## Project Structure

```
src/
  routes/        Route definitions (pages)
  components/    React components (ui/ for shadcn)
  lib/           Utilities, API client
  hooks/         Custom React hooks
  types/         Shared TypeScript types
```

The `@/` import alias maps to `src/`.
