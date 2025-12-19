# Feature Specification: BiliBili AccessToken Debug Flows

**Feature Branch**: `004-bilibili-access-token`  
**Created**: 2025-02-14  
**Status**: Draft  
**Input**: User description: "Create specification for BiliBili AccessToken debug flows as described in docs/develop/access-token/bilibili/debug.md"

## Clarifications

### Session 2025-12-18

- Q: Flow 数据（Redis/内存）的默认保留时长是多少？ → A: 24 小时（86,400 秒） TTL，用于覆盖典型调试/支持周期。
- Q: AccessToken 调试服务默认应监听哪个地址？ → A: 默认 127.0.0.1，仅在显式设置 `ACCESSTOKEN_LISTEN_ADDR` 时才对外暴露。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Stand up debug service (Priority: P1)

As an operator, I can boot the AccessToken debug server with the prescribed environment variables so that I can expose the `/debug` console for BiliBili investigations.

**Why this priority**: Without a reliable way to start the service there is no path to the rest of the workflow.

**Independent Test**: Follow the startup commands, observe the log line `accesstoken-server: listening...`, and hit `http://127.0.0.1:7071/debug` to confirm availability.

**Acceptance Scenarios**:

1. **Given** MediaX repo root and required env variables, **When** the operator runs the documented `go run ./cmd/accesstoken/server ...` command, **Then** the server logs the listening message and the browser reaches `/debug` without TLS or auth failures.
2. **Given** Redis variables are unset, **When** the operator starts the server, **Then** the service still boots using in-memory storage while warning that Flow data is volatile.

---

### User Story 2 - Capture & inspect OAuth tokens (Priority: P2)

As a QA engineer, I can walk through the `/debug` UI to authorize a BiliBili account, log Flow records, and view masked `access_token` metadata so I can share validated credentials with other tools.

**Why this priority**: Token capture with clear provenance is the core value of the debug surface.

**Independent Test**: Complete an OAuth cycle from the UI, verify the callback table populates, and call `/accesstoken/token` to inspect the stored token.

**Acceptance Scenarios**:

1. **Given** the debug page with Provider/App/Mode selectors, **When** QA chooses the configured BiliBili template and clicks “发起授权”, **Then** the returned Flow record contains `flow_id`, `expire_at`, `masked_token`, `token_source`, and `oauth_key` fields.
2. **Given** an existing Flow record, **When** QA posts to `/accesstoken/token` without specifying `access_token`, **Then** the API responds with the latest stored token plus TTL and config metadata.

---

### User Story 3 - Retrieve or share Flow records (Priority: P3)

As a support engineer, I can reload expired browser sessions by entering a Flow ID or querying Redis so that I can reconstruct past authorizations without rerunning the full OAuth dance.

**Why this priority**: Flow reuse shortens incident response time and avoids repeated user interaction.

**Independent Test**: Capture a Flow ID, restart the service (or open a new browser), use the “Flow ID 回填” workflow, and confirm `/accesstoken/token` now references the reloaded record.

**Acceptance Scenarios**:

1. **Given** a stored Flow ID, **When** support supplies it through the UI or API, **Then** the system fetches `accesstoken:oauth:flow:<id>` and hydrates the main record without manual JSON edits.
2. **Given** Redis retention is enabled, **When** support queries the key manually (per the documented CLI snippet), **Then** the JSON matches the UI output, ensuring parity for external consumers.

---

### Edge Cases

- When the configured `redirect_url` does not match the actual server address, the OAuth callback must fail gracefully with actionable instructions instead of leaving the UI hanging.
- When Redis is unreachable at startup or during Flow persistence, the system should fall back to in-memory storage while clearly labeling that Flow data will not survive restarts.
- When a Flow ID is invalid or expired, the UI and API should return a structured error explaining whether the Flow key was missing or the record TTL elapsed, avoiding silent failures.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Operators MUST be able to start the AccessToken debug server via documented CLI or Make targets, supplying `ACCESSTOKEN_CONFIG`, `ACCESSTOKEN_LISTEN_ADDR`, and related env vars.
- **FR-002**: The system MUST expose a `/debug` UI that lists Provider/App/Auth-mode combinations sourced from `access_token_providers` for BiliBili entries, including field autofill.
- **FR-003**: The UI and API MUST allow generating OAuth authorize URLs, capturing callbacks, and displaying Flow records (with masked tokens, TTL, oauth_key, config path, and provider metadata).
- **FR-004**: `/accesstoken/token` MUST resolve an access token from one of: request body, env overrides, config defaults, or the latest Flow record, and return token provenance plus TTL.
- **FR-005**: The system MUST offer Flow ID reuse: accept Flow IDs from UI/API, fetch data from memory or Redis (`accesstoken:oauth:*`), and repopulate the debug payloads.
- **FR-006**: Redis integration MUST be optional: when Redis variables are set, Flow records persist using the documented key scheme; when absent, the UI warns that data is volatile but remains functional.
- **FR-007**: The documentation MUST outline curl samples for `/accesstoken/token` and `/accesstoken/oauth/start`, plus troubleshooting steps for common HTTP/redirect/token errors.

### Key Entities

- **Provider Selection**: Represents the chosen provider/app/auth-mode triple plus config path; drives both UI templates and payloads for API calls.
- **Access Token Record**: Stores masked token info, TTL, token source (env/config/flow), oauth_key, and timestamps, ensuring operators can reason about credential freshness.
- **Flow Record**: Captures Flow ID, authorize URL state, callback payload, Redis key references, and expiration, enabling replay and sharing across teams; each record carries一个固定的 24 小时（86,400 秒） TTL，在 Redis 与内存实现中保持一致。
- **Debug Service Binding**: AccessToken debug server 默认监听 `127.0.0.1`（来自 `ACCESSTOKEN_LISTEN_ADDR` 默认值），除非操作者显式设置该环境变量，否则不会暴露到外部网络。

## MediaX Architecture Guardrails *(must reference Constitution sections)*

- **Provider Adapter Parity (Constitution §Provider-Adapters)**: This feature reads BiliBili entries from the shared `access_token_providers` tree; no new adapter code is introduced, but the spec mandates that any future BiliBili adapter additions register via existing `MediaX` factories to keep parity across providers.
- **Config-Layered Security (Constitution §Config-Hardening)**: All secrets (client_id/client_secret/API tokens) remain in env variables referenced by `config.yaml`. Documentation reiterates that YAML files must only contain `${ENV}` placeholders, satisfying layered security expectations.
- **Token Lifecycle Discipline (Constitution §Token-Lifecycle)**: `/accesstoken/token` and Flow reuse flows must continue masking tokens, respect TTL surfaced by `kernel.TokenHandler`, and never bypass the cache-path (Redis or in-memory). Refresh decisions remain centralized in `GetOAuthToken` hooks.
- **Observability & Error Traceability (Constitution §Observability)**: Logs must retain the structured `provider`, `provider_app`, `action`, `flow_id`, and `token_source` fields referenced in the doc so support staff can correlate UI actions with server-side traces while masking sensitive values.
- **Testable Modularity & SessionToken Readiness (Constitution §Modularity)**: Functional slices—startup validation, OAuth capture, Flow replay—are isolated so unit tests can simulate env/Redis combinations. Although SessionToken isn’t directly touched, the spec enforces the same module-level discipline for future convergence.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Operators can start the debug server and reach `/debug` in under 2 minutes following the documented command sequence (measured across 5 trial runs).
- **SC-002**: At least 95% of OAuth authorizations initiated from the UI produce a stored Flow record with populated `flow_id`, `masked_token`, and TTL fields, as observed in logs over a week of use.
- **SC-003**: 99% of Flow ID reuse attempts (UI or API) succeed in retrieving the exact record from Redis or memory, demonstrated via automated or manual regression tests.
- **SC-004**: Support cases tied to “missing access token” errors drop by 50% after the guide ships, indicating the troubleshooting section resolves common misconfigurations.
