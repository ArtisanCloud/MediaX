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
   export SESSIONTOKEN_ZHIHU_API_VERSION="v4"
   ```
   > 插件/Studio 必须与 MediaX 服务共用上述值，避免 BaseURL/API Token 不一致导致 401。
1. **拉取依赖**
   ```bash
   go mod tidy
   ```
2. **配置示例** (`config.yaml`，可直接复制 `config.example.yaml`)
   ```yaml
   session_token_providers:
     providers:
       - code: zhihu
         name: 知乎
         apps:
           - code: web
             name: 知乎 Web
             provider_code: zhihu_sessiontoken
             auth_modes:
               - key: default
                 label: 默认
                 zhihu_session_token_config:
                   strategy: zhihu_pc_v1
                   service:
                     base_url: http://127.0.0.1:7070
                     api_token: dev-session-token
                     timeout: 30
                     http_debug: false
                     api_version: v4
                   authenticator:
                     entries:
                       - type: pc
                         url: https://www.zhihu.com/signin?next=%2F
                     default_user_agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.58 Safari/537.36
                     script_ids:
                       - sessiontoken.zhihu.auth.pc.v1
                     captcha_strategy: auto
                   harvester:
                     watch_cookies:
                       - SESSIONID
                       - JOID
                       - osd
                       - q_c1
                       - d_c0
                       - unlock_ticket
                       - z_c0
                     watch_headers:
                       - X-XSRF-TOKEN
                     harvest_script_id: sessiontoken.zhihu.harvest.pc.v1
                   callback:
                     secret: mediax-sessiontoken-callback
                     max_retry: 3
                     retry_backoff: [2, 4, 8]
                   network:
                     proxy: ""
                     proxy_pool: zhihu-default
                     ip_strategy: china_rotating
                     request_timeout: 60
                   api:
                     retry_backoff: [2, 4, 8]
   ```
   > 建议在首次启动前执行 `cp config.example.yaml config.yaml`（或调用 `make sessiontoken-bootstrap`）自动生成配置，再按环境需要覆盖 `callback_secret`/Redis 等少数字段。
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
   go test ./pkg/client/sessionToken/... ./server/zhihu/sessionToken/...
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
8. **调试台（/debug）**
   - 浏览器打开 `http://127.0.0.1:7070/debug`，选择 Provider/App 模板（默认为 Zhihu Web），页面会自动填充 API Token 与 Callback（默认 `/debug/callback`），并在 `localStorage` 记忆上一次的输入。
   - 点击“创建 Flow”后，输出区域会展示完整响应，同时写入 Flow ID、authorize_url；`copyAuthorizeURL` 可直接复制登录地址；`pollFlow` 按钮用于查看最新 `metadata.session_token`。
   - 勾选“复用 Cookie”后，metadata 将包含 `reuse_session=true`，服务会先尝试使用同租户/账号最近一次成功 Flow 的 SessionToken；若缓存存在则立即 `succeeded`，否则继续提示用户登录。
   - “Mock 回调”区域会实时显示最近 20 条 `/debug/callback` 记录，也可点击“刷新/清空”按钮或通过 `POST /debug/callback` 手动重放 payload。
   - 若希望自动拉起浏览器并回写 Cookie，可在仓库根目录执行 `pnpm install && npx playwright install chromium`，随后运行 `pnpm sessiontoken:browser --flow $FLOW_ID [--base http://127.0.0.1:7070]`。脚本会打开 Chromium、指导你登录知乎，并在完成后自动调用 `POST /debug/flows/<flow_id>/metadata`，同时在终端输出 `/debug/callback` 记录。
   - “API 调试（Beta）”面板会根据 Provider 及所选 API 版本（如 v4）自动加载可测试的 Zhihu API，支持输入 `channel_id=xxx&limit=10` 这类 query/path 参数；点击“从 Flow 填充 SessionToken”即可自 `GET /session-token/flows/<id>` 自动写入 `X-SessionToken`，再点“发送请求”就能用当前凭证直接访问 `/zhihu/v1/*`，结果会在页面底部实时展示，方便一站式验证“抓 Cookie → 调 API → 模拟失效”。
   - 若需要复用已存在的 Flow，可直接在本机 `redis-cli --raw keys 'sessionToken:flow:*'` 列出 key，再用 `redis-cli --raw GET "sessionToken:flow:<id>" | jq '.'` 查看 metadata；把 Flow ID 粘回 `/debug` 即可继续调试。
9. **快速调用 Zhihu API / 查看 Flow**
   - 使用提供的脚本一键操作：
     ```bash
     scripts/sessiontoken-debug.sh flow $FLOW_ID
     scripts/sessiontoken-debug.sh followings "$SESSION_TOKEN"
     scripts/sessiontoken-debug.sh channels "$SESSION_TOKEN" zhihu_column_id 10 0
     scripts/sessiontoken-debug.sh sanity "$SESSION_TOKEN"
     ```
     其中 `SESSION_TOKEN` 取自回调或 `GET /session-token/flows/{id}` 的 `metadata.session_token` 字段；脚本会自动读取 `POWERX_SESSION_TOKEN_BASE_URL` 及 API Token。
10. **查看日志指标**
   - Flow 创建/查询会输出 `sessiontoken_metric: action=create_flow provider=zhihu provider_app=zhihu_article tenant_uuid=tenant_x account_id=acct_demo state=ui-flow flow_id=stf_xxx status=pending latency_ms=12 retry=0`。
   - 回调成功时会看到 `sessiontoken_callback: success provider=zhihu provider_app=zhihu_article tenant_uuid=tenant_x flow_id=stf_xxx state=ui-flow flow_status=succeeded retry=0 http_status=200 latency_ms=5`；若失败则会有 `sessiontoken_callback: failed ... retry=2 latency_ms=900 error=...`。
   - 使用 `tail -f logs/sessiontoken-info.log | rg 'sessiontoken_(metric|callback)'` 可实时观察上述指标，或在生产中接入日志/指标系统。

## Troubleshooting
- **回调签名失败**：确认 callback secret 与插件端配置一致，并检查 payload 中 timestamp/nonce 是否在允许窗口内。
- **Flow 永远 pending**：确认浏览器容器可访问 authorize_url，或检查 Authenticator 脚本日志。
- **本地代理拦截**：如果本机装有 Surge/Charles，记得在 curl 中添加 `--noproxy "*"`.
- **Redis 连接失败**：检查 `MediaXConfig.Redis` 是否指向正确实例，并确保 Flow key `sessionToken:flow:*` 能被创建。
