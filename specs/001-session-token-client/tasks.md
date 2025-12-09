---

description: "Task list template for feature implementation"
---

# Tasks: SessionTokenClient Provider Architecture

**Input**: Design documents from `/specs/001-session-token-client/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: 仅在明确要求时添加测试任务；本功能涉及状态机、签名与安全逻辑，需包含关键单测。

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Go SDK monorepo：核心代码位于 `pkg/`，配置在 `pkg/client/config/`，文档在 `docs/` 与 `specs/`
- 新增 SessionToken 共享逻辑在 `pkg/client/sessiontoken/`
- Provider 适配器位于 `pkg/client/<provider>/sessionToken/`

## Constitution Guardrails Checklist *(apply to every story)*

- Include explicit tasks for provider/product directory wiring (`core/client.go`, `accessTokenClient/`, etc.) when adapters change.
- Capture config/schema edits under `pkg/client/config`, including yaml/json tags, env sourcing, and TTL enforcement.
- Add tasks for token lifecycle handling (BaseClient hooks, cache refresh, SessionToken flow persistence) plus the accompanying unit tests.
- Ensure observability tasks cover logging fields (`provider`, `api`, `flow_id`, `tenant_uuid`), retry counters, and sensitive-data masking。
- Add standalone test/documentation tasks for pure logic modules, SessionToken state transitions, callback签名，以及必需的授权/安全校验。

## Phase 1: Setup (Shared Infrastructure)

- [ ] T001 确认 `config.yaml` 示例中新增 Zhihu `sessionToken` 段并引用环境变量 (`config.yaml`)
- [ ] T002 拉取依赖并验证 go modules 版本（`go mod tidy`）
- [ ] T003 [P] 创建 `pkg/client/sessiontoken/` 目录框架（manager、flow、callback、storage 子目录）

## Phase 2: Foundational (Blocking Prerequisites)

- [ ] T004 设计 Flow 数据结构与 Redis key 前缀（`pkg/client/sessiontoken/flow.go`）
- [ ] T005 实现 Redis 存储驱动（读写 Flow、TTL、幂等）（`pkg/client/sessiontoken/storage/redis/storage.go`）
- [ ] T006 定义 `SessionTokenClient`、`Authenticator`、`CredentialHarvester`、`CallbackDispatcher` 接口（`pkg/client/sessiontoken/interfaces.go`）
- [ ] T007 实现 `sessiontoken.Manager` 基础骨架（依赖 BaseClient、logger、cache）（`pkg/client/sessiontoken/manager.go`）
- [ ] T008 [P] 编写 Flow 状态机与事件枚举（pending→authorizing→succeeded/failed）（`pkg/client/sessiontoken/state_machine.go`）
- [ ] T009 [P] 添加 Flow 状态机单测覆盖正常与异常分支（`pkg/client/sessiontoken/state_machine_test.go`）
- [ ] T010 定义 HMAC 签名工具与回调 dispatcher 接口（`pkg/client/sessiontoken/callback/dispatcher.go`）
- [ ] T011 [P] 编写签名与重试逻辑的单元测试（`pkg/client/sessiontoken/callback/dispatcher_test.go`）
- [ ] T012 完成日志脱敏/敏感字段工具（`pkg/client/sessiontoken/sanitizer/sanitizer.go` + tests）
- [ ] T013 实现 SessionToken API Bearer Token 验证中间件/拦截器，并复用配置中的 `api_token`（`pkg/server/middleware/sessiontoken_auth.go`）
- [ ] T014 [P] 为授权中间件添加单测，覆盖缺失/错误 token → 401（`pkg/server/middleware/sessiontoken_auth_test.go`）

## Phase 3: User Story 1 - 插件创建 SessionToken Flow 并拉起登录 (Priority: P1) 🎯 MVP

**Goal**: 提供 `POST /session-token/flows` API，创建 Flow、生成 authorize_url、保持 pending 状态。

**Independent Test**: 使用 contract (`contracts/session-token-flows.yaml`) 与 curl 示例验证 Flow 创建返回 flow_id/authorize_url 并存入 Redis。

### Implementation

- [ ] T015 [US1] 实现 Flow 创建请求 DTO + 校验逻辑（`pkg/client/sessiontoken/api/create_flow_request.go`）
- [ ] T016 [US1] 在 `sessiontoken.Manager` 实现 CreateFlow 方法（含 authorize_url 构建钩子）（`pkg/client/sessiontoken/manager.go`）
- [ ] T017 [US1] 集成 Redis 存储写入、重复 state 幂等处理（`pkg/client/sessiontoken/manager.go`）
- [ ] T018 [US1] 在 HTTP 层注册 `POST /session-token/flows` handler 并接入 Bearer 中间件（`pkg/server/handlers/sessiontoken_flow_create.go`）
- [ ] T019 [US1] 输出结构化日志（provider/api/tenant_uuid/flow_id）与错误包装（`pkg/server/handlers/sessiontoken_flow_create.go`）
- [ ] T020 [P] [US1] 编写 handler + manager 的单元测试（`pkg/server/handlers/sessiontoken_flow_create_test.go`、`pkg/client/sessiontoken/manager_create_test.go`）

## Phase 4: User Story 2 - 插件轮询 Flow 状态 (Priority: P1)

**Goal**: 提供 `GET /session-token/flows/{flow_id}`，返回状态、TTL、last_error/result。

**Independent Test**: 通过 contract 与 curl 验证查询接口可展示状态、错误，并在 TTL 内可重复读取。

### Implementation

- [ ] T021 [US2] 实现查询请求解析/响应 DTO（`pkg/client/sessiontoken/api/get_flow_response.go`）
- [ ] T022 [US2] 在 `sessiontoken.Manager` 添加 `GetFlow`，支持 expired/not_found 分支（`pkg/client/sessiontoken/manager.go`）
- [ ] T023 [US2] wiring HTTP handler `GET /session-token/flows/{flow_id}` 并挂载 Bearer 中间件（`pkg/server/handlers/sessiontoken_flow_get.go`）
- [ ] T024 [US2] 追加 last_error/result 序列化与脱敏逻辑（`pkg/server/handlers/sessiontoken_flow_get.go`）
- [ ] T025 [P] [US2] 编写查询接口单元测试，覆盖 pending/succeeded/failed/expired/401 未授权（`pkg/server/handlers/sessiontoken_flow_get_test.go`）

## Phase 5: User Story 3 - 凭证采集与回调 (Priority: P2)

**Goal**: Harvester 产出凭证后更新 Flow、触发 HMAC 回调，失败可重试 3 次并记录 last_error。

**Independent Test**: 使用 fake Harvester + mock callback server 验证成功/失败路径与重试次数。

### Tests (必要)
- [ ] T026 [P] [US3] 构建 callback dispatcher 集成测试，模拟签名通过/失败（`pkg/client/sessiontoken/callback/dispatcher_integration_test.go`）

### Implementation
- [ ] T027 [US3] 在 `sessiontoken.Manager` 增加 `CompleteFlowSuccess`/`CompleteFlowFailed` API（`pkg/client/sessiontoken/manager.go`）
- [ ] T028 [US3] 实现凭证脱敏存储与 Redis 更新 result（`pkg/client/sessiontoken/manager.go`）
- [ ] T029 [US3] 在 dispatcher 中实现 2s→4s→8s 指数退避逻辑并记录 retry（`pkg/client/sessiontoken/callback/dispatcher.go`）
- [ ] T030 [US3] 添加回调 HMAC 签名与 payload 结构（`pkg/client/sessiontoken/callback/payload.go`）
- [ ] T031 [US3] 在日志中记录 callback 成功/失败、flow_id、retry 次数（`pkg/client/sessiontoken/callback/dispatcher.go`）

## Phase 6: User Story 4 - 知乎适配器与配置注入 (Priority: P2)

**Goal**: 在 Zhihu provider 中实现 SessionToken adapter，配置 Authenticator/Harvester/Callback，并通过 MediaX 工厂注入。

**Independent Test**: 通过 fake 脚本验证 adapter 生效；运行 `go test ./pkg/client/zhihu/sessionToken/...`。

### Implementation
- [ ] T032 [US4] 在 `pkg/client/config/zhihu.go` 添加 `ZhihuSessionTokenConfig` 结构、yaml/json tags、默认值（`pkg/client/config/zhihu.go`）
- [ ] T033 [US4] 在 `config.yaml`/README 添加 Zhihu sessionToken 配置示例（`config.yaml`, `README.md`）
- [ ] T034 [US4] 新建 `pkg/client/zhihu/sessionToken/core/client.go`，复用 BaseClient 并注入配置/logger/cache（`pkg/client/zhihu/sessionToken/core/client.go`）
- [ ] T035 [US4] 实现 Zhihu Authenticator（入口 URL、UA、脚本选择）（`pkg/client/zhihu/sessionToken/authenticator/authenticator.go`）
- [ ] T036 [US4] 实现 Zhihu Harvester（抓取 cookies/headers）（`pkg/client/zhihu/sessionToken/harvester/harvester.go`）
- [ ] T037 [US4] 实现 Zhihu CallbackDispatcher wrapper，复用共享 dispatcher 并添加 provider metadata（`pkg/client/zhihu/sessionToken/callback/dispatcher.go`）
- [ ] T038 [US4] 在 `pkg/client/mediaX.go` 注册 `CreateZhihuSessionTokenClient`，并确保工厂注入日志/缓存（`pkg/client/mediaX.go`）
- [ ] T039 [P] [US4] 编写 Zhihu adapter 单元测试（配置解析 + Authenticator/Harvester 行为）（`pkg/client/zhihu/sessionToken/..._test.go`）

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T040 更新文档：`docs/plan/session_token_client.md` 与 `quickstart.md`，说明新的 API、授权与配置（`docs/plan/session_token_client.md`, `specs/001-session-token-client/quickstart.md`）
- [ ] T041 校验 `contracts/session-token-flows.yaml` 与实现一致并生成示例（`specs/001-session-token-client/contracts/session-token-flows.yaml`）
- [ ] T042 `go test ./...` + `golangci-lint`（若配置）全量回归，确保无敏感日志（项目根）
- [ ] T043 运行示例 curl/quickstart，验证 Flow 创建→查询→回调闭环（`specs/001-session-token-client/quickstart.md` 指引）
- [ ] T044 为 Flow 创建/查询与回调添加 latency/成功率指标上报或日志量化（`pkg/client/sessiontoken/manager.go`, `pkg/client/sessiontoken/callback/dispatcher.go`）
- [ ] T045 根据指标结果更新 README 或运维指引，记录 SLA 检查/调优方法（`README.md`, `docs/plan/session_token_client.md`）

## Dependencies & Execution Order

- Setup → Foundational → US1 → US2 → US3 → US4 → Polish
- US1/US2 为 MVP（Flow 创建 + 查询）；US3、US4 可在 Foundational 完成后并行启动，但回调 dispatcher 与授权中间件需先完成。

## Parallel Opportunities

- T003 与 T001/T002 可并行。
- Foundational 中 T008/T009、T011/T012、T013/T014 可并行。
- US1 handler/manager 与 US2 查询在基础能力完成后可并行推进。
- US3 dispatcher 实现与 US4 配置/adapter 可并行，前提是 manager 扩展已完成。
- Polish 中的文档/指标任务（T040–T045）可按职责分摊并发执行。

## Implementation Strategy

1. 优先完成 Foundational + US1/US2，形成可创建/查询 Flow 且具备授权校验的 MVP。
2. 在 MVP 基础上扩展 US3（凭证回调）以满足闭环，再落地 Zhihu adapter（US4）。
3. 最后完善文档、契约与性能/成功率指标，确保 SLA 可量化验证。
