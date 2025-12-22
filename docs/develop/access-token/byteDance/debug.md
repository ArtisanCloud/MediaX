# DouYin AccessToken 调试指南

> 与 Google/BiliBili 的 `/debug` 调试台一致，只是 Provider/App/Mode 换成 “字节跳动 / 抖音 Douyin”。本指南聚焦本地调试流程、Redis Flow 管理与常见问题。

## 1. 启动服务

```bash
cd /private/var/www/html/ArtisanCloud/X/MediaX/core/MediaX
export ACCESSTOKEN_CONFIG=${ACCESSTOKEN_CONFIG:-config.yaml}
export ACCESSTOKEN_LISTEN_ADDR=${ACCESSTOKEN_LISTEN_ADDR:-127.0.0.1:7071}
export ACCESSTOKEN_API_TOKEN=${ACCESSTOKEN_API_TOKEN:-dev-accesstoken}
go run ./cmd/accesstoken/server -config "$ACCESSTOKEN_CONFIG"
```

日志应看到：`accesstoken-server: event=listening ... storage_backend=redis|memory provider=byte_dance`。

## 2. 配置要求

- `config.yaml` 必须包含 DouYin Provider，详见 `docs/develop/access-token/byteDance/develop.md`。
- `redirect_url` 必须指向服务监听地址的 `/debug/callback`。
- 若使用 `${DOUYIN_*}` 占位符，记得在 shell 中 export。

## 3. 调试台操作

访问 <http://127.0.0.1:7071/debug?api_token=dev-accesstoken>：

1. **选择 Provider**：`Provider → 字节跳动`，`App → 抖音 Douyin`，`授权模式 → default`。
2. **同步模板**：点击“同步模板”按钮，右侧 JSON 会自动填入 `provider_code=byte_dance_douyin` 等字段。
3. **发起授权**：点击“发起授权”，浏览器会跳到 DouYin 登录/授权页；登录后回调 `/debug/callback`。
4. **查看 Flow**：底部“授权记录”会看到 `flow_id=oauth-xxx`，可点击“填充”把 token 写回 JSON；也可在“Flow ID 回填”输入 `oauth-xxx` → “加载 Flow ID” → 读取 Redis。若 Redis 不可用，界面会高亮 `storage_backend=memory` 提示重新授权。
5. **解析 Token**：在 “AccessToken 解析” 区块点击“解析”即可查看 token 来源（request/env/config/flow）与 TTL，便于排障。
6. **API 调试**：手动编辑 `action/payload` 并点击“执行 API”，即可触发 `/accesstoken/call`。响应面板会显示 `retry_count`、`last_backoff_ms`、`flow_id` 等信息，便于确认是否发生刷新或重试。

> Flow 列表中若出现 “需重新授权” 标记或 `status=need reauth`，表示自动刷新失败，必须重新发起 OAuth。

## 4. Redis Flow 结构

- 主 key：`accesstoken:oauth:byte_dance_douyin:default:<mode>`。
- 索引 key：`accesstoken:oauth:flow:<flow_id>`。
- TTL：默认 86,400 秒，可通过 `ACCESSTOKEN_FLOW_TTL_SECONDS` 调整。

```bash
redis-cli --raw GET accesstoken:oauth:flow:oauth-abc | xargs redis-cli --raw GET
```

## 5. API 调试 & 速率限制

1. 每个 DouYin action 配置了 1 QPS 的 `rate.Limiter`。若同一 action 连续调用，两次请求之间间隔 <1s，第二次会立即返回 `429`，UI toast 会提示“达到 1 QPS 限制”。
2. `429/5xx` 会触发最多 3 次指数退避重试（200ms 起倍增），响应 JSON 中的 `retry_count`、`last_backoff_ms` 以及日志 `event=token.call retry_count=...` 可用于确认实际重试次数。
3. 如果 AccessToken 剩余时间 ≤5 分钟，调试台执行 API 前会自动刷新 Flow。刷新成功时响应仍为 200，Flow 列表的 `last_refresh_at` 与 `status=refreshed` 会更新。
4. 刷新失败或 Flow 缺失时，响应返回 `{"error":"need reauth"}`，调试台会弹出提示并在 Flow 列表标注“需重新授权”。

## 6. 常见问题

| 场景 | 解决办法 |
| --- | --- |
| 页面显示 `${DOUYIN_CLIENT_ID}` | 未设置环境变量，或 `config.yaml` 未展开，确保 shell 已 export 对应变量后重启服务。 |
| 授权跳转到登录页 | 首次 OAuth 需要先登录 DouYin 主体账号，登录成功后就会继续授权，这是官方流程。 |
| 提示 `unsupported provider` | 检查 `provider_code` 是否填 `byte_dance_douyin`，以及 `auth_modes` 是否包含默认 key。 |
| `storage_backend=memory` 警告 | Redis 未连上，Flow 仅存在内存；设置 `ACCESSTOKEN_REDIS_ADDR` 后重启即可。 |
| 授权后 Flow 消失 / need reauth | Flow 过期或 refresh_token 失效，再次发起 OAuth 即可；可通过 Flow 表的 `flow_id` 交叉检查 `redis-cli`。 |
| 429 频繁出现 | 降低相同 action 的调用频率，或在脚本中增加 `sleep`，因为调试服务每 action 仅允许 1 QPS。 |

> 更多环境/配置细节请参阅 `docs/develop/access-token/byteDance/develop.md`。倘若需要 CLI 调试，可执行 `go run ./cmd/accesstoken -provider byte_dance -provider-app douyin ...`，调试台生成的 JSON 也可直接复制到 CLI 中使用。

## 7. SC-004 走查提示

- 走查开始前，请为被访者提供 `specs/006-bytedance-access-token/quickstart.md` 与本调试文档的链接，其他帮助暂时隐藏。
- 观察其是否能独立完成：启动服务 → `/debug` 发起授权 → 查看 Flow → 复制 JSON → 执行 `/accesstoken/call`。
- 完成后，将耗时、遇到的问题与截图补充到 `docs/develop/access-token/byteDance/develop.md` 的“文档走查记录”，必要时在 Quickstart 中追加排障说明。
