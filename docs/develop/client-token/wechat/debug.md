# WeChat ClientToken Debug 指南

> 目标：在 MediaX 仓库中本地启动 `cmd/clienttoken/server`（默认 `:7072`）、通过 `/debug` 页面与 `scripts/clienttoken-debug.sh` 完成 Token 刷新、API 代理、消息验签，并在第一时间拿到微信返回的原始错误/回调日志，方便排查 `invalid ip`、`invalid grant_type` 等问题。

## 1. 环境与准备

1. Go 1.18+、Redis（推荐 `redis:7-alpine`）。若无本地实例，可执行 `docker run --rm -p 6379:6379 redis:7-alpine`。
2. 工作目录：`/private/var/www/html/ArtisanCloud/X/MediaX/core/MediaX`。
3. 建议在 shell profile 中写入：
   ```bash
   export CLIENTTOKEN_CONFIG="${PWD}/config.yaml"
   export CLIENTTOKEN_LISTEN_ADDR=":7072"
   export CLIENTTOKEN_API_TOKEN="dev-clienttoken"
   export CLIENTTOKEN_REDIS_ADDR="127.0.0.1:6379"
   export CLIENTTOKEN_REDIS_DB="0"
   ```
4. `config.yaml` 可由 `config.example.yaml` 拷贝；若只调试公众号，可保留 `client_token_providers.wechat` 节点并把 `client_token` 下的 `appid/appsecret/message_token` 等字段替换成环境变量。

> **提示**：ClientToken 进程在启动前会主动 ping Redis，若连不上会直接退出，因此务必先确认 `redis-cli ping` 成功。

## 2. 启动 ClientToken Server

```bash
make clienttoken
# 或
GO111MODULE=on go run ./cmd/clienttoken/server \
  -config "${CLIENTTOKEN_CONFIG:-config.yaml}" \
  -listen "${CLIENTTOKEN_LISTEN_ADDR:-:7072}"
```

出现 `clienttoken: server listening addr=:7072 api_token=dev-clienttoken config=config.yaml` 即代表正常。若日志中包含 `failed to cache access token`、`invalid grant_type` 等，则说明服务已经把微信原始响应透传出来，可直接按本文 5.2 小节排查。

## 3. 调试页面 `/debug`

浏览器访问 <http://127.0.0.1:7072/debug> 可获得完整的可视化工具，常见功能：

- **Provider/App/Mode 级联**：读取 `client_token_providers` 模板，选中后会自动填充 `provider_code/provider_app/provider_auth_mode/config_path`，无需重复输入。
- **Token 卡片**：显示最近一次刷新结果（masked token、TTL、来源 `memory|redis|refresh`、Redis Key）。点击“刷新 Token”会调用 `POST /client-token/token`，Redis 及日志会同步更新。
- **API 调试**：填写 `action=cgi-bin/user/get`、`method=GET/POST` 与 Query/Body，调试页会自动拼出请求体并调用 `/client-token/call`，在面板中展示微信的响应或错误码。
- **消息验签**：输入 `signature/timestamp/nonce/echostr`，点击按钮即可检验签名与 EchoStr，适合联调公众号 URL 验证。
- **回调日志**：`/client-token/message/callback` 所有请求会缓存在调试页最下方，支持清空/导出，便于第三方复现。

调试流程示例：
1. 选中 `wechat → official_account → default`，点击“刷新 Token”。若失败，日志会出现微信原始 JSON（例如 `{"errcode":40164,...}`）。
2. 在“API 调试”中填写 `action=cgi-bin/getcallbackip` → “调用 API”→ 观察响应；若 Token 过期，服务会自动刷新，并在结果里标记 `token_source=refresh`。
3. “消息验签”输入公众号后台提供的签名数据，确认 `valid=true`。若 `message_token` 配置错误，此处会立即失败。
4. 如需模拟公众号推送，`POST /client-token/message/callback`（或直接在页面输入 JSON）即可生成回调记录，后端/QA 可以直接查看。

## 4. CLI：`scripts/clienttoken-debug.sh`

脚本封装了常用 API，所有请求默认加 `--noproxy '*'`，防止被系统代理截获。

```bash
# 手动刷新 Token
CLIENTTOKEN_PROVIDER_CODE=wechat \
CLIENTTOKEN_APP_CODE=official_account \
CLIENTTOKEN_AUTH_MODE=default \
scripts/clienttoken-debug.sh token

# 查看缓存命中与 TTL
scripts/clienttoken-debug.sh cache wechat official_account default config.yaml

# 调用公众号接口：GET cgi-bin/getcallbackip
scripts/clienttoken-debug.sh call cgi-bin/getcallbackip GET "" ""

# POST 并读取文件 body
scripts/clienttoken-debug.sh call cgi-bin/message/custom/send POST "" @payload.json

# 验签
scripts/clienttoken-debug.sh validate <signature> $(date +%s) $RANDOM "hello"

# 查询/清空回调日志
scripts/clienttoken-debug.sh callbacks list
scripts/clienttoken-debug.sh callbacks clear
```

常用环境变量：

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `CLIENTTOKEN_BASE_URL` | `http://127.0.0.1:7072` | Server 地址 |
| `CLIENTTOKEN_API_TOKEN` | `dev-clienttoken` | HTTP 鉴权 |
| `CLIENTTOKEN_PROVIDER_CODE` | `wechat` | Provider code |
| `CLIENTTOKEN_APP_CODE` | `official_account` | App code |
| `CLIENTTOKEN_AUTH_MODE` | `default` | 授权模式 |
| `CLIENTTOKEN_CONFIG_PATH` | `config.yaml` | 读取配置文件的路径 |

脚本内部 JSON 采用 `python3` 编码，因此缺少 Python 会直接退出。

## 5. 验证请求链路

### 5.1 手动 `curl`（不依赖脚本）

```bash
curl --noproxy '*' -sS -X POST "http://127.0.0.1:7072/client-token/token" \
  -H "Authorization: Bearer ${CLIENTTOKEN_API_TOKEN}" \
  -H 'Content-Type: application/json' \
  -d '{
    "provider_code":"wechat",
    "provider_app":"official_account",
    "provider_auth_mode":"default",
    "config_path":"config.yaml"
  }' | jq '.'
```

该接口会在失败时直接返回微信错误体，便于 QA/外部应用展示。例如当 IP 不在白名单时，HTTP 200 但 body 为 `{"errcode":40164,"errmsg":"invalid ip ..."}`，终端即可看到完整 JSON。

### 5.2 Redis 与缓存命中

- 默认 Key：`clientToken:wechat:<appid>`。可通过 `cache.redis_key` 覆盖。
- 查看内容：
  ```bash
  redis-cli --raw GET "clientToken:wechat:${WECHAT_OFFICIAL_APP_ID}" | jq '.'
  ```
- 若 `token_source=memory`，说明命中内存缓存；`refresh` 表示刚向微信请求；`redis` 表示从 Redis 读取。
- 当 TTL 低于 `refresh_before_seconds`（默认 600 秒）时，`/client-token/call` 会自动先刷新 Token，再向微信发起真实 API 请求。

## 6. 日志与排查

- 实时查看：
  ```bash
  tail -f logs/clienttoken-info.log | rg 'clienttoken'
  tail -f logs/clienttoken-error.log
  ```
- 常见日志：
  - `clienttoken_metric: action=token ... status=error ... error=access token empty: {"errcode":40002,...}` → 微信直接返回 4xx，需要检查 `grant_type`、AppID/Secret。
  - `clienttoken_metric: action=cache ... status=miss ...` → 缓存未命中，服务会继续刷新。
  - `clienttoken_metric: action=call ... status=success ... token_source=refresh` → API 调用前触发了刷新。
  - `clienttoken_metric: action=message_callback ...` → 表示收到公众号回调，可在 `/debug` 页面查看 payload。
- 若看到 `invalid ip ... not in whitelist`，说明当前出口 IP 未加入公众号安全域。此时服务仍会把错误原文返回给调用方，外部系统无需再去日志截取。

## 7. 与外部应用联调建议

1. **确保 API Token 一致**：外部应用请求 `http://127.0.0.1:7072/client-token/token` 时需携带 `Authorization: Bearer <CLIENTTOKEN_API_TOKEN>`。若返回 401，请先在日志里确认打印的 token。
2. **错误透传**：现在 `tokenHandler` 会校验 `access_token` 是否为空，若微信返回 `errcode/errmsg`，HTTP 200 的 body 会直接包含该 JSON。外部页面可直接展示，不必再组合缓存结构。
3. **多配置切换**：
   - `/debug` 页面与 CLI 都允许动态修改 `config_path`，因此可以在同一进程下测试多个 YAML 文件。
   - 若要对比生产、预发环境，只需 `config_path=conf/staging.yaml` 并在命令里覆盖即可。
4. **排查调用链**：当 `/client-token/call` 返回 403/400 时，结合日志里的 `clienttoken_metric action=token`、`action=call` 两条记录即可判断是刷新失败还是 API 本身拒绝。必要时对照 Redis 中的 `masked_token` 与 `expire_at`，确认 Token 已被刷新。

通过上述步骤，可以在本地快速完成“启动服务 → 刷新 Token → 验签 → 调用 `cgi-bin/*` → 观察日志/回调”的闭环，并把微信原始错误透明地暴露给上游系统。
