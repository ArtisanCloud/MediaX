# Feature Specification: WeChat ClientToken Server

**Feature Branch**: `003-wechat-client-token`  
**Created**: 2025-12-17  
**Status**: Draft  
**Input**: docs/plan/wechat/client_token_client.md – “WeChat ClientToken 模式实施计划”

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 – 独立 ClientToken Server 启动 (Priority: P1)
作为 SDK 维护者，我需要一个独立的 `cmd/clienttoken` 服务，`go run ./cmd/clienttoken -config clienttoken.yaml` 即可启动，默认监听 `:7072` 并提供 API Token 鉴权及 `/debug` 页面，以便 QA/产品在浏览器中调试公众号接口。

**Independent Test**: 不依赖 AccessToken/SessionToken server；单独运行 `go run ./cmd/clienttoken -config clienttoken.yaml`，浏览器访问 `http://127.0.0.1:7072/debug` 验证 UI。

**Acceptance Scenarios**
1. **Given** 准备好的 `clienttoken.yaml`，**When** 执行 `go run ./cmd/clienttoken -config clienttoken.yaml`，**Then** 日志打印 `listening :7072`，`GET /healthz` 返回 200。
2. **Given** 浏览器访问调试页，**When** 未携带 `Authorization: Bearer <api_token>` 调接口，**Then** 后端返回 401；携带正确 token 时接口成功响应。

---

### User Story 2 – Client Credential Token 生命周期 (Priority: P1)
作为运维，我希望在调试页或 REST API 中一键刷新公众平台 `access_token`，并缓存到 Redis（key `clientToken:wechat:<appid>`），以便复用同一 access_token 调试多种接口。

**Independent Test**: 使用 `POST /client-token/token` 接口即可验证，无需实际调用业务接口。

**Acceptance Scenarios**
1. **Given** 配置了 `appid/appsecret`，**When** 调用 `POST /client-token/token`，**Then** 返回 JSON `{"access_token":"...","ttl":7200,"cached_at":"..."}`，并在 Redis 中写入记录（TTL 至少 7000 秒）。
2. **Given** Redis 中已有未过期 token，**When** 调用 `GET /client-token/cache?appid=xxx`，**Then** 直接返回缓存命中结果且不会重复刷新。
3. **Given** Token 即将过期（TTL < 600 秒），**When** 再次调用刷新接口，**Then** 系统主动重新请求微信 API，并更新缓存。

---

### User Story 3 – 浏览器调试页面与常用 API (Priority: P2)
作为 QA，我希望在 `/debug` 页面选择公众号配置、查看当前 token、填写 `action/method/query/body` 并点击“调用 API”，以便快速验证粉丝列表、消息发送、素材操作等常见场景；同时可模拟微信消息验证与回调。

**Independent Test**: 仅需启动 ClientToken server，打开调试页操作即可，不依赖 CLI。

**Acceptance Scenarios**
1. **Given** 调试页列出配置列表，**When** 切换到不同 app，**Then** 页面立即更新 AppID、当前 token、缓存状态，并能点击“刷新 Token”获取新值。
2. **Given** 在 API 表单填写 `action=cgi-bin/user/get`，**When** 点击“调用 API”，**Then** 后端调用 `pkg/client/wechat/officialAccount/clientTokenClient` 并把响应 JSON 展示在页面下方，同时输出 `clienttoken_metric` 日志。
3. **Given** QA 在“消息验证”区域输入 `signature/timestamp/nonce/echostr`，**When** 点击“验证”，**Then** 服务按配置 `message_token` 校验签名并返回 echostr；如果配置了 `message_aes_key`，还可在“Mock 消息回调”区上传加密 XML/JSON 并查看解密结果。

---

## Requirements *(mandatory)*

### Functional Requirements
1. **FR-001**: 新增 `clienttoken.yaml` 配置块（或 `client_token_providers` 节点），字段必须覆盖 `appid`, `appsecret`, `message_token`, `message_aes_key`, `api_token`, `cache.redis` 等，且与 `pkg/client/config` 的结构保持同步。
2. **FR-002**: `cmd/clienttoken` MUST 支持 `-config` 参数、`CLIENTTOKEN_CONFIG` 环境变量，并在日志中提示监听端口、API token。
3. **FR-003**: Server MUST 暴露 `POST /client-token/token`、`GET /client-token/cache`、`POST /client-token/call`、`POST /client-token/message/validate`、`POST /client-token/message/callback`、`GET /debug`、`GET /healthz`。
4. **FR-004**: 所有 REST 接口 MUST 通过 `Authorization: Bearer <API_TOKEN>` 或 `X-API-Token` 鉴权，默认 token `dev-clienttoken` 可通过配置覆盖。
5. **FR-005**: Cache MUST 采用 Redis （key `clientToken:wechat:<appid>`） + 进程内内存，默认 TTL 7000 秒；刷新失败时应保留旧值并返回错误信息。
6. **FR-006**: 调试页 MUST 用浏览器就绪（HTML/JS），可显示配置列表、当前 token、缓存 TTL、刷新按钮，并支持 API 调试/消息验证/回调日志等交互。
7. **FR-007**: API 调用 MUST 复用 `pkg/client/wechat/officialAccount/clientTokenClient`，禁止在 handler 中直接拼装 HTTP 请求。
8. **FR-008**: `clienttoken_metric` 日志 MUST 包含 `action/provider/appid/status/latency_ms/token_source`，并对 token 进行脱敏。
9. **FR-009**: CLI/脚本 示例 MUST 覆盖 `go run ./cmd/clienttoken -config clienttoken.yaml`、`scripts/clienttoken-refresh.sh` 等命令，确保复制即可运行。
10. **FR-010**: 消息验证 MUST 按微信官方算法校验 `signature`，若配置了 `message_aes_key`，Mock 回调 MUST 支持加解密并展示结果。

### Key Entities
- **ClientTokenProviderConfig**: 描述公众号的 AppID/AppSecret、消息 Token/AESKey、缓存策略、API token。
- **ClientToken Cache Record**: `{appid, access_token, ttl_seconds, stored_at, expire_at, source}` – 存在 Redis 与内存中。
- **ClientToken Request Payload**: `{"action":"cgi-bin/user/get","method":"GET","query":"next_openid=xxx","body":"{}"}` – 由调试页/调用方提交，Server 负责拼接 URL 与注入 `access_token`。
- **Message Validation Payload**: `{"appid":"...","signature":"...","timestamp":"...","nonce":"...","echostr":"..."}` – 用于模拟微信服务器验证。

---

## MediaX Architecture Guardrails *(reference Constitution)*
- **Provider Adapter Parity**: 所有接口调用必须通过 `pkg/client/wechat/officialAccount/clientTokenClient`，不得直连微信 API；保持与 MediaX provider 工厂一致。
- **Config-Layered Security**: 不得在仓库中硬编码真实 AppID/AppSecret，示例需使用占位符；敏感配置从 `clienttoken.yaml` 或环境变量注入。
- **Cache & Token Discipline**: Token 获取必须走 `clientTokenClient`，缓存逻辑要遵循 TTL/提前刷新策略；日志/响应中需对 token 脱敏。
- **Observability**: 依赖 `MediaXCore/pkg/logger` 打印结构化日志，并在 `/debug` 页面呈现回调/调用信息；错误需附带上下文（provider/appid/action/error_reason）。
- **Testability**: `cmd/clienttoken` 架构应与现有服务一致（`main.go` 调用 `run()`，handler 可单测），调试页纯静态文件，便于独立验证。

---

## Success Criteria *(mandatory)*
1. **SC-001**: 新成员在 10 分钟内按照文档完成 `clienttoken.yaml` 配置、运行 `go run ./cmd/clienttoken -config clienttoken.yaml`，即可访问调试页。
2. **SC-002**: 使用 `/client-token/token` + `/client-token/cache` 接口时，Redis 中必然存在 `clientToken:wechat:<appid>` 记录，刷新 TTL 正常。
3. **SC-003**: 调试页调用 `cgi-bin/user/get`、`message/custom/send` 至少两个接口全部成功，响应与日志均可见。
4. **SC-004**: 消息验证/回调模拟功能可在 5 分钟内完成演示，包含 signature 校验、AES 解密、回显 payload。
5. **SC-005**: CLI/脚本示例经实测可直接运行，无额外依赖；文档中的所有命令成功率 ≥ 90%。

---

## Assumptions & Dependencies
- 可用的微信公众平台测试号（AppID/AppSecret/消息 Token/AESKey）。
- Redis 服务可访问（默认 `127.0.0.1:6379`，可通过环境变量覆盖）。
- `pkg/client/wechat/officialAccount/clientTokenClient` 已具备调用基础；如需扩展接口，应先在该包提供方法。
- 调试页沿用 SessionToken/AccessToken 的样式与构建方式（纯静态模板 + 内嵌 JS）。
- 网络需访问 `https://api.weixin.qq.com`；若需代理，server/客户端支持读取 `HTTP_PROXY` 等配置。
