# DouYin ClientToken 开发指南

> 目标：在 MediaX 内使用 `cmd/clienttoken/server`（与 WeChat ClientToken 一致）完成抖音服务端 token 刷新、Redis 缓存与 API 调试。

## 1. 适用场景

- 只需要 ClientToken 的服务端接口（不依赖用户授权），如内容模板、任务、沙箱、票据等。
- 希望与 WeChat ClientToken 相同的步骤：配置 `client_token_providers` → 启动 `cmd/clienttoken/server` → 浏览器 `/client-token/debug` 操作。

## 2. 目录与依赖

| 目录 | 说明 |
| --- | --- |
| `pkg/client/byteDance/douYin/clientTokenClient` | 抖音 ClientToken 客户端，包含 `content/`、`tools/` 模块。 |
| `pkg/client/config/douyin.go` | `ByteDanceDouYinConfig`，在 ClientToken 模式下需要填写 `client_token` 相关字段。 |
| `cmd/clienttoken/server` | ClientToken 调试服务；Provider 列表来自 `client_token_providers`。 |
| `docs/develop/client-token/byteDance/debug.md` | 调试台操作指南。 |

依赖：Go 1.21、Redis（缓存 client_token）、`MediaXCore`。

## 3. 环境变量

| 变量 | 说明 |
| --- | --- |
| `DOUYIN_CLIENT_KEY` / `DOUYIN_CLIENT_SECRET` | DouYin ClientToken client_key/client_secret。 |
| `CLIENTTOKEN_CONFIG` | `cmd/clienttoken/server` 的配置文件路径（默认 `config.yaml`）。 |
| `CLIENTTOKEN_LISTEN_ADDR` | 调试服务监听地址，默认 `:7072`。 |
| `CLIENTTOKEN_API_TOKEN` | 调试页面/HTTP API 使用的 Bearer token。 |
| `CLIENTTOKEN_REDIS_ADDR/DB/PASS` | Redis 缓存配置。 |

## 4. `client_token_providers` 片段

```yaml
client_token_providers:
  providers:
    - code: "byte_dance"
      name: "字节跳动"
      apps:
        - code: "douyin_service"
          name: "抖音服务端"
          provider_code: "byte_dance_douyin_clienttoken"
          auth_modes:
            - key: "default"
              label: "默认"
              byte_dance_douyin_config:
                api_url: "https://open.douyin.com"
                timeout: 5
                http_debug: false
                client_token:
                  client_key: "${DOUYIN_CLIENT_KEY}"
                  client_secret: "${DOUYIN_CLIENT_SECRET}"
                cache:
                  redis_key: "clientToken:douyin:${DOUYIN_CLIENT_KEY}"
                  ttl_seconds: 7000
                  refresh_before_seconds: 600
                redis:
                  addr: "${CLIENTTOKEN_REDIS_ADDR:-127.0.0.1:6379}"
                  db: "${CLIENTTOKEN_REDIS_DB:-0}"
                api_token: "${CLIENTTOKEN_API_TOKEN:-dev-clienttoken}"
```

> 所有凭证需使用环境变量，避免提交到仓库。

## 5. 启动服务

```bash
export CLIENTTOKEN_CONFIG=${CLIENTTOKEN_CONFIG:-config.yaml}
export CLIENTTOKEN_LISTEN_ADDR=${CLIENTTOKEN_LISTEN_ADDR:-:7072}
export CLIENTTOKEN_API_TOKEN=${CLIENTTOKEN_API_TOKEN:-dev-clienttoken}
go run ./cmd/clienttoken/server -config "$CLIENTTOKEN_CONFIG"
```

日志应输出 `clienttoken: server listening addr=:7072 ... provider=byte_dance`。

## 6. Go 示例

```go
mx := mediax.NewMediaX(&config.MediaXConfig{}, cacheStore)
cfg := &config.ByteDanceDouYinConfig{
    ClientConfig: &config.ClientConfig{
        ClientToken: &config.ClientTokenConfig{
            ClientKey:    os.Getenv("DOUYIN_CLIENT_KEY"),
            ClientSecret: os.Getenv("DOUYIN_CLIENT_SECRET"),
        },
    },
}
client, _ := mx.CreateByteDanceDouYinCTClient(cfg)
resp, err := client.GetContentTaskClient().List(ctx, &task.ListReq{})
```

## 7. 与 `/client-token/debug` 联动

详见 `docs/develop/client-token/byteDance/debug.md`：页面提供 Token 刷新、缓存查看、API 调试与回调日志，流程与 WeChat ClientToken 相同。

## 8. 常见问题

| 问题 | 解决办法 |
| --- | --- |
| 缓存 Key 不更新 | 确认 `clientToken:douyin:<client_key>` 是否存在，或设置 `CLIENTTOKEN_REDIS_ADDR`。 |
| 页面提示 `provider code not found` | 检查 `client_token_providers` 配置是否包含 `byte_dance/douyin_service`。 |
| refresh 失败 | DouYin client_key/secret 无效，或 API 返回签名错误，查看日志 `clienttoken: refresh`。 |

> 抖音 ClientToken 只适用于服务端接口，若需要用户授权的 AccessToken，请参考 `docs/develop/access-token/byteDance/` 中的指南。
