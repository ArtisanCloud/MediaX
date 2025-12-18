# Tasks: WeChat ClientToken Server

**Input**: Design documents from `/specs/003-wechat-client-token/`
**Prerequisites**: spec.md, plan.md, docs/plan/wechat/client_token_client.md

---

## Phase 0: Setup

- [x] T000 `go test ./...` 确认当前仓库无回归问题，记录 baseline

---

## Phase 1: 配置与工厂 (Blocking)

**Goal**: 在 `config.yaml` 中新增 `client_token_providers`，并提供对应的 config 结构/工厂。

- [x] T101 [CFG] 扩展 `pkg/client/config`，定义 `ClientTokenProviderConfig` & `WechatClientTokenConfig`（字段：appid/appsecret/message_token/message_aes_key/api_token/cache/redis）
- [x] T102 [CFG] 更新 `pkg/client/mediaX.go`，新增 `CreateWechatClientTokenClient`，复用 `pkg/client/wechat/officialAccount/clientTokenClient`
- [x] T103 [CFG] 更新 `config.example.yaml` 与 `config.yaml` 注释，示例化 `client_token_providers`（使用占位符）
- [x] T104 [CFG] 文档 `docs/plan/wechat/client_token_client.md` / `docs/develop/wechat/client_token.md` 补充配置说明与环境变量覆盖顺序
- [x] T105 [CFG] 校验 `config` 加载逻辑（新 Node/flag），确保 `go run ./cmd/clienttoken/server -config config.yaml` 能解析

---

## Phase 2: ClientToken Server 基础 (P1 🎯)

**Goal**: 交付 `cmd/clienttoken/server` 主体，提供 REST API + Redis 缓存 + API Token 鉴权。

- [ ] T201 [SRV] 新建 `cmd/clienttoken/server/main.go`（flag：config/port），引用 `run()`
- [ ] T202 [SRV] `run()`：加载配置→初始化 Redis（失败直接退出）→ 初始化 Logger/MediaX client→注册 HTTP 路由
- [ ] T203 [SRV] 实现 API Token 中间件（默认 `dev-clienttoken`，可配置环境变量覆盖）
- [ ] T204 [SRV] Handlers：
  - `POST /client-token/token`（刷新 token，写入 Redis + 内存缓存）
  - `GET /client-token/cache`（读取缓存并返回 TTL）
  - `POST /client-token/call`（通用 API 调用：action/method/query/body）
  - `POST /client-token/message/validate`（signature 验证，含 echostr 返回）
  - `POST /client-token/message/callback`（记录/可选 AES 解密）
  - `GET /healthz`
- [ ] T205 [SRV] 缓存模块：封装 `clientToken:wechat:<appid>` key，包含 TTL、提前刷新逻辑；Redis 必须可用
- [ ] T206 [SRV] 结构化日志 `clienttoken_metric`（action/provider/appid/status/latency/token_source），日志脱敏

---

## Phase 3: 调试页 & 回调日志 (P1 🎯)

**Goal**: 提供与 SessionToken/AccessToken 相同体验的 `/debug` 页面。

- [ ] T301 [UI] `cmd/clienttoken/server/debug_page.go`（Go 模板）实现：
  - 配置选择（Provider/App）
  - Token 卡片（来源/TTL/刷新按钮）
  - API 调试表单（action/method/query/body）
  - 消息验证/回调区域
  - 回调日志列表
- [ ] T302 [UI] `callback_store` 实现（复用 AccessToken 结构）：记录最近 50 条消息（时间/Query/Headers/Body）
- [ ] T303 [UI] 前端 JS 调用 `/client-token/token`、`/client-token/cache`、`/client-token/call`、`/client-token/message/*`，携带 API Token
- [ ] T304 [UI] 生成 `http://127.0.0.1:7072/debug` 页面，确保跨浏览器可用

---

## Phase 4: CLI & 文档 (P1)

**Goal**: 一键启动/调试体验；提供脚本/QS 文档。

- [ ] T401 [DOC] 在 `docs/develop/wechat/client_token.md` 编写 Quickstart（配置→`go run ./cmd/clienttoken/server -config config.yaml`→调试步骤）
- [ ] T402 [DOC] 更新 `docs/plan/wechat/client_token_client.md`、`specs/003-.../plan.md` 的交付/流程图
- [ ] T403 [CLI] 提供 `make clienttoken` 或 `scripts/clienttoken-refresh.sh`：示例刷新 token/调用 API
- [ ] T404 [DOC] 在 README/根目录文档中说明 ClientToken Server 的作用与启动命令

---

## Phase 5: 验收 & 测试 (P1)

- [ ] T501 [TEST] `go test ./cmd/clienttoken/...`（配置加载/缓存 helper/消息验证算法）
- [ ] T502 [MANUAL] 启动服务 + 浏览器访问 `/debug`，完成 token 刷新、`cgi-bin/user/get` 调用、消息验证流程，附截图/日志
- [ ] T503 [QA] 使用 `curl` 验证 `/client-token/token`、`/client-token/cache`、`/client-token/call`、`/client-token/message/validate`、`/client-token/message/callback`
- [ ] T504 [DOC REVIEW] 校对 spec/plan/quickstart，确保命令与配置一致（`go run ./cmd/clienttoken/server -config config.yaml`）

---

## Execution Notes

- 阶段顺序：Phase1 → Phase2 → Phase3 → Phase4 → Phase5（调试页依赖 server，文档最终与实现对齐）
- 并行机会：配置/工厂（Phase1）可与调试文档部分同步推进；调试页 UI 与 handler 可并行开发，只需约定 API 格式；CLI 脚本/文档在 server 可运行后再产出。
- 交付标准：调试页在浏览器可运行、Redis 中存在 `clientToken:wechat:<appid>` 缓存、所有 REST API 均有日志/鉴权/错误提示。
