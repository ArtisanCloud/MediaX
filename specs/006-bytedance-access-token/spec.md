# Feature Specification: DouYin AccessToken Debug Integration

**Feature Branch**: `[006-bytedance-access-token]`  
**Created**: 2025-12-20  
**Status**: Draft  
**Input**: `docs/plan/byteDance/access_token_client.md`, `docs/develop/access-token/byteDance/{develop,debug}.md`

> **Design Principle**: 在不破坏现有 Provider（Google/BiliBili/RedBook 等）调试体验的前提下，把抖音 DouYin AccessToken 能力植入 `cmd/accesstoken` 服务，遵循相同的配置、Flow、日志、Redis 规范。

## Clarifications

### Session 2025-12-20

- Q: DouYin Flow 里需要怎样保存 refresh_token 与自动续期策略？ → A: 保存 access_token + refresh_token，并在 `/accesstoken/call` 检测到将过期或已过期时自动刷新并更新 Flow。
- Q: DouYin Provider 是否需要在 `config.yaml` 支持多套应用（不同 client_key/client_secret）同时调试？ → A: 仅支持单个 DouYin 应用，`<app>` 固定为 `default`。
- Q: DouYin Flow 的 TTL/Redis 过期时间要如何计算？ → A: 使用 DouYin 返回的 `expires_in` 值作为 TTL，并在自动刷新成功后重置 TTL。
- Q: `/accesstoken/call` 调 DouYin action 时需要怎样控制速率和重试？ → A: 每 action 1 QPS 的速率限制 + 针对 429/5xx 的最多 3 次指数退避重试。
- Q: 当自动刷新 DouYin token 失败（如 refresh_token 失效或返回不可恢复错误）时，Flow 应该如何处理？ → A: 将 Flow 标记为失效并从内存/Redis 删除，`/accesstoken/call` 返回“需重新授权”，`/debug` 提醒重新发起 OAuth。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - `/debug` 支持 DouYin OAuth (Priority: P1)

作为调试 `cmd/accesstoken` 的工程师，我可以在 `/debug` 页面看到 “字节跳动 / 抖音 Douyin” Provider，点击“发起授权”完成 DouYin OAuth，并在授权记录表中查看 Flow ID/过期时间，以便 CLI 与脚本复用。

**Why**: 没有 OAuth 就无法获取 AccessToken，也就无法调用任何 DouYin API。

**Independent Test**: `go run ./cmd/accesstoken/server -config config.yaml` → `/debug` 选择 DouYin Provider → 发起授权 → 浏览器跳至 DouYin 登录/授权页 → 回调成功后在 “授权记录” 中出现新的 `provider_code=byte_dance_douyin` Flow。

**Acceptance Scenarios**:
1. **Given** `config.yaml` 填写 `byte_dance_douyin_config`，**When** 在 `/debug` 发起授权，**Then** 返回有效授权 URL、日志输出 `provider=byte_dance_douyin`、UI 展示 Flow ID/expire_at。
2. **Given** DouYin 回调 `/debug/callback`，**When** 交换 access token 成功，**Then** Flow 同时写入内存与可选 Redis，并能在授权记录列表中点击“填充”回写 JSON。

---

### User Story 2 - DouYin API 可通过 `/accesstoken/call` 调试 (Priority: P2)

作为后端调试者，我希望在调试台或 cURL 调用 `/accesstoken/call`，指定 `provider_code=byte_dance_douyin` 与 DouYin action（例如 `douyin.video.list`、`douyin.im.message.send`），系统复用 SDK 客户端与最新 Flow Token，并脱敏日志/错误。

**Why**: 只有实际调用接口才能验证 AccessToken 生命周期与 SDK 封装可用，且能与其它 Provider 形成统一体验。

**Independent Test**: 在 Redis 中存在有效 DouYin Flow 后执行 `curl -H "Authorization: Bearer dev-accesstoken" -d '{"provider_code":"byte_dance_douyin","action":"douyin.video.list",...}' http://127.0.0.1:7071/accesstoken/call`，输出成功 JSON，日志含 `event=token.call provider=byte_dance_douyin action=douyin.video.list`。

**Acceptance Scenarios**:
1. **Given** Flow 可用，**When** 通过 `/accesstoken/call` 发起 DouYin action，**Then** Handler 使用 `ByteDanceDouYinACClient`，脱敏 AccessToken，返回结构化 JSON 与 `provider/provider_app/action` 元信息。
2. **Given** 请求缺少 `DOUYIN_SCOPE` 或 token，**When** 调试台调用，**Then** 返回解释性错误（如 `missing access token`、`OAuth scope 未配置`），并保持其它 Provider 不受影响。

---

### User Story 3 - 文档可独立复现 (Priority: P3)

作为 SDK/运营伙伴，我只需阅读 `docs/develop/access-token/byteDance/{develop,debug}.md`，就能完成环境变量、配置、授权、Flow 回填与 API 调试整个流程，且与 Google/BiliBili 文档格式一致。

**Independent Test**: 让两位未参与开发的同事只看该文档，按照说明在 1 小时内完成配置→授权→调用示例 API → 记录 Flow 证据。

**Acceptance Scenarios**:
1. **Given** 使用者仅阅读文档，**When** 根据步骤操作，**Then** 能成功获取 Flow 并通过 `/accesstoken/call` 调 DouYin 示例接口；若有问题能通过文档列出的排障项解决。

### Edge Cases

- `config.yaml` 中未配置 DouYin Provider 时 `/debug` 不应出现字节跳动入口，避免空配置干扰其它 Provider。
- 出现 `client_id`/`client_secret`/`DOUYIN_SCOPE` 缺失时必须沿用 FR-002 的校验分支，返回指明缺失字段的错误并记录日志，避免误导其它 Provider。
- Redis 不可用时需自动降级至内存，顶部显示 `storage_backend=memory` 提示，且 Flow 回填 API 返回 `flow_not_found` 时文档需说明重新授权。
- **回归约束**：任何新增逻辑不得破坏既有 Provider（Google/BiliBili/RedBook）的授权、调用与日志结构。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `/debug` 必须读取 `byte_dance_douyin_config` 并展示 Provider 卡片（含 `provider_code`、`api_version`、`oauth_key`、`config_path`）。
- **FR-002**: `/accesstoken/oauth/start` 对 `provider_code=byte_dance_douyin` 时需校验必填字段（client_id/secret/scope/redirect_url），缺失时返回描述性错误并记录日志。
- **FR-003**: `/debug/callback` 完成 DouYin token 交换后要写入内存与可选 Redis，键名遵循 `accesstoken:oauth:byte_dance_douyin:<app>:<mode>`，并保留 `flow_id` 索引。
- **FR-004**: `/accesstoken/call` 在 DouYin action 下必须复用 `ByteDanceDouYinACClient`，记录 `provider/provider_app/action/flow_id/token_source`，并通过 `mask` 工具脱敏 token。
- **FR-005**: 文档与 `config.example.yaml` 需补充 DouYin 配置片段与环境变量说明，确保复制即可使用。
- **FR-006**: `/accesstoken/flows`、`/api/oauth/tokens` 支持 `provider_code=byte_dance_douyin` 过滤与 Flow 回填。
- **FR-007**: 所有新增逻辑必须保持与现有 Provider 相同的日志字段与错误码，不可引入破坏性更改。
- **FR-008**: DouYin Flow 必须同时保存 access_token 与 refresh_token，`/accesstoken/call` 检测 token 将过期或已过期时自动使用 refresh_token 续期并同步到内存/Redis Flow 记录。
- **FR-009**: `byte_dance_douyin_config` 仅支持单个应用配置，`<app>` 固定 `default`（对应 `provider_app=douyin`），CLI/UI 必须提示“单实例/单 app”以避免误解。
- **FR-010**: Flow TTL 直接取 DouYin OAuth 响应 `expires_in`，自动续期成功后必须刷新 Redis/内存 TTL，并将新 `expires_in` 写入 Flow 记录。
- **FR-011**: `/accesstoken/call` 对 DouYin action 需实现 per-action 1 QPS 速率限制，并在收到 429/5xx 时触发最多 3 次指数退避重试；重试过程需记录 `retry_count` 与最终结果。
- **FR-012**: 若自动刷新失败（如 refresh_token 失效/不可恢复错误），需立即删除对应 Flow（内存与 Redis），并在 `/accesstoken/call` 返回“need reauth”错误码；`/debug` Flow 列表需提示“需重新授权”。

### Key Entities

- **DouYinProviderMeta**：`provider_code=byte_dance_douyin`、`app_code=douyin`（映射 `<app>=default`）、`auth_mode=default`、`oauth_key`、`api_version`。
- **DouYinOAuthFlowRecord**：Flow ID（`oauth-<state>`）、token 来源、TTL（等于 DouYin `expires_in`）、`storage_backend`、`provider_app`、`masked_token`、`refresh_token` 及最近续期时间；结构必须与现有 Flow 记录一致以便复用 UI/CLI。

## MediaX Architecture Guardrails

- **Provider Adapter Parity** (§Provider Adapter Parity): 只能复用 `pkg/client/byteDance/douYin/accessTokenClient` 与 `MediaX.CreateByteDanceDouYinACClient`，禁止新增平行实现。
- **Config-Layered Security** (§Config Layered Security): `config.yaml` 仅保留环境变量占位符；文档要求通过 `DOUYIN_*` env 注入，避免硬编码。
- **Token Lifecycle Discipline** (§Token Lifecycle Discipline): OAuth/Flow 写入沿用 `kernel.BaseClient` + cache 机制，Redis 键名遵循 `accesstoken:oauth:*`，不得自定义未经审计的缓存策略。
- **Observability & Error Traceability** (§Observability & Error Traceability): 日志记录 `provider`, `provider_app`, `flow_id`, `token_source`, `config_path`, `action`，并使用 `internal/accesstoken/handler/mask` 处理敏感值。
- **Testable Modularity & SessionToken Readiness** (§Testable Modularity): 新增合约测试验证 DouYin OAuth Start、Flow 保存、`/accesstoken/call` 调用；必须保持其它 Provider 测试通过，禁止回归。

## Success Criteria *(mandatory)*

- **SC-001**: 5 次迭代内有 ≥4 次在 `/debug` 完成 DouYin OAuth 与 Flow 保存（成功率 ≥80%），日志中无 panic。
- **SC-002**: Redis 与内存模式均可通过 `/api/oauth/tokens?flow_id=` 回填 DouYin Flow，随机抽取 20 个 Flow，成功率 ≥95%。
- **SC-003**: 选定至少 1 个 DouYin action（如 `douyin.video.list`）在 `/accesstoken/call` 调试成功，p95 响应时间 ≤2s，日志无明文 token。
- **SC-004**: 两名未参与开发的使用者仅凭文档即可在 1 小时内完成配置→授权→API 调试→截取 Redis 证据，反馈“无阻碍”；若失败需补充文档并复测。
