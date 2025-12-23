# Implementation Plan: DouYin ClientToken Mode Onboarding

**Branch**: `[007-douyin-client-token]` | **Date**: 2025-12-22 | **Spec**: specs/007-douyin-client-token/spec.md
**Input**: Feature specification from `/specs/007-douyin-client-token/spec.md`

## Summary

- 扩展 `cmd/clienttoken/server`、CLI 与脚本，使抖音 DouYin ClientToken 可通过统一 Provider 配置刷新、缓存、调试 API；Redis 为主缓存，内存自动降级。
- 更新 `pkg/client/config/douyin.go`、`config.example.yaml`、`docs/develop/client-token/byteDance/{develop,debug}.md`，确保运维/集成商按文档即可复制流程。
- 新增 CLI/脚本与日志字段，覆盖刷新→缓存→调用→清除的可观测链路，所有迈向 DouYin 的调用都复用 BaseClient 与 `clientTokenClient`。

## Technical Context

**Language/Version**: Go 1.18（仓库 go.mod）  
**Primary Dependencies**: `github.com/ArtisanCloud/MediaXCore`、`github.com/redis/go-redis/v9`、既有 `pkg/client/byteDance/douYin/clientTokenClient`  
**Storage**: Redis（主缓存，key `clientToken:douyin:<client_key>`）+ 进程内内存降级  
**Testing**: `go test ./...`（新增 pkg/client 与 cmd/clienttoken 子包单测 + CLI 集成测试）  
**Target Platform**: Linux/macOS 开发环境 + 部署在 MediaX CLI/调试服务  
**Project Type**: 单一 Go 仓库（cmd + pkg + docs）  
**Performance Goals**: 刷新/调用接口 p95 ≤ 1.5s；`refresh_before_seconds` 触发后 1 次 API 内完成续期  
**Constraints**: 所有凭证通过 env 注入；不得产生新 HTTP 客户端；日志需脱敏；多实例共享 Redis 键，避免缓存抖动  
**Scale/Scope**: 单 Provider / App / Mode（`byte_dance`/`douyin_service`/`default`），面向多租户 CLI/插件共用。

## Constitution Check

- **Provider Adapter Parity**: 方案沿用 `pkg/client/byteDance/douYin/clientTokenClient` + `MediaX.CreateByteDanceDouYinCTClient`，在 `cmd/clienttoken/server` 通过统一 Provider factory 暴露，不新增平行目录。
- **Config-Layered Security**: 计划在 `pkg/client/config/douyin.go` 中扩展 `ByteDanceDouYinConfig.ClientToken` 字段，并在 `config.example.yaml`/文档中仅引用 `${ENV}` 占位；所有秘密（client_key/secret、api_token、Redis）均取自环境变量。
- **Token Lifecycle Discipline**: 刷新逻辑复用 `core.ByteDanceTokenHandler` 与缓存接口；新增 TTL 检测（读取 token 前检查 `refresh_before_seconds`）并写回 Redis/内存。
- **Observability & Error Traceability**: `clienttoken` 日志统一输出 `event, provider, provider_app, action, redis_key, ttl_remaining, token_source`，并使用 `internal/accesstoken/handler/mask` 脱敏；CLI 与 HTTP API 都重用这些字段。
- **Testable Modularity & SessionToken Readiness**: 计划新增 `pkg/client/byteDance/douYin/clientTokenClient` 单测（刷新、TTL 阈值）、`cmd/clienttoken` 集成测试（刷新→调用→清除）, 并确保 SessionToken 模块加载配置不受影响。

*Gate Status*: 所有宪章约束均已在计划中覆盖，无需豁免。

## Project Structure

```text
specs/007-douyin-client-token/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── clienttoken-douyin.yaml
└── checklists/
    └── requirements.md

cmd/clienttoken/
├── server/
│   ├── main.go
│   ├── internal/handler/
│   └── internal/cache/
├── templates/
└── debug/ui/

pkg/client/
└── byteDance/
    └── douYin/
        └── clientTokenClient/
            ├── client.go
            ├── cache.go
            └── logger.go

docs/
└── develop/client-token/
    └── byteDance/
        ├── develop.md
        └── debug.md

scripts/
└── clienttoken-douyin.sh (new)
```

**Structure Decision**: 单体 Go 仓库；核心改动集中在 `cmd/clienttoken`（UI/API）、`pkg/client/byteDance/douYin/clientTokenClient`（刷新与缓存）、`pkg/client/config/douyin.go`（配置）以及文档/脚本目录。

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| *None* | | |
