# Implementation Plan: Google YouTube AccessToken Enablement

**Branch**: `002-youtube-access-token` | **Date**: 2025-03-07 | **Spec**: specs/002-youtube-access-token/spec.md
**Input**: Feature specification from `/specs/002-youtube-access-token/spec.md`

## Summary

本特性需要把 `google_youtube_config` 模板、`cmd/accesstoken` CLI、Playground 示例与全新的 Web 调试服务整理成统一的开发/调试体验。计划包括：
1. 更新配置文档与示例，确保 `config.yaml`/环境变量都能驱动 AccessToken 客户端。
2. 设计并实现 CLI（action/参数/错误处理）与 Playground 演练流程，形成快速验证路径，并提供“订阅 → 视频 → 发布 → 评论”闭环用例。
3. 交付 AccessToken 调试服务（默认端口 `:7070`），页面体验复用 SessionToken 调试台，可在浏览器中完成授权回调与 API 调用。
4. 产出 Quickstart、数据模型与 API/CLI/调试服务契约，便于后续在 README 或 docs 中引用。

## Technical Context

**Language/Version**: Go 1.18（仓库 `go.mod` 已指定）  
**Primary Dependencies**: `github.com/ArtisanCloud/MediaXCore`（BaseClient/Logger）、`github.com/redis/go-redis/v9`（可选 Token 缓存）、`gopkg.in/yaml.v3`（配置解析）  
**Storage**: Redis（用于 AccessToken cache，可选；本特性主要读取本地 YAML + 环境变量）  
**Testing**: `go test ./...`（覆盖 CLI 参数解析与 helper），手动 `make accesstoken`/`go run ./main.go` 验证  
**Target Platform**: Linux/macOS 开发机、CI 环境（CLI + Go 工程）  
**Project Type**: 单一 Go 仓库（含 `cmd/`, `pkg/`, `docs/`）  
**Performance Goals**: CLI 返回时间 <5s（取决于 Google API）；日志输出不阻塞；配置加载常量时间  
**Constraints**: 不得暴露敏感 token；重用 BaseClient 日志；CLI 不引入额外网络依赖；新增调试服务需沿用 SessionToken 的 HTTP/HTML 结构并遵守宪章约束  
**Scale/Scope**: 文档 + CLI 工具 + 示例；使用者规模预估为团队内 10~20 名开发/运维

## Constitution Check

- **Provider Adapter Parity**: 仅触及现有 `pkg/client/google/youtube/accessTokenClient` 与 `client/mediaX.go` 的工厂，计划中不会新增新目录结构。CLI/文档均引用该适配器，确保合规。✅
- **Config-Layered Security**: `google_youtube_config` 继续位于 `pkg/client/config/google.go`，计划中所有示例都强调使用环境变量/Secret Store 注入敏感值，禁止写死 token。✅
- **Token Lifecycle Discipline**: 文档会解释 `GoogleAccessTokenHandler`、`kernel.BaseClient`、`GetCustomToken` 的关系；CLI 通过回调注入 token 并允许缓存 TTL，自始至终复用既有机制。✅
- **Observability & Error Traceability**: 计划中 CLI 与 Playground 都复用 `MediaXCore/pkg/logger`，要求输出 `provider=google`、`action`, `channel_id/ids` 摘要，并提示脱敏 AccessToken。✅
- **Testable Modularity & SessionToken Readiness**: CLI 参数解析与文档逻辑属于纯 Go 代码，可加 `cmd/accesstoken/main_test.go`；不影响 SessionToken 流，但会在计划中列出需新增的测试与 Quickstart 步骤。✅

## Project Structure

### Documentation (this feature)

```text
specs/002-youtube-access-token/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── cli.md
└── tasks.md (由 /speckit.tasks 生成)
```

### Source Code (repository root)

```text
cmd/
├── accesstoken/           # CLI 入口（新增，可测试）
└── sessiontoken/

pkg/
├── client/
│   ├── config/            # google_youtube_config 定义
│   └── google/youtube/... # AccessToken 客户端
└── ...

docs/
├── develop/access-token/google/  # 调试/开发文档
└── plan/google/                  # 功能介绍

specs/002-youtube-access-token/   # 本特性文档
```

**Structure Decision**: 使用单一 Go 仓库结构；新增/修改集中在 `cmd/accesstoken`、`docs/develop/...` 与 `config*.yaml` 示例，不引入额外项目。 

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| _None_ | | |

## Additional Scope: AccessToken 调试服务

为满足最新需求，需要在现有 CLI/Playground 之外新增一套 Web 调试服务，沿用 SessionToken 架构：

- **入口形式**：`cmd/accesstoken/server`（或 `cmd/accesstoken/main.go` 中新增 `serve` 子命令）读取 `config.yaml`，默认监听 `:7070`，通过 `SESSIONTOKEN` 同款 `registerDebugPage` 思路提供 HTML 页面。
- **核心路由**：
  1. `/debug/accesstoken`：静态页面，包含 Provider/App/API 版本选择、AccessToken/RefreshToken 操作、API 调试表单、回调日志展示。
  2. `/accesstoken/token`：POST，根据 `oauth_key` 或自定义凭证刷新 AccessToken，底层调用 `GoogleYouTubeACClient` 的 `GetOAuthToken` 回调。
  3. `/accesstoken/call`：POST，接收 action/参数，复用 CLI 的 `execVideosList/SearchList/PlaylistsList` 逻辑返回 JSON。
  4. `/debug/callback`：用于 OAuth 回调录制，与 SessionToken 共用 `callbackLogStore`。
- **共享能力**：沿用 `MediaX.New` + `cache.NewMemoryCache/RedisCache`、`loggerconfig.LogConfig`、`resolveConfigPath` 等 helper，减少重复实现；API Token 验证逻辑可参照 `sessionhandler.RegisterSessionTokenFlow*`.
- **交互要求**：页面需显示 token 来源（flag/env/config/refresh）、最近回调、请求/响应摘要，并允许复制 curl 命令，实现“本地沙盒”体验。

后续任务会在 `tasks.md` 中补充，具体实现时需参考 SessionToken 现有的 debug 页面代码结构（`cmd/sessiontoken/debug_page.go`）。
