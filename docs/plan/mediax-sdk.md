# MediaX SDK 对接规划（频道订阅 & 平台账号）

> 目的：集中记录 MediaX SDK/MediaXCore 在频道订阅与平台账号矩阵场景下需要实现/扩展的能力，便于与插件团队对齐依赖与分工。

## 1. 模拟登录（SessionTokenClient）

- **服务进程**：MediaXCore 需提供独立的 SessionToken HTTP 服务（默认监听 `:7070`），在 `cmd/sessiontoken/main.go` 中启动 HTTP server，加载配置并注册 `/session-token/flows` 路由。
- **API Handler**：复用 `pkg/server/handlers/session_token` 下的 `RegisterSessionTokenFlowCreateRoute` 与 `RegisterSessionTokenFlowGetRoute`，并在中间件中校验 Bearer Token（默认 `dev-session-token`，可配置）。
- **Flow 管理器**：构建 `sessiontoken.Manager`，注入 Redis FlowStore、provider Authenticator/Harvester/CallbackDispatcher、回调派发器、日志脱敏与指标（`sessiontoken_metric`）。
- **Provider 适配器**：按 `docs/plan/session_token_client.md` 实现各平台的 Authenticator/Harvester（入口 URL、脚本、代理策略、凭证标准化）。
- **回调**：服务在 Flow 成功后向插件传入的 `callback_url`（`/admin/platforms/session-token/callback`）POST 标准化凭证；支持签名校验、重试与日志。
- **配置共享**：插件端通过 `POWERX_SESSION_TOKEN_*` 指向该服务；MediaX 端需保证 BaseURL/API Token/CallbackURL 一致，详见 `docs/plan/creative/channels.md#2.1`。

## 2. OAuth 授权（AuthManager/SessionManager）

- 插件要求 OAuth 账号统一走 MediaX AuthManager/SessionManager，因此 SDK 需暴露可重用的授权 URL 生成、state/session 存储、回调处理能力。
- 建议在 MediaXCore 提供瘦封装：读取 PowerSocialite 配置、自动处理 refresh_token、token 加密、租户隔离，并提供 `auth.GenerateURL`、`auth.HandleCallback`、`session.StoreToken` 等接口供插件调用。
- 需要配套文档/示例（Quickstart）说明如何在 plugin 内调用这些接口，以及如何共享 redis/cache。

## 3. 频道同步 & 任务执行

- 虽然频道刷新/发现主要由插件实现，但 MediaX SDK 可提供辅助能力：
  - 官方账号调用封装（例如 YouTube/Bilibili API client，带 token 管理与配额限流）。
  - 视频解析/下载 SDK（与 video-parser 模块共用）。
  - RPA/Playwright runner（供直链解析/模拟浏览器使用）。
- 若 MediaXCore 侧已有 worker/scheduler，可暴露任务接口让插件提交刷新任务，由 MediaX 执行并回调 `ingest`。

## 4. 配置与部署

- MediaX SDK 应在 `config.yaml` 或等效文件中提供 `session_token`、`oauth`、`provider`、`network/redis` 配置模板，保证与插件默认配置兼容（BaseURL、token、回调地址）。
- 提供 `make sessiontoken` / `go run ./cmd/sessiontoken` 等命令，方便 MediaX Studio 或本地开发启动 SessionToken 服务。
- 在 README/Quickstart 中记录启动顺序：先启动 MediaX SessionToken 服务，再启动 PowerX 插件（指向相同的 BaseURL/token），最后打开 `/publish/platforms` 触发模拟登录。

## 5. 对齐文档

- `docs/plan/session_token_client.md`（MediaX 仓库）定义了详细设计；本文件补充插件侧依赖。后续若接口或部署方式有变化，需同步更新：
  1. MediaX 仓库的 plan/spec 文档；
  2. 本仓库的 `docs/plan/creative/channels.md` 和 `specs/003-channels-subscription`。

---

如需进一步拆分任务，请参考 `specs/003-channels-subscription/plan.md#additional-notes` 与 `tasks.md` 中关于 SessionToken（T052/T053）的说明。
