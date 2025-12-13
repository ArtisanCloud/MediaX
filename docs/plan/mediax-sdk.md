# MediaX SDK 对接规划（频道订阅 & 平台账号）

> 目的：集中记录 MediaX SDK/MediaXCore 在频道订阅与平台账号矩阵场景下需要实现/扩展的能力，便于与插件团队对齐依赖与分工。

## 1. 模拟登录（SessionTokenClient）

- **服务进程**：MediaXCore 需提供独立的 SessionToken HTTP 服务（默认监听 `:7070`），在 `cmd/sessiontoken/main.go` 中启动 HTTP server，加载配置并注册 `/session-token/flows` 路由。
- **API Handler**：复用 `server/handlers/session_token` 下的 `RegisterSessionTokenFlowCreateRoute` 与 `RegisterSessionTokenFlowGetRoute`，并在中间件中校验 Bearer Token（默认 `dev-session-token`，可配置）。
- **Flow 管理器**：构建 `sessiontoken.Manager`，注入 Redis FlowStore、provider Authenticator/Harvester/CallbackDispatcher、回调派发器、日志脱敏与指标（`sessiontoken_metric`）。
- **Provider 适配器**：按 `docs/plan/session_token_client.md` 实现各平台的 Authenticator/Harvester（入口 URL、脚本、代理策略、凭证标准化）。
- **回调**：服务在 Flow 成功后向插件传入的 `callback_url`（`/admin/platforms/session-token/callback`）POST 标准化凭证；支持签名校验、重试与日志。
- **配置共享**：插件端通过 `POWERX_SESSION_TOKEN_*` 指向该服务；MediaX 端需保证 BaseURL/API Token/CallbackURL 一致，详见 `docs/plan/creative/channels.md#21-环境变量映射`。
- **默认配置模板**：MediaX 仓库必须维护可直接使用的 `config.example.yaml`（或等效模板），在 `zhihu_config.sessionToken` 等段落中填写“版本配套”的默认入口 URL、脚本 ID、UA、回调 secret/Redis 模板，而不是 `${ZH_*}` 占位符；插件可在启动 SessionToken 服务前自动复制该模板，避免开发者手工配置。模板中已包含 `service.api_version`、`harvester.watch_cookies`、`callback.retry_backoff`、`network.proxy` 等字段。
- **版本化实现**：`pkg/client/zhihu/web/sessionTokenClient/` 入口负责根据 `service.api_version`（或 `SESSIONTOKEN_ZHIHU_API_VERSION`）加载对应版本目录，例如当前默认 `v4/*`；新增版本只需要落地新的目录并在入口注册，调试页元数据会自动带上 `api_version`，日志 (`sessiontoken_api`) 亦输出 `version=...` 供监控。
- **默认策略模板（`zhihu_pc_v1`）**：`config.example.yaml` 已提供 MediaX Studio 同步验证过的默认参数，可直接复制为 `config.yaml` 使用：
  - `service`：`base_url=http://127.0.0.1:7070`、`api_token=dev-session-token`、`timeout=30s`、`http_debug=false`；与 `make sessiontoken` 默认监听一致。
  - `authenticator`：入口指向 `https://www.zhihu.com/signin?next=%2F`，`default_user_agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 14_3) ... Chrome/123.0.6312.58`，`script_ids=[sessiontoken.zhihu.auth.pc.v1]`，`captcha_strategy=auto`，并在顶层显式标识 `strategy: zhihu_pc_v1` 方便插件识别策略版本。
  - `harvester`：默认监听 `z_c0/zst_82` Cookie 与 `X-XSRF-TOKEN` Header，抓取脚本为 `sessiontoken.zhihu.harvest.pc.v1`。
  - `callback`/`network`：`callback_secret=mediax-sessiontoken-callback`、`max_retry=3`、`retry_backoff=[2,4,8]`，`proxy_pool=zhihu-default`、`ip_strategy=china_rotating`、`request_timeout=60s`。
- **模板复制流程**：MediaX Studio/插件应先检测 `config.yaml` 是否存在，不存在则复制 `config.example.yaml`（仓库提供 `make sessiontoken-bootstrap` / `scripts/sessiontoken-bootstrap.sh` 可直接完成），再按环境覆写少量差异（如 `callback_secret`、Redis 地址）。同时提供 `.env.example`，集中列出 `SESSIONTOKEN_*`/`POWERX_SESSION_TOKEN_*` 变量，便于统一管理。如此即可确保“拿到仓库即可直接跑”，版本升级只需同步模板即可。
- **调试设施**：内置 `/debug` 页面支持 Provider → App → API 版本级联、Mock callback ring buffer、全局 `window.createFlow`/`pollFlow` 函数及默认回调地址 `/debug/callback`。外部应用在无插件环境下也能按照文档创建 Flow、登录并回放回调；脚本 `scripts/sessiontoken-debug.sh` 供命令行快速调用 Flow/zhihu API。
- **Metadata 等待**：FlowOrchestrator 在触发 Harvester 前会每 2 秒轮询一次 `metadata.session_token`，默认最长等待 2 分钟，确保浏览器脚本/人工登录有时间写回 Cookie。仅在超时或脚本写入空值时才会失败（日志会输出 `metadata session_token not ready before timeout`），避免“创建 Flow 后立刻失败”的体验。
- **会话复用**：会在 `CompleteFlowSuccess` 时把凭证写入 `sessionToken:reuse:<hash>` 缓存，勾选 `metadata.reuse_session=true`（调试页对应“复用 Cookie”）即可直接复用上次成功的 SessionToken，减少重复登录。

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
- 发布流程需同步更新 `config.example.yaml` 与相关模板：当 provider 登录策略、脚本或回调策略变更时，MediaX 版本必须更新示例配置并在 release note 中说明；插件侧自动复制模板即可得到最新默认值，无需额外手工配置。
- SessionToken 服务需提供 `make sessiontoken-bootstrap`（或等效脚本）检测/生成 `config.yaml`：若文件缺失则复制模板，若仍存在 `${...}` 占位符需 fail-fast 并提示；插件可以在自动启动前调用该脚本，确保所有开发者免填配置即可使用。

## 5. 对齐文档

- `docs/plan/session_token_client.md`（MediaX 仓库）定义了详细设计；本文件补充插件侧依赖。后续若接口或部署方式有变化，需同步更新：
  1. MediaX 仓库的 plan/spec 文档；
  2. 本仓库的 `docs/plan/creative/mediax-sdk.md`、`docs/plan/creative/channels.md`，以及 `specs/003-channels-subscription`。
- `docs/plan/creative/mediax-sdk.md` 记录了 Creative 插件对 MediaX SDK 的能力依赖（配置模板、SessionToken 策略、登录态控制等），供 channels-subscription 任务（如 T054/T055）引用。

---

如需进一步拆分任务，请参考 `specs/003-channels-subscription/plan.md#additional-notes` 与 `tasks.md` 中关于 SessionToken（T052/T053）的说明。
