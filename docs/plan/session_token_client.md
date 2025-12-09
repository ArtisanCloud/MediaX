# SessionTokenClient 需求实现计划

## 1. 背景
- MediaX SDK 需要统一管理多平台（知乎、小红书、抖音等）的模拟登录流程，以便插件端能够拉起登录、抓取凭证、回调写库。
- 现有 `pkg/client/sessiontoken` 仅提供一个通用 HTTP 客户端，且唯一的配置 `SessionTokenClientConfig` 与各 provider 无关，无法覆盖多入口登录 URL、代理策略、凭证解析等差异化需求。
- 插件端已经定义了 `/session-token/flows` 接口契约，需要 SDK 侧适配器实现 Flow 生命周期管理、凭证采集和回调投递。

## 2. 目标
1. 在 MediaX SDK 内建立 SessionTokenClient 模式，支持 provider 级别的 Authenticator、CredentialHarvester、CallbackDispatcher。
2. 完成知乎 provider 的 SessionToken 流程实现，作为首个适配器，满足账号密码/扫码/手机号等多入口登录。
3. 调整配置体系：每个 provider 拥有独立的 SessionToken 配置块（含登录入口、代理、回调签名、脚本策略等）。
4. 更新 MediaX 客户端工厂，能够按 provider 创建对应 SessionTokenClient，并暴露 Flow 创建/查询能力给插件。
5. 提供安全、日志脱敏、回调签名等通用能力。

## 3. 现状问题
- `config.SessionTokenClientConfig` 无法表达 provider 差异。
- `MediaX.CreateSessionTokenClient` 只返回统一客户端，无法指定 provider 登录策略。
- 缺少 Flow 状态模型、状态持久化和生命周期管理实现。
- 日志与回调未做脱敏/签名。

## 4. 功能范围
| 模块 | 内容 |
| --- | --- |
| SDK 接口 | `POST/GET /session-token/flows` 请求结构、Flow 状态模型、结果标准化 |
| SessionTokenClient 抽象 | 定义 `Authenticator`、`CredentialHarvester`、`CallbackDispatcher` 接口，负责 URL 构建、凭证抓取、回调投递 |
| Provider 适配器 | 每个 provider（知乎等）实现上述接口，配置代理、入口类型、脚本注入 |
| 状态管理 | Flow 存储 `{flow_id, provider_code, status, authorize_url, expires_at, last_error, result}`；状态迁移 pending→authorizing→succeeded/failed |
| 安全 | API token 校验、回调签名/鉴权、日志脱敏、代理/IP 池支持 |

## 5. 接口设计
### 5.1 Flow 创建
- 请求：
```json
{
  "provider_code": "zhihu",
  "provider_app_code": "zhihu_article",
  "account_id": 123,
  "tenant_uuid": "...",
  "state": "<plugin-state>",
  "callback_url": "https://plugin/api/v1/admin/platforms/session-token/callback",
  "metadata": {"login_mode": "password"}
}
```
- 响应：
```json
{
  "flow": {
    "flow_id": "stf_abc123",
    "status": "pending",
    "authorize_url": "https://www.zhihu.com/signin?sessionflow=xxx",
    "expires_at": "2025-12-08T12:35:00Z"
  }
}
```

### 5.2 Flow 查询
- `GET /session-token/flows/{flow_id}`，返回 `status/message/result`，供插件轮询。

### 5.3 标准化凭证
```json
{
  "session_token": "...",
  "cookies_json": [{"name":"z_c0","value":"..."}],
  "headers_json": {"User-Agent":"...", "X-XSRF-TOKEN":"..."},
  "expires_at": "2025-12-09T00:00:00Z",
  "note": "手机号登陆",
  "captured_at": "2025-12-08T11:33:21Z"
}
```

## 6. 技术方案
### 6.1 配置结构调整
- 在每个 provider 配置文件（例：`pkg/client/config/zhihu.go`、`redbook.go` 等）增加 `SessionTokenConfig`：
  - `Service`：与 SessionToken 服务交互的 base_url、api_token。
  - `Authenticator`：登录入口列表（pc/h5/扫码）、默认 UA、脚本注入资源、验证码策略。
  - `Harvester`：需要监听的 Cookie/Header/Storage key、抓取脚本描述。
  - `Callback`：目标 URL、签名 secret、重试策略。
  - `Network`：代理池、IP 池、超时。

### 6.2 SDK 抽象
- 定义：
  ```go
  type SessionTokenClient interface {
      CreateFlow(ctx context.Context, req *CreateFlowOptions) (*Flow, error)
      GetFlow(ctx context.Context, flowID string) (*Flow, error)
  }
  type Authenticator interface { BuildAuthorizeURL(*FlowContext) (string, error) }
  type CredentialHarvester interface { Watch(*FlowContext) (*Credentials, error) }
  type CallbackDispatcher interface { Dispatch(ctx context.Context, payload CallbackPayload) error }
  ```
- `sessiontoken.Manager`：负责 Flow 状态机、持久化、调用 provider adapter。

### 6.3 Flow 状态机
- `pending`：Flow 创建，生成 `authorize_url`。
- `authorizing`：Authenticator 驱动浏览器容器进行登录。
- `succeeded`：CredentialHarvester 抓取凭证后写入 `result`，触发回调。
- `failed`：登录失败、超时或手工终止，记录 `message/last_error`。

### 6.3.1 Flow 存储与生命周期
- Flow 记录的是“某次模拟登录会话”的上下文，字段包含 `flow_id`、provider、租户、登录模式、状态、结果、回调与错误信息；需要存入可持久化的存储（Redis/DB）以便多实例共享、重启恢复以及插件轮询查询。
- Flow 生命周期通常跨多个 HTTP 请求（插件创建 → 浏览器容器交互 → 凭证回调），因此存储层需支持 TTL/状态索引，默认在 `succeeded/failed` 后保留一段时间供审计/回查，然后按策略清理。
- Flow 存储不负责长期保存 session/cookie：当 Flow 成功时，凭证由插件/业务方写入自有 session 表或缓存；凭证过期后，业务再调用 `POST /session-token/flows` 创建新的 Flow 以触发下一次模拟登录。
- 由于 Flow 具有可恢复的状态机，任何实例都可以根据存储中的状态继续回调重试、记录事件或查询结果，避免耦合到单个 SDK 进程。

### 6.4 安全 & 日志
- API token 必填；请求头 `Authorization: Bearer` 验证。
- callback payload：`{state, flowId, status, credentials, metadata}`，附签名（HMAC-SHA256）。
- 日志脱敏：`session_token`、`cookies_json` 仅展示前后若干字符。
- 支持可选代理池配置，避免同 IP 被识别。

## 7. 里程碑
1. **配置与抽象**（Day 1-2）
   - 删除旧 `SessionTokenClientConfig`，新增 provider 专属配置。
   - 定义 SessionToken 抽象接口与 Flow 模型。
2. **知乎适配器**（Day 3-5）
   - 实现 `ZhihuSessionTokenClient`（Authenticator/Harvester/Callback）。
   - 接入代理、登录入口配置。
3. **API 层与存储**（Day 6-7）
   - 实现 Flow 存储、状态机、API token 校验、日志脱敏。
   - 完成回调签名与重试逻辑。
4. **测试 & 文档**（Day 8）
   - 单元/集成测试（Flow 状态、回调 payload）。
   - 文档更新（README/示例）。

## 8. 风险与缓解
- **浏览器容器稳定性**：引入守护与超时，失败后自动更新 Flow 状态。
- **代理池可用性**：加健康检查与回退机制。
- **凭证敏感信息泄露**：日志与事件统一脱敏；回调仅通过 HTTPS。

## 9. 成功标准
- 插件端能够通过 `provider_code=zhihu` 发起 Flow，成功拉起登录并在 95% 情况下获取凭证。
- Flow 状态 api 支持轮询，能准确反映 pending/authorizing/failed/succeeded。
- 回调 payload 满足插件校验（state & signature）。
