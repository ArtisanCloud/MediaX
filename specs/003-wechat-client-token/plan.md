# Implementation Plan: WeChat ClientToken Server

**Branch**: `003-wechat-client-token` | **Date**: 2025-12-17 | **Spec**: specs/003-wechat-client-token/spec.md  
**Input**: `docs/plan/wechat/client_token_client.md`

## Summary

目标是交付一个与 SessionToken/AccessToken 结构一致的独立 ClientToken Server（`cmd/clienttoken/server`，默认监听 `:7072`），用于管理微信公众号的 client credential token、API 调试和消息验证模拟。主要工作：

1. 扩展配置体系：在 `config.yaml` 中新增 `client_token_providers`，提供 appid/appsecret/消息 token/aes key/API token 等字段，并更新 `config.example.yaml`。
2. 新建 `cmd/clienttoken/server`：加载配置 → 初始化 Redis/Logger/MediaX 客户端 → 注册 HTTP 路由（`/debug` 页面、`/client-token/token`、`/client-token/cache`、`/client-token/call`、`/client-token/message/*`、`/healthz` 等）。
3. 构建调试页 HTML/JS，沿用 SessionToken/AccessToken 的风格，支持配置切换、token 刷新、API 调试、消息验证及回调日志查看。
4. 提供 CLI/脚本示例（`scripts/clienttoken-debug.sh` 或 `make clienttoken`），确保本地可以通过 `go run ./cmd/clienttoken/server -config config.yaml` 启动并调试。

## Technical Context

- **语言/版本**：Go 1.18（仓库 `go.mod`）  
- **依赖**：`github.com/ArtisanCloud/MediaXCore`（BaseClient/Logger/Cache）、`github.com/redis/go-redis/v9`（Token 缓存）、`gopkg.in/yaml.v3`（配置解析）  
- **存储**：Redis（ClientToken 缓存，强依赖，需要成功连接）  
- **服务入口**：`cmd/clienttoken/server`（与现有 `cmd/sessiontoken`/`cmd/accesstoken/server` 命令格式一致）  
- **调试页面**：Go 模板内嵌 HTML/JS（类似 `cmd/sessiontoken/debug_page.go`）  
- **日志**：`MediaXCore/pkg/logger`，输出 `clienttoken_metric`、`clienttoken_debug` 等结构化日志  
- **测试**：`go test ./cmd/clienttoken/...` + 手动访问 `/debug`，Redis 可使用本地 `127.0.0.1:6379`

## Constitution Check

- ✅ **Provider Adapter Parity**：仅复用 `pkg/client/wechat/officialAccount/clientTokenClient`，所有 API 调用走该适配器，不直连 HTTP。  
- ✅ **Config-Layered Security**：配置集中在 `config.yaml`，敏感字段通过环境变量覆盖，示例使用占位符，符合安全要求。  
- ✅ **Token Lifecycle Discipline**：Token 刷新由 clientTokenClient 完成，缓存策略统一在 Redis/内存层实现，遵循 TTL 规则。  
- ✅ **Observability & Error Traceability**：服务将复用 Logger/回调日志机制，`clienttoken_metric` 输出 action/provider/appid/status/latency，调试页显示请求/响应摘要。  
- ✅ **Testable Modularity**：`cmd/clienttoken/server` 结构与其他服务一致，handler 可单测；`debug_page.go` 为纯模板，可通过 e2e 访问验证。

## Project Structure Impact

```
cmd/
 ├─ clienttoken/
 │   ├─ server/main.go        # 新增：入口，解析 flag，调用 run()
 │   ├─ server/run.go         # 新增：加载配置、初始化 Redis/Logger
 │   └─ server/debug_page.go  # 新增：静态调试页模板
 ├─ accesstoken/...
 └─ sessiontoken/...

pkg/
 └─ client/config/            # 新增 ClientToken 配置结构

docs/
 └─ develop/wechat/client_token.md  # 调试/使用文档

specs/003-wechat-client-token/      # Spec/Plan/Tasks
```

## Detailed Plan & Tasks

### Phase 1 – 配置与工厂（约 1 天）
1. **config 定义**：在 `pkg/client/config` 新增 `ClientTokenProviderConfig`、`WechatClientTokenConfig`，支持 `appid/appsecret/message_token/message_aes_key/api_token/cache` 等字段，加载自 `config.yaml` 的 `client_token_providers`。
2. **MediaX 工厂**：在 `pkg/client/mediaX.go` 添加 `CreateWechatClientTokenClient`，返回 `pkg/client/wechat/officialAccount/clientTokenClient`，注入 Logger/Cache。
3. **config.example.yaml**：加入示例配置（使用占位符），并在 README/plan 中引用。

### Phase 2 – ClientToken Server 基础（约 2 天）
1. 新建 `cmd/clienttoken/server` 目录：`main.go` 解析 `-config/-port` 等 flag，调用 `run()`.
2. `run()`：加载配置 → 初始化 Redis（连接失败直接退出）→ 创建 Logger/Cache → 实例化 MediaX ClientTokenClient。
3. HTTP 路由：
   - `GET /healthz`
   - `GET /debug`（后续与调试页绑定）
   - `POST /client-token/token`
   - `GET /client-token/cache`
   - `POST /client-token/call`
   - `POST /client-token/message/validate`
   - `POST /client-token/message/callback`
4. API Token 鉴权中间件（默认 `dev-clienttoken`，可配置）。
5. 缓存模块：封装 `clientToken:wechat:<appid>` 的读写逻辑，包含 TTL、提前刷新策略。

### Phase 3 – 调试页 & CLI（约 2 天）
1. **调试页模板**（`debug_page.go`）：提供配置列表、Token 卡片、API 调试表单、消息验证区域、回调日志展示，交互逻辑用原生 JS（复用 SessionToken 页面模式）。
2. **回调日志存储**：参考 `cmd/accesstoken/server/callback_store.go`，记录最近 N 条消息（时间/Query/Headers/Body）。
3. **CLI/脚本示例**：新增 `scripts/clienttoken-debug.sh`（可由 `make clienttoken` 调用）并在文档中示例刷新 token/查看缓存/调用接口/验签。
4. **文档**：`docs/develop/wechat/client_token.md` 记录配置、启动命令 `go run ./cmd/clienttoken/server -config config.yaml`、调试步骤与常见问题。

### Phase 4 – 验收与打磨（约 1 天）
1. 联调微信公众号测试号：刷新 token、调用 `cgi-bin/user/get`、`message/custom/send` 等至少两个接口。
2. 验证消息验证/回调模拟流程（signature/AES）。
3. 审核日志脱敏、错误提示、配置示例一致性。

## Testing Strategy

- `go test ./cmd/clienttoken/...`：覆盖配置加载、缓存 helper、消息验证算法。
- 手动测试：
  1. 启动服务：`go run ./cmd/clienttoken/server -config config.yaml`
  2. `curl` 调用 `/client-token/token`、`/client-token/cache`、`/client-token/call`
  3. 浏览器访问 `/debug`，完成 token 刷新、API 调用、消息验证、回调日志查看。
- Redis：使用本地 `redis-server`（默认 db0），通过 `redis-cli` 检查 `clientToken:wechat:*` key。

## Risks & Mitigations

| 风险 | 说明 | Mitigation |
| --- | --- | --- |
| Redis 不可用 | ClientToken 依赖 Redis，启动失败会阻塞调试 | 启动前检测，提供明确错误信息；文档中说明需本地运行 Redis。 |
| 微信 API 限流/凭证失效 | 测试号 token 易过期 | 提供 CLI/脚本快速刷新；调试页展示 TTL/来源，便于定位。 |
| 调试页交互复杂 | UI 元素多（配置、API、消息验证） | 复用 SessionToken/AccessToken 的布局/组件，减少重复设计。 |

## Deliverables

1. `cmd/clienttoken/server` 代码（main/run/handlers/debug 页面/存储）。
2. `pkg/client/config` & `config.example.yaml` 更新、MediaX 工厂函数。
3. 调试页 HTML/JS + 回调日志组件。
4. CLI/脚本示例、开发文档 `docs/develop/wechat/client_token.md`。
5. Spec/Plan/Tasks 文档（`specs/003-wechat-client-token/`）。
