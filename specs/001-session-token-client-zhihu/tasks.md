---
description: "Task list for Zhihu SessionToken Cookie Orchestrator"
---

# Tasks: Zhihu SessionToken Cookie Orchestrator

**Input**: `docs/develop/session-token/zhihu/develop.md`, `docs/develop/session-token/zhihu/debug.md`, updated [spec.md](./spec.md) & [plan.md](./plan.md)
**Prerequisites**: Redis 可用，本地可运行 `make sessiontoken`
**Tests**: 每个核心模块（Flow/Callback/API/中间件/调试页）必须有单测或可执行脚本；PRD 要求的闭环需在 quickstart 中验证。

## Format: `[ID] [P?] [Story] Description`

- **[Story]** 对应 spec 中的 US-01～US-04
- **[P]** 表示可与其它任务并行（不同文件/模块）
- 引用具体文件/目录，便于 code search

## Phase 0 – Flow & Metadata 升级（US-01）

- [x] T100 [US01] 更新 `pkg/client/sessionToken/flow.go`、`storage/redis/store.go`，新增 `provider_app_code`、`metadata.session_token`、`metadata.cookie_*`、`credentials_note`、`credentials_expires_hint`、`code/message/last_failed_api` 字段及编解码逻辑。
- [x] T101 [US01] 扩展 `pkg/client/sessionToken/manager.go` (`CreateFlow`/`GetFlow`/`CompleteFlowSuccess`/`CompleteFlowFailed`) 以写入/读取新字段，并保证 TTL、state 幂等。
- [x] T102 [US01] 调整 `server/handlers/session_token/flow_create.go` 与 `flow_get.go` 的响应 DTO/contract，输出新的 metadata 字段；同步更新 `specs/001-session-token-client-zhihu/contracts/session-token-flows.yaml`。
- [x] T103 [P][US01] 单测：Flow encode/decode、manager 成功/失败路径、handler 响应脱敏（`pkg/client/sessionToken/..._test.go`, `server/handlers/session_token/..._test.go`）。

## Phase 1 – 回调、错误码与失效逻辑（US-01/US-03）

- [x] T110 [US01] 在 `pkg/client/sessionToken/callback/payload.go` 定义 `credentials`、`code/message`、`metadata.last_failed_api`，并在 dispatcher 注入 HMAC 签名、2s→4s→8s retry 与结构化日志。
- [x] T111 [US01] 扩展 `pkg/client/sessionToken/callback/dispatcher.go` 记录 `sessiontoken_callback` 字段（provider/api/flow_id/status/code/retry_count/latency_ms），并在 `/debug/callback`（`cmd/sessiontoken/debug_page.go`）存储最近 20 条记录。
- [x] T112 [US03] 实现 `server/zhihu/sessionToken/middleware/session_token.go`（或扩展现有中间件）解析 Header `X-SessionToken`，缺失/格式错误返回 401 `ZH_COOKIE_EXPIRED` 并写审计日志。
- [x] T113 [US03] 在 manager 增加 `MarkTokenInvalid(flow_id, code, message, api)`（或复用 `CompleteFlowFailed`），供 API/心跳调用触发 Flow fail + 回调。
- [x] T114 [P][US03] 单测：callback payload/签名/retry、X-SessionToken 中间件、API 触发失效 → Flow fail + 回调（可用 fake callback server）。

## Phase 2 – Zhihu API 网关（US-02/US-03）

- [x] T120 [US02] 构建 `pkg/client/zhihu/web/sessionTokenClient/client.go`（封装 BaseClient 调用 + 错误码映射），支持代理/UA/header 注入。
- [x] T121 [US02] 实现 `GET /zhihu/v1/me/followings` 与 `GET /zhihu/v1/channels/{channel_id}/articles`（`followings.go`, `channels_articles.go`），支持 `limit/offset`，附带 `meta.source`。
- [x] T122 [US02] 实现 `GET /zhihu/v1/articles/{id}`（`articles_get.go`），遵循同样的错误码策略。
- [x] T123 [US03] 实现 `POST /zhihu/v1/articles` 与 `POST /zhihu/v1/sanity/check`（`articles_post.go`, `sanity_check.go`），前者校验标题/内容，后者调用 `https://www.zhihu.com/api/v4/me` 并在 401/403 时触发 Flow fail。
- [x] T124 [P][US02] 为所有 handler 添加单测，覆盖成功/400/401/403/5xx、缺失 `X-SessionToken`、风控返回等路径；可使用 httptest + fake upstream。

## Phase 3 – 调试台与文档（US-01/US-04）

- [x] T130 [US04] 改造 `cmd/sessiontoken/debug_page.go`：Provider→App 级联（含自定义）、`window.createFlow/pollFlow/copyAuthorizeURL` global 暴露、默认 Callback URL `/debug/callback`、LocalStorage 记忆、Mock callback 日志区域；确保点击按钮有反馈。
- [x] T131 [US04] `/debug/callback` handler（`cmd/sessiontoken/debug_page.go` 或单独文件）使用 ring buffer 保存 20 条记录，提供 GET（列表）/DELETE（清空），并输出 `sessiontoken_debug` 日志。
- [x] T134 [US04] 暴露 `/debug/flows/<flow_id>/metadata`（GET/POST/PUT），允许浏览器授权页/脚本直接写回 Cookie、credentials_note 等字段，取代手动编辑 Redis；写入成功需打 `sessiontoken_debug: metadata_update` 日志。
- [x] T132 [US04] 根据最新调试策略更新 `docs/develop/session-token/zhihu/debug.md`、`specs/001-session-token-client-zhihu/quickstart.md`、README，写明“启动 Redis → make sessiontoken → /debug → 浏览器登录 → 抓取 Cookie → POST metadata → Mock callback → 调用 API → 模拟失效”的步骤。
- [x] T133 [P][US04] 编写脚本/命令示例：`scripts/sessiontoken-debug.sh` 或 README 片段，支持快速读取 Redis metadata 或调用 `zhihu/v1` API；文档应引用 `tail -f logs/sessiontoken-info.log | rg sessiontoken_(metric|callback)`。
- [x] T135 [US04] 提供 Playwright/Puppeteer 驱动的调试器脚本（`tools/sessiontoken-browser/`），支持 `pnpm sessiontoken:browser --flow <id>`：脚本会打开 Chromium 加载 authorize_url、提示用户登录、自动抓取 Cookie 并调用 `/debug/flows/<flow_id>/metadata`，完成后在终端输出 `/debug/callback` 返回结果；在 `/debug` 页面增加 CLI 指南以便复制命令。

## Phase 4 – 配置与测试覆盖（跨故事）

- [x] T140 [P][US01] 扩展 `pkg/client/config/zhihu.go` 与 `config.yaml` 示例，新增 `watch_cookies`, `callback.secret`, `network.proxy`, `api.retry_backoff`, `SESSIONTOKEN_DEBUG_CALLBACK_URL` 等字段；在 `README.md`/`.env.example` 中给出环境变量映射。
- [x] T141 [P][US01] 更新 `specs/001-session-token-client-zhihu/data-model.md`、`contracts/` 等文档，描述新字段/回调协议。
- [x] T142 [P][All] `go test ./...`、lint、`make sessiontoken` + Playbook 跑通：创建 Flow → 浏览器登录（或模拟 metadata）→ `/debug/callback` 观察 → 调用 `zhihu/v1` API → 触发 `ZH_COOKIE_EXPIRED` 并确认回调。

## Phase 5 – Zhihu API 版本切换（US-02 扩展）

- [x] T150 [US02] 重构 `pkg/client/zhihu/web/sessionTokenClient` 目录，按版本拆分（如 `v4/`, `v5/`），将现有 handler/测试迁移至 `v4`，并保留公共中间件/日志在上层。
- [x] T151 [US02] 在 `pkg/client/config/zhihu.go` 与 `config.yaml`/`config.example.yaml`/`.env.example` 新增 `service.api_version` 配置；`client.go` 根据版本动态注册对应 handler，默认使用 `v4`。
- [x] T152 [P][US02] 更新 `/debug` 页面 Provider App 模板、`README.md`、`docs/develop/session-token/zhihu/develop.md`/`debug.md`、`specs/001-session-token-client-zhihu/quickstart.md`，描述 API 版本切换流程，并在 Flow metadata/日志中输出当前版本；补充单测覆盖 `api_version` 选择。

## Notes

- 每个阶段都要确保日志脱敏（只输出 Cookie 前后缀）。
- `tasks.md` 应随着实现进度勾选完成，便于追踪。
