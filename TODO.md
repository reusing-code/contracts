# TODO

## Frontend: migrate from localStorage to API

The frontend repositories (`lib/category-repository.ts`, `lib/contract-repository.ts`) currently read/write through `lib/storage.ts` which uses localStorage. The repository functions are already async, so migration is straightforward: replace the implementations with `fetch()` calls to `/api/v1/...` endpoints.

- [ ] Create an API client helper (`lib/api.ts`) with typed fetch wrappers and error handling
- [ ] Rewrite `category-repository.ts` to call the backend API
- [ ] Rewrite `contract-repository.ts` to call the backend API
- [ ] Remove `storage.ts` (localStorage layer)
- [ ] Update React Query hooks if response shapes change (they shouldn't — JSON fields match the Zod schemas)

## Backend: tests

- [ ] Store tests — CRUD operations, cascade delete, index consistency
- [ ] Handler tests — request validation, error responses, status codes
- [ ] Integration test — start server, run requests against real endpoints

## Authentication

Storage layer is multi-user ready (all store methods accept `userID`). Handlers currently pass a hardcoded default user ID.

- [ ] Choose auth strategy (session cookies, JWT, OAuth/OIDC)
- [ ] Add auth middleware that extracts user ID from request context
- [ ] Replace `defaultUserID` in handlers with user ID from context
- [ ] Add user registration/login endpoints or external provider integration

## Business logic

- [ ] Contract cancellation date calculation (based on notice period, minimum duration, extension)
- [ ] Renewal tracking / upcoming renewal notifications
- [ ] Cost summary / dashboard aggregations on the backend
