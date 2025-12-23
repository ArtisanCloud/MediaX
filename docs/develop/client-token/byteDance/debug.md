# DouYin ClientToken 调试指南

> 目标：在 `cmd/clienttoken/server` 的 `/client-token/debug` 页面完成抖音 ClientToken 刷新、缓存检查与 API 调试，体验与微信 ClientToken 调试台一致。

## 1. 启动服务

```bash
cd /private/var/www/html/ArtisanCloud/X/MediaX/core/MediaX
export CLIENTTOKEN_CONFIG=${CLIENTTOKEN_CONFIG:-config.yaml}
export CLIENTTOKEN_LISTEN_ADDR=${CLIENTTOKEN_LISTEN_ADDR:-:7072}
export CLIENTTOKEN_API_TOKEN=${CLIENTTOKEN_API_TOKEN:-dev-clienttoken}
go run ./cmd/clienttoken/server -config "$CLIENTTOKEN_CONFIG"
```

日志应包含 `clienttoken: server listening ... provider=byte_dance`。确保 `CLIENTTOKEN_REDIS_ADDR` 指向可访问的 Redis，否则会降级为内存缓存。

## 2. 浏览器操作

访问 <http://127.0.0.1:7072/debug?api_token=dev-clienttoken>：

1. **选择 Provider/App/Mode**：`Provider → 字节跳动`、`App → 抖音服务端`、`模式 → default`。
2. **同步模板**：点击“同步模板”填充 JSON（包含 `client_token.redis_key`、`client_key` 等信息）。
3. **刷新 Token**：点击“刷新 Token” → 服务端调用 DouYin client_token API，并将结果写入 Redis（key 形如 `clientToken:douyin:<client_key>`）。
4. **查看缓存**：点击“查看缓存”可看到 `access_token`、TTL、来源（redis/memory）。
5. **调用 API**：在 API 调试面板填写 `action`（如 `open_api/1/content/video/list/`）、请求方式与 JSON Body，点击“调用 API”后会自动附带最新 client_token。若页面提示 “Douyin API 模板已禁用”，说明 `device_id`/`risk_info` 尚未同时配置，可在 `config.yaml` 或环境变量中补齐后重启服务。
6. **回调日志**：若外部服务向 `/client-token/callback` 发送事件，页面最下方会显示收集到的请求，可用于调试第三方回调逻辑。

> **Risk 字段提示**：为了避免误用需风控字段的接口，调试台会自动检测 `byte_dance_douyin_config.device_id` 与 `risk_info`。仅当两个字段均已设置时，Douyin API 模板才会启用，否则会显示禁用提示。

## 3. Redis 缓存

- 默认 key：`clientToken:douyin:<client_key>`。
- TTL：配置中的 `ttl_seconds`（示例 7000 秒），在剩余 `refresh_before_seconds` 时会自动刷新。
- 检查命令：`redis-cli --raw GET clientToken:douyin:<client_key> | jq '.'`。

## 4. CLI/HTTP 验证

脚本 `scripts/clienttoken-douyin.sh` 帮你封装了常用操作：

```bash
CLIENTTOKEN_API_TOKEN=dev-clienttoken ./scripts/clienttoken-douyin.sh refresh
./scripts/clienttoken-douyin.sh cache show
./scripts/clienttoken-douyin.sh cache clear
./scripts/clienttoken-douyin.sh call open_api/1/content/video/list/ POST "" '{"page":1,"size":10}'
```

脚本输出与页面一致，若遇到 401/502，可结合 `logs/clienttoken-server-info.log` 中的 `clienttoken_event` 日志定位问题。

## 5. 常见问题

| 场景 | 解决办法 |
| --- | --- |
| 页面没有字节跳动 Provider | 检查 `client_token_providers` 是否包含对应配置，或重启服务重载配置。 |
| 刷新失败 `invalid client_key` | 确认环境变量提供的 key/secret 与 DouYin 后台一致。 |
| 缓存不落地 | Redis 未连接，日志会提示 fallback memory；设置 `CLIENTTOKEN_REDIS_*` 后重启。 |
| API 返回 401 | ClientToken 过期或接口需要 AccessToken 模式，核对文档或重新刷新。 |
| 提示 `token refresh failed` | TTL 低于阈值且 Douyin 刷新失败，检查 `clienttoken_event` 日志与 Douyin 平台响应。 |
| Douyin 模板禁用 | 未同时设置 `DOUYIN_DEVICE_ID` 与 `DOUYIN_RISK_INFO`。若需要这些接口，请补齐两个环境变量后重启服务。 |

> 与 AccessToken 调试不同，ClientToken 不涉及 Flow/授权回调，因此页面不会出现 Flow 表格。若需要用户授权，请参考 `docs/develop/access-token/byteDance/`。
