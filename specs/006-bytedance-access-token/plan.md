# Implementation Plan: DouYin AccessToken Debug Integration

**Branch**: `[006-bytedance-access-token]` | **Date**: 2025-12-20 | **Spec**: `specs/006-bytedance-access-token/spec.md`
**Input**: Feature specification + clarifications stored in `specs/006-bytedance-access-token/spec.md`

## Summary

把 DouYin（`provider_code=byte_dance_douyin`）接入 `cmd/accesstoken` 服务：在 `/debug` 暴露 Provider 卡片并完成 OAuth、在 `/accesstoken/call` 里复用 `ByteDanceDouYinACClient` 调试 action，Flow 存储遵循 `accesstoken:oauth:*` 规范（access_token+refresh_token、TTL=`expires_in`、刷新失败立即失效），并补齐限流/日志/文档，确保 CLI、脚本与 UI 共用一套流程。

## Technical Context

**Language/Version**: Go 1.18（`go.mod`）  
**Primary Dependencies**: `github.com/ArtisanCloud/MediaXCore`（BaseClient/Logger/Cache）、`github.com/redis/go-redis/v9`（可选 Flow 缓存）、`pkg/client/byteDance/douYin/accessTokenClient`（DouYin SDK）、`internal/accesstoken/handler/mask`  
**Storage**: 进程内 map + 可选 Redis（key `accesstoken:oauth:byte_dance_douyin:default:<mode>`，flow 索引 `accesstoken:oauth:flow:<flow_id>`）  
**Testing**: `go test ./cmd/accesstoken/...`（含 `oauth_contract_test.go`、`token_contract_test.go`、`flow_contract_test.go`）+ 针对 DouYin 新增用例，以及必要的 `pkg/client/byteDance/douYin/...` 单测  
**Target Platform**: `cmd/accesstoken` HTTP 服务（macOS/Linux dev、容器部署）  
**Project Type**: Go 单仓（CLI + HTTP 服务）；本特性聚焦 `cmd/accesstoken`  
**Performance Goals**: `/accesstoken/call` DouYin action p95 ≤ 2s；速率限制 1 QPS/action + 针对 429/5xx 的 3 次指数退避重试  
**Constraints**: 只能复用既有 Adapter/Factory；单 app（`<app>=default`）；Flow TTL 来自 DouYin `expires_in`；自动刷新失败必须删除 Flow 并提示重新授权；日志字段与错误码需保持与其他 Provider 一致  
**Scale/Scope**: 单租户数十条调试 Flow，低并发但需处理 Redis 不可用降级场景

## Constitution Check（Gate）

- **Provider Adapter Parity** ✅：仅扩展 `pkg/client/byteDance/douYin/accessTokenClient`（refresh/call helper）与 `MediaX.CreateByteDanceDouYinACClient` 注入点，在 `cmd/accesstoken/server` 里调用 `mode.ByteDanceDouYinConfig`，不新增平行客户端或目录。
- **Config-Layered Security** ✅：更新 `pkg/client/config/douyin.go`、`pkg/client/config/config.go`、`config.example.yaml` 与文档，让 DouYin OAuth 字段只引用 `DOUYIN_*` 环境变量，解释 TTL/逐出策略及单 app 限制，仓库不落地明文。
- **Token Lifecycle Discipline** ✅：`handleCall` 检测 Flow TTL，触发 `ByteDanceDouYinACClient.RefreshToken`（复用 BaseClient HTTP）；成功后调用 `saveOAuthTokenRecord` 更新 access/refresh/TTL；刷新失败调用统一 `invalidateFlow` 删除 Redis/内存并返回 `need reauth`。
- **Observability & Error Traceability** ✅：沿用 `app.LogInvocation` 与 `mask.MaskField` 记录 `provider/provider_app/action/flow_id/token_source/retry_count/storage_backend`，新增速率/重试日志字段，多跳错误用 `fmt.Errorf("douyin: %w", err)` 包裹。
- **Testable Modularity & SessionToken Readiness** ✅：列出将补充的单测（OAuth 授权成功/失败、Flow TTL=expires_in、自动刷新成功/失败、速率限制），并保持 Flow 结构字段（`refresh_token`、`flow_expire_at`）可被 SessionToken 流程复用。

## Project Structure

### Documentation（本特性）

```text
specs/006-bytedance-access-token/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
└── tasks.md (后续由 /speckit.tasks 生成)
```

### Source Code（仓库根）

```text
cmd/
└── accesstoken/
    ├── main.go
    ├── internal/
    └── server/
        ├── server.go               # Provider wiring / Flow 生命周期
        ├── handlers.go             # /accesstoken/call, /debug, OAuth handlers
        ├── oauth_contract_test.go
        ├── token_contract_test.go
        └── flow_contract_test.go

internal/
└── accesstoken/
    └── handler/                    # mask、日志字段工具

pkg/
└── client/
    ├── byteDance/
    │   ├── core/
    │   └── douYin/accessTokenClient/
    ├── config/
    └── mediaX.go                   # CreateByteDanceDouYinACClient 工厂

docs/
└── develop/access-token/byteDance/{develop,debug}.md
```

**Structure Decision**: 继续沿用上述 Go 单仓结构，修改集中在 `cmd/accesstoken` 服务、`pkg/client/byteDance/douYin` 与配置/文档目录，不引入新模块或子仓。

## Complexity Tracking

当前方案未触发宪章例外，无需额外条目。

## Implementation Strategy

1. **配置与 Provider 元数据**
   - 在 `pkg/client/config/douyin.go` 定义完整 OAuth 字段、单 app（`app_code=douyin`）提示，以及 refresh/TTL 注释；更新 `pkg/client/config/config.go` 解析逻辑与 `config.example.yaml`、`docs/.../byteDance` 示例，提醒仅支持 `default` 模式（FR-009）。
   - `cmd/accesstoken/server/server.go` 的 `initProviders`/`renderProviderMeta` 中补充 DouYin 卡片内容：`provider_code`、`api_version`、`oauth_key`、`config_path`。

2. **OAuth 授权链路**
   - `handleOAuthStart`/`buildOAuthAuthorizeURL` 检查 DouYin 必填字段（client_id/secret/scope/oauth_url/access_token_url/redirect_url），缺失时给出描述性错误（FR-001~FR-003）。
   - `handleDebugCallback`/`completeOAuthFlow` 将 DouYin exchange 结果写入 `oauthTokenRecord`：保存 `refresh_token`、`expires_in`、`FlowTTLSeconds=expires_in`，日志记录 `provider=byte_dance_douyin`。

3. **Flow 保存 & 自动刷新**
   - 在 `saveOAuthTokenRecord` 之前调用 `enrichFlowMetadata` 补 `refresh_token` 字段；新增 `maybeRefreshDouYinToken(rec *oauthTokenRecord)`：
     - 当 `TokenExpireAt-Now ≤ 300s` 或 `AccessToken` 已过期时，用 `MediaX.CreateByteDanceDouYinACClient(ctx.Douyin)` 的 refresh API 获取新 token。
     - 成功：更新 `rec.AccessToken/RefreshToken/ExpiresIn/FlowTTLSeconds`，写回缓存与 Redis，记录 `event=token.refresh`。
     - 失败：调用 `invalidateFlow(rec, reason)`（删除内存/Redis + log），在 `/accesstoken/call` 返回 `need reauth`（FR-012）。

4. **`/accesstoken/call` 执行 DouYin action**
   - `handleCall` 增加 DouYin 分支：确保 `opts.Validate()`（必填 action/provider），在执行前调用上一步自动刷新。
   - 新建 `executeDouyinCall(ctx context.Context, w http.ResponseWriter, pctx *providerContext, opts *app.Options, payload map[string]any, storedRec *oauthTokenRecord)`：
     - 维护 `map[string]*rate.Limiter`（key=`action`）实现 1 QPS；命中限流时返回 429 提示。
     - 调用 `ByteDanceDouYinACClient` 的 `Call(action string, payload map[string]any)`（若缺失则补足），对 429/5xx 响应重试 3 次，指数退避（例如 200ms→400ms→800ms），并在日志中记录 `retry_count/backoff_ms`。
     - 响应结构包含 `provider/provider_app/action/flow_id/token_source`，敏感字段使用 `mask.MaskField`。

5. **Flow 列表与回填**
   - `handleListFlows`、`handleListOAuthTokens` 支持 `provider_code=byte_dance_douyin` 过滤，展示 `refresh_token`、`flow_expire_at`、`storage_backend`，回填 JSON 可直接复制到 CLI。
   - `flowIndex` API 在回填 DouYin Flow 时写出 `token_source`=`flow`、`token_source_detail`=`flow_id`，满足 FR-006。

6. **文档 & Quickstart**
   - `docs/develop/access-token/byteDance/{develop,debug}.md`：写入环境变量表（DOUYIN_CLIENT_ID/SECRET/SCOPE/REDIRECT_URL/OAUTH_URL/TOKEN_URL）、授权步骤、速率限制/刷新失败排障。
   - `quickstart.md`：一步步演示配置→授权→`/accesstoken/call`→Flow 回填，附日志/Redis 示例。

7. **测试计划**
   - `cmd/accesstoken/server/oauth_contract_test.go`: mock DouYin token endpoint，断言缺失字段报错与 Flow TTL=expires_in。
   - `token_contract_test.go`: 构造 near-expiration Flow，验证自动刷新成功路径与失败后 Flow 删除/返回 need reauth。
   - `flow_contract_test.go`: 验证 `/api/oauth/tokens` 新增字段（refresh_token/ttl）与过滤。
   - 速率限制/重试逻辑：通过 fake client 或注入接口测试 1 QPS 与 3 次退避行为。

8. **Agent Context**
   - 计划完成 Phase 1 后运行 `.specify/scripts/bash/update-agent-context.sh codex`，把新增 DouYin 技术信息同步至 Agent 记忆，供后续自动化参考。

9. **Success Criteria Validation**
   - 为 SC-001~SC-004 准备脚本/演练：统计 `/debug` OAuth 成功率、Redis/内存 Flow 回填命中率、`/accesstoken/call` DouYin action p95 性能、双人文档走查反馈，并在 README/Docs 中记录方法与结果，确保上线前指标可量化复验。

以上产物准备好后，可进入 `/speckit.tasks` 拆解实现任务。
