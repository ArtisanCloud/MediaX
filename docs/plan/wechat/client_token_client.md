# WeChat ClientToken 模式实施计划

## 1. 背景
- MediaX AccessToken 调试台目前只覆盖 OAuth 场景（以 Google YouTube 为首），需要人工走 OAuth Flow + 回调，流程长、数据依赖 Redis Flow。
- 对于微信公众平台（Official Account），常见调用只需要 AppID + AppSecret 即可换取 `access_token`，无需用户交互。SDK 已有 `pkg/client/wechat/officialAccount/clientTokenClient`，但 AccessToken 调试台、配置和缓存流程尚未接入。
- 业务需求：运维或 QA 在本地输入一组公众号配置（AppID/AppSecret/消息加解密参数）即可实时生成 AccessToken、调试消息推送、调用公众号开放接口。

## 2. 目标
1. 新增独立的 **ClientToken Server**（默认 `cmd/clienttoken/server`，监听 `:7072`），支持微信公众平台官方的 Client Credential Flow，与 AccessToken Server 解耦。
2. 配置沿用现有 `config.yaml`，在其中新增 `client_token_providers` 节点为多个公众号配置 AppID、AppSecret、消息验证 token、消息 aes key，调试台可以直接切换（启动命令与其他服务保持一致：`go run ./cmd/clienttoken/server -config config.yaml`）。
3. ClientToken Server 提供 REST API/调试页用于：
   - 解析当前 access_token（来源、TTL、最近缓存时间）
   - 主动刷新/缓存 access_token（Redis 缓存 key 规范）
   - 调试公众号接口，例如自定义菜单、模板消息发送、素材管理等。
4. 保留 Redis 缓存、日志脱敏、API token 校验等既有基础设施。

## 3. 现状问题
- AccessToken 服务的 provider metadata 没有 “wechat_official_account” 模式，也无法渲染对应配置。
- `clientTokenClient` 只能在 pkg 层单独使用；没有 HTTP 服务/CLI 接入，也没有缓存 TTL、重试机制。
- 调试 UI 只适配 Google OAuth，缺少“直接用 AppSecret 换 token”的入口。

## 4. 功能范围
| 模块 | 内容 |
| --- | --- |
| 配置 | 独立的 `clienttoken.yaml`（或 `config.yaml` 的 `client_token_providers`）中新增 `provider_code=wechat_official_account`，app 节点包含 `appid/appsecret/messagetoken/messageaeskey` 等字段。 |
| SDK | 在 `client/providers.go` 增加“创建 Wechat ClientToken Client”的工厂，使用 `pkg/client/wechat/officialAccount/clientTokenClient`。 |
| 服务层 | 新建 `cmd/clienttoken` 服务，提供 `POST /client-token/token`（主动刷新）、`GET /client-token/cache`（查看缓存）、`POST /client-token/call`（调试公众号 API）、`POST /client-token/message/validate` 等接口。 |
| 缓存 | Redis key 规范：`clientToken:wechat:<appid>` 保存 token 与过期时间；支持缓存时长 > 2 小时，提前 10 分钟刷新。 |
| 调试页 | ClientToken Server 自带 `/debug` 页面：展示配置下拉、显示当前 token、有刷新按钮，以及常用 API（例如获取粉丝列表、发送文本消息等）。 |
| 安全&日志 | 复用 SessionToken/AccessToken server 的 API Token 鉴权方案；日志记录 `clienttoken_metric`（action/provider/app/status/latency）；避免打印完整 token。 |

## 5. 配置示例
```yaml
# config.yaml 示例
client_token_providers:
  - provider_code: wechat_official_account
    name: 示例公众号
    apps:
      - code: demo_official
        name: Demo Official
        default_mode: client_credential
        client_token:
          appid: wx_demo_appid_123456
          appsecret: demo_appsecret_abcdef1234567890
          message_token: demoMessageToken123
          message_aes_key: demoMessageAESKey1234567890abcdef
        cache:
          redis_key: clientToken:wechat:wx_demo_appid_123456
          ttl_seconds: 7000
```

## 6. 技术方案
### 6.1 Client 工厂
- 在 `pkg/client/mediaX.go` 增加 `CreateWechatClientTokenClient(cfg *config.WechatOfficialAccountConfig)`，内部复用 `pkg/client/wechat/officialAccount/clientTokenClient`。
- 统一提供 logger / cache / http proxy 设置。

### 6.2 ClientToken Server
- 新建 `cmd/clienttoken` 进程，结构参考 AccessToken/SessionToken server：启动时加载配置、初始化 Redis+logger+MediaX 客户端，并注册 HTTP 路由。
- 主要路由：
  - `GET /debug`：静态调试页（纯前端 HTML/JS）。
  - `POST /client-token/token`：根据 provider/app，从配置中读取 `appid/appsecret`，调用 clientTokenClient 拉取新 token，写入 Redis，并返回 `access_token`、`ttl`、`cached_at`。
  - `GET /client-token/cache`：读取 Redis key，查看当前 token 及剩余有效期（若无则尝试内存缓存）。
  - `POST /client-token/call`：请求体包含 `action`（例如 `cgi-bin/user/get`）、`method`、`query/body`，服务端注入 `access_token` 并通过 clientTokenClient 发起请求。
  - `POST /client-token/message/validate`：输入 `signature/timestamp/nonce/echostr`，按配置的 `message_token` 验签，返回 echo 字符串。
  - `POST /client-token/message/callback`：模拟微信服务器推送事件/消息，可结合消息 AESKey 做加解密展示。
- Redis 操作复用现有 cache 封装，key 规范 `clientToken:wechat:<appid>`；服务启动时必须成功连接 Redis（不可 fallback 到内存），以保证本地调试与生产缓存策略一致。

### 6.3 调试页 UI
- 结构：
  1. 配置选择：Provider/App 列表 + 当前 AppID。
  2. Token 信息卡片：显示 `access_token`（脱敏）、`来源`（cache/refresh）、`过期时间`、刷新按钮。
  3. API 调试表单：预置常用接口（粉丝列表、群发消息、素材上传等），支持自定义 `path/query/body`。
  4. 消息验证：输入 `signature/timestamp/nonce/echostr` 可模拟系统验证。
- 调用接口均携带 `Authorization: Bearer <api_token>`，复用现有 `apiBase`。

### 6.4 CLI/脚本
- 提供 `make clienttoken ARGS='...'` 或 `scripts/clienttoken-refresh.sh`：读取 `config.yaml` 主动刷新 token、打印 TTL；支持 `--call cgi-bin/user/get` 等一键调用示例（命令默认引用 `go run ./cmd/clienttoken/server -config config.yaml` 同一配置）。

### 6.5 安全与观测
- Token 刷新、API 调用均打印结构化日志 `clienttoken_metric`：
  ```
  clienttoken_metric: action=refresh provider=wechat_official account=demo_official status=success latency_ms=120 ttl=7000
  ```
- 错误日志脱敏：`token=wx***a2a`。

### 6.6 启动与调试
- 开发阶段推荐直接使用 `go run` 启动：
  ```bash
  go run ./cmd/clienttoken/server -config config.yaml
  ```
  默认监听 `:7072`，调试页访问 `http://127.0.0.1:7072/debug`。
  如需编译二进制，可执行 `go build -o bin/clienttoken ./cmd/clienttoken/server` 后运行 `./bin/clienttoken -config config.yaml`。

## 7. 交付项
1. **配置**：`config.example.yaml` 中新增 `client_token_providers` 示例。
2. **SDK**：新增创建 WeChat ClientToken Client 的工厂、缓存工具。
3. **服务端**：`cmd/clienttoken`（main/handler/缓存/log）。
4. **前端**：ClientToken Server `/debug` 页面（HTML/JS/CSS）。
5. **文档**：`docs/develop/access-token/wechat/README.md`（说明如何启动/调用）、`docs/develop/wechat/client_token.md`（调试步骤）。
6. **脚本/Makefile**：一键刷新命令或 curl 示例。

## 8. 里程碑
| 阶段 | 说明 | 预估 |
| --- | --- | --- |
| M1 配置与工厂 | 定义配置、提供 `CreateWechatClientTokenClient` | 1 天 |
| M2 服务端接口 | Handler、Redis 缓存、API Token 校验 | 2 天 |
| M3 调试页 & CLI | Web UI、脚本、文档 | 2 天 |
| M4 验收 | 联调公众号测试号、验证消息推送/接口调用 | 1 天 |
