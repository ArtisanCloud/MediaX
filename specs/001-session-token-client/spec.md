# Feature Specification: SessionTokenClient Provider Architecture

**Feature Branch**: `001-session-token-client`  
**Created**: 2025-12-09  
**Status**: Draft  
**Input**: User description: "Implement provider-scoped SessionTokenClient with Zhihu adapter, flow storage, and signed callbacks per docs/plan/session_token_client.md"

## Clarifications

### Session 2025-12-09
- Q: Flow 持久化的主存储应选择哪种方案以同时满足 TTL、幂等和多实例共享？ → A: 统一使用 Redis 作为主存储，依托 go-redis 管理 TTL/索引，并在需要长期保留时再导出到业务存储。
- Q: 回调失败时的默认重试策略需要怎么设定？ → A: 默认最多重试 3 次，并采用 2s→4s→8s 的指数退避，超出后标记失败并记录 last_error。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 插件创建 SessionToken Flow 并拉起登录 (Priority: P1)
插件运营人员希望为指定 provider（首期知乎）发起一次模拟登录，SDK 需要生成 authorize_url、flow_id，并在 TTL 内保持 pending 状态，供插件交互式登录。

**Why this priority**: 这是 SessionToken 模式的入口，没有 Flow 创建就无法驱动登录与凭证采集。

**Independent Test**: 通过调用 `POST /session-token/flows`，验证返回 flow_id/authorize_url，并在无浏览器交互时维持 `pending`，可完全回归测试。

**Acceptance Scenarios**:

1. **Given** 插件提供 provider_code、tenant_uuid、callback_url 等必填字段，**When** SDK 验证请求并存储 Flow，**Then** 返回带有 `flow_id`、`authorize_url`、`expires_at` 的 `pending` 状态响应。
2. **Given** 同一租户重复提交 state，**When** SDK 接收到重复 Flow 请求，**Then** 通过幂等策略返回已有 Flow 或提示冲突，避免并发创建。

---

### User Story 2 - 插件轮询 Flow 状态 (Priority: P1)
插件需要实时了解 Flow 是否进入 authorizing、成功或失败，以便提示运营人员继续登录或重试。

**Why this priority**: 没有可查询 API，就无法驱动交互式登录流程并判断何时结束。

**Independent Test**: 通过 `GET /session-token/flows/{flow_id}` 验证状态迁移、错误描述与结果 payload 可独立测试。

**Acceptance Scenarios**:

1. **Given** Flow 已处于 `authorizing`，**When** 插件轮询查询，**Then** 返回状态、authorize_url、剩余 TTL 以及最近的 `last_error`（若有）。
2. **Given** Flow 进入 `succeeded` 或 `failed`，**When** 插件查询，**Then** 返回最终状态、结果/错误描述，并在审计 TTL 内可重复读取。

---

### User Story 3 - 凭证采集与回调 (Priority: P2)
Authenticator/Harvester 完成登录后，SDK 需要脱敏存储结果、触发签名回调，把 session_token/cookie 返回给插件。

**Why this priority**: 没有凭证回传就无法把 Flow 结果写入业务，会导致流程失效。

**Independent Test**: 通过模拟 Harvester 输出和回调验证，确保签名校验、日志脱敏与重试策略可单独回归。

**Acceptance Scenarios**:

1. **Given** Harvester 捕获凭证，**When** SDK 更新 Flow 为 `succeeded`，**Then** 触发回调，payload 包含脱敏字段、state、签名，通过 HMAC 校验。
2. **Given** 回调 5xx 或超时，**When** SDK retry policy 生效，**Then** 记录 retry 次数并在超出阈值后将 Flow 标记为 `failed`、写入 `last_error`。

---

### User Story 4 - 知乎适配器与配置注入 (Priority: P2)
平台团队需要为知乎配置多种登录入口、代理、脚本策略并在 MediaX 工厂中注入 `ZhihuSessionTokenClient`。

**Why this priority**: 没有首个 provider 落地，抽象无法验证，插件也无法在真实场景验证流程。

**Independent Test**: 通过构造 Zhihu 配置与模拟登录脚本，验证工厂创建、配置读取、Authenticator/Harvester 组合逻辑。

**Acceptance Scenarios**:

1. **Given** `config.zhihu.SessionTokenConfig` 已在 `config.yaml` 中赋值，**When** MediaX 工厂创建 Zhihu SessionTokenClient，**Then** 生成的 Authenticator/Harvester 能读取登录入口、代理与脚本 ID。
2. **Given** 配置缺失或 secret 无效，**When** 构造器运行，**Then** 返回结构化错误并阻止客户端实例化。

### Edge Cases

- Flow 创建参数缺失或 provider_code 不被支持时，返回 4xx 并记录原因。
- `authorizing` 状态长时间无进展（超时/浏览器崩溃）需自动转 `failed` 并写入 `last_error`。
- 回调目标不可达或签名校验失败时，需要重试并在最终失败后阻止重复派发。
- Redis/持久化短暂不可用时，需要保证 Flow 写入具备重试/幂等，避免生成多个 flow_id。
- 多入口（密码/扫码/手机号）之间切换需确保 authorize_url/metadata 一致，避免插件拿到错误入口。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: SDK MUST 暴露 `POST /session-token/flows` 与 `GET /session-token/flows/{flow_id}` API，使插件可创建/查询 Flow，响应结构与 docs/plan/session_token_client.md 中示例保持一致。
- **FR-002**: 系统 MUST 定义 `Flow` 状态机（pending → authorizing → succeeded/failed），并把 Flow 记录持久化到共享 Redis 存储（默认 go-redis，实现 TTL/索引/幂等），字段包含 `flow_id`、provider、租户、状态、authorize_url、expires_at、last_error、result；如需长期审计再由业务侧异步导出。
- **FR-003**: 系统 MUST 在 `pkg/client/sessionToken` 中定义 `SessionTokenClient`、`Authenticator`、`CredentialHarvester`、`CallbackDispatcher` 接口，并在 `sessionToken.Manager` 内编排状态迁移与流程回调。
- **FR-004**: MediaX 工厂 MUST 暴露 `Create<Provider>SessionTokenClient`，自动注入 Logger、Cache、HTTP Helper 与 provider 配置，拒绝缺失配置的实例化。
- **FR-005**: 每个 provider 的配置文件 MUST 新增 `SessionTokenConfig`，包含 `Service`、`Authenticator`、`Harvester`、`Callback`、`Network` 五大块，并可通过 `yaml/json` tag 映射至 `config.yaml`。
- **FR-006**: 系统 MUST 实现 Zhihu 适配器（Authenticator/Harvester/CallbackDispatcher），支持账号密码、扫码、手机号多入口，允许配置代理池与脚本 ID/URL。
- **FR-007**: Flow 成功后 MUST 脱敏存储 session_token/cookies，仅输出前后缀到日志，并向配置中的 callback_url 发送含 state、flow_id、status、credentials 的 HMAC-SHA256 签名 payload，默认重试 3 次（2s→4s→8s 指数退避），超出后标记 `failed` 并写入 `last_error`，同时允许 provider 配置覆盖该策略。
- **FR-008**: API 层 MUST 校验 `Authorization: Bearer <api_token>`，拒绝未授权请求，并对所有外发/持久化操作记录 `provider`、`api`, `tenant_uuid`, `account_id`, `flow_id`。
- **FR-009**: 系统 MUST 提供可配置 TTL，确保 Flow 在完成后保留审计窗口并自动清理，TTL 到期后需返回明确的 `not_found/expired` 响应。

### Key Entities *(include if feature involves data)*

- **Flow**: 表示一次 SessionToken 登录会话，字段含 `flow_id`, `provider_code`, `tenant_uuid`, `account_id`, `state`, `status`, `authorize_url`, `expires_at`, `last_error`, `result`, `callback_url`, `metadata`。状态机驱动 API 与回调。
- **SessionTokenConfig**: 每个 provider 的配置结构，划分 `Service`（base_url、api_token）、`Authenticator`（入口 URL、UA、脚本、验证码策略）、`Harvester`（需抓取的 cookie/header、脚本描述）、`Callback`（URL、签名 secret、重试阈值）、`Network`（代理/超时）。
- **CredentialPayload**: `session_token`, `cookies_json`, `headers_json`, `expires_at`, `note`, `captured_at` 的封装体，供回调与查询返回使用。

## MediaX Architecture Guardrails *(must reference Constitution sections)*

- **Provider Adapter Parity**: 新增 `pkg/client/zhihu/sessionToken`（或等效路径）并在 provider 目录下保持 `core/`, `authenticator/`, `harvester/`, `callback/` 等结构，同时在 `pkg/client/mediaX.go` 注册 `CreateZhihuSessionTokenClient`，遵循宪章对目录和工厂模式的统一要求。
- **Config-Layered Security**: 在 `pkg/client/config/zhihu.go`（及其他 provider 配置文件）加入带 `yaml/json` tag 与文档注释的 `SessionTokenConfig`，所有 secret（api_token、callback secret、脚本凭证）必须通过环境变量或秘密管理注入并声明 TTL/逐出策略。
- **Token Lifecycle Discipline**: SessionTokenClient 内部仍复用 `kernel.BaseClient` 的 HttpHelper、重试与刷新钩子；Flow 创建/回调涉及的 HTTP 请求禁止绕过 BaseClient，并在需要覆写行为（如请求重放或代理）时给出明确注释。
- **Observability & Error Traceability**: 日志/事件必须携带 `provider`, `api`, `tenant_uuid`, `account_id`, `flow_id`，敏感凭证脱敏；回调 payload 强制 HMAC-SHA256 + timestamp/nonce 签名，并为 retry/replay 记录序号与触发条件。
- **Testable Modularity & SessionToken Readiness**: `sessionToken.Manager`、状态机与签名/脱敏工具需提供 `*_test.go`，覆盖 pending→authorizing→succeeded/failed 分支以及回调签名验证，同时提供 Zhihu 配置解析与入口选择的单测。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 插件使用 `POST /session-token/flows` 创建 Flow，在 200ms 内获得 authorize_url 和 flow_id，失败率 < 1%。
- **SC-002**: Flow 状态查询可在 100ms 内返回，并准确反映四种状态，错误信息具备可操作性（含 last_error）。
- **SC-003**: 95% 的 Zhihu Flow 在 5 分钟内完成并触发成功回调，失败 Flow 均写入 last_error 并通过回调/查询可见。
- **SC-004**: 所有回调 payload 100% 通过签名校验测试，且日志未出现明文 token/cookie；安全审计可复现完整 Flow 记录直至 TTL 结束。
