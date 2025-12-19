# Tasks: BiliBili AccessToken Debug Flows

**Input**: Design documents from `/specs/004-bilibili-access-token/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

## Phase 1: Setup (Shared Infrastructure)

- [X] T001 Sync Go module dependencies referenced in `go.mod`/`go.sum` via `go mod download`
- [X] T002 Update environment defaults and startup instructions in `docs/develop/access-token/bilibili/debug.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

- [X] T003 Implement config loader defaults for `ACCESSTOKEN_LISTEN_ADDR` and TTL in `internal/accesstoken/config/config.go`
- [X] T004 Add startup validation + structured log fields (listen addr, storage backend) in `cmd/accesstoken/server/main.go`
- [X] T005 [P] Build Redis client factory with memory fallback in `internal/accesstoken/flow/store.go`
- [X] T006 [P] Create store unit tests covering Redis/memory modes in `internal/accesstoken/flow/store_test.go`
- [X] T007 Wire shared logging helpers that inject `provider`, `provider_app`, `flow_id`, `token_source` in `internal/accesstoken/handler/middleware.go`
- [X] T008 [P] Reference new quickstart (`specs/004-bilibili-access-token/quickstart.md`) from `README.md`

---

## Phase 3: User Story 1 - Stand up debug service (Priority: P1) 🎯 MVP

**Goal**: Operators can boot the server with documented env vars, hit `/debug`, and see warnings when Redis is absent.
**Independent Test**: Run `ACCESSTOKEN_LISTEN_ADDR=127.0.0.1:7071 go run ./cmd/accesstoken/server -config config.yaml`, observe log line with listen addr/storage backend, then load `http://127.0.0.1:7071/debug` without auth failures.

### Implementation
- [X] T009 [US1] Split listen addr/port CLI flags & bind defaults inside `cmd/accesstoken/server/main.go`
- [X] T010 [P] [US1] Render API-token dialog + volatility banner in `internal/accesstoken/ui/templates/debug.html`
- [X] T011 [US1] Secure `/debug` handler with Bearer token + static asset serving in `internal/accesstoken/handler/debug.go`
- [X] T012 [US1] Document Redis-optional startup + warning screenshots in `docs/develop/access-token/bilibili/debug.md`
- [X] T013 [US1] Add env validation tests (listen addr, API token) in `cmd/accesstoken/server/main_test.go`

---

## Phase 4: User Story 2 - Capture & inspect OAuth tokens (Priority: P2)

**Goal**: QA can authorize via UI, persist Flow records (24h TTL), and inspect masked tokens through `/accesstoken/token`.
**Independent Test**: Complete OAuth using `/debug`, confirm Flow table entry with `flow_id/expire_at`, then call `/accesstoken/token` (no token in body) and verify response fields + TTL.

### Implementation
- [X] T014 [US2] Implement Flow creation + TTL enforcement in `internal/accesstoken/flow/store.go`
- [X] T015 [P] [US2] Build `/accesstoken/oauth/start` handler per contract spec in `internal/accesstoken/handler/oauth.go`
- [X] T016 [P] [US2] Handle `/debug/callback` payload masking + Flow persistence in `internal/accesstoken/handler/callback.go`
- [X] T017 [US2] Implement `/accesstoken/token` resolver precedence (request/env/config/Flow) in `internal/accesstoken/handler/token.go`
- [X] T018 [P] [US2] Add contract tests for `/accesstoken/token` scenarios in `tests/contract/accesstoken_token_test.go`
- [X] T019 [US2] Emit structured OAuth logs (provider/app/flow/token_source/listen_addr) inside `internal/accesstoken/handler/oauth.go`

---

## Phase 5: User Story 3 - Retrieve or share Flow records (Priority: P3)

**Goal**: Support engineers reload Flow IDs via UI/API (memory or Redis) and understand TTL/expiry reasons.
**Independent Test**: Restart server, call `/accesstoken/flow/replay` with stored Flow ID, receive `status=authorized` + storage backend, then consume `/accesstoken/token` to confirm hydration.

### Implementation
- [X] T020 [US3] Implement `/accesstoken/flow/replay` handler returning status + payload in `internal/accesstoken/handler/flow.go`
- [X] T021 [P] [US3] Expose Flow listing endpoint `/accesstoken/flows` with pagination + masked account fields in `internal/accesstoken/handler/flow.go`
- [X] T022 [US3] Add Flow ID 回填 controls + API wiring in `internal/accesstoken/ui/templates/debug.html`
- [X] T023 [US3] Extend `internal/accesstoken/flow/store.go` to fetch by Flow ID and surface expired/missing errors
- [X] T024 [P] [US3] Write contract tests for `/accesstoken/flow/replay` and `/accesstoken/flows` in `tests/contract/accesstoken_flow_test.go`

---

## Phase 6: Polish & Cross-Cutting

- [X] T025 Document observability/log fields + troubleshooting table updates in `docs/develop/access-token/bilibili/debug.md`
- [X] T026 [P] Refactor masked token helpers for reuse in `internal/accesstoken/handler/mask.go`
- [X] T027 Align quickstart command samples with final behavior in `specs/004-bilibili-access-token/quickstart.md`
- [X] T028 Capture verification evidence (curl outputs + screenshots) inside `specs/004-bilibili-access-token/checklists/verification.md`

---

## Dependencies & Execution Order

1. Phase 1 → Phase 2 (blocking for all stories)
2. After Phase 2, User Stories can proceed in priority order (US1 → US2 → US3) or in parallel if dependencies satisfied
3. Phase 6 (Polish) executes once targeted stories are complete

---

## Parallel Opportunities

- T005/T006/T008 can run concurrently once config defaults are defined
- In US1, UI work (T010) can proceed parallel to handler updates (T011) after config wiring (T009)
- In US2, Flow store (T014) must finish before handlers consume it, but tests (T018) can stub against contracts earlier
- In US3, UI integration (T022) can run once API contracts (T020/T021) are outlined

---

## Independent Tests per Story

- **US1**: Manual startup hitting `/debug` + log verification
- **US2**: OAuth round-trip verifying Flow table + `/accesstoken/token` response
- **US3**: Flow replay after restart via `/accesstoken/flow/replay` and `/accesstoken/token`

---

## Suggested MVP Scope

- Deliver through Phase 3 (User Story 1) for MVP; subsequent stories layer additional capabilities independently.
