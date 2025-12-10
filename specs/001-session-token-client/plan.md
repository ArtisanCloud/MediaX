# Implementation Plan: SessionTokenClient Provider Architecture

**Branch**: `001-session-token-client` | **Date**: 2025-12-09 | **Spec**: [specs/001-session-token-client/spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-session-token-client/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

MediaX 需要把 SessionToken 流程升级为 provider 级适配器。此次迭代将：

- 拆解 `pkg/client/sessionToken`，提供 `SessionTokenClient`、`sessionToken.Manager` 与 Flow 状态机，统一接入 `kernel.BaseClient`。
- 构建 `cmd/sessiontoken` HTTP 服务与 `make sessiontoken` 命令，默认监听 `:7070` 并注册 `/session-token/flows`，使插件可通过 `POWERX_SESSION_TOKEN_*` 指向统一 BaseURL/API Token。
- 为知乎实现首个 SessionToken 适配器（Authenticator/Harvester/CallbackDispatcher），并在 `pkg/client/mediaX.go` 中注册创建函数。
- 重构配置体系：每个 provider (`pkg/client/config/<provider>.go`) 新增 `SessionTokenConfig`，并在 `config.yaml` 里示例如何注入。
- 构建 Redis 持久化 Flow、回调签名与重试、日志脱敏以及 CLI/Quickstart 样例，满足插件侧 Flow 创建、轮询和凭证回传的 SLA。
- 在 API 层提供 Bearer Token 验证中间件与性能/成功率指标采集，确保 FR-008 及性能目标可被验证。
- 按 `docs/plan/mediax-sdk.md` 同步 README/Quickstart/`docs/plan/creative/channels.md`，说明插件如何共享 SessionToken 服务、如何联动 OAuth/AuthManager 与频道任务接口。

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.18（仓库 `go.mod`）  
**Primary Dependencies**: `github.com/ArtisanCloud/MediaXCore`（日志/HTTP/BaseClient）、`github.com/redis/go-redis/v9`（Flow 持久化）、`gopkg.in/yaml.v3`（配置解析）  
**Storage**: Redis（Flow/状态机共享存储，包含 TTL/索引）  
**Testing**: `go test ./...`；新增单元测试覆盖状态机、配置加载、回调签名，必要处使用 `httptest`/fake transport  
**Target Platform**: Linux/macOS 开发环境 + Linux 容器/服务器运行（MediaX SDK）  
**Project Type**: Go monorepo SDK，provider 适配器 + factory  结构  
**Performance Goals**: Flow 创建 p95 < 200ms；Flow 查询 p95 < 100ms；Zhihu Flow 成功率 ≥95% 且 5 分钟内完成；回调签名 100% 通过  
**Constraints**: 必须脱敏日志、HMAC 签名、Redis 缓存 TTL、`Authorization: Bearer` 校验；禁止硬编码 secret；Flow 完成后仍需在审计窗口内可读  
**Scale/Scope**: 首期面向单 provider（知乎）与多租户插件；架构需可扩展到其他 provider，Flow 量级预期数千级并发

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Provider Adapter Parity**: Plan MUST document provider/product directory layout, required `core/client.go` location, and factory injection so no team invents a new topology.
- 满足方式：`pkg/client/sessionToken` 承载共享 manager/Flow 模型；`pkg/client/zhihu/core/client.go` + `sessionToken/authenticator`、`harvester`、`callback` 负责产品细节。在 `pkg/client/mediaX.go` 增加 `CreateZhihuSessionTokenClient`，并在 README/quickstart 里注明如何注入。新增 provider 时沿用同一层级结构。
- **Config-Layered Security**: Plan MUST enumerate configs added or modified under `pkg/client/config` plus how secrets/TTL are sourced without hardcoding.
- 满足方式：在 `pkg/client/config/zhihu.go` 定义 `ZhihuSessionTokenConfig`（含 `Service/Auth/Harvester/Callback/Network` 五块），并在 `config.yaml`/README 中展示 env var 引用。Flow TTL、callback secret、代理池参数全部声明在配置体内，避免散落常量。
- **Token Lifecycle Discipline**: Plan MUST explain how BaseClient helpers, refresh hooks, and cache interfaces are reused (or explicitly extended) for Access/Client/Session tokens.
- 满足方式：`sessionToken.Manager` 使用 `kernel.BaseClient` 提供的 HttpHelper & Logger；Zhihu adapter 仅封装特定 payload，HTTP 调用仍走 BaseClient；Redis 缓存通过现有 cache interface 注入，无自建 client。
- **Observability & Error Traceability**: Plan MUST specify structured logging fields, retry instrumentation, and desensitization strategy for the feature scope.
- 满足方式：所有 Flow/回调日志均由 `client.Logger.WithContext` 输出，字段含 `provider/api/tenant_uuid/account_id/flow_id/retry`。凭证脱敏由 `sessionToken/sanitizer` 辅助函数统一执行，回调请求记录签名校验结果与重试次数；API 层强制 Bearer Token 校验并输出授权失败日志，同时记录 Flow 创建/查询 latency 指标以验证 SLA。
- **Testable Modularity & SessionToken Readiness**: Plan MUST call out the pure-logic components that will ship with unit tests and any SessionToken flow/state-machine coverage.
- 满足方式：新增 `sessionToken/manager_test.go` 覆盖 pending→authorizing→succeeded/failed；`sessionToken/callback/signature_test.go` 覆盖 HMAC；`config/zhihu_sessionToken_test.go` 确保配置解析；Zhihu adapter 内部逻辑通过 fake Authenticator/Harvester 单测。文档在 quickstart 中列出如何运行 `go test ./pkg/client/sessionToken/...`。

## Project Structure

### Documentation (this feature)

```text
specs/001-session-token-client/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── session-token-flows.yaml
└── tasks.md (由 /speckit.tasks 生成)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
pkg/
├── client/
│   ├── mediaX.go                    # 工厂：新增 CreateZhihuSessionTokenClient
│   ├── config/
│   │   ├── zhihu.go                 # 添加 ZhihuSessionTokenConfig + yaml/json tag
│   │   └── sessionToken.go          # （若需要）公共 SessionToken 配置片段
│   ├── sessionToken/
│   │   ├── manager.go               # Flow 状态机 & API 编排
│   │   ├── flow.go / storage/redis/ # Flow 模型与 Redis 持久化
│   │   ├── callback/dispatcher.go   # HMAC 签名 & 重试
│   │   └── sanitizer/
│   └── zhihu/
│       └── sessionToken/
│           ├── core/client.go       # 复用 BaseClient
│           ├── authenticator/
│           ├── harvester/
│           └── callback/
├── utils/...

docs/
└── plan/session_token_client.md     # 现有背景文档

cmd/
└── sessiontoken/
    └── main.go                      # 独立 HTTP 服务入口，注册 SessionToken 路由

Makefile
└── sessiontoken                     # 启动命令：`make sessiontoken`/`go run ./cmd/sessiontoken`
```

**Structure Decision**: 使用现有 Go monorepo；SessionToken 共享逻辑集中在 `pkg/client/sessionToken`，provider 实现放在 `pkg/client/<provider>/sessionToken/`，与宪章要求的 `core/` + adapter 子目录保持一致；配置仍位于 `pkg/client/config` 并通过工厂注入。

## MediaX SDK 对齐事项

- **SessionToken 服务**：遵循 `docs/plan/mediax-sdk.md#1` 要求，`cmd/sessiontoken/main.go` 需在启动时加载配置、注册 `/session-token/flows`，并输出 `sessiontoken_metric/sessiontoken_callback` 日志。
- **OAuth/AuthManager 触点**：在 plan 与 README 中指向 MediaX AuthManager/SessionManager 的复用方式，说明插件如需 OAuth 账号仍可通过 SDK 共享 state/session 存储。
- **频道任务依赖**：文档需补充 SessionToken Flow 完成后如何向 MediaX worker/scheduler 提交刷新任务，保持与 “频道同步 & 任务执行” 章节一致。
- **配置/部署**：`config.yaml`、`config.example.yaml`、Quickstart与 `make sessiontoken` 需要展示 `POWERX_SESSION_TOKEN_*` 的默认值和启动顺序（先 MediaX SessionToken 服务再启动插件）。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
