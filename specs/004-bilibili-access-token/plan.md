# Implementation Plan: BiliBili AccessToken Debug Flows

**Branch**: `004-bilibili-access-token` | **Date**: 2025-12-18 | **Spec**: specs/004-bilibili-access-token/spec.md  
**Input**: Feature specification from `/specs/004-bilibili-access-token/spec.md`

## Summary
扩展 MediaX AccessToken 调试服务以专注 BiliBili Provider：标准化 `/debug` UI 模板、实现 Flow ID 回填与 Redis 可选持久化、统一 `/accesstoken/token` 与 `/accesstoken/oauth/start` 的回溯日志，同时确保默认仅监听 `127.0.0.1` 并将 Flow TTL 固定为 24 小时。计划复用 `MediaXCore` 的 BaseClient、日志与配置设施，新增文档/合约帮助运营与 QA 在 2 分钟内完成授权与 Token 校验。

## Technical Context

**Language/Version**: Go 1.18（仓库 go.mod）  
**Primary Dependencies**: `github.com/ArtisanCloud/MediaXCore`（BaseClient/日志/配置绑定）、`github.com/redis/go-redis/v9`（可选 Flow 缓存）、`gopkg.in/yaml.v3`（解析 access_token_providers）  
**Storage**: 进程内内存 map + 可选 Redis (key 前缀 `accesstoken:oauth:*`, TTL=24h)  
**Testing**: `go test ./cmd/accesstoken/... ./internal/accesstoken/...` 使用 Go `testing` + `testify` helpers；加上文档化的 curl/self-test  
**Target Platform**: 本地或 CI Linux/macOS 实例，运行 `cmd/accesstoken/server`，默认 127.0.0.1:7071  
**Project Type**: Go CLI + HTTP 服务（单仓库，多 cmd/internal 包）  
**Performance Goals**: 服务启动 <2 分钟（SC-001）、OAuth 写入 Flow 成功率 ≥95%（SC-002）、Flow 回填成功率 ≥99%（SC-003）、API 响应 p95 < 300ms 在 10 RPS 以内  
**Constraints**: 所有敏感值来自 env；UI 与 API 必须遮罩 token；默认不暴露公网；Flow 数据仅保留 24h；必须兼容无 Redis 情形并提供可观测日志  
**Scale/Scope**: 单调试节点、并发 <5 人、每日 Flow <200，主要供运营/QA 支持场景

## Constitution Check

- **Provider Adapter Parity**: 仅消费既有 `access_token_providers` 声明与 `pkg/client/bilibili/accessTokenClient`，不新增目录；通过 `cmd/accesstoken/server` 注册 provider factory，保持 parity。
- **Config-Layered Security**: 所有 BiliBili OAuth 参数继续读取 `pkg/client/config/bilibili.go` 的强类型结构；文档强调 `${ENV}` 占位与 `ACCESSTOKEN_LISTEN_ADDR`/`ACCESSTOKEN_API_TOKEN` 的默认值与可配置 TTL，避免明文。
- **Token Lifecycle Discipline**: `/accesstoken/token`、Flow 复用逻辑复用 `kernel.TokenHandler` 与 BaseClient 缓存接口；Refresh/回填时仅通过已有 handler 与 Redis accessor，不绕过缓存或直接暴露 token。
- **Observability & Error Traceability**: 计划要求在 Flow 写入/回填/Redis fallback 时记录 `provider`, `provider_app`, `flow_id`, `token_source`, `listen_addr`，并沿用结构化 logger + 错误 wrap 规范；UI/API 返回带可操作错误。
- **Testable Modularity & SessionToken Readiness**: 将 Flow 存储、Redis 入口、UI handler、token 解析拆分到 `internal/accesstoken/{flow,handler,ui}` 中的可单测函数，新增 `*_test.go` 模拟 env/Redis 组合；SessionToken 约束（状态与脱敏）在 docs/quickstart 中重申。

## Project Structure

```text
specs/004-bilibili-access-token/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
└── contracts/

cmd/
└── accesstoken/
    └── server/          # main.go, CLI flag/env wiring

internal/
└── accesstoken/
    ├── flow/            # Flow persistence, TTL, Redis fallback
    ├── handler/         # HTTP handlers for /accesstoken/*
    ├── ui/              # /debug assets + provider/app templates
    └── config/          # env binding helpers for listen addr, tokens

pkg/
└── client/
    ├── config/          # access_token_providers schema (YAML -> struct)
    └── bilibili/
        └── accessTokenClient/  # OAuth token helpers, reused by debug server

docs/
└── develop/access-token/bilibili/debug.md  # 主调试文档

tests/
└── accesstoken/
    ├── flow/
    └── handler/
```

**Structure Decision**: 采用现有 `cmd + internal + pkg` Go 单仓结构；本特性只在 `cmd/accesstoken/server` 与 `internal/accesstoken/*` 下扩展逻辑与测试，并在 `specs/004-bilibili-access-token/` 追加研究/设计文档。

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | - | - |
