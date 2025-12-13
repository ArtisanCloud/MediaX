# Implementation Plan: Zhihu SessionToken Cookie Orchestrator

**Branch**: `001-session-token-client-zhihu` | **Spec**: [spec.md](./spec.md) | **Sources**: `docs/develop/session-token/zhihu/develop.md`, `docs/develop/session-token/zhihu/debug.md`

## Summary

为了配合新的 PRD，本次迭代需要把 SessionToken 服务从“仅能创建/查询 Flow”扩展为“可模拟登录、聚合 Cookie、标准化回调并暴露 Zhihu API”的产品化能力。核心目标：

1. **Flow & Metadata 升级**：Flow 实体要保存 Provider→App、`session_token` 聚合串、拆分后的 `cookie_*` 字段、`credentials_note`、`credentials_expires_hint` 以及 `code/message/last_failed_api`。Harvester/manager/Redis 层要支持这些字段的写入/读取/TTL 并在 `/debug`、回调、API 查询中展示。
2. **回调与失效策略**：扩展 callback dispatcher，输出 `ZH_COOKIE_EMPTY/EXPIRED/RISK` 等错误码、HMAC 签名、retry 计数；当 API 调用检测到 401/403、风控等事件时自动触发 Flow fail + 回调。
3. **API 网关**：实现五条 Zhihu API (`me/followings`、`channels/.../articles`、`articles/{id}`、`articles` POST、`sanity/check`)，所有请求通过 `X-SessionToken`＋BaseClient 调用上游，输出统一错误码并记录指标。
4. **调试台与文档**：`cmd/sessiontoken/debug_page.go` 需要支持 Provider→App 模板、Mock callback 日志、默认 Callback URL、本地记忆，并暴露 `/debug/flows/<flow_id>/metadata` 等调试接口，方便浏览器授权页写回 Cookie。额外提供 Playwright/Puppeteer 驱动的桌面脚本，在本地自动拉起 Chromium→登录知乎→捕获 `document.cookie`→POST metadata→刷新 `/debug/callback`，把整个流程固化到命令行/按钮里。`docs/.../debug.md`/`develop.md` 的流程要体现在 quickstart/spec/tasks 里，覆盖“启动脚本→登录→写回→Mock 回调→调用 API”。

## Technical Context

- **Language**: Go 1.18
- **存储**: Redis（Flow 状态、metadata、TTL）
- **关键组件**:
  - `pkg/client/sessionToken/manager.go`: 状态机/Flow 完成逻辑
  - `pkg/client/sessionToken/flow.go` & `storage/redis`: Flow 实体、Redis 编解码
  - `pkg/client/sessionToken/callback/*`: 签名、retry、日志
  - `server/zhihu/sessionToken/*`: Authenticator/Harvester/Callback adapter
  - `pkg/client/zhihu/web/sessionTokenClient/*`: Zhihu Web API 封装（对外调用入口）
  - `cmd/sessiontoken/debug_page.go`: 调试 UI
- **Testing**: `go test ./pkg/client/sessionToken/... ./server/...` + 集成测试（mock callback server、fake Zhihu API）。
- **Dependencies**: `github.com/redis/go-redis/v9`, `github.com/ArtisanCloud/MediaXCore` BaseClient、logger；需要 MockHTTP/httptest 用于 API 代理与 sanity check。

## Constitution Guardrails

| Guardrail | Plan 对应措施 |
| --- | --- |
| Provider Adapter Parity | 保持 Zhihu adapter 在 `server/zhihu/sessionToken/{authenticator,harvester,callback}`，并由 `pkg/client/mediaX.go` 暴露 `CreateZhihuSessionTokenClient`。配置仍由 `pkg/client/config/zhihu.go` 注入。 |
| Config-Layered Security | `ZhihuSessionTokenConfig` 将新增 `watch_cookies`, `callback.secret`, `network.proxy`, `api.retry_backoff` 等字段，所有 secret 通过 env；`config.yaml` / README / quickstart 更新示例。 |
| Token Lifecycle Discipline | Flow 完成逻辑仍依赖 manager + BaseClient；新增 API 通过 BaseClient HttpHelper 代理知乎，上游失败触发 Flow fail。Redis TTL/幂等继续由 FlowStore 管理。 |
| Observability & Error Traceability | 扩展 `sessiontoken_metric`、`sessiontoken_callback` 日志字段；API handler 输出 `provider/api/flow_id/tenant_uuid/latency_ms` 并统一脱敏 `session_token`。`/debug/callback` 也写入 logger，帮助定位。 |
| Testable Modularity & SessionToken Readiness | Flow/metadata 编解码、回调 payload、错误码映射、`X-SessionToken` 中间件、Zhihu API handlers 均补齐单测。`docs/.../debug.md` 的 Playbook 需要在 quickstart 中转换为可执行步骤。 |
| Operational Reproducibility | `make sessiontoken` 仍是权威入口，文档要指导如何同时启动 Redis + 服务 + `/debug` 页，如何 tail `logs/sessiontoken-info.log`，以及如何使用 Mock callback。 |

## Project Structure（关键文件）

- `cmd/sessiontoken/main.go`：注册 `/session-token/flows`、`/zhihu/v1/*`、`/debug`，加载 `SESSIONTOKEN_DEBUG_CALLBACK_URL`。
- `cmd/sessiontoken/debug_page.go`：Provider/App 模板、Mock 回调日志、`window.createFlow` 等函数。
- `pkg/client/sessionToken/flow.go`、`storage/redis/store.go`: Flow struct + metadata。
- `pkg/client/sessionToken/manager.go`: `CreateFlow`、`GetFlow`、`CompleteFlowSuccess/Failed`、`MarkTokenInvalid`。
- `pkg/client/sessionToken/callback/*`: payload、签名、retry、日志。
- `server/zhihu/sessionToken/*`: authenticator（登录入口）、harvester（Cookie 抓取）、callback（provider metadata）。
- `pkg/client/zhihu/web/sessionTokenClient/v4/*`: REST API 封装（按版本拆分）+ `X-SessionToken` 校验，入口包负责路由版本化。
- `specs/001-session-token-client-zhihu/quickstart.md`: Playbook。

## Phase Breakdown

### Phase 0 – Flow & Metadata 升级
1. 更新 `pkg/client/sessionToken/flow.go` + Redis encoder，添加 `provider_app_code`, `metadata.session_token`, `metadata.cookie_*`, `code/message`, `last_failed_api`, `credentials_note`, `credentials_expires_hint`。
2. 扩展 `sessionToken.Manager`（`CreateFlow`, `GetFlow`, `CompleteFlowSuccess/Failed`）以处理新字段；`CompleteFlowFailed` 接受 `code/message/last_failed_api`。
3. 更新 `/session-token/flows` handlers & contract 输出，确保响应/查询可见新字段。
4. 单测：Flow encode/decode、manager 完成逻辑、handler 响应包含 metadata。

### Phase 1 – Callback & 失效策略
1. 在 `pkg/client/sessionToken/callback/payload.go` 中定义 `credentials` 块、`code/message`、`metadata.last_failed_api`；dispatcher 根据 `code` 选择日志级别。
2. 扩展 Callback retry（2s→4s→8s），并记录 `sessiontoken_callback` 日志字段 `code`/`retry_count`；`/debug/callback` 存储 20 条记录。
3. 新建 `server/zhihu/sessionToken/middleware/session_token.go` 或复用现有逻辑，解析 Header `X-SessionToken`，缺失→401 `ZH_COOKIE_EXPIRED`。
4. 在 API 调用失败（401/403/风控）时调用 manager `MarkFlowFailed`（带 code/message/last_failed_api）并回调。
5. 单测：dispatcher payload、Middleware、API 401 → Flow fail。

### Phase 2 – Zhihu API 网关
1. `pkg/client/zhihu/web/sessionTokenClient/followings.go`：调用 `https://www.zhihu.com/api/v4/people/{uid}/following-columns`（uid 从 session_token metadata/配置获得），返回 JSON + `meta.source` 字段。
2. `pkg/client/zhihu/web/sessionTokenClient/channels_articles.go`：代理 feed API，支持 `limit/offset`。
3. `pkg/client/zhihu/web/sessionTokenClient/articles_get.go`：获取文章详情。
4. `pkg/client/zhihu/web/sessionTokenClient/articles_post.go`：发布文章，body 需校验标题/内容。
5. `pkg/client/zhihu/web/sessionTokenClient/sanity_check.go`：调用 `/api/v4/me` 检测 Cookie；若失败 → Flow fail + 回调。
6. 公共助手：`pkg/client/zhihu/web/sessionTokenClient/client.go`（封装 BaseClient 调用 + error mapping）。
7. 单测：使用 httptest/fake upstream 覆盖成功/401/403/5xx、参数校验、错误码映射。

### Phase 3 – 调试台 & 文档
1. `cmd/sessiontoken/debug_page.go`: 
   - Provider/App 模板配置来自 `providerCatalog`（Zhihu/RedBook/Google），支持自定义 app code。
   - `window.createFlow/pollFlow/copyAuthorizeURL` 暴露到 global，避免按钮无响应。
   - 添加 Mock callback 区、默认 Callback URL= `/debug/callback`、LocalStorage 记忆、state 随机器。
   - 暴露 `/debug/flows/<flow_id>/metadata`（GET/POST/PUT），便于 QA 在浏览器授权页 `document.cookie` 抓取后直接写回 Flow metadata，而无需手动改 Redis。
2. 新增 Playwright/Puppeteer 调试器：
   - `tools/sessiontoken-browser/`（或同级目录）封装 CLI：`pnpm sessiontoken:browser --flow <id> --base http://127.0.0.1:7070`。
   - CLI 根据 Flow ID 拉取 authorize_url，使用 Playwright 打开 Chromium；提示用户在自动拉起的窗口完成登录。
   - 登录完成后脚本通过 CDP 读取 Cookie 并调用 `/debug/flows/<flow_id>/metadata`；结果/回调输出在终端（或写回 `/debug` via WebSocket）。
   - `/debug` 页面可增加“Launch Playwright”按钮，点击后提示运行命令或通过自定义协议触发脚本。
3. `/debug/callback`：内存 ring buffer（20 条），GET/DELETE API + 日志 `sessiontoken_debug`。
4. 文档：
   - `docs/develop/session-token/zhihu/debug.md` → 已更新，需同步引用到 `specs/.../quickstart.md` 与 README。
   - `specs/.../quickstart.md`：写出“启动 Redis → make sessiontoken → 打开 /debug → 浏览器登录 → tail logs → 调用 zhihu/v1 API → 模拟失效”的步骤。
   - `tasks.md`/`plan.md`/`spec.md` 对应 PRD 的更新。
4. Quick verification 脚本：`scripts/sessiontoken-debug.sh`（可选）帮助 QA 读取 metadata。

## Testing & Verification

- 单元测试：Flow 状态机、callback、X-SessionToken 中间件、Zhihu handlers。
- 合同测试：`contracts/session-token-flows.yaml` + 新增 `contracts/zhihu-api.yaml`（若需要）。
- E2E：`make sessiontoken` + 浏览器登录 + Mock callback + curl API。
- 日志检查：`tail -f logs/sessiontoken-info.log | rg 'sessiontoken_(metric|callback)'` + `/debug/callback` JSON。

## Risks & Mitigations

- **Cookie 字段遗漏**：通过 config `watch_cookies` + 单测 + `/debug` Metadata viewer 及时发现。
- **回调失败影响 SLA**：重试 + 指标 `retry_count`，并使用 `/debug/callback` 记录失败 payload。
- **API 代理受限**：统一在 BaseClient 层管理代理/UA/headers，与 `ZhihuSessionTokenConfig.Network` 对齐，并允许在 `.env` 里修改。
- **调试体验差**：调试页/文档/Playbook 统一更新，默认 Callback URL 直达 `/debug/callback`，QA 无需插件即可复现。

## Deliverables

1. 更新后的代码（Flow/Manager/Callback/API/调试页）。
2. 文档：`spec.md`、`plan.md`、`tasks.md`、`quickstart.md`、README/Makefile 片段，与 `docs/develop/session-token/zhihu/*.md` 对齐。
3. 测试报告：`go test ./...` 输出 + curl/Playbook 验证结果。
4. 日志示例：`sessiontoken_metric`、`sessiontoken_callback`、`sessiontoken_debug`。
