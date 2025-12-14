# Implementation Plan: Google YouTube AccessToken Enablement

**Branch**: `002-youtube-access-token` | **Date**: 2025-03-07 | **Spec**: specs/002-youtube-access-token/spec.md
**Input**: Feature specification from `/specs/002-youtube-access-token/spec.md`

## Summary

本特性需要把 `google_youtube_config` 模板、`cmd/accesstoken` CLI 以及 Playground 示例整理成统一的开发/调试体验。计划包括：
1. 更新配置文档与示例，确保 `config.yaml`/环境变量都能驱动 AccessToken 客户端。
2. 设计并实现 CLI（action/参数/错误处理）与 Playground 演练流程，形成快速验证路径，并提供“订阅 → 视频 → 发布 → 评论”闭环用例。
3. 产出 Quickstart、数据模型与 API/CLI 契约，便于后续在 README 或 docs 中引用。

## Technical Context

**Language/Version**: Go 1.18（仓库 `go.mod` 已指定）  
**Primary Dependencies**: `github.com/ArtisanCloud/MediaXCore`（BaseClient/Logger）、`github.com/redis/go-redis/v9`（可选 Token 缓存）、`gopkg.in/yaml.v3`（配置解析）  
**Storage**: Redis（用于 AccessToken cache，可选；本特性主要读取本地 YAML + 环境变量）  
**Testing**: `go test ./...`（覆盖 CLI 参数解析与 helper），手动 `make accesstoken`/`go run ./main.go` 验证  
**Target Platform**: Linux/macOS 开发机、CI 环境（CLI + Go 工程）  
**Project Type**: 单一 Go 仓库（含 `cmd/`, `pkg/`, `docs/`）  
**Performance Goals**: CLI 返回时间 <5s（取决于 Google API）；日志输出不阻塞；配置加载常量时间  
**Constraints**: 不得暴露敏感 token；重用 BaseClient 日志；CLI 不引入额外网络依赖；遵守宪章约束  
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
