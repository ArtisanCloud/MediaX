# Tasks: Google YouTube AccessToken Enablement

**Input**: Design documents from `/specs/002-youtube-access-token/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

## Phase 1: Setup (Shared Infrastructure)

- [x] T001 Execute baseline regression tests to confirm clean starting point (`go test ./...`)

---

## Phase 2: Foundational (Blocking Prerequisites)

- [x] T002 Add `accesstoken` make target wiring到 `cmd/accesstoken` 入口，确保整个团队有统一命令 (`Makefile`)

---

## Phase 3: User Story 1 - 统一配置模板 (Priority: P1) 🎯 MVP

**Goal**: 为 Google YouTube AccessToken 客户端提供统一的配置模板与文档，保证本地/CI 可快速加载凭证。
**Independent Test**: 新人复制 `config.example.yaml` → `config.yaml`，填入凭证后运行 `make accesstoken -action videos.list -part snippet -ids <id>` 可成功返回数据。

### Implementation

- [x] T003 [US1] 扩充 `google_youtube_config` 模板到 `config.example.yaml`，包含 API URL、OAuth 字段与注释
- [x] T004 [US1] 同步 `config.yaml` 示例中的 `google_youtube_config`，并在注释里提醒使用环境变量注入敏感值
- [x] T005 [US1] 在 `docs/develop/access-token/google/develop.md` 中编写配置章节（模板拷贝、oauth_key、环境变量覆盖顺序）
- [x] T006 [US1] 更新 `docs/develop/access-token/google/debug.md`，描述准备配置、代理参数、AccessToken 缺失/过期处理以及 Redis 缓存清理命令
- [x] T007 [US1] 在 `docs/plan/google/youtube_access_token_client.md` 的能力矩阵内追加配置能力与引用链接，确保 SDK 使用者能定位模板

---

## Phase 4: User Story 2 - CLI 调试体验 (Priority: P1)

**Goal**: 提供 `cmd/accesstoken` CLI，支持 videos/search/playlists 调试并输出标准 JSON。
**Independent Test**: `make accesstoken ARGS='-action search.list -part snippet -query MediaX -max-results 3 -access-token <token>'` 返回 200 且日志包含 provider/action。

### Implementation

- [x] T008 [US2] 在 `cmd/accesstoken/main.go` 实现 flag 解析、配置加载与 AccessToken 解析（flag > env > config）
- [x] T009 [US2] 完成 `execVideosList`/`execSearchList`/`execPlaylistsList` 请求构造、结构化日志输出（provider/action/参数摘要）与敏感字段脱敏 (`cmd/accesstoken/main.go`)
- [x] T010 [US2] 编写 CLI 单元测试，校验缺失参数/互斥参数/AccessToken 优先级 (`cmd/accesstoken/main_test.go`)
- [x] T011 [US2] 扩展 `docs/develop/access-token/google/debug.md`，加入 CLI 命令示例、常见错误（invalid_grant/quotaExceeded/insufficientPermissions/401/403）、代理设置、`-mine`/`-ids` 互斥说明与排查指引
- [x] T012 [US2] 在 `docs/develop/access-token/google/develop.md` 增加 CLI 使用流程与样例输出，映射到配置章节
- [x] T021 [US2] 在 `docs/develop/access-token/google/debug.md` / `develop.md` 中新增“订阅→视频→发布→评论”闭环示例表格，列出对应 CLI 命令或 Playground API

---

## Phase 5: User Story 4 - AccessToken 调试服务 (Priority: P1)

**Goal**: 为产品/QA/合作伙伴提供可视化调试入口，复用 SessionToken 调试台的交互，在浏览器中完成授权回调观测、AccessToken 刷新与 API 调用。
**Independent Test**: `make accesstoken-serve`（或 `go run ./cmd/accesstoken/server`）后访问 `http://127.0.0.1:7070/debug/accesstoken`，无需 CLI/Playground 即可完成最小链路。

### Implementation

- [ ] T022 [US4] 新增 `cmd/accesstoken/server`（或 `cmd/accesstoken/main.go` 中引入 `serve` 子命令），整合配置加载、日志、缓存与 API Token 解析，默认监听 `:7070`。
- [ ] T023 [US4] 复用/改写 `cmd/sessiontoken/debug_page.go`，实现 AccessToken 专用调试页，覆盖 Provider/App/API 版本选择、授权 URL/AccessToken 操作、API 调试表单、回调日志表格。
- [ ] T024 [US4] 暴露 RESTful 接口（`/accesstoken/token`、`/accesstoken/call`、`/debug/callback`），内部复用 CLI 的执行逻辑，并输出 JSON/错误结构，日志需脱敏 token。
- [ ] T025 [US4] 在 `docs/develop/access-token/google/develop.md` / `debug.md` / `quickstart.md` 等文档补充调试服务的启动方式、页面说明、API token 配置与常见排查步骤。

---

## Phase 6: User Story 3 - Playground 示例验证 (Priority: P2)

**Goal**: 通过 Playground 演练展示 `MediaX.CreateGoogleYouTubeACClient` 与自定义 token 回调的完整链路。
**Independent Test**: `go run ./main.go`（启用 `playground.PlayGoogleYouTube`）成功输出 videoList，并在 `logs/info.log` 中看到 provider/action/参数摘要。

### Implementation

- [x] T013 [US3] 调整 `playground/google.go` 示例，使其说明如何利用 `google_youtube_config` 与 `GetOAuthToken` 回调（含注释/日志）
- [x] T014 [US3] 更新 `main.go`，添加启用 Playground 的注释/步骤，补充 BaseClient 日志字段/脱敏示例，并演示如何切换 Redis vs 内存缓存
- [x] T015 [US3] 在 `docs/develop/access-token/google/develop.md` 增补 Playground 小节（启用方式、日志位置、Token 回调示例）

---

## Phase 7: Polish & Cross-Cutting Concerns

- [x] T016 [P] 完成 `docs/develop/access-token/google/debug.md` 与 `develop.md` 的交叉校对，确保 CLI/Playground/配置相互引用
- [x] T017 [P] 运行 `go test ./cmd/accesstoken` 并把示例命令记录到 `docs/develop/access-token/google/debug.md` 的“验证”小节
- [x] T018 在 `docs/plan/google/youtube_access_token_client.md` 中记录 CLI/Playground 与 AccessToken Handler 的关系图示
- [x] T019 [P] 更新 `specs/002-youtube-access-token/quickstart.md`，确保配置/CLI/Playground 步骤与最新实现一致，并包含闭环演练提示
- [x] T020 [P] 校准 `specs/002-youtube-access-token/contracts/cli.md`，对齐最终参数集与示例输出，涵盖闭环相关 action

---

## Dependencies & Execution Order

1. **Setup → Foundational**: T001 必须先跑，随后 T002（Makefile 入口）解锁所有用户故事。
2. **User Stories**: 先完成 US1 (Phase 3) 与 US2 (Phase 4) 以稳定配置和 CLI，再实现新加入的 US4 (Phase 5) —— 调试服务依赖 CLI 的执行逻辑与配置模板，完成后再进入 P2 的 Playground（US3，Phase 6）。
3. **Polish**: T016–T018 需在所有用户故事完成后进行，以保证文档一致性与最终测试。

## Parallel Opportunities

- T003–T007（配置文档）可由一人集中完成；CLI 开发（T008–T010）可并行于文档更新（T011–T012），前提是 Makefile 目标已就绪。
- AccessToken 调试服务（T022–T025）可以在 CLI 核心逻辑稳定后启动，与 Playground（T013–T015）互不阻塞；页面/后端可由不同成员并行实现。
- Playground 更新（T013–T015）与 CLI 单测（T010）可同时进行，因为两者改动不同文件。
- Polish 阶段的 T016 与 T017 可并行执行；T018 需在前两项完成后更新能力矩阵。

## Implementation Strategy

1. **MVP**: 完成 US1 + US2（配置模板 + CLI），即可交付命令行调试体验。
2. **Incremental**: 基于 CLI 逻辑实现 US4（调试服务），随后追加 US3（Playground 演练），最后进行文档交叉检查与测试收尾。
3. **Testing**: 每个阶段结束后执行 `go test ./cmd/accesstoken` + CLI/HTTP 调试命令（`make accesstoken`、`make accesstoken-serve`），确保对应 User Story 可独立验证。
