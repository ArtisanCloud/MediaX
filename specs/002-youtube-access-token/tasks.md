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

- [ ] T008 [US2] 在 `cmd/accesstoken/main.go` 实现 flag 解析、配置加载与 AccessToken 解析（flag > env > config）
- [ ] T009 [US2] 完成 `execVideosList`/`execSearchList`/`execPlaylistsList` 请求构造、结构化日志输出（provider/action/参数摘要）与敏感字段脱敏 (`cmd/accesstoken/main.go`)
- [ ] T010 [US2] 编写 CLI 单元测试，校验缺失参数/互斥参数/AccessToken 优先级 (`cmd/accesstoken/main_test.go`)
- [ ] T011 [US2] 扩展 `docs/develop/access-token/google/debug.md`，加入 CLI 命令示例、常见错误（invalid_grant/quotaExceeded/insufficientPermissions/401/403）、代理设置、`-mine`/`-ids` 互斥说明与排查指引
- [ ] T012 [US2] 在 `docs/develop/access-token/google/develop.md` 增加 CLI 使用流程与样例输出，映射到配置章节
- [ ] T021 [US2] 在 `docs/develop/access-token/google/debug.md` / `develop.md` 中新增“订阅→视频→发布→评论”闭环示例表格，列出对应 CLI 命令或 Playground API

---

## Phase 5: User Story 3 - Playground 示例验证 (Priority: P2)

**Goal**: 通过 Playground 演练展示 `MediaX.CreateGoogleYouTubeACClient` 与自定义 token 回调的完整链路。
**Independent Test**: `go run ./main.go`（启用 `playground.PlayGoogleYouTube`）成功输出 videoList，并在 `logs/info.log` 中看到 provider/action/参数摘要。

### Implementation

- [ ] T013 [US3] 调整 `playground/google.go` 示例，使其说明如何利用 `google_youtube_config` 与 `GetOAuthToken` 回调（含注释/日志）
- [ ] T014 [US3] 更新 `main.go`，添加启用 Playground 的注释/步骤，补充 BaseClient 日志字段/脱敏示例，并演示如何切换 Redis vs 内存缓存
- [ ] T015 [US3] 在 `docs/develop/access-token/google/develop.md` 增补 Playground 小节（启用方式、日志位置、Token 回调示例）

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T016 [P] 完成 `docs/develop/access-token/google/debug.md` 与 `develop.md` 的交叉校对，确保 CLI/Playground/配置相互引用
- [ ] T017 [P] 运行 `go test ./cmd/accesstoken` 并把示例命令记录到 `docs/develop/access-token/google/debug.md` 的“验证”小节
- [ ] T018 在 `docs/plan/google/youtube_access_token_client.md` 中记录 CLI/Playground 与 AccessToken Handler 的关系图示
- [ ] T019 [P] 更新 `specs/002-youtube-access-token/quickstart.md`，确保配置/CLI/Playground 步骤与最新实现一致，并包含闭环演练提示
- [ ] T020 [P] 校准 `specs/002-youtube-access-token/contracts/cli.md`，对齐最终参数集与示例输出，涵盖闭环相关 action

---

## Dependencies & Execution Order

1. **Setup → Foundational**: T001 必须先跑，随后 T002（Makefile 入口）解锁所有用户故事。
2. **User Stories**: US1 (Phase 3) 与 US2 (Phase 4) 均为 P1，应先完成功能配置与 CLI，再实现 P2 的 Playground（US3）。
3. **Polish**: T016–T018 需在所有用户故事完成后进行，以保证文档一致性与最终测试。

## Parallel Opportunities

- T003–T007（配置文档）可由一人集中完成；CLI 开发（T008–T010）可并行于文档更新（T011–T012），前提是 Makefile 目标已就绪。
- Playground 更新（T013–T015）与 CLI 单测（T010）可同时进行，因为两者改动不同文件。
- Polish 阶段的 T016 与 T017 可并行执行；T018 需在前两项完成后更新能力矩阵。

## Implementation Strategy

1. **MVP**: 完成 US1 + US2（配置模板 + CLI），即可交付可运行的调试体验。
2. **Incremental**: 在 MVP 稳定后追加 US3（Playground 演练），最后进行文档交叉检查与测试收尾。
3. **Testing**: 每个阶段结束后执行 `go test ./cmd/accesstoken` + CLI 实测命令，确保对应 User Story 可独立验证。
