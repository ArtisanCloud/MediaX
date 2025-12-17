# Feature Specification: Zhihu SessionToken Cookie Orchestrator

**Feature Branch**: `001-session-token-client-zhihu`  
**Created**: 2025-12-09（更新：2026-02-XX）  
**Status**: Draft → Aligning with docs/develop/session-token/zhihu  
**Input**: `docs/develop/session-token/zhihu/develop.md`（PRD）+ `docs/develop/session-token/zhihu/debug.md`（调试指南）

## Clarifications & Alignment

1. SessionToken 服务需要完整覆盖“创建 Flow → 浏览器模拟 → Cookie 聚合 → 回调告警 → API 复用”闭环；`/debug` 页面被视为官方调试入口，必须支持 Provider→App 选择、Mock Callback、元数据观察。
2. `metadata.session_token` 并非真实 Cookie 名称，而是聚合字符串（`SESSIONID=...; JOID=...; ...`），任何对外 API 都要求调用方通过 `X-SessionToken` 回传该串，服务端负责拆解为 Cookie；此概念需要在 spec/plan 中明确。
3. PRD 定义的错误码 `ZH_COOKIE_EMPTY`、`ZH_COOKIE_EXPIRED`、`ZH_COOKIE_RISK` 等需体现在 Flow 状态机与回调 payload 中，调试页面和日志必须能观察到。
4. 需要新增一批 REST API（频道/文章/发帖/心跳），所有请求共用 `X-SessionToken`，失败会触发 Flow 失败及回调。
5. `/debug` 页面要配合本地 7070 服务、Redis，并默认提供 Mock callback（`/debug/callback`），用于独立调试，无需依赖插件环境。

## User Scenarios & Testing

| ID | Priority | User | 场景 | 独立测试 |
| --- | --- | --- | --- | --- |
| US-01 | P1 | 平台运营 | 在 `/debug` 或插件上创建 Zhihu Flow、复制 `authorize_url`、完成浏览器登录，SDK 聚合 Cookie 并通过回调提供 `session_token` | 通过 `/session-token/flows` + `/debug` + `/debug/callback` 验证 Flow 生命周期、metadata 字段、回调 payload |
| US-02 | P1 | 内容采集服务 | 持有 `session_token` 后调用 `GET /zhihu/v1/me/followings`、`/channels/{id}/articles`、`/articles/{id}` 等 API 获取订阅内容 | 使用 curl/contract 调用 API 并检查 `X-SessionToken` 校验、上游错误码映射、日志脱敏 |
| US-03 | P1 | 发布服务 | 使用 `POST /zhihu/v1/articles` 发布内容以及 `POST /zhihu/v1/sanity/check` 进行心跳，失效时触发 `ZH_COOKIE_EXPIRED` 回调 | 单测 + e2e：伪造 session_token 调用 API，验证 401/403 → Flow 标记 failed + 回调 |
| US-04 | P2 | 运维/QA | 本地使用 `/debug` + Mock callback + Redis 工具定位 `metadata.session_token is empty`/脚本失败等问题 | 通过 `docs/develop/session-token/zhihu/debug.md` 的 Playbook 重复流程，检查日志/回调记录 |

### Acceptance criteria examples
- US-01: Flow 成功时 `metadata` 必须包含 `session_token` 以及拆分的 `cookie_sessionid`/`cookie_joid`/...，并透传到回调；失败时 `code`/`message` 显示 `ZH_COOKIE_EMPTY` 等。
- US-02: 每个 API 都要求 Header `X-SessionToken`，并在日志中打印 `provider=zhihu api=channels.articles flow_id=? tenant_uuid=? latency_ms`；缺失 header 返回 401。
- US-03: `POST /zhihu/v1/sanity/check` 会访问 `https://www.zhihu.com/api/v4/me`，若 401/403 则调用 `CompleteFlowFailed`、写入 `metadata.last_failed_api` 并回调。
- US-04: `/debug` 页面“创建 Flow/查询 Flow/Mock 回调日志”均可操作，默认 Callback URL= `http://127.0.0.1:7070/debug/callback`，日志中可见 `sessiontoken_callback` 或 `sessiontoken_callback: failed ... metadata.session_token is empty`。

## Functional Requirements

1. **FR-001 Flow 管理**：`POST /session-token/flows`、`GET /session-token/flows/{id}` 必须覆盖 Provider→App、state 幂等、TTL、作者 metadata，状态机扩展为 `pending → authorizing → succeeded/failed` 并记录 `code/message/last_failed_api`。
2. **FR-002 Cookie 聚合**：Harvester 必须按 PRD 定义收集 `SESSIONID/JOID/osd/q_c1/d_c0/unlock_ticket/...`，生成 `metadata.session_token`（完整拼接）及 `metadata.cookie_*` 字段，`credentials.session_token` 与 metadata 保持一致。
3. **FR-003 回调协议**：回调 payload 需要包含 `credentials` 块、`code/message`、`metadata.last_failed_api`；失败时 `credentials=null`，code 取 `ZH_COOKIE_*`；所有 payload 通过 HMAC-SHA256 + timestamp/nonce 签名，失败重试 3 次（2s→4s→8s）。
4. **FR-004 调试控制台**：`cmd/sessiontoken/debug_page.go` 提供 Provider→App 模板、Mock callback 日志、LocalStorage 记忆、默认 Callback URL、State 随机工具；按钮必须调用全局函数，确保无响应问题。
5. **FR-005 API 封装**：实现 `GET /zhihu/v1/me/followings`、`GET /zhihu/v1/channels/{id}/articles`、`GET /zhihu/v1/articles/{id}`、`POST /zhihu/v1/articles`、`POST /zhihu/v1/sanity/check`，统一读取 `X-SessionToken`，通过 BaseClient 代理至知乎 API，原样/轻度裁剪响应，并映射错误码（`ZH_BAD_REQUEST`、`ZH_COOKIE_EXPIRED`、`ZH_RISK_BLOCK`、`ZH_UPSTREAM_ERROR`）。
6. **FR-006 Token 校验**：所有 `zhihu/v1/*` API 未包含 `X-SessionToken` 或解析失败时返回 401 + `ZH_COOKIE_EXPIRED`，并触发 Flow 失败/回调；成功调用需记录 `captured_at`/`credentials_note`。
7. **FR-007 监控与日志**：新增/保持 `sessiontoken_metric`、`sessiontoken_callback` 日志，字段包含 `action/provider/api/tenant_uuid/account_id/flow_id/status/latency_ms/retry/code`；调试页的 Mock callback 也要写日志 `sessiontoken_debug: callback ...`。
8. **FR-008 文档交付**：`docs/develop/session-token/zhihu/debug.md`、`develop.md` 的流程要在 spec/plan/quickstart 中反映；`README.md` 与 `specs/.../quickstart.md` 必须说明本地调试步骤、`X-SessionToken`、API 列表。

## Non-functional & Guardrails

- **Config-Layered Security**：`pkg/client/config/zhihu.go` 中的 `ZhihuSessionTokenConfig` 需要提供 `watch_cookies`、`callback.secret`、`network.proxy` 等字段，全部可通过 env 覆盖；`config.yaml` 示例需提及 `POWERX_SESSION_TOKEN_*`、`SESSIONTOKEN_DEBUG_CALLBACK_URL`。
- **Observability**：日志禁止输出完整 Cookie，使用 `sanitizer.MaskToken`；`/debug` Mock 回调存储最近 20 条记录，便于 QA 复查。
- **Testable Modularity**：Flow 状态机、回调签名、API proxy、`X-SessionToken` 中间件必须拥有独立单测，覆盖 `ZH_COOKIE_EMPTY`/`EXPIRED` 等分支。
- **Operational Reproducibility**：`make sessiontoken` + `/debug` 指南必须可让新同学在 10 分钟内完成 Flow 创建、浏览器登录、Mock 回调调试。

## Success Criteria

1. 95% Zhihu Flow 能在 5 分钟内完成并回调 `session_token`，失败 Flow 均包含 `code/message` 与 `last_failed_api`。
2. 所有 `zhihu/v1/*` API p95 延迟 < 300ms（代理上游除外），`X-SessionToken` 缺失/失效立即返回 401 并触发 Flow fail 回调。
3. `/debug` 页面和 Mock callback 在本地环境可操作，日志中能检索到 `sessiontoken_callback` 或 `sessiontoken_callback: failed ... metadata.session_token is empty`。
4. Quickstart 按文档步骤可复现授权→调试→API 调用→失效告警全链路。
