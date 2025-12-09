# Quickstart: SessionTokenClient Provider Architecture

## Prerequisites
- Go 1.18+
- Redis 实例（默认 `redis://localhost:6379/0`，可通过 `config.MediaXConfig.Redis` 配置）
- `config.yaml` 里为知乎添加 `SessionTokenConfig`，包含 Service/Auth/Harvester/Callback/Network 字段，secret 使用环境变量引用

## Setup Steps
1. **拉取依赖**
   ```bash
   go mod tidy
   ```
2. **配置示例** (`config.yaml`)
   ```yaml
   zhihu_config:
     sessionToken:
       service:
         base_url: ${ZH_ST_BASE_URL}
         api_token: ${ZH_ST_API_TOKEN}
       authenticator:
         entries:
           - type: pc
             url: https://www.zhihu.com/signin
         default_user_agent: Mozilla/5.0 ...
       harvester:
         watch_cookies: [z_c0]
       callback:
         callback_secret: ${ZH_ST_CALLBACK_SECRET}
         max_retry: 3
       network:
         proxy_pool: zhihu-default
   ```
3. **运行单元测试**
   ```bash
   go test ./pkg/client/sessionToken/... ./pkg/client/zhihu/sessionToken/...
   ```
4. **示例：创建 Flow**
   ```bash
   curl -X POST http://localhost:8080/session-token/flows \
     -H "Authorization: Bearer $SESSIONTOKEN_API_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
           "provider_code": "zhihu",
           "provider_app_code": "zhihu_article",
           "account_id": "acct_1",
           "tenant_uuid": "tenant_x",
           "state": "ui-flow-123",
           "callback_url": "https://plugin/api/callback",
           "metadata": {"login_mode": "password"}
         }'
   ```
5. **轮询 Flow**
   ```bash
   curl -H "Authorization: Bearer $SESSIONTOKEN_API_TOKEN" \
     http://localhost:8080/session-token/flows/stf_abc123
   ```
6. **查看回调日志**
   - 日志中应包含 `provider=zhihu api=session_token.callback flow_id=... retry=0`。
   - 若回调失败，可在日志内看到 retry=1/2/3 与 `last_error`。

## Troubleshooting
- **回调签名失败**：确认 callback secret 与插件端配置一致，并检查 payload 中 timestamp/nonce 是否在允许窗口内。
- **Flow 永远 pending**：确认浏览器容器可访问 authorize_url，或检查 Authenticator 脚本日志。
- **Redis 连接失败**：检查 `MediaXConfig.Redis` 是否指向正确实例，并确保 Flow key `sessionToken:flow:*` 能被创建。
