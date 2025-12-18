# WeChat ClientToken 调试指南

> 目标：在 MediaX 仓库内独立启动 `cmd/clienttoken/server`（默认 `:7072`），通过调试页与 CLI 刷新公众号 `access_token`、校验消息签名、代理调用 `cgi-bin/*` 接口，并配合 Redis 缓存快速排查 Token 生命周期问题。

## 1. 依赖与基础环境

1. Go 1.18+、Redis（默认 `127.0.0.1:6379`）。若没有本地 Redis，可先运行 `docker run --rm -p 6379:6379 redis:7-alpine`。
2. 工作目录：`/private/var/www/html/ArtisanCloud/X/MediaX/core/MediaX`。
3. 核心环境变量（可写入 shell profile）：
   ```bash
   export CLIENTTOKEN_CONFIG="${PWD}/config.yaml"
   export CLIENTTOKEN_LISTEN_ADDR=":7072"
   export CLIENTTOKEN_API_TOKEN="dev-clienttoken"
   export CLIENTTOKEN_REDIS_ADDR="127.0.0.1:6379"
   export CLIENTTOKEN_REDIS_DB="0"
   ```
4. 没有 `config.yaml` 时，可直接复制 `config.example.yaml` 并替换为占位符，也可新建最小配置（见下）。

### 1.1 `client_token_providers` 最小配置

```yaml
client_token_providers:
  providers:
    - code: "wechat"
      name: "微信"
      apps:
        - code: "official_account"
          name: "公众号"
          provider_code: "wechat_official_account"
          auth_modes:
            - key: "default"
              label: "默认"
              wechat_official_account_config:
                api_url: "https://api.weixin.qq.com"
                proxy_api_url: ""
                timeout: 5
                http_debug: false
                oauth:
                  access_token_url: "https://api.weixin.qq.com/cgi-bin/token"
                  client_id: "${WECHAT_OFFICIAL_APP_ID}"
                  client_secret: "${WECHAT_OFFICIAL_APP_SECRET}"
                client_token:
                  appid: "${WECHAT_OFFICIAL_APP_ID}"
                  appsecret: "${WECHAT_OFFICIAL_APP_SECRET}"
                  message_token: "${WECHAT_OFFICIAL_MESSAGE_TOKEN}"
                  message_aes_key: "${WECHAT_OFFICIAL_MESSAGE_AES_KEY}"
                api_token: "${WECHAT_CLIENTTOKEN_API_TOKEN:-dev-clienttoken}"
                cache:
                  redis_key: "clientToken:wechat:${WECHAT_OFFICIAL_APP_ID}"
                  ttl_seconds: 7000
                  refresh_before_seconds: 600
                redis:
                  addr: "${WECHAT_CLIENTTOKEN_REDIS_ADDR:-127.0.0.1:6379}"
                  db: "${WECHAT_CLIENTTOKEN_REDIS_DB:-0}"
```

> 建议所有敏感字段使用环境变量覆盖，避免把真实 `AppSecret` 写入仓库。

## 2. 启动 ClientToken Server

1. 确保 Redis 正常。可执行 `redis-cli ping` 检查。
2. 在仓库根目录运行：
   ```bash
   go run ./cmd/clienttoken/server -config "${CLIENTTOKEN_CONFIG:-config.yaml}"
   # 或使用 Makefile
   make clienttoken
   ```
3. 日志出现 `clienttoken: server listening addr=:7072 api_token=dev-clienttoken` 代表启动成功。若 Redis 不可达，进程会立即退出，需要先排查网络/凭证。

## 3. 内置调试页面（/debug）

1. 打开浏览器访问 <http://127.0.0.1:7072/debug>。
2. 页面提供与 SessionToken 调试台一致的体验：
   - **Provider/App/Mode**：直接读取 `client_token_providers`，可切换不同公众号或授权模式，表单会记忆上次输入。
   - **Token 卡片**：展示 `access_token`（脱敏）、来源（缓存/刷新）、TTL、Redis Key，可一键触发“刷新缓存”。
   - **API 调试**：输入 `cgi-bin/user/get` 等 action，支持 GET/POST，Query/Body 均为纯文本。点击“调用 API”后会自动注入最新 Token 并返回格式化 JSON。
   - **消息验签**：填写 `signature/timestamp/nonce/echostr`，服务端会根据配置中的 `message_token` 计算期望签名并给出结果。
   - **回调日志**：`POST /client-token/message/callback` 的所有请求都会记录在页面最下方，方便联调第三方回调逻辑，可随时清空。
3. 常见操作流程：
   1. 进入页面后选择 Provider 模板 → 点击“刷新 Token” → 观察日志/Redis 中是否写入 `clientToken:wechat:<appid>`。
   2. 在 API 调试区填写 `action=cgi-bin/getcallbackip`，点击“调用 API”验证正式接口。
   3. 在“消息验证”区域输入真实签名或指定 timestamp/nonce，确认 `valid=true`。
   4. 如需模拟公众号推送，直接 `POST /client-token/message/callback` 或在调试页粘贴 JSON，该记录会显示在列表中供前端或测试人员查看。

## 4. Redis 缓存与 Token 生命周期

- 缓存 Key 默认为 `clientToken:wechat:<appid>`，可通过 `cache.redis_key` 覆盖。
- 记录结构：
  ```json
  {
    "appid": "wx_demo",
    "access_token": "ya29...",
    "stored_at": "2025-12-18T11:40:00+08:00",
    "expire_at": "2025-12-18T13:40:00+08:00",
    "source": "refresh|memory"
  }
  ```
- 逻辑：
  1. `POST /client-token/token` 会无条件刷新，成功后写入 Redis 与内存 map。
  2. `GET /client-token/cache` 优先读内存，再读 Redis，返回 `token_source` 与剩余 TTL。
  3. `POST /client-token/call` 若发现 Token 即将过期（`ttl_seconds < refresh_before_seconds`），会提前调用刷新逻辑，确保 API 调用永远拿到有效 Token。
- 排查建议：
  - `redis-cli --raw GET "clientToken:wechat:xxx" | jq '.'` 查看真实缓存。
  - 注意 TTL 是否与配置一致（默认 7000 秒，提前 600 秒刷新）；可在调试页随时查看。

## 5. CLI：`scripts/clienttoken-debug.sh`

为减少重复输入，仓库附带脚本支持以下命令：

```bash
# 刷新 Token
CLIENTTOKEN_PROVIDER_CODE=wechat CLIENTTOKEN_APP_CODE=official_account \
scripts/clienttoken-debug.sh token

# 查看缓存
scripts/clienttoken-debug.sh cache

# 调用 API（GET cgi-bin/user/get）
scripts/clienttoken-debug.sh call cgi-bin/user/get GET "next_openid=" ""

# 以文件作为 POST body
scripts/clienttoken-debug.sh call cgi-bin/message/custom/send POST "" @payload.json

# 验签
scripts/clienttoken-debug.sh validate "<signature>" "$(date +%s)" "$RANDOM" "hello"

# 查看/清空回调日志
scripts/clienttoken-debug.sh callbacks list
scripts/clienttoken-debug.sh callbacks clear
```

- 脚本默认读取以下环境变量：
  | 变量 | 默认 | 说明 |
  | --- | --- | --- |
  | `CLIENTTOKEN_BASE_URL` | `http://127.0.0.1:7072` | Server 地址 |
  | `CLIENTTOKEN_API_TOKEN` | `dev-clienttoken` | HTTP 鉴权 |
  | `CLIENTTOKEN_PROVIDER_CODE` | `wechat` | Provider code |
  | `CLIENTTOKEN_APP_CODE` | `official_account` | App code |
  | `CLIENTTOKEN_AUTH_MODE` | `default` | 授权模式 |
  | `CLIENTTOKEN_CONFIG_PATH` | `config.yaml` | 当请求中需要切换配置文件时使用 |
- 所有 `curl` 请求强制加 `--noproxy '*'`，避免被系统代理篡改。

## 6. 常见问题与排查

| 场景 | 现象 | 排查步骤 |
| --- | --- | --- |
| Redis 未启动 | 进程启动即退出，日志 `ping redis` 失败 | 确认 `CLIENTTOKEN_REDIS_ADDR`、本地 Redis 服务；临时可 `docker run redis` |
| API Token 不匹配 | `/debug` 或 CLI 返回 `401 unauthorized` | 检查服务端日志打印的 `api_token`，与请求头 `Authorization: Bearer ...` 是否一致 |
| 调用公众号接口报 `40001/invalid credential` | Token 失效或 AppSecret 错误 | 重新刷新 Token，检查配置中的 AppID/AppSecret 是否真实可用 |
| 调试页无法列出 Provider | `client_token_providers` 为空或 YAML 路径错误 | 确保 `-config` 指向含新节点的文件，可在页面上改 `config_path` 字段临时切换 |
| `scripts/clienttoken-debug.sh` 报 `python3 is required` | Mac/Linux 缺少 Python3 | 安装 `python3` 或修改脚本逻辑 |
| 消息验签一直失败 | `message_token` 未配置/不匹配 | 检查配置中 `client_token.message_token`，确保与公众号后台一致 |

## 7. 与 SessionToken 调试页的差异

- ClientToken 不涉及 Flow/回调链路，所有操作都是幂等 HTTP 接口；Redis 是强依赖，未连接成功直接退出。
- `/debug` 页面以 Provider/App 列表为核心，默认端口 `7072`（SessionToken 为 `7070`，AccessToken 为 `7071`），可以同时运行互不影响。
- CLI 主要围绕 Token 刷新、缓存检查、API 调试与消息验签，相比 SessionToken CLI，不需要管理 Flow ID，也无需写回 metadata。

借助本文档步骤，可以在 5 分钟内完成以下闭环：启动服务 → 刷新 Token → 在调试页或脚本调用公众号接口 → 校验消息签名 → 将日志/缓存交给 QA 或第三方系统复现。若需要进一步扩展（例如 CLI 一键刷新全部 App、或将 `/debug` 接口接入 CI），可在 `cmd/clienttoken/server` 与 `scripts/clienttoken-debug.sh` 基础上继续演进。
