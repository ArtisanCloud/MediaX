# Tasks: RedBook JuGuang AccessToken Debug Integration

**Input**: Design artifacts from `/specs/005-redbook-access-token/`
**Prerequisites**: spec.md, plan.md, research.md, data-model.md, contracts/, quickstart.md

## Phase 1: Provider & Config Foundations

- [x] T101 Update `pkg/client/config/redbook.go` 与 `config.example.yaml`，补齐 `oauth_url`、`access_token_url`、`scope`、`redirect_url` 并加入 env 说明。
- [x] T102 调整 `cmd/accesstoken/server/server.go` Provider 初始化，读取 `redbook_juguang_config` 并在 `/debug` Provider JSON 中展示聚光元数据。
- [x] T103 在 `/debug` 前端模板与静态资源中新增聚光 Provider 卡片（标题、API 版本、默认模板），保持与现有 Provider UI 一致。
- [x] T104 扩充启动/配置文档（`docs/develop/access-token/redbook/develop.md` / `debug.md` 初稿），描述环境变量、Redis 依赖与警告提示。

---

## Phase 2: OAuth Flow 接入 (User Story 1)

- [x] T201 在 `/accesstoken/oauth/start` handler 中处理 `provider_code=redbook_juguang`，验证配置字段并返回授权 URL。
- [x] T202 实现 `/debug/callback` 对聚光回调的解析与脱敏，复用 Flow 存储逻辑（memory/Redis）。
- [x] T203 完善 `/accesstoken/flows` 与 `/api/oauth/tokens`，确保可以根据 provider 过滤并回填聚光 Flow。
- [x] T204 添加契约/单元测试覆盖聚光 OAuth Start + Flow 存储（含 scope 缺失与 Redis fallback）。

---

## Phase 3: 聚光 API 调试入口 (User Story 2)

- [x] T301 新增 `/accesstoken/redbook/account/balance` endpoint：读取 Flow/Config，调用 `JuGuangAccountClient.GetAccountBalance`，返回脱敏 JSON。
- [x] T302 在调试台 UI/示例 JSON 中加入账户余额模板（advertiser_id），支持一键填充。
- [x] T303 根据 `contracts/accesstoken_redbook.yaml` 为新 endpoint 添加 contract tests（成功、缺少 token、配置错误场景）。
- [x] T304 扩充日志与错误处理：为聚光 API 调用输出 `provider`, `provider_app`, `flow_id`, `token_source` 字段并复用掩码工具。

---

## Phase 4: 文档 & Quickstart (User Story 3 / Polish)

- [x] T401 将 quickstart.md 中的步骤同步到 `docs/develop/access-token/redbook/{develop,debug}.md`，包含命令、截图/日志示例。
- [x] T402 在仓库 README 与 `docs/develop/access-token/google/debug.md` 等入口增加聚光章节链接。
- [x] T403 记录验证证据（curl 输出、Redis key、Flow 回填截图）到 `specs/005-redbook-access-token/checklists/verification.md`。
- [x] T404 `go test ./cmd/accesstoken/server` + 手动调试回归，确保 Google/Bilibili 流程不受到影响。

## Dependencies & Notes

1. Phase 1 必须先完成以提供配置/Provider 元数据。
2. Phase 2 基于 Phase 1 的 Provider wiring，可与 Phase 3 并行起步但必须先完成 OAuth Flow。
3. Phase 3 endpoint 与 contract tests 完成后，方可进入 Phase 4 的文档更新与验证。
4. 全程需保持日志脱敏与 Flow 键命名规则，与现有 Provider 对齐。
