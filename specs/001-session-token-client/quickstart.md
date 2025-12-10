# Quickstart: SessionTokenClient Provider Architecture

## Prerequisites
- Go 1.18+
- Redis 实例（默认 `redis://localhost:6379/0`，可通过 `config.MediaXConfig.Redis` 配置）
- `config.yaml` 里为知乎添加 `SessionTokenConfig`，包含 Service/Auth/Harvester/Callback/Network 字段，secret 使用环境变量引用

## Setup Steps
0. **准备环境变量（服务与插件共用）**
   ```bash
   export POWERX_SESSION_TOKEN_BASE_URL="http://127.0.0.1:7070"
   export POWERX_SESSION_TOKEN_API_TOKEN="dev-session-token"
   export POWERX_SESSION_TOKEN_CALLBACK_URL="https://plugin.local/api/v1/admin/platforms/session-token/callback"
   export SESSIONTOKEN_REDIS_ADDR="127.0.0.1:6379"
   ```
   > 插件/Studio 必须与 MediaX 服务共用上述值，避免 BaseURL/API Token 不一致导致 401。
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
3. **启动 SessionToken 服务**
   ```bash
   make sessiontoken
   # 或者
   go run ./cmd/sessiontoken -config config.yaml
   ```
   先在 MediaX 仓库启动服务，待日志输出 `sessiontoken: server listening` 后，再启动插件或浏览器容器。

4. **（可选）挂载 HTTP Handler**
   ```go
   mux := http.NewServeMux()
   store := redisstore.NewFlowStore(redisClient, mediaX.Logger)
   mgr := sessiontoken.NewManager(baseClient, mediaX.Logger, cache, store,
     sessiontoken.WithAuthenticator(auth),
     sessiontoken.WithCallbackDispatcher(dispatcher),
   )
   token := os.Getenv("POWERX_SESSION_TOKEN_API_TOKEN")
   if token == "" {
     token = os.Getenv("SESSIONTOKEN_API_TOKEN")
   }
   sessiontokenhandler.RegisterSessionTokenFlowCreateRoute(mux, mgr, token, mediaX.Logger)
   sessiontokenhandler.RegisterSessionTokenFlowGetRoute(mux, mgr, token, mediaX.Logger)
   ```
5. **运行单元测试**
   ```bash
   go test ./pkg/client/sessionToken/... ./pkg/client/zhihu/sessionToken/...
   ```
6. **示例：创建 Flow**
   ```bash
   FLOW_RESPONSE=$(curl --noproxy "*" -s -X POST http://localhost:8080/session-token/flows \
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
         }')
   echo "$FLOW_RESPONSE"
   export FLOW_ID=$(echo "$FLOW_RESPONSE" | jq -r '.flow.flow_id')
   ```
7. **轮询 Flow + 查看凭证**
   ```bash
   curl --noproxy "*" -H "Authorization: Bearer $SESSIONTOKEN_API_TOKEN" \
     http://localhost:8080/session-token/flows/$FLOW_ID | jq '.flow.status,.flow.result'
   ```
8. **查看日志指标**
   - Flow 创建/查询会输出 `sessiontoken_metric: action=create_flow provider=zhihu ... latency_ms=12`，可据此评估延迟。
   - 回调成功时会看到 `sessiontoken_callback: success provider=zhihu tenant_uuid=tenant_x flow_id=... retry=0 latency_ms=5`；若失败则会有 `sessiontoken_callback: failed ... retry=2 latency_ms=900 error=...`，便于排查。

## Troubleshooting
- **回调签名失败**：确认 callback secret 与插件端配置一致，并检查 payload 中 timestamp/nonce 是否在允许窗口内。
- **Flow 永远 pending**：确认浏览器容器可访问 authorize_url，或检查 Authenticator 脚本日志。
- **本地代理拦截**：如果本机装有 Surge/Charles，记得在 curl 中添加 `--noproxy "*"`.
- **Redis 连接失败**：检查 `MediaXConfig.Redis` 是否指向正确实例，并确保 Flow key `sessionToken:flow:*` 能被创建。
