# Feature Specification: Google YouTube AccessToken Enablement

**Feature Branch**: `002-youtube-access-token`  
**Created**: 2025-03-07  
**Status**: Draft  
**Input**: User description: "根据 docs/develop/access-token/google 下的开发与调试指南，为 Google YouTube AccessToken 客户端产出规范化的 spec 文档，覆盖本地配置、CLI 调试、使用示例等"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 统一配置模板 (Priority: P1)

作为 MediaX SDK 维护者，我需要一套明确的 `google_youtube_config` 模板与变量说明，使任何本地/CI 环境都能快速加载 YouTube OAuth 凭证并复用 AccessToken 调试工具。

**Why this priority**: 没有标准化配置即无法实例化客户端或挂载 CLI，属于所有后续场景的前提。

**Independent Test**: 仅验证配置文档与示例即可——把模板复制到新环境、填入凭证后 `make accesstoken -action videos.list` 能成功返回数据。

**Acceptance Scenarios**:

1. **Given** 新 clone 的仓库，**When** 按文档复制 `config.example.yaml` 并填写 `google_youtube_config` 字段，**Then** CLI 可以解析并使用这些字段发起 API 调用。
2. **Given** 仅提供环境变量（如 `GOOGLE_YOUTUBE_ACCESS_TOKEN`），**When** 跳过 YAML 中的 access_token，**Then** CLI 仍可读取 token 并完成调用。

---

### User Story 2 - CLI 调试体验 (Priority: P1)

作为 SDK 使用者，我希望有一个 `cmd/accesstoken` CLI，能根据命令行参数调 YouTube Video/Search/Playlist 接口，并输出结构化 JSON，方便验证 AccessToken 是否可用以及排查权限问题。

**Why this priority**: CLI 直接承载调试/验证主流程，是开发指南的核心价值。

**Independent Test**: 单独执行 CLI 命令即可完成端到端测试，无需依赖其他组件。

**Acceptance Scenarios**:

1. **Given** 已配置 OAuth 信息并提供 AccessToken，**When** 运行 `make accesstoken ARGS='-action search.list -query test -part snippet'`，**Then** CLI 返回格式化 JSON 结果且日志清晰标识请求参数。
2. **Given** 传入无效 token，**When** 调用 CLI，**Then** 能看到 Google API 的错误信息并在文档中找到故障排查步骤（如 quota、invalid_grant）。
3. **Given** 使用 CLI/Playground 文档提供的“订阅 → 视频 → 发布 → 评论”闭环指引，**When** 依次执行订阅列表、频道视频列表、视频发布、评论获取与回复操作，**Then** 每个操作都能在本地完成调试并输出预期结果。

---

### User Story 3 - Playground 示例验证 (Priority: P2)

作为平台演示者，我需要在 `main.go` 中启用 `playground.PlayGoogleYouTube`，让团队通过同一配置文件快速复现 SDK 调用栈，验证 `GetOAuthToken` 回调、BaseClient 日志等行为。

**Why this priority**: Playground 属于次要渠道，但能帮助了解 SDK 与 CLI 一致性。

**Independent Test**: 仅执行 `go run ./main.go` 并观察日志即可，不依赖 CLI。

**Acceptance Scenarios**:

1. **Given** 在 `main.go` 打开示例，**When** 执行程序，**Then** 控制台打印 `youtube#videoListResponse`，并在 `logs/info.log` 中看到 HTTP 请求/响应。
2. **Given** 调整 `GetOAuthToken` 回调读取外部存储，**When** Playground 运行，**Then** 不需要改 SDK 代码即可替换 token 获取策略。

---

### User Story 4 - AccessToken 调试服务 (Priority: P1)

作为外部应用开发者，我希望像 SessionToken 调试台一样，能够在本地启动一套默认监听 `:7070` 的 Web 服务，页面上可以切换 Provider/App/API 版本、快速复现授权与接口调试，从而在浏览器里模拟真实业务如何消费授权账号。

**Why this priority**: CLI/Playground 对非工程人员门槛较高，缺少“沙盒式”交互会阻碍演示与联调；统一的调试服务能沿用 SessionToken 的交互范式并支持 OAuth 回调观察，是交付给产品/QA/合作伙伴的必要能力。

**Independent Test**: 只需启动服务并访问 `/debug/accesstoken` 页面即可独立验证，不依赖 CLI/Playground。

**Acceptance Scenarios**:

1. **Given** 在仓库根目录执行 `make accesstoken-serve`（或 `go run ./cmd/accesstoken/server`），**When** 进程启动后，**Then** 控制台输出 `listening on :7070`，浏览器访问 `http://127.0.0.1:7070/debug/accesstoken` 能看到 Provider/App/API 版本的选择控件。
2. **Given** 调试页选择 `provider=google`、`app=youtube.default` 并沿用配置里的 `oauth_key`，**When** 点击“生成授权链接/刷新 AccessToken”，**Then** 服务会调用与 CLI 相同的逻辑输出 JSON（含脱敏 `token_source`），同时 `/debug/callback` 记录 OAuth 回调 Body 与 Query。
3. **Given** 在调试页配置 `videos.list`／`search.list` 参数，**When** 点击“调用 API”，**Then** 服务器通过 `GoogleYouTubeACClient` 发起请求并返回结构化响应，同时页面展示请求日志、耗时与错误信息，方便 QA 直接验证。

---

### Edge Cases

- AccessToken 缺失或过期：文档需指引通过环境变量/刷新 token 解决，并提示清理 Redis `mediax.access_token.*` 缓存。
- 网络受限或需要代理：需说明 `proxy_api_url` 的配置位置及 CLI 如何继承。
- CLI 参数冲突：明确 `-mine`、`-channel-id`、`-ids` 互斥关系以及缺失 `part` 时的报错行为。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 文档 MUST 提供完整的 `google_youtube_config` 字段说明（API URL、timeout、http_debug、oauth block、可选 AccessToken/RefreshToken、oauth_key），并指向示例文件。 
- **FR-002**: 文档 MUST 描述如何使用环境变量覆盖配置（如 `MEDIA_X_CONFIG`, `GOOGLE_YOUTUBE_ACCESS_TOKEN`）以及哪些字段支持覆盖。 
- **FR-003**: CLI MUST 支持 `videos.list`, `search.list`, `playlists.list` 三种 action，且参数映射（`part`, `ids`, `query`, `channel-id`, `max-results`, `page-token`, `mine`, `search-mine`, `search-type`）。 
- **FR-004**: CLI MUST 能接受命令行 AccessToken 或配置文件中的 token/refresh_token，并在输出中提供 JSON 响应或报错信息。 
- **FR-005**: 文档 MUST 提供故障排查指南（invalid_grant、quotaExceeded、insufficientPermissions、401/403）和推荐的清理/重试步骤。 
- **FR-006**: 文档 MUST 指导如何通过 `playground.PlayGoogleYouTube` 与 CLI 互相验证，保证两个入口共享配置与日志策略。 
- **FR-007**: 文档 MUST 强调敏感配置的管理方式（仅本地/Secret Manager），并提醒在提交前移除临时 token。 
- **FR-008**: 文档 MUST reference `docs/plan/google/youtube_access_token_client.md` 以说明可调用的子客户端列表和能力范围。 
- **FR-009**: 文档 MUST 描述“订阅 → 视频 → 发布 → 评论”的完整示例流程，列出各步骤对应的 CLI 命令或 Playground API，确保外部应用能据此调试闭环能力。
- **FR-010**: 新增的 AccessToken 调试服务 MUST 默认监听 `:7070`，支持从 `config.yaml`/环境变量读取 `google_youtube_config`，并提供 API Token（默认 `dev-accesstoken`）校验。
- **FR-011**: 调试服务 MUST 提供 HTML 调试页（Provider/App/API 选择、授权 URL/AccessToken 操作、API 参数输入）以及 REST API（如 `/accesstoken/token`、`/accesstoken/call`、`/debug/callback`），功能与 SessionToken 调试台一致，所有请求都复用 `GoogleYouTubeACClient`。
- **FR-012**: 调试服务 MUST 记录回调日志（时间、FlowID、Query、Headers、Body），以列表形式展示在页面上，并允许用户清空记录；日志需脱敏 token。

### Key Entities *(include if feature involves data)*

- **GoogleYouTubeConfig**: 描述 API URL、超时、调试开关、OAuth 凭证、AccessToken/RefreshToken、oauth_key；与 `config.yaml`、环境变量以及 CLI 参数映射。 
- **AccessToken CLI Invocation**: 由 action、核心参数、AccessToken 来源构成，输出 JSON 结果或错误；记录日志并遵循代理/超时配置。 
- **Playground Execution Context**: 使用 `MediaX.CreateGoogleYouTubeACClient`、`GetOAuthToken` 回调和 `BaseClient` 日志；与 CLI 共享配置、用于展示完整 SDK 行为。
- **AccessToken Debug Service**: 常驻监听 `:7070` 的 HTTP 服务，集成 HTML 调试页、REST 接口与 `/debug/callback` 日志，可供产品/QA 在浏览器中选择 Provider/App/API 并发起授权或 API 调用，背后依然调用 `GoogleYouTubeACClient`。

## MediaX Architecture Guardrails *(must reference Constitution sections)*

- **Provider Adapter Parity**: 该文档聚焦 Google/YouTube AccessToken 客户端，明确所有调用都走 `pkg/client/google/youtube/accessTokenClient/*`，保持与 `MediaX.CreateGoogleYouTubeACClient` 工厂的绑定，避免引入额外 provider 实现。 
- **Config-Layered Security**: `google_youtube_config` 需继续使用 `pkg/client/config` 中的结构体，Yaml/JSON 标签不可变化；凭证来源（client_secret、refresh_token、access_token）需通过环境变量或 Secret Store 注入，禁止硬编码。 
- **Token Lifecycle Discipline**: 文档需说明 AccessToken 如何注入 `core.GoogleAccessTokenHandler` → `kernel.BaseClient`，以及 `GetCustomToken` 回调/缓存 TTL 的使用方式，确保不会绕过既有刷新流程。 
- **Observability & Error Traceability**: CLI 与 Playground 均需复用 `MediaXCore/pkg/logger`，输出 provider、action、参数摘要与 Google 错误 reason，遵守日志脱敏（隐藏 token）。 
- **Testable Modularity & SessionToken Readiness**: CLI 用例与 Playground 示例都不影响 SessionToken 状态机，但需提供可独立运行的测试/演练步骤，确保未来可添加单元测试验证参数解析逻辑。 

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 新成员在 30 分钟内可按照文档完成本地配置并通过 CLI 拉取至少一个视频列表（人工演练反馈 ≥90% 成功率）。 
- **SC-002**: CLI 调用覆盖视频/搜索/播放列表三种 action，文档列举的示例命令全部可直接复制执行并返回预期响应。 
- **SC-003**: 常见错误（invalid_grant、quotaExceeded、insufficientPermissions、401/403）均提供排查路径，内部支持团队可在 10 分钟内定位问题，无需升级。 
- **SC-004**: playground 演练步骤可在 15 分钟内完成，且日志能展示 AccessToken 注入→HTTP 请求→响应的完整链路。 
- **SC-005**: 调试服务可在 5 分钟内完成启动与页面访问，产品/QA 不需要 Go/CLI，即可在浏览器上完成授权回调、AccessToken 注入与 API 调用（至少覆盖 `videos.list`/`search.list`/`playlists.list`），并在页面上看到最新回调记录。

## Assumptions & Dependencies

- 开发者已有可用的 Google Cloud 项目与 OAuth Client（桌面应用或 Web）。
- `config.example.yaml` 已同步至最新，包含 `google_youtube_config`。 
- 网络能够访问 `https://www.googleapis.com`；若需代理，相关配置在 CLI/Playground 中自动加载。 
- Redis 仅在需要缓存 token 时使用，CLI 默认为内存缓存即可满足调试（不强制依赖）。 
- 调试服务将借鉴 `cmd/sessiontoken` 的架构（Redis/HTTP/Debug Page），并尽量与现有 CLI 共享 helper，避免重复实现 AccessToken 注入逻辑。
