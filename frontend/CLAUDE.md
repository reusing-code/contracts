# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

- `bun dev` — Start dev server (Vite, port 5173)
- `bun run build` — Type-check with `tsc -b` then Vite production build
- `bun run lint` — ESLint
- `bun run preview` — Preview production build locally
- `npx shadcn@latest add <component>` — Add a shadcn/ui component (e.g. button, table, dialog)

## Architecture

React 19 + TypeScript SPA built with Vite. No SSR — this is a CRUD/business app behind auth.

**Routing:** TanStack Router with type-safe route definitions in `src/routes/`. The router is created in `router.ts`, root layout in `__root.tsx`. Register the router type via the `Register` interface in `router.ts`.

**Data fetching:** TanStack Query. `QueryClientProvider` wraps the app in `App.tsx`.

**Forms:** React Hook Form + Zod for validation via `@hookform/resolvers`.

**Styling:** Tailwind CSS v4 with the Vite plugin. shadcn/ui for pre-built components (config in `components.json`, components go in `src/components/ui/`). Use the `cn()` helper from `@/lib/utils` for conditional classNames.

**Path alias:** `@/` maps to `src/` (configured in both tsconfig and vite.config).

## Key directories

- `src/routes/` — Route definitions (pages)
- `src/components/` — React components; `ui/` subdirectory for shadcn/ui
- `src/lib/` — Utilities, API client
- `src/hooks/` — Custom React hooks
- `src/types/` — Shared TypeScript types
