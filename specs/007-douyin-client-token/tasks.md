# Tasks: DouYin ClientToken Mode Onboarding

**Input**: Design documents from `/specs/007-douyin-client-token/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 准备开发环境与本地调试配置，确保后续修改可在本机复现。

- [x] T001 在 `go.mod`/`go.sum` 上运行 `go mod tidy`，确保 DouYin ClientToken 相关依赖可解析。
- [x] T002 在根目录 `config.yaml` 中复制 `client_token_providers` 的字节跳动占位段，并填入本地 `DOUYIN_*` 与 `CLIENTTOKEN_*` 环境变量，供调试台自测。

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 扩展配置层以支持 DouYin ClientToken，所有用户故事依赖这些配置。

- [x] T003 在 `pkg/client/config/douyin.go` 与 `pkg/client/config/config.go` 中新增 `byte_dance_douyin_config.client_token/cache/redis` 字段与 YAML/JSON tag，映射 `client_key`、`client_secret`、TTL 及 Redis 设定。
- [x] T004 [P] 更新 `config.example.yaml`，追加 `provider_code=byte_dance_douyin_clienttoken` 的示例块与环境变量注释，保持与文档示例一致。
- [x] T005 [P] 在 `pkg/client/config/config_test.go` 新增用例，校验 `ClientTokenProviders.Resolve...` 可解析 DouYin 配置并返回正确的 `redis_key` 与 `refresh_before_seconds`。

**Checkpoint**: 配置层可加载 DouYin ClientToken，后续实现可安全读取。

---

## Phase 3: User Story 1 - `/client-token/debug` 刷新 DouYin ClientToken (Priority: P1) 🎯 MVP

**Goal**: 调试台可列出 DouYin Provider，并完成刷新/调试 API 的端到端流程。
**Independent Test**: `go run ./cmd/clienttoken/server -config config.yaml` → 浏览器进入 `/client-token/debug` → DouYin Provider 可选、点击“刷新 Token”成功写入缓存并可调用 API。

### Implementation for User Story 1

- [x] T006 [US1] 更新 `cmd/clienttoken/server/server.go` 的 `resolveProviderConfig`，识别 `provider_code=byte_dance_douyin_clienttoken`，并从 `byte_dance_douyin_config` 派生缓存 key/TTL。
- [x] T007 [P] [US1] 在 `cmd/clienttoken/server/server.go`/`server.go:handleToken` 中调用 `MediaX.CreateByteDanceDouYinCTClient`，实现 DouYin client_token 刷新与 `tokenStore` 写入（含脱敏字段）。
- [x] T008 [P] [US1] 扩展 `cmd/clienttoken/server/debug_page.go` 与相关模板，增加 DouYin Provider 选项、默认 JSON 模板及 `action` 下拉示例。
- [x] T009 [US1] 在 `cmd/clienttoken/server/server_test.go` 添加集成测试，验证 `/client-token/token` 与 `/client-token/call` 对 DouYin Provider 返回结构化日志与 200 响应。

**Checkpoint**: `/client-token/debug` 可刷新 DouYin Token 并调用 API。

---

## Phase 4: User Story 2 - Redis 缓存与 CLI/脚本复用 (Priority: P2)

**Goal**: 客户端与 CLI 共享 Redis 缓存，自动续期，并提供脚本工具。
**Independent Test**: 配置 Redis → `./scripts/clienttoken-douyin.sh refresh` → Redis 记录 `clientToken:douyin:<client_key>`，TTL 逼近阈值时自动续期，脚本可查看/清除缓存。

### Implementation for User Story 2

- [x] T010 [US2] 在 `pkg/client/byteDance/douYin/clientTokenClient/client.go`（或新增缓存辅助文件）实现读取 token 时的 TTL 检测，低于 `refresh_before_seconds` 时调用 DouYin API 自动刷新并回写。
- [x] T011 [P] [US2] 扩展 `cmd/clienttoken/server/cache_store.go`，统一 Redis Key `clientToken:douyin:<client_key>` 的读写逻辑，记录 `source`, `ttl_remaining`, `refreshed_at`。
- [x] T012 [P] [US2] 改造 `cmd/clienttoken/server/server.go:handleCache`（含 DELETE）以返回 TTL/来源/最近刷新时间，并允许清空缓存供复测。
- [x] T013 [US2] 新建 `scripts/clienttoken-douyin.sh`，实现 `refresh`/`cache`/`call` 子命令并复用 `CLIENTTOKEN_API_TOKEN` 发送 HTTP 请求。
- [x] T014 [P] [US2] 在 `pkg/client/byteDance/douYin/clientTokenClient/client_test.go` 编写单测，覆盖 TTL 阈值触发刷新与 Redis/内存降级行为。
- [x] T015 [P] [US2] 增强 `cmd/clienttoken/server/server.go` 日志，输出 `event`, `provider`, `provider_code`, `app_code`, `action`, `redis_key`, `ttl_remaining`, `storage_backend`, `token_source`, `api_token_subject` 等字段，并统一使用 `internal/accesstoken/handler/mask` 脱敏。
- [x] T016 [US2] 在 `cmd/clienttoken/server/server.go` 的刷新/调用流程中加入“TTL 低于阈值必须等待刷新”分支，失败时返回 `token refresh failed` 并记录剩余 TTL。
- [x] T017 [US2] 在 `cmd/clienttoken/server/server.go`、`debug_page.go` 与配置模板中检测 `device_id`、`risk_info` 等可选字段，默认禁用相关 API 模板；若启用需验证配置齐全并给出显式提示，确保与文档声明一致。

**Checkpoint**: CLI 与 Redis 缓存协同工作，可自动续期且日志完备。

---

## Phase 5: User Story 3 - 文档与配置可复制 (Priority: P3)

**Goal**: 仅靠文档即可在 1 小时内完成 DouYin ClientToken 接入。
**Independent Test**: 新同事按文档设置 env→复制配置→刷新 Token→调试 API，全程依靠文档排障。

### Implementation for User Story 3

- [x] T018 [US3] 更新 `docs/develop/client-token/byteDance/develop.md`，加入环境变量、启动命令、CLI/脚本使用及常见问题，并新增 `device_id`/`risk_info` 等可选字段的说明及默认禁用策略。
- [x] T019 [P] [US3] 更新 `docs/develop/client-token/byteDance/debug.md`，补充调试台步骤、模板同步、缓存查看与自动刷新提示，标注当启用额外参数时的风险提示。
- [x] T020 [US3] 在根 `config.yaml`（或 `config.example.yaml` 相应部分）添加注释，说明各 `DOUYIN_*`、Redis 变量以及额外字段的默认值与作用，确保复制即可运行。
- [x] T021 [P] [US3] 在 `docs/plan/byteDance/client_token_client.md` 与 `specs/007-douyin-client-token/quickstart.md` 中同步 CLI/Redis 场景、额外字段排障案例及验证步骤，方便运维复查。

**Checkpoint**: 文档/配置示例完整，外部团队可自助复现。

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 统一编码风格、回归测试并同步计划文档。

- [x] T022 [P] 对 `cmd/clienttoken/server/*.go` 与 `pkg/client/byteDance/douYin/**/*.go` 运行 `gofmt`/`goimports`，保证代码风格一致。
- [x] T023 运行 `go test ./cmd/clienttoken/... ./pkg/client/byteDance/douYin/...`，确认刷新/缓存逻辑全部通过。
- [x] T024 在 `docs/plan/byteDance/client_token_client.md` 记录 DouYin ClientToken 接入说明与后续待办。
- [x] T025 依照 `docs/develop/client-token/byteDance/develop.md` 步骤复跑一次 Quickstart，并在 `specs/007-douyin-client-token/quickstart.md` 标注最新验证日期与结果。

---

## Dependencies & Execution Order

1. **Phase 1 → Phase 2**：先准备本地配置与依赖，再扩展配置层。
2. **Phase 2 → Phase 3/4/5**：所有用户故事都依赖 DouYin 配置结构，必须完成后方可进入实现。
3. **User Stories**：
   - US1 解锁基本刷新/调用能力，是 MVP。
   - US2 构建共享缓存与 CLI，可在 US1 完成后并行推进（部分代码同文件需串行 review）。
   - US3 主要是文档，可在 US1 稳定后开始，与 US2 并行。
4. **Polish**：最后统一格式、跑测试、同步 quickstart。

## Parallel Opportunities

- Setup 阶段任务彼此独立，可并行完成。
- Foundational 中的 T004/T005 互不依赖，可并行；T003 完成后再执行与配置读取相关的用户故事。
- US1 中 T007/T008 可并行（分别改 server 逻辑与 UI），最终由 T006/T009 串联验证。
- US2 中 T011/T012/T013/T015 可由不同成员并行，T010/T014/T016 依赖核心 TTL 逻辑调整。
- 文档（T017–T020）可由文档负责人并行推进，与开发并不冲突。
- Polish 阶段 T021 与 T022/T023/T024 可由不同成员并行执行。

## Implementation Strategy

1. **MVP (US1)**：完成 Phase 1-3，使 `/client-token/debug` 能刷新 DouYin Token 并调用 API。
2. **Incremental**：
   - Phase 4 引入自动续期 + CLI，满足批量任务需求。
   - Phase 5 完整文档与配置示例，支持外部团队自助。
3. **Validation**：每个 User Story 完成后依据其 Independent Test 进行验收，可随时停在任意“Checkpoint”交付阶段成果。
