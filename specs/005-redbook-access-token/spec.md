# Feature Specification: RedBook JuGuang AccessToken Debug Integration

**Feature Branch**: `[005-redbook-access-token]`  
**Created**: 2025-12-19  
**Status**: Draft  
**Input**: User description: "基于 docs/plan/redBook/access_token_client.md，将小红书聚光 AccessToken 客户端像 Google/Bilibili 一样接入现有 cmd/accesstoken 本地调试环境与文档体系。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 调试台支持聚光 OAuth (Priority: P1)

作为调试 cmd/accesstoken 的工程师，我可以在 `/debug` 页面选择“小红书 / 聚光” Provider，发起 OAuth 授权并在回调页看到 Flow 记录，从而像 Google/Bilibili 一样复用 Redis Flow ID 进行 CLI 或 API 调试。

**Why this priority**: 无法授权就无法使用任何聚光接口，因此这是整个接入最核心的价值。

**Independent Test**: 启动 `go run ./cmd/accesstoken/server -config config.yaml`，在 `/debug` 选择聚光 Provider → 点击“发起授权” → 浏览器跳转到聚光授权页 → 回调后在“授权记录”中看到新 Flow。

**Acceptance Scenarios**:

1. **Given** config.yaml 已配置 `redbook_juguang_config`，**When** `/debug` 选择聚光 Provider 并点击“发起授权”，**Then** 服务返回有效授权 URL 并显示 Flow ID/过期时间。
2. **Given** 聚光 OAuth 回调到 `/debug/callback`，**When** 服务用 code 换取 token，**Then** Flow 存入内存和 Redis，可在“授权记录”列表中查看并复制到 CLI。

---

### User Story 2 - 聚光 API 可通过 cmd/accesstoken 调用 (Priority: P2)

作为后端调试者，我希望在调试台或带 API token 的 HTTP 请求中调用至少一个聚光 REST API（如账户余额），系统会复用聚光客户端并脱敏日志，以验证 AccessToken 生命周期与真实业务能力。

**Why this priority**: 只有完成 API 调用才能证明授权 token 可用，支撑对外说明“聚光已入驻调试环境”。

**Independent Test**: 使用 `curl -H "Authorization: Bearer dev-accesstoken" -d '{"provider_code":"redbook_juguang", "action":"redbook.account.balance", ...}' http://127.0.0.1:7071/accesstoken/call`，观察成功 JSON 与日志。

**Acceptance Scenarios**:

1. **Given** 聚光 Flow 已在 Redis 中，**When** 通过 /accesstoken/call 发起 `action=redbook.account.balance` 请求，**Then** 服务从 Flow 读取 token、调用聚光客户端并返回余额 JSON，同时日志包含 `provider=redbook_juguang` 且 token 已脱敏。

---

### User Story 3 - 文档指导聚光入驻 (Priority: P3)

作为 SDK/运营伙伴，我可以查阅 `docs/develop/access-token/redbook/{develop,debug}.md`，复制配置、了解 OAuth/Flow 操作、调试 API 示例，确保与 Google/Bilibili 文档体验一致。

**Why this priority**: 配套文档是让外部团队能够自助调试的唯一途径，可以减少支持成本。

**Independent Test**: 仅阅读新文档即可完成从配置、环境变量到调试页面/CLI 的整个流程。

**Acceptance Scenarios**:

1. **Given** 新成员仅阅读文档，**When** 按步骤配置环境变量、运行调试台并调用样例 API，**Then** 全流程可独立完成且文档示例与 UI/接口保持一致。

### Edge Cases

- 当 `redbook_juguang_config` 缺少 `scope` 或 `oauth_url` 时，`/accesstoken/oauth/start` 必须返回解释性错误并在 `/debug` 显示提示。
- 若 Redis 未连接，Flow 需降级至内存并在日志中提示 `storage_backend=memory`，用户仍可调试但需重新授权。
- 如果 config.yaml 中未启用聚光 Provider，则 `/debug` 不应显示该入口，避免空模式引发混淆。
- Flow TTL 到期或被清理后，Flow ID 回填要返回 `flow_not_found` 且文档需说明重新授权步骤。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `/debug` 页面必须在读取 config.yaml 后展示“小红书 / 聚光” Provider、默认 App 与授权模式，并支持发起 `/accesstoken/oauth/start` 请求。
- **FR-002**: `cmd/accesstoken/server` 解析 `redbook_juguang_config` 时，缺少 `client_id`/`scope` 等字段要给出明确错误，并保持 Google/Bilibili Provider 不回归。
- **FR-003**: OAuth 回调完成后需将聚光 Flow（含 `provider_code`、`storage_backend`、`flow_ttl_seconds`、token 来源）写入内存缓存与 Redis（若可用），并允许通过 Flow ID 回填。
- **FR-004**: 新增的调试 API（如账户余额）必须复用 `RedBookJuGuangACClient`，脱敏日志，返回结构化 JSON，并在请求/响应中标明 Provider。
- **FR-005**: 文档与配置模板需补充聚光条目，包括 `config.example.yaml`、`docs/develop/access-token/redbook/{develop,debug}.md`，并强调使用 `go run ./cmd/accesstoken/server -config config.yaml` 的统一命令。
- **FR-006**: `/accesstoken/flows` API 必须支持 `provider_code=redbook_juguang` 的过滤、分页与 Flow 回填，使 CLI/调试台共享同一份数据。
- **FR-007**: 所有聚光相关日志字段需遵循现有脱敏与结构化规范（`provider`, `provider_app`, `flow_id`, `token_source`）。

### Key Entities *(include if feature involves data)*

- **Provider Metadata**: 包含 Provider、App、授权模式、`config_path`、`oauth_key`、`api_version`；聚光条目需加入现有调试台 JSON 中供 UI 渲染。
- **OAuth Flow Record**: 存储 Flow ID、`provider_code=redbook_juguang`、`provider_auth_mode`、`storage_backend`、`token_source`、TTL、脱敏 token 片段；在内存/Redis 之间共享。

## MediaX Architecture Guardrails *(must reference Constitution sections)

- **Provider Adapter Parity** (Constitution §Provider Adapter Parity): 仅在 `pkg/client/redBook/juGuang/accessTokenClient` 现有适配器与 `MediaX.CreateRedBookJuGuangACClient` 工厂中注册，不增加平行实现，保持与 Google/Bilibili 相同的 cmd/accesstoken 入驻模式。
- **Config-Layered Security** (Constitution §Config Layered Security): 扩充 `pkg/client/config/redbook.go` 与配置模板，所有 `client_secret`/`scope` 均通过环境变量注入；文档禁止硬编码敏感值。
- **Token Lifecycle Discipline** (Constitution §Token Lifecycle Discipline): OAuth Flow 继续使用 `kernel.BaseClient` + `cache.ICache` 的 TTL/Flow-Index 设计，不改变其它 Provider 的刷新策略；Flow/Redis键名遵循 `accesstoken:oauth:*` 规范。
- **Observability & Error Traceability** (Constitution §Observability & Error Traceability): 调试 API 与 OAuth Handler 需记录 `provider`, `provider_app`, `flow_id`, `token_source`, `config_path` 字段，并通过 `internal/accesstoken/handler/mask` 实现脱敏。
- **Testable Modularity & SessionToken Readiness** (Constitution §Testable Modularity): 新增契约测试覆盖聚光 OAuth Start、Flow 列表、API 调用，严格限制在纯逻辑单元，可独立执行并保护其它 Provider 的 SessionToken 行为不受影响。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 在单台开发机上，运行 `go run ./cmd/accesstoken/server -config config.yaml` 后 `/debug` 可列出聚光 Provider，并能在 5 分钟内完成 OAuth 授权（以 5 次试验为样本，成功率 ≥ 80%）。
- **SC-002**: 聚光 Flow 记录支持 Redis 与内存两种模式，服务重启后通过 Flow ID 回填的成功率 ≥ 95%（20 次回填实验）。
- **SC-003**: 至少 1 条聚光 API 通过新增调试 endpoint 成功调用，平均响应时间 ≤ 2 秒（本地网络），日志中无明文 token。
- **SC-004**: 新增 `docs/develop/access-token/redbook/{develop,debug}.md`，由两名内部使用者仅靠文档即可完成配置到调试的闭环，成功率 100%。
