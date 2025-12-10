# Creative Channels - SessionToken 联调指引

> 参考：`docs/plan/mediax-sdk.md` 中的「SessionTokenClient 对接规划」。本文件聚焦于 Creative/Channels 插件如何与 MediaX SessionToken 服务共享配置与部署。

## 1. 服务依赖与入口

- SessionToken 服务由 `cmd/sessiontoken/main.go` 驱动，默认监听 `:7070`。
- 启动命令：
  ```bash
  make sessiontoken             # 默认读取 config.yaml
  # 或
  go run ./cmd/sessiontoken -config config.yaml
  ```
- 注册路由：`/session-token/flows` (POST/GET)，全部挂载在 `RegisterSessionTokenFlowCreateRoute` 与 `RegisterSessionTokenFlowGetRoute` 中，强制校验 Bearer Token。
- 服务启动后会输出 `sessiontoken_metric` 和 `sessiontoken_callback` 日志，用于观察 Flow 创建/查询/回调状态。

## 2. 配置对齐

### 2.1 环境变量映射

| 变量 | 示例 | 说明 |
| --- | --- | --- |
| `POWERX_SESSION_TOKEN_BASE_URL` | `http://127.0.0.1:7070` | 插件访问 SessionToken 服务的 BaseURL，需与服务监听地址一致。 |
| `POWERX_SESSION_TOKEN_API_TOKEN` | `dev-session-token` | 插件调用 `/session-token/flows` 所用的 Bearer Token；服务端 `SESSIONTOKEN_API_TOKEN`/`POWERX_SESSION_TOKEN_API_TOKEN` 读取该值。 |
| `POWERX_SESSION_TOKEN_CALLBACK_URL` | `https://plugin.local/api/v1/admin/platforms/session-token/callback` | 插件暴露的凭证回调入口，Flow 创建时写入 `callback_url` 字段。 |
| `SESSIONTOKEN_LISTEN_ADDR` | `:7070` | （可选）覆盖服务监听地址，默认 `:7070`。 |
| `SESSIONTOKEN_REDIS_ADDR` | `127.0.0.1:6379` | FlowStore 使用的 Redis 地址，服务与插件必须可访问该实例。 |
| `SESSIONTOKEN_LOG_LEVEL` | `info` | 控制 SessionToken 服务日志级别。 |

> **约定**：PowerX/插件侧仅依赖 `POWERX_SESSION_TOKEN_*` 变量；MediaX 服务可读取 `SESSIONTOKEN_*` 或相同的 `POWERX_*` 值，但最终必须与插件保持一致。

### 2.2 其他配置

- `config.yaml` 中的 `zhihu_config.sessionToken.service.api_token` 与 `POWERX_SESSION_TOKEN_API_TOKEN` 保持一致。
- `callback_url` 由插件构造，建议读取 `POWERX_SESSION_TOKEN_CALLBACK_URL` 形成 `https://<plugin-domain>/admin/platforms/session-token/callback`。

## 3. 启动顺序

1. **启动 Redis**：确保 `SESSIONTOKEN_REDIS_ADDR` 可用，并具备 TTL 支持。
2. **启动 SessionToken 服务**：运行 `make sessiontoken`，确认日志中出现 `sessiontoken: server listening addr=...`。
3. **启动插件 / 浏览器容器**：加载 `POWERX_SESSION_TOKEN_*` 后启动 Creative/Channels 插件或 MediaX Studio。
4. **触发模拟登录**：在插件 UI（如 `/publish/platforms`）发起 Flow 创建，随后可在 `sessiontoken_metric` 日志中看到 `action=create_flow` 记录。

若插件先于服务启动，会导致 5xx/401；因此务必遵循上述顺序。

## 4. 回调与监控

- Flow 成功/失败后，服务会根据 `callback_url` 向插件 POST 标准化凭证，payload 包含 `state`、`flow_id`、`status`、`credentials` 以及 HMAC-SHA256 签名。
- 日志样例：
  ```
  sessiontoken_metric: action=create_flow provider=zhihu tenant_uuid=tenant_x flow_id=stf_xxx status=pending latency_ms=14 retry=0
  sessiontoken_callback: success provider=zhihu tenant_uuid=tenant_x flow_id=stf_xxx state=ui-flow retry=0 http_status=200 latency_ms=6
  ```
- 插件应将回调结果写入自身 Session 存储，以便频道刷新/任务执行模块消费（参见 `docs/plan/mediax-sdk.md#3`）。

## 5. 故障排查

- `401 Unauthorized`：确认插件与服务的 API Token 完全一致，并在请求头中携带 `Authorization: Bearer <token>`。
- `flow not found`：Flow TTL 默认 30 分钟 + 审计保留 6 小时，若超时需重新创建。
- 回调超时：检查插件 `POWERX_SESSION_TOKEN_CALLBACK_URL` 是否可通过 HTTPS 访问，并确保返回 2xx；服务默认 3 次（2s→4s→8s）重试。

> 若文档或部署方式更新，请同步修改此文件与 `docs/plan/mediax-sdk.md`，保持插件/SDK 对齐。
