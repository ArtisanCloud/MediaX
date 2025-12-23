# Feature Specification: DouYin ClientToken Mode Onboarding

**Feature Branch**: `[007-douyin-client-token]`  
**Created**: 2025-12-22  
**Status**: Draft  
**Input**: `docs/plan/byteDance/client_token_client.md`, `docs/develop/client-token/byteDance/{develop,debug}.md`

> **Design Principle**：将抖音 DouYin ClientToken 能力与 WeChat ClientToken 相同体验接入 `cmd/clienttoken/server`、CLI 及三方插件，配置全部来自 `client_token_providers` 与环境变量，缓存遵循统一 Redis/内存规范，文档调试流程对齐 `docs/develop/client-token` 约定。

## Clarifications

### Session 2025-12-22

- Q: DouYin ClientToken 自动刷新应如何触发？ → A: 在任意读取/调用 client_token 时检测 TTL，若低于 `refresh_before_seconds` 立即刷新

## User Scenarios & Testing *(mandatory)*

### User Story 1 - `/client-token/debug` 刷新 DouYin ClientToken (Priority: P1)

作为运维/开发者，我可以在 `cmd/clienttoken/server` 启动的 `/client-token/debug` 页面选择 “字节跳动 / 抖音服务端 / default” 模式，填充配置模板并一键刷新 DouYin ClientToken，缓存写入 Redis（或内存降级），并立即用于 API 调试。

**Why this priority**: 没有可视化刷新与调试入口，就无法在 MediaX 中完成 ClientToken 入驻，也无法向 CLI/插件暴露统一能力。

**Independent Test**: 启动 `cmd/clienttoken/server` → 访问 `/client-token/debug?api_token=...` → 同步模板 → 刷新 Token → 查看日志/页面确认 `provider_code=byte_dance_douyin_clienttoken` 成功写入缓存。

**Acceptance Scenarios**:

1. **Given** `config.yaml` 按计划新增 `byte_dance/douyin_service` 配置，**When** 在 `/client-token/debug` 选择对应 Provider 并点击“刷新 Token”，**Then** 成功获得 DouYin client_token JSON，页面展示 TTL、来源与 Redis key。
2. **Given** 用户在调试台中选择 API 模板（如 `content/video/list`），**When** 发起 API 调试，**Then** 调用自动附带最新 client_token，并在日志中记录 `event=clienttoken.call provider_code=byte_dance_douyin_clienttoken`。

---

### User Story 2 - Redis 缓存与 CLI/脚本复用 (Priority: P2)

作为需要批量刷新 ClientToken 的工程师，我希望缓存策略、Redis Key（`clientToken:douyin:<client_key>`）、TTL/`refresh_before_seconds` 被清晰定义，并可通过 CLI 或 `scripts/clienttoken-douyin.sh` 复用调试台同一配置，避免重复登录。

**Why this priority**: 客户端任务通常在无人值守环境运行，需要稳定的缓存命名与 TTL，以便监控和自动刷新。

**Independent Test**: 设置 `CLIENTTOKEN_REDIS_*` → 调用 CLI/脚本刷新 client_token → 使用 `redis-cli` 检查 key/TTL → CLI 再次调用复用缓存。

**Acceptance Scenarios**:

1. **Given** Redis 可用且配置 `refresh_before_seconds=600`，**When** token 剩余 TTL < 600 秒时触发刷新，**Then** 新 token 写入同一 Redis Key，调试台“查看缓存”显示 `source=redis refreshed=true`。
2. **Given** CLI 或脚本执行 `./scripts/clienttoken-douyin.sh refresh`, **When** 提供必需环境变量，**Then** 命令返回成功并在日志中输出与调试台一致的结构化字段。

---

### User Story 3 - 文档与配置可复制 (Priority: P3)

作为 SDK 集成商/合作伙伴，我仅阅读 `docs/develop/client-token/byteDance/{develop,debug}.md` 与 `config.example.yaml`，即可在 1 小时内完成环境变量、Provider 配置、服务启动、页面操作与常见问题排查，体验与微信 ClientToken 文档一致。

**Why this priority**: 文档是外部团队唯一入口，必须保证自助可复现，减少支持成本。

**Independent Test**: 由未参与开发的成员按照文档步骤，从复制配置到完成 API 调试全程无额外指导即可通过。

**Acceptance Scenarios**:

1. **Given** 用户复制 `client_token_providers` 示例到自己的 `config.yaml`，**When** 依照文档设置所有 `DOUYIN_*` 及 `CLIENTTOKEN_*` 环境变量，**Then** 能成功启动服务并通过 `/client-token/debug` 刷新 Token。
2. **Given** 文档列出 `provider code not found`、Redis 降级等问题，**When** 用户遇到这些问题，**Then** 可按照排障表解决且无需额外支持。

### Edge Cases

- 缺少 `DOUYIN_CLIENT_KEY/SECRET` 或 `CLIENTTOKEN_API_TOKEN` 时必须在启动与刷新流程中明确报错，并阻止调试台展示 DouYin Provider，避免误操作。
- Redis 不可用时需自动降级为进程内缓存，同时在 UI/日志展示 `storage_backend=memory`，并提示仅适合单实例调试。
- 若 DouYin API 额外要求 `device_id`、`risk_info` 等参数，缺失时需在配置/文档里标记并默认禁用相关 API 模板。
- CLI/脚本调用时认证失败（如 Bearer token 不符）需返回 401 并记录 `api_token` 校验日志，确保脚本/调试台共用安全策略。
- 当缓存命中但 TTL 小于 `refresh_before_seconds`，必须阻塞当前调用直至刷新成功，否则返回“token refresh failed”并在日志中指出阈值与剩余 TTL，防止使用将过期 token。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `client_token_providers` 必须新增 `provider_code=byte_dance_douyin_clienttoken` 结构，支持 `apps[].auth_modes[].byte_dance_douyin_config` 字段，并全部通过环境变量注入 `client_key/client_secret/api_token/redis`。
- **FR-002**: `cmd/clienttoken/server` 在启动时需加载上述 Provider，并在 `/client-token/debug` 下拉列出 Provider→App→Mode；若缺乏配置则隐藏此项并输出清晰报错。
- **FR-003**: 刷新 client_token 成功后必须写入缓存（Redis 优先、内存备选），key 命名固定为 `clientToken:douyin:<client_key>`，TTL 取 `ttl_seconds`，并在调试台/CLI/API 读取 token 时检测剩余 TTL，若低于 `refresh_before_seconds` 即刻刷新并回写缓存。
- **FR-004**: 调试台和 HTTP API (`/client-token/call`) 必须自动带上最新 token，记录 `provider_code`, `provider_app`, `action`, `token_source`, `ttl_remaining` 字段，并使用 `internal/accesstoken/handler/mask` 脱敏输出。
- **FR-005**: `docs/develop/client-token/byteDance/develop.md` 与 `debug.md` 需更新为 DouYin 版本，结构/章节顺序与微信 ClientToken 文档一致，覆盖环境变量、启动命令、页面操作与常见问题。
- **FR-006**: `config.example.yaml` 与 `config.yaml` 模板必须包含完整 DouYin Provider 节点，附带注释说明环境变量及默认值，确保复制即可运行。
- **FR-007**: 新增或扩展 CLI/脚本（例如 `scripts/clienttoken-douyin.sh`）以支持刷新、查看缓存、调用 API 三个指令；脚本需读取同一配置并复用 `CLIENTTOKEN_API_TOKEN` 校验。
- **FR-008**: 日志需遵循 `clienttoken` 命名规范，包含 `event`、`provider`、`app_code`、`redis_key`、`refresh_before_seconds`，并在 Redis 连接失败时输出 fallback 说明供排障。
- **FR-009**: 若份外字段（`device_id`、`risk_info` 等）被 DouYin API 要求，必须在配置/文档/模板中声明，默认关闭相应 API 模板，避免用户误用；若字段被启用需验证后才允许保存。
- **FR-010**: `/client-token/debug` 的“查看缓存”与 CLI 查询需展示 TTL、来源（`redis`/`memory`）、最近刷新时间，并允许一键清除缓存以便复测。

### Key Entities *(include if feature involves data)*

- **DouYinClientTokenProvider**：描述 Provider/App/Mode 组合，包含 `provider_code`、`app_code`（`douyin_service`）、`auth_mode`（`default`）、`byte_dance_douyin_config`（api_url、timeout、http_debug、client_token、cache、redis、api_token）。用于调试台枚举。
- **ClientTokenCacheRecord**：对应单个 DouYin client_token 的缓存结构，含 `redis_key`、`token_value`（脱敏）、`ttl_seconds`、`refresh_before_seconds`、`source`（redis/memory）、`updated_at`、`provider_code`，供调试台和 CLI 显示。
- **ClientTokenDebugAction**：用户在调试台或 CLI 触发的操作（刷新、查看缓存、调用 API），记录 `action_name`、`request_payload`、`api_endpoint`、`result_status`、`error_hint`，用于调试日志和追踪。

## MediaX Architecture Guardrails *(must reference Constitution sections)*

- **Provider Adapter Parity (§Provider Adapter Parity)**：仅能复用 `pkg/client/byteDance/douYin/clientTokenClient` 及 `MediaXCore` 已有的 BaseClient/Logger，禁止新增平行 HTTP 客户端；`cmd/clienttoken/server` 必须通过既有 factory 注入该 Provider，保证与其它 ClientToken 客户端一致。
- **Config-Layered Security (§Config-Layered Security)**：新增字段全部放在 `pkg/client/config/douyin.go` 的 `ByteDanceDouYinConfig` 与 `ClientTokenConfig` 下，并在 YAML 示例中使用 `${ENV}` 占位；`client_key/client_secret/api_token/redis` 只能来自环境变量或 secret store，严禁写入源码。
- **Token Lifecycle Discipline (§Token Lifecycle Discipline)**：刷新逻辑需复用 `core.ByteDanceTokenHandler` 与既有缓存接口；Redis key 固定 `clientToken:douyin:*`，TTL/`refresh_before` 由配置控制，自动续期要通过统一的缓存层，并确保 `/client-token/debug`、CLI、三方插件共享缓存。
- **Observability & Error Traceability (§Observability & Error Traceability)**：日志与调试 UI 必须输出 `provider`, `provider_code`, `action`, `redis_key`, `ttl`, `storage_backend`, `api_token_subject` 等字段；调用 DouYin API 时必须使用 `mask` 包处理 `client_token`，并记录 `retry_count`、错误码、HTTP 状态用于排障。
- **Testable Modularity & SessionToken Readiness (§Testable Modularity & SessionToken Readiness)**：`pkg/client/byteDance/douYin/clientTokenClient` 新增方法需具备单元测试；`cmd/clienttoken` 需添加集成测试覆盖刷新→缓存→调用→清除的完整流程；虽不涉及 SessionToken 流程，也需验证不会影响现有 SessionToken/AccessToken 模块的配置加载。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 使用新文档与配置，5 次独立演练中至少 4 次可在 10 分钟内从启动 `cmd/clienttoken/server` 到完成 DouYin ClientToken 刷新与 API 调试（≥80% 成功率）。
- **SC-002**: 设置 Redis 后，连续 24 小时内自动刷新不少于 12 次，日志中无 `storage_backend=fallback`，`clientToken:douyin:<client_key>` TTL 始终大于配置的 `refresh_before_seconds`。
- **SC-003**: CLI/脚本在 3 种操作（刷新/查看/调用）上通过回归测试，响应时间 p95 ≤ 1.5s，且日志无明文凭证；测试需在 Redis 与纯内存两种模式各执行一次。
- **SC-004**: 两名未参与开发的同事仅凭 `docs/develop/client-token/byteDance/{develop,debug}.md` 即可复制配置并完成整套流程，问卷反馈“无需额外帮助”，若失败需补充文档后复测。
