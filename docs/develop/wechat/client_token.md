# 微信 ClientToken Server Quickstart

本指南面向需要在本地或测试环境调试微信公众号 Client Credential Flow 的同学，帮助你在 5 分钟内跑通 `cmd/clienttoken/server`、刷新 Token、调用接口以及模拟消息回调。

## 1. 前置条件

- Go 1.18+（确保 `go version` 输出 ≥ 1.18）
- 已准备 `config.yaml`，其中包含 `client_token_providers` 节点（可复制 `config.example.yaml` 并替换占位符）
- 可访问的 Redis（本地默认 `127.0.0.1:6379`，或通过 `CLIENTTOKEN_REDIS_ADDR` 覆盖）

常用环境变量：

| 变量 | 说明 | 默认 |
| --- | --- | --- |
| `CLIENTTOKEN_CONFIG` | 配置文件路径 | `config.yaml` |
| `CLIENTTOKEN_LISTEN_ADDR` | HTTP 监听地址 | `:7072` |
| `CLIENTTOKEN_API_TOKEN` | API 鉴权 Token | `dev-clienttoken` |
| `CLIENTTOKEN_REDIS_ADDR` / `CLIENTTOKEN_REDIS_DB` | Token 缓存 Redis | `127.0.0.1:6379` / `0` |
| `CLIENTTOKEN_PROVIDER_CODE` / `CLIENTTOKEN_APP_CODE` / `CLIENTTOKEN_AUTH_MODE` | CLI/脚本默认 Provider App | `wechat` / `official_account` / `default` |

## 2. 配置示例

`client_token_providers` 核心字段：

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
                timeout: 5.0
                http_debug: true
                oauth:
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

> `client_token` 节点用于 Token 换取；`api_token` 控制 `/client-token/*` 与 `/debug` 鉴权；`cache/redis` 可按环境覆盖。

## 3. 启动 ClientToken Server

```bash
go run ./cmd/clienttoken/server -config config.yaml -port 7072
```

- 若需要自定义端口或配置路径，可直接设置环境变量：
  ```bash
  export CLIENTTOKEN_CONFIG=/path/to/config.yaml
  export CLIENTTOKEN_LISTEN_ADDR=:17072
  export CLIENTTOKEN_API_TOKEN=my-secret
  go run ./cmd/clienttoken/server
  ```
- 启动成功后日志会输出 `clienttoken: server listening addr=:7072 api_token=**`，并在 Redis 中预留 `clientToken:wechat:<appid>` 键。

## 4. 使用 `/debug` 页面

1. 打开浏览器访问 <http://127.0.0.1:7072/debug>。
2. 在“Provider 选择”区域选中需要调试的公众号配置，可覆盖 `config_path` 与 API Token。
3. “Token 缓存”卡片可查看当前缓存状态、刷新 Token；命中后会显示 TTL、存储时间以及 Redis Key。
4. “API 调试”中填写微信公众号 API 路径、Method、Query/Body，即可点击“调用 API”。默认使用 `cgi-bin/getcallbackip` 以验证网络与 Token。
5. “消息验证 / 回调”支持：
   - 根据配置 `message_token` 计算签名，模拟微信验签；
   - `POST /client-token/message/callback` 可写入任意 payload，日志会出现在下方列表（配合外部系统测试时非常方便）。

## 5. 常用 REST 接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/client-token/token` | 强制刷新 Token 并写入 Redis，返回脱敏值/TTL/来源 |
| `GET` | `/client-token/cache` | 仅查看缓存命中情况，不触发刷新 |
| `POST` | `/client-token/call` | 自定义 action/method/query/body，代理调用微信公众号 API |
| `POST` | `/client-token/message/validate` | 根据 `message_token` 验签，返回 `expected` 与 `valid` |
| `POST` | `/client-token/message/callback` | 记录一条调试回调日志（配合 `/debug` 展示） |
| `GET/DELETE` | `/client-token/message/callbacks` | 获取或清空调试回调日志 |
| `GET` | `/healthz` | 健康检查 |

示例（刷新 Token）
```bash
curl -sS -X POST \
  -H "Authorization: Bearer dev-clienttoken" \
  -H "Content-Type: application/json" \
  -d '{
        "provider_code": "wechat",
        "provider_app": "official_account",
        "provider_auth_mode": "default",
        "config_path": "config.yaml"
      }' \
  http://127.0.0.1:7072/client-token/token | jq '.'
```

## 6. CLI / 脚本辅助

仓库提供 `scripts/clienttoken-debug.sh`，可在终端快速刷新 Token、查看缓存、调用接口或验证回调：

```bash
# 刷新 Token（默认读取 config.yaml 中的第一个 Provider）
CLIENTTOKEN_API_TOKEN=dev-clienttoken \
CLIENTTOKEN_PROVIDER_CODE=wechat \
scripts/clienttoken-debug.sh token

# 查看 Redis 缓存命中
scripts/clienttoken-debug.sh cache

# 调用 API（示例：粉丝列表 GET cgi-bin/user/get）
scripts/clienttoken-debug.sh call cgi-bin/user/get GET "next_openid=" ""

# 提交 JSON 请求体（可使用 @file 引用）
scripts/clienttoken-debug.sh call cgi-bin/message/custom/send POST "" @payload.json

# 验证消息签名
ts=$(date +%s)
nonce=$RANDOM
scripts/clienttoken-debug.sh validate "<expected_signature>" "$ts" "$nonce" "hello"
```

该脚本默认读取以下环境变量：
- `CLIENTTOKEN_BASE_URL`（默认 `http://127.0.0.1:7072`）
- `CLIENTTOKEN_API_TOKEN`
- `CLIENTTOKEN_PROVIDER_CODE` / `CLIENTTOKEN_APP_CODE` / `CLIENTTOKEN_AUTH_MODE`
- `CLIENTTOKEN_CONFIG_PATH`

执行 `scripts/clienttoken-debug.sh --help` 可查看所有子命令与参数。

## 7. 故障排查

| 问题 | 排查思路 |
| --- | --- |
| `client_token_providers 未配置` | 确认 `config.yaml` 中存在该节点，或设置 `CLIENTTOKEN_CONFIG` 指向正确文件。 |
| `redis: dial tcp` / `ping redis` 失败 | 检查 `CLIENTTOKEN_REDIS_ADDR`、Redis 服务是否启动。 |
| 调试页提示 `unauthorized` | API Token 不一致，确保浏览器携带的 Token 与服务端 `CLIENTTOKEN_API_TOKEN` 一致。 |
| 调用微信公众号接口返回 `40001/40002` | Token 失效或 AppSecret 错误，重新刷新 Token，必要时检查微信控制台。 |

如仍无法定位，可在 `logs/clienttoken-server-error.log` 中查看详细日志，或开启 `http_debug: true` 观测 HTTP 请求/响应。
