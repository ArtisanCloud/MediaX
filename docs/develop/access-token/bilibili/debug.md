# BiliBili AccessToken 调试指南

> 目标：通过 `go run ./cmd/accesstoken/server` 的调试页，快速完成 B 站 OAuth 授权、查看 AccessToken/Flow 记录、调用 `/accesstoken/token` 验证配置。一些 API 面板目前仅支持 Google YouTube，因此本文侧重 Token 获取与排障。

## 1. 启动调试服务

```bash
cd /private/var/www/html/ArtisanCloud/X/MediaX/core/MediaX
export ACCESSTOKEN_CONFIG=${ACCESSTOKEN_CONFIG:-config.yaml}
export ACCESSTOKEN_LISTEN_ADDR=${ACCESSTOKEN_LISTEN_ADDR:-127.0.0.1:7071}
go run ./cmd/accesstoken/server \
  -config "$ACCESSTOKEN_CONFIG" \
  -port "${ACCESSTOKEN_LISTEN_ADDR##*:}"
```

日志若输出 `accesstoken-server: listening addr=127.0.0.1:7071 api_token=dev-accesstoken config=config.yaml` 即代表成功；若需要开放到其他主机，只需覆盖 `ACCESSTOKEN_LISTEN_ADDR`（如 `0.0.0.0:8080`）。也可以使用 `make accesstoken-serve`（支持 `ARGS='-port 8080'`，端口仍会根据 `ACCESSTOKEN_LISTEN_ADDR` 自动选择）。

必备环境变量：
- `ACCESSTOKEN_API_TOKEN`（默认 `dev-accesstoken`）保护调试接口；浏览器访问 `/debug` 时需在 URL 末尾追加 `?api_token=<值>` 或点击右上角“更新 API Token”按钮写入浏览器存储，后续请求会自动附带 `Authorization: Bearer ...`。
- `ACCESSTOKEN_LISTEN_ADDR`（默认 `127.0.0.1:7071`）决定 HTTP 监听地址；默认仅绑定本机防止调试界面暴露到公网。
- `ACCESSTOKEN_REDIS_ADDR/DB/PASS`：如果希望授权记录可跨进程复用，需要提前在本地启动 Redis（或指向测试 Redis）。未设置时记录只保存在进程内存，重启即失效。
- `GOOGLE_*`/`BILIBILI_*` 等变量：`config.yaml` 若使用 `${BILIBILI_CLIENT_ID}` 占位符，需要提前注入。
- Flow TTL 与记录保留：无论使用 Redis 或内存，Flow 记录默认保留 24 小时（86,400 秒，可通过 `ACCESSTOKEN_FLOW_TTL_SECONDS` 覆盖），超时后需要重新授权。

## 2. 配置要求

`config.yaml` 必须包含 `access_token_providers` 并声明 BiliBili Provider。例如：

```yaml
access_token_providers:
  - code: bilibili
    name: BiliBili
    apps:
      - code: content_center
        provider_code: bilibili_arcopen
        auth_modes:
          - key: default
            bilbili_config:
              api_url: "https://member.bilibili.com"
              timeout: 10
              http_debug: false
              oauth:
                client_id: "${BILIBILI_CLIENT_ID}"
                client_secret: "${BILIBILI_CLIENT_SECRET}"
                redirect_url: "http://localhost:7071/debug/callback"
                scope: "ARC_BASE,ARTICLE_BASE,LIVE_ROOM_DATA"
```

注意：
- `redirect_url` 必须指向本地调试服务的回调 `/debug/callback`，否则授权成功后无法写入记录。
- `scope` 需包含业务所需权限；至少要有 `ARC_BASE` 才能调用视频稿件相关接口。
- 同一 Provider 下可以配置多个 App 或授权模式（`auth_modes`），调试页会自动读取并生成对应下拉列表。

## 3. 浏览器调试流程

访问 <http://127.0.0.1:7071/debug?api_token=dev-accesstoken>（或任意自定义 Token）：

> 提示：首次访问推荐在 URL 中附带 `?api_token=`，页面加载后会自动将 Token 写入浏览器 LocalStorage 并清理地址栏，也可以通过页面右上角的“更新 API Token”按钮随时更换。

1. **选择 Provider/App/Mode**：左上角三个下拉分别对应 `code`/`app code`/`auth_mode key`。选择后可点击“同步模板”让 JSON 区块快速填充。
2. **解析 AccessToken**：`POST /accesstoken/token`，输入 JSON 包含 `provider_code/provider_app/provider_auth_mode/config_path`，可留空 `access_token` 以使用最新授权记录。返回结果会显示：
   - `masked_token`、`token_source`（env/config/flow）
   - `access_token_ttl`
   - `oauth_key`、`config_path`、`provider` 信息
   适合排查“为什么外部调用拿不到 token”这类问题。
3. **发起授权**：点击“发起授权”按钮 → 浏览器会跳到 B 站登录页。登录完成后 B 站会将 `code/state` 回调到 `/debug/callback`，调试页底部的“授权回调记录”表格会新增一条记录（含 `flow_id`、`expire_at`、账号信息）。
4. **回填 Flow**：如果调试页刷新或服务重启，可在“Flow ID 回填”处输入 `oauth-xxxx` 并点击“加载 Flow ID”，服务会从内存/Redis (`accesstoken:oauth:flow:<id>`) 中读取记录并填入 JSON，继续解析即可。
5. **授权记录列表**：点击“刷新授权记录”会列出当前 Provider/App/Mode 关联的所有历史记录，来源可包括浏览器授权、CLI 触发、或自定义回调。支持清空、复制 JSON、查看 `provider_code/provider_app/oauth_key` 等字段。

> 当前版本 `/accesstoken/call` 的 API 执行只开放 Google YouTube。若要调用 B 站接口，请使用 Go SDK 或后续定制脚本。调试台仍然提供 Token 解析/Flow 管理能力，方便与外部系统协作。

> 若顶部出现“Flow 数据当前存储在内存中”的黄色横幅，表示 Redis 未启用，Flow/Token 记录将在服务重启后清空；如需持久化请配置 `ACCESSTOKEN_REDIS_*`。当监听地址为 `0.0.0.0` 等公网可访问地址时，还会出现红色安全提醒。

## 4. Redis 键空间

授权成功后会写入两个 key（主键与 flow 索引都会设置 24 小时 / 86,400 秒 TTL）：
- `accesstoken:oauth:bilibili:<app>:<mode>`：保存最新授权结果（JSON 含 `access_token/refresh_token/expires_in/flow_id`），调试页的“刷新授权记录”即读取此 key。
- `accesstoken:oauth:flow:<flow_id>`：索引 key，方便用户复制 Flow ID 后在其他环境复用。内容为字符串，指向主 key。

示例：
```bash
redis-cli --raw GET accesstoken:oauth:flow:oauth-123456789
# → accesstoken:oauth:bilibili:content_center:default
redis-cli --raw GET accesstoken:oauth:bilibili:content_center:default | jq '.'
```
若希望默认不写 Redis，只需不设置 `ACCESSTOKEN_REDIS_ADDR`；调试页仍可在单次进程内使用授权记录。

## 5. cURL/脚本示例

### 5.1 手动解析 Token
```bash
curl --noproxy '*' -sS -X POST 'http://127.0.0.1:7071/accesstoken/token' \
  -H 'Authorization: Bearer dev-accesstoken' \
  -H 'Content-Type: application/json' \
  -d '{
    "provider_code": "bilibili",
    "provider_app": "content_center",
    "provider_auth_mode": "default",
    "config_path": "config.yaml",
    "access_token": "${BILIBILI_ACCESS_TOKEN:-}"
  }' | jq '.'
```
- 若 body 中 `access_token` 为空，服务会尝试读取环境变量（`BILIBILI_ACCESS_TOKEN`、`ACCESSTOKEN_ACCESS_TOKEN` 等）或最新 Flow 记录。
- 若返回 `missing access token`，说明没有任何来源可用，需要重新授权或显式传入 token。

### 5.2 开启授权流程（API）
```bash
curl -sS -X POST 'http://127.0.0.1:7071/accesstoken/oauth/start' \
  -H 'Authorization: Bearer dev-accesstoken' \
  -H 'Content-Type: application/json' \
  -d '{
    "provider_code": "bilibili",
    "provider_app": "content_center",
    "provider_auth_mode": "default",
    "config_path": "config.yaml"
  }' | jq '.authorize_url'
```
复制返回的 URL 在浏览器打开即可完成授权。CLI/UI 在内部也是调用该接口。

## 6. 常见问题排查
| 现象 | 处理建议 |
| --- | --- |
| 浏览器 401 Unauthorized | 确认访问 `/debug` 时附带 `?api_token=...` 或已通过右上角按钮写入 Token。Token 值需与 `ACCESSTOKEN_API_TOKEN` 一致。|
| 授权页回跳 404 | `redirect_url` 未指向 `http://localhost:7071/debug/callback`。修改配置后需重启服务并在 B 站后台同步回调域。|
| `missing access token` | 未执行授权或 Redis 缺记录。确认 `accesstoken:oauth:...` 是否存在；若没有需重新发起授权。|
| Flow ID 查询为 `[]` | Flow 已过期或 key 被清理。默认 TTL 为 24 小时（86,400 秒），若需要更长保留请在 24 小时内重新授权或自建延长策略。|
| `/accesstoken/flow/replay` 返回 `FLOW_NOT_FOUND` | Flow ID 填写错误或尚未触发授权。可在“授权记录”表或 `redis-cli KEYS accesstoken:oauth:flow:*` 中确认真实 ID。|
| `/accesstoken/flow/replay` 返回 `FLOW_EXPIRED` | Flow 超过 TTL 已被清理。重新点击“发起授权”或让业务方再授权一遍。|
| 仍想调用 B 站 API | 当前 `/accesstoken/call` 会拒绝非 Google Provider。请直接在 Go 代码中调用 `BiliBiliACClient`，或自行编写 CLI。|

## 7. 观测与日志
- **Info/Warning 日志**：`logs/accesstoken-server-info.log`、`logs/accesstoken-server-error.log` 分别收集运行态信息。每条日志都包含 `event` 字段（如 `oauth.start`、`oauth.complete`、`token.resolve`、`flow.replay`）以及以下关键字段：`provider_code`、`provider_app`、`provider_auth_mode`、`flow_id`、`token_source`、`storage_backend`、`listen_addr`。可以使用 `rg 'event=flow' logs/accesstoken-server-info.log` 追踪整条授权链路。
- **Flow 追踪**：UI/CLI 中显示的 `status` 与 API 返回值一致：`pending`（等待回调）、`authorized`（已拿到 token）、`expired`（TTL 超时）。`flow_state` 字段会保留原始 OAuth `state`，方便排查多个授权混淆的问题。
- **回调与 masking**：`/debug/callback` 存储的 `callback.body/query` 已自动替换 `code/state/access_token/refresh_token/client_secret` 字段为 `***`，导出时不会泄露原始值。
- **诊断技巧**：
  - 监听 `accesstoken-server: event=oauth.start.error`，常见 detail 包括 `invalid_json`、`resolve_provider`、`build_authorize_url`。
  - 当 Redis 不可用或回退到内存模式时，会打印 `storage_backend=memory` 且 UI 显示黄色横幅，可提醒操作人重启后需重新授权。
  - `flow.replay` / `flow.list` 请求同样会写日志，可结合 Flow ID 与 `listen_addr`、`token_source_detail` 检查是谁发起的回填。

## 8. 与外部应用协同
- 若外部服务只需要 AccessToken，可让其调用 `POST /accesstoken/token`，由本地调试服务统一返回 `masked_token/token_source/ttl` 等字段。
- Flow ID 可以共享给外部系统，帮助他们直接在 Redis 中获取授权记录，避免多次登录。确保 Flow ID 在日志中脱敏（服务端默认打 `oauth-****`）。
- 建议使用 `logs/accesstoken-server-info.log` + `rg 'accesstoken'` 观察 OAuth/Token 行为，配合 `tail -f` 可以实时查看浏览器操作结果。

通过该调试流程，可以在几分钟内完成“配置 Provider → 启动服务 → 发起 OAuth → 解析 Token → 提供给业务代码”这一闭环，确保 B 站 AccessToken 能力在接入前充分验证。
