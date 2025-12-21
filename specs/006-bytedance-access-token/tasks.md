# Tasks: DouYin AccessToken Debug Integration

**Input**: Design documents from `/specs/006-bytedance-access-token/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

## Phase 1: Setup (Shared Infrastructure)

目的：准备所有用户故事都会依赖的配置/示例。

- [x] T001 扩展 DouYin 配置结构体，补齐 OAuth URL、client_id/secret、scope、redirect、TTL 注释与 yaml/json tag（pkg/client/config/douyin.go）
- [x] T002 将 `byte_dance_douyin_config` 正式挂载到 AccessTokenProviders：解析 YAML、默认 `<app>=default`、`FirstSelection` 与 MediaX 工厂（pkg/client/config/config.go）
- [x] T003 [P] 在 `config.example.yaml` 的 `access_token_providers.byte_dance.douyin` 节点添加 DouYin 示例与 `DOUYIN_*` 环境变量占位（config.example.yaml）

---

## Phase 2: Foundational (Blocking Prerequisites)

目的：为所有故事提供 Provider 元数据、缓存模型与 SDK 钩子。完成前不得进入用户故事。

- [x] T004 注入 DouYin 的 OAuth env key 映射，确保 `/debug` 能根据 `app.EnvKeysForProvider` 读取 AccessToken（cmd/accesstoken/internal/app/providers.go）
- [x] T005 在 `initProviders`/`providerMeta`/`providerContext` 中注册 DouYin 卡片信息（code、name、api_version、oauth_key、config_path），并在缺少 `byte_dance_douyin_config` 时直接禁用/隐藏该 Provider，保证 UI/日志只在配置完整时才暴露入口（cmd/accesstoken/server/server.go）
- [x] T006 扩展 `oauthTokenRecord` 及 Flow 缓存结构，加入 `refresh_token`、`flow_ttl_seconds=expires_in`、`status`、`last_refresh_at`、`invalidateFlow` 等元数据（cmd/accesstoken/server/server.go）
- [x] T007 [P] 在 DouYin SDK 中补充 `Call`/`RefreshToken` 辅助方法，复用 `MediaX.CreateByteDanceDouYinACClient` 并暴露错误包装（pkg/client/byteDance/douYin/accessTokenClient/client.go）

---

## Phase 3: User Story 1 - `/debug` 支持 DouYin OAuth (Priority: P1) 🎯 MVP

**Goal**: `/debug` 可展示 DouYin Provider，完成 OAuth 授权并把 Flow 保存到内存/Redis（含 refresh_token/TTL）。

**Independent Test**: `go run ./cmd/accesstoken/server -config config.yaml` → `/debug` 中选择 DouYin → 点击发起授权 → 浏览器跳至 DouYin 登录 → 回调成功后在 Flow 列表出现 `provider_code=byte_dance_douyin`、正确的 `expire_at` 和 `storage_backend`。

### Implementation & Tests

- [x] T008 [US1] 编写 DouYin OAuth 合约测试：覆盖缺少 client_id/scope 报错、`expires_in` 写入 Flow TTL 的断言（cmd/accesstoken/server/oauth_contract_test.go）
- [x] T009 [US1] 在调试页模板/静态资源中新增 “字节跳动 / DouYin” 卡片，仅在服务端判定配置齐全时渲染，并在卡片上提示 `provider_code`/`oauth_key`/`storage_backend` 以及“单实例/Single app（<app>=default）”标签（cmd/accesstoken/server/debug_page.go, cmd/accesstoken/server/static/*）
- [x] T010 [US1] 在 `/accesstoken/oauth/start` 分支校验 DouYin 必填字段并构造授权 URL/state（cmd/accesstoken/server/handlers.go）
- [x] T011 [US1] 实现 `/debug/callback` 交换 DouYin token，写入内存与 Redis，记录 `flow_id`、`refresh_token`、`flow_ttl_seconds=expires_in`（cmd/accesstoken/server/handlers.go, cmd/accesstoken/server/server.go）
- [x] T012 [US1] 扩展 `/api/oauth/tokens` 与 `/accesstoken/flows` 输出，包含 `refresh_token` mask、`status`、`flow_expire_at`，并允许 `provider_code=byte_dance_douyin` 过滤（cmd/accesstoken/server/server.go）
- [x] T013 [US1] 在 Flow 列表/回填 UI 与日志中处理 Redis 降级：显式提示 `storage_backend=memory`、`flow_not_found`，并在 `/debug` 提示需重新授权（cmd/accesstoken/server/debug_page.go, cmd/accesstoken/server/server.go）
- [x] T014 [US1] 合约测试覆盖 Redis 不可用/flow_not_found 分支，以及 `byte_dance_douyin_config` 缺失时 `/debug`/`/api` 隐藏或拒绝 DouYin Provider 的 gating 逻辑，验证文案与 HTTP 响应（cmd/accesstoken/server/flow_contract_test.go）

**Checkpoint**: 完成后 `/debug` 可独立完成 DouYin OAuth 并在 UI/接口中看到 Flow。

---

## Phase 4: User Story 2 - DouYin API 可通过 `/accesstoken/call` 调试 (Priority: P2)

**Goal**: `/accesstoken/call` 可使用最新 DouYin Flow，具备自动刷新、per-action 1 QPS、429/5xx 指数退避与全链路脱敏日志。

**Independent Test**: Redis 中存在有效 DouYin Flow → 执行 `curl -H "Authorization: Bearer dev-accesstoken" -H "Content-Type: application/json" -d '{"provider_code":"byte_dance_douyin","action":"douyin.video.list"}' http://127.0.0.1:7071/accesstoken/call` → 响应成功 JSON，日志含 `event=token.call provider=byte_dance_douyin action=douyin.video.list retry_count=0`；若 token 过期则自动刷新；若 refresh 失败返回 `need reauth`。

### Implementation & Tests

- [x] T015 [US2] 补充 `/accesstoken/call` DouYin 场景的合约测试：验证自动刷新成功/失败、`need reauth` 错误、记录 `retry_count`（cmd/accesstoken/server/token_contract_test.go）
- [x] T016 [US2] 在 `handleCall` 中新增 DouYin 分支与 `executeDouyinCall`，封装 `ByteDanceDouYinACClient` action dispatch + payload 透传（cmd/accesstoken/server/handlers.go）
- [x] T017 [US2] 实现 `maybeRefreshDouYinToken`/`invalidateFlow`，检测 token 将过期时刷新并同步 TTL，失败时删除 Flow 并抛出 `need reauth`（cmd/accesstoken/server/server.go）
- [x] T018 [US2] 引入 per-action `rate.Limiter` 与 429/5xx 指数退避重试（含日志 `retry_count/backoff_ms`），超限返回 429 提示（cmd/accesstoken/server/server.go）
- [x] T019 [US2] 扩展日志与脱敏：`mask` 工具处理 DouYin access/refresh token，日志记录 `provider/provider_app/action/flow_id/token_source/retry_count/storage_backend`（internal/accesstoken/handler/mask, cmd/accesstoken/server/logging.go）
- [x] T020 [US2] 对比 DouYin 错误码/日志结构与现有 Provider，确保错误码、HTTP Status、字段名称保持一致，并补充必要的单测/校验（cmd/accesstoken/server/token_contract_test.go, docs/develop/access-token/byteDance/debug.md）
- [x] T021 [US2] 为缺少 access_token 或 `DOUYIN_SCOPE` 的请求编写合约测试与 handler 校验，确保返回 `missing access token`/`OAuth scope 未配置` 的解释性错误并与其他 Provider 一致（cmd/accesstoken/server/token_contract_test.go, cmd/accesstoken/server/handlers.go）

**Checkpoint**: `/accesstoken/call` DouYin action 可长期复用 Flow，且自动刷新/限流/重试逻辑可观测。

---

## Phase 5: User Story 3 - 文档可独立复现 (Priority: P3)

**Goal**: 仅阅读 `docs/develop/access-token/byteDance/{develop,debug}.md` 即可 1 小时内完成配置→授权→API 调试→Flow 回填。

**Independent Test**: 新同事按照文档步骤配置 env、启动服务、完成 OAuth 和 `/accesstoken/call` 示例；若受阻可在文档排障章节找到解决方案。

### Implementation

- [x] T022 [US3] 在 `docs/develop/access-token/byteDance/develop.md` 撰写环境变量表、配置示例、CLI/脚本调用说明、`need reauth` 排障，并加粗说明“仅支持单实例 `<app>=default`，如需多 app 需多实例部署”的提示（docs/develop/access-token/byteDance/develop.md）
- [x] T023 [US3] 在 `docs/develop/access-token/byteDance/debug.md` 记录 `/debug` 授权步骤、Flow “填充”演练、Redis/内存降级及速率限制提示（docs/develop/access-token/byteDance/debug.md）
- [x] T024 [US3] 将 `quickstart.md` 的端到端步骤同步到文档并在 Quickstart 中补充 rate-limit/refresh FAQ（specs/006-bytedance-access-token/quickstart.md, docs/develop/access-token/byteDance/develop.md）

**Checkpoint**: 文档读者可完全按图复现并定位常见错误。

---

## Phase 6: Polish & Cross-Cutting

- [x] T025 [P] 运行 `gofmt`/`go mod tidy` 并执行 `go test ./cmd/accesstoken/...` 确认 OAuth/Call/Flow 合约测试通过
- [x] T026 汇总 DouYin 接入说明到 `README.md` 或相关顶层文档（README.md, docs/plan/byteDance/access_token_client.md）
- [x] T027 [P] 依照 `quickstart.md` 实际走一遍从配置到 `/accesstoken/call` 的流程，记录剩余问题并更新文档/日志示例
- [x] T028 运行 Google/BiliBili/RedBook 等既有 Provider 的关键授权/调用回归测试，确认新增 DouYin 逻辑未破坏日志/错误码/流程（cmd/accesstoken/server/token_contract_test.go, docs/develop/access-token/*）

---

## Phase 7: Success Criteria Validation

- [x] T029 [P] SC-001：编写脚本统计 5 次 `/debug` DouYin OAuth 成功率，记录日志且未达 80% 时输出改进建议（scripts/, logs/）
- [x] T030 [P] SC-002：在 Redis 与内存模式下回填 20 条 Flow，生成统计报告（cmd/accesstoken/server/flow_contract_test.go, docs/develop/access-token/byteDance/develop.md）
- [x] T031 [P] SC-003：对 `/accesstoken/call` DouYin action 做 p95 ≤2s 的性能基准测试并在 README/文档中记录方法与结果（cmd/accesstoken/server/token_contract_test.go, docs/develop/access-token/byteDance/develop.md）
- [x] T032 [P] SC-004：组织 2 位未参与开发的同事依照文档完成端到端演练，收集反馈并更新排障章节（docs/develop/access-token/byteDance/{develop,debug}.md, quickstart.md）

---

## Dependencies & Execution Order

1. Phase 1 → Phase 2（配置与元数据必须先完成）
2. Phase 2 完成后方可进入各用户故事；US1 是 MVP，建议先完成
3. US2 依赖 US1 的 Flow 写入/字段；US3 可与 US2 并行但需引用最终接口形态
4. Polish 完成后，再进行 Success Criteria 验证以打包交付

## Parallel Opportunities

- [P] 任务可在不同文件/模块上并行推进（如配置示例、SDK helper、Go fmt/test、成功准则验证）
- 完成 Phase 2 后，US1/US2/US3 可由不同成员并行，但需保持接口契约同步
- Success Criteria 验证任务可与 Polish 并行开展，只需保证核心功能稳定

## Implementation Strategy

1. **MVP**：完成 Phase 1-3 → `/debug` DouYin OAuth + Flow 保存上线，可立刻被 CLI 复用
2. **迭代**：在 MVP 基础上落地 US2（调用/刷新/限流 + 错误码校验）→ 可验证 token 生命周期与跨 Provider 兼容
3. **完备**：补齐 US3 文档、Polish 以及 Phase 7 成功准则验证，确保指标、文档与排障全部闭环
