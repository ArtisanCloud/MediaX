# Implementation Plan: RedBook JuGuang AccessToken Debug Integration

**Branch**: `[005-redbook-access-token]` | **Date**: 2025-12-19 | **Spec**: specs/005-redbook-access-token/spec.md
**Input**: Feature specification from `/specs/005-redbook-access-token/spec.md`

## Summary
- 在已有 `cmd/accesstoken/server` 与 `/debug` 页面中注册“小红书 / 聚光” Provider，沿用 Google/BiliBili 的 OAuth/Flow 管线。
- 扩充 `config.yaml` / `config.example.yaml` 的 `redbook_juguang_config`（补齐 `oauth_url`,`access_token_url`,`scope`），并在文档中提供 Quickstart。
- 增加一个最小聚光业务接口（账户余额）作为 API 调试示例，提供新的 `/accesstoken/redbook/account/balance` endpoint 与 OpenAPI 合同。
- 加强日志脱敏与 Flow 持久化，确保 Redis/内存双模式下都能回填 Flow。

## Technical Context

**Language/Version**: Go 1.21 (仓库 go.mod)  
**Primary Dependencies**: MediaXCore BaseClient、`github.com/redis/go-redis/v9`、聚光 AccessToken 客户端  
**Storage**: Redis（首选）+ 进程内内存 fallback，用于 Flow 与调试记录  
**Testing**: `go test ./cmd/accesstoken/server` + 针对聚光 handler 的 contract tests  
**Target Platform**: Linux/macOS 开发机上运行的 CLI/HTTP 服务  
**Project Type**: 单一 Go 仓库（CLI + HTTP Server + 文档）  
**Performance Goals**: 调试 API 本地 p95 ≤ 2s；OAuth Flow 写入 Redis 成功率 ≥ 95%  
**Constraints**: 遵守统一脱敏日志、禁止 hardcode 凭证；保持 `/debug` UI 行为一致  
**Scale/Scope**: 单人可完成的增量（新增 Provider 元数据 + 1 个 API Endpoint + 2 篇文档）

## Constitution Check

| 原则 | 设计落实 |
| --- | --- |
| Provider Adapter Parity | 仅复用 `pkg/client/redBook/juGuang/accessTokenClient` 与 `MediaX.CreateRedBookJuGuangACClient`，在 `cmd/accesstoken/server` 的 provider 列表中注册，不创建新拓扑。 |
| Config-Layered Security | 在 `pkg/client/config/redbook.go` / `config.example.yaml` 扩充 OAuth 字段，文档要求通过 env 注入 `client_secret`/`scope`，禁止明文。 |
| Token Lifecycle Discipline | 复用 BaseClient + TokenHandler，Flow 写入沿用现有缓存键（`accesstoken:oauth:*`），不得绕过缓存接口。 |
| Observability & Error Traceability | OAuth/API 调用日志记录 `provider`, `provider_app`, `flow_id`, `token_source`, `config_path`，并借助 `internal/accesstoken/handler/mask` 脱敏。 |
| Testable Modularity & SessionToken Readiness | 新增 `oauth_contract_test.go` / `flow_contract_test.go` 聚光用例，确保纯逻辑代码可单测；SessionToken 状态机无变更，仅在调试服务添加 contract tests。 |

> Gate 通过：所有原则均已纳入设计约束，Phase 1/2 需再次自检以防偏离。

## Project Structure

### Documentation (this feature)

```text
specs/005-redbook-access-token/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── accesstoken_redbook.yaml
└── tasks.md   # 由 /speckit.tasks 生成
```

### Source Code (repository root)

```text
cmd/
└── accesstoken/
    └── server/
        ├── handlers.go          # 注册聚光 Provider、API endpoint
        ├── server.go            # Provider 初始化、Flow 存储
        └── static/              # /debug 前端模板资源

docs/
└── develop/access-token/redbook/
    ├── develop.md
    └── debug.md

pkg/
└── client/
    └── redBook/juGuang/
        ├── accessTokenClient/
        └── config/             # `redbook_juguang_config`
```

**Structure Decision**: 继续沿用单仓库 Go 项目结构，只在现有目录中新增/修改代码与文档。

## Phase 0: Outline & Research

- 结论详见 `research.md`（复用 cmd/accesstoken、首选账户余额 API、沿用既有配置层级）。
- 当前无未解的 NEEDS CLARIFICATION，Phase 0 完成。

## Phase 1: Design & Contracts

- `data-model.md` 描述 ProviderMetadata / OAuthFlowRecord / RedBookJuGuangConfig 扩展字段。
- `contracts/accesstoken_redbook.yaml` 定义 `/accesstoken/redbook/account/balance` OpenAPI 草案。
- `quickstart.md` 提供统一的配置与调试步骤。
- 已运行 `.specify/scripts/bash/update-agent-context.sh codex`，同步本特性涉及的技术上下文。

## Constitution Re-check

- Provider 注册、配置与日志规范均在计划和数据模型中体现，新文档也强调 env 注入 → 原则仍满足。
- 进一步实现时若新增更多 API，需重复上述检查确保日志/缓存一致。

## Phase 2 (Next Steps)

- 使用 `/speckit.tasks` 将上述工作拆解为编码/测试/文档任务。
- 实施完成后再执行 `/speckit.checklist` 及常规 PR 流程。
