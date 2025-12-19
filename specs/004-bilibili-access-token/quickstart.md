# Quickstart: BiliBili AccessToken Debug Flows

## 1. Prerequisites
- Go 1.18+
- MediaX dependencies installed via `go mod download`
- `config.yaml` already contains `access_token_providers` entry for `bilibili` with `${BILIBILI_CLIENT_ID}`-style placeholders
- Required environment variables:
  ```bash
  export ACCESSTOKEN_CONFIG=${ACCESSTOKEN_CONFIG:-config.yaml}
  export ACCESSTOKEN_API_TOKEN=${ACCESSTOKEN_API_TOKEN:-dev-accesstoken}
  export ACCESSTOKEN_LISTEN_ADDR=${ACCESSTOKEN_LISTEN_ADDR:-127.0.0.1:7071}
  # Optional Redis cache
  export ACCESSTOKEN_REDIS_ADDR=127.0.0.1:6379
  export ACCESSTOKEN_REDIS_DB=0
  export ACCESSTOKEN_REDIS_PASS=""
  ```

## 2. Start the debug service
```bash
cd /private/var/www/html/ArtisanCloud/X/MediaX/core/MediaX
make accesstoken-serve ARGS="-config $ACCESSTOKEN_CONFIG -port ${ACCESSTOKEN_LISTEN_ADDR##*:}"
# or
ACCESSTOKEN_LISTEN_ADDR=$ACCESSTOKEN_LISTEN_ADDR go run ./cmd/accesstoken/server \
  -config "$ACCESSTOKEN_CONFIG" \
  -port "${ACCESSTOKEN_LISTEN_ADDR##*:}"
```
Expected log: `accesstoken-server: listening addr=127.0.0.1:7071 api_token=*** token_cache=redis|memory`.

## 3. Use the browser console
1. 打开 `http://127.0.0.1:7071/debug`（可在首次 URL 中附带 `?api_token=...` 自动写入 Token）
2. 右上角齿轮输入 `ACCESSTOKEN_API_TOKEN`
3. 选择 Provider/App/Auth Mode 后点击 **同步模板** 让 JSON 表单填充默认值
4. 点击 **发起授权** 完成 OAuth；授权成功后面板最下方的“授权记录”会出现新的 `flow_id` 行
5. 点击 **刷新授权记录** 或 **加载 Flow ID** 可按当前 Provider/App 回填历史 Flow；多次点击“加载 Flow 列表”会根据 `cursor` 翻页 Redis 中的记录
6. 需要手动回填时，可在“Flow ID 回填”输入框中粘贴 `oauth-xxx`，调试台会调用 `/accesstoken/flow/replay` 自动填充 AccessToken
7. 最后发送 `/accesstoken/token` 请求块检查 Token 来源、TTL、`storage_backend`（`redis` 表示可跨进程复用，`memory` 说明重启即丢失）

## 4. API smoke tests

### 4.1 解析最新 AccessToken
```bash
curl --noproxy '*' -sS -X POST 'http://127.0.0.1:7071/accesstoken/token' \
  -H "Authorization: Bearer $ACCESSTOKEN_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "provider_code": "bilibili",
    "provider_app": "content_center",
    "provider_auth_mode": "default",
    "config_path": "config.yaml"
  }' | jq '.'
```
- `token_source=flow` 表示已经回填成功；`storage_backend=redis` 则代表可跨进程复用

### 4.2 列出 Flow 记录
```bash
curl -sS 'http://127.0.0.1:7071/accesstoken/flows?provider_code=bilibili&provider_app=content_center&provider_auth_mode=default' \
  -H "Authorization: Bearer $ACCESSTOKEN_API_TOKEN" | jq '.flows[] | {flow_id,status,expire_at,masked_account}'
```
返回的 `status` 会是 `pending`/`authorized`/`expired`，并附带 `masked_account` 帮助定位具体 provider。

### 4.3 Flow 回放
```bash
curl -sS -X POST 'http://127.0.0.1:7071/accesstoken/flow/replay' \
  -H "Authorization: Bearer $ACCESSTOKEN_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"flow_id": "oauth-123"}' | jq '.'
```
`status=authorized` 表示可以直接调用 `/accesstoken/token`，若返回 `FLOW_NOT_FOUND/FLOW_EXPIRED` 则需要重新授权。

## 5. Flow TTL & replay
- Every Flow + primary record TTL = 24 hours (86,400 seconds)
- Redis keys:
  - `accesstoken:oauth:bilibili:<app>:<mode>` for latest record
  - `accesstoken:oauth:flow:<flow_id>` pointer
- Replay via API:
  ```bash
  curl -sS -X POST 'http://127.0.0.1:7071/accesstoken/flow/replay' \
    -H "Authorization: Bearer $ACCESSTOKEN_API_TOKEN" \
    -H 'Content-Type: application/json' \
    -d '{"flow_id": "oauth-123"}' | jq '.'
  ```
  Response `status=authorized` means `/accesstoken/token` can now resolve without re-login.

## 6. Memory-mode fallback
- Missing Redis env or connection failure logs `flow_store=memory` and shows a yellow UI banner announcing volatility
- `storage_backend` field returns `memory`
- Support docs instruct operators to re-run OAuth after restart

## 7. Logging & masking checklist
- Log fields: `provider`, `provider_app`, `provider_auth_mode`, `flow_id`, `token_source`, `storage_backend`, `listen_addr`
- Masked token keeps first/last four chars only
- Callback payload strips PII before persistence/export

Follow these steps to satisfy the Success Criteria before shipping downstream automation or `/speckit.tasks` planning.
