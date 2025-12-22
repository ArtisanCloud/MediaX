# DouYin AccessToken 开发指南

> 目标：让抖音开放平台的 AccessToken 客户端像 Google/BiliBili 一样在 `cmd/accesstoken` 调试台中完成 OAuth、Flow 管理与 API 调试。本文面向 SDK/运营同学，说明配置、环境变量与常用命令。

> ⚠️ **仅支持单实例 `<app>=default`。** 如果需要多套 DouYin 应用，请在不同实例/配置文件中分别部署；同一 `config.yaml` 下的 DouYin AccessToken Flow 将始终绑定 `provider_app=douyin`、`auth_mode=default`。

## 1. 适用场景

- 使用 MediaX SDK 访问抖音开放平台的用户授权接口（视频上传、粉丝画像、IM 消息等）。
- 希望沿用 `go run ./cmd/accesstoken/server` 的 `/debug` 页面完成 OAuth/Flow/Token 操作，而不是单独写脚本。
- 已申请 DouYin 开放平台应用，并可将授权回调指向 `http://127.0.0.1:7071/debug/callback` 进行本地调试。

## 2. 目录与依赖

| 目录 | 说明 |
| --- | --- |
| `pkg/client/byteDance/douYin/accessTokenClient` | 抖音 AccessToken SDK，涵盖 `content/`、`connection/`、`im/`、`market/` 等模块。 |
| `pkg/client/config/douyin.go` | `ByteDanceDouYinConfig` 定义，继承 `ClientConfig`，可指定 `oauth_key`、代理、超时等。 |
| `cmd/accesstoken/server` | 调试服务与 `/debug` 页面，Provider 元数据从 `config.yaml` 读取。 |
| `docs/develop/access-token/byteDance/debug.md` | 与本文配套的调试台使用说明。 |

依赖：Go 1.21、`github.com/redis/go-redis/v9`（可选，用于 Flow 持久化）、`github.com/ArtisanCloud/MediaXCore`。

## 3. 环境变量一览

| 变量 | 说明 |
| --- | --- |
| `DOUYIN_CLIENT_ID` / `DOUYIN_CLIENT_SECRET` | OAuth client 凭证（``client_key``/``client_secret``）。 |
| `DOUYIN_REDIRECT_URL` | OAuth 回调地址，调试时建议 `http://127.0.0.1:7071/debug/callback`。 |
| `DOUYIN_SCOPE` | 授权范围，例如 `user_info,video.create`。 |
| `DOUYIN_OAUTH_KEY` | 用于区分租户/账号的 key，会显示在 `/debug` Provider 卡片与 Flow 记录中，默认 `default`。 |
| `DOUYIN_ACCESS_TOKEN` | （可选）直接注入 AccessToken，便于绕过 OAuth 测试。 |
| `ACCESSTOKEN_REDIS_*` | 调试服务的 Redis 地址/DB/密码，未设置时默认尝试 `127.0.0.1:6379` 并在失败后回退内存。 |

## 4. `config.yaml` 配置示例

```yaml
access_token_providers:
  providers:
    - code: "byte_dance"
      name: "字节跳动"
      apps:
        - code: "douyin"
          name: "抖音 Douyin"
          provider_code: "byte_dance_douyin"
          api_version: "v1"
          auth_modes:
            - key: "default"
              label: "默认 OAuth"
              byte_dance_douyin_config:
                api_url: "https://open.douyin.com"
                timeout: 5
                http_debug: false
                oauth:
                  authorize_url: "${DOUYIN_AUTHORIZE_URL:-https://open.douyin.com/platform/oauth/connect/}"
                  access_token_url: "${DOUYIN_ACCESS_TOKEN_URL:-https://open.douyin.com/oauth/access_token/}"
                  client_id: "${DOUYIN_CLIENT_ID}"
                  client_secret: "${DOUYIN_CLIENT_SECRET}"
                  redirect_url: "${DOUYIN_REDIRECT_URL:-http://127.0.0.1:7071/debug/callback}"
                  scope: "${DOUYIN_SCOPE:-user_info,video.create}"
                oauth_key: "${DOUYIN_OAUTH_KEY:-default}"
```

> 所有敏感字段请使用环境变量覆盖。若需要多环境/多租户，可在 `auth_modes` 下追加更多 `key`。

## 5. Go/CLI 示例

```go
package main

import (
    "context"
    mediax "github.com/ArtisanCloud/MediaX/pkg/client"
    "github.com/ArtisanCloud/MediaX/pkg/client/config"
    "github.com/ArtisanCloud/MediaXCore/pkg/cache/driver/redis"
)

func main() {
    cacheStore, _ := redis.NewRedis(map[string]interface{}{"Addr": "127.0.0.1:6379"})
    mx := mediax.NewMediaX(&config.MediaXConfig{}, cacheStore)
    cfg := &config.ByteDanceDouYinConfig{ClientConfig: &config.ClientConfig{}}
    cfg.GetOAuthToken = func(key string, refresh bool) object.HashMap {
        return object.HashMap{"access_token": os.Getenv("DOUYIN_ACCESS_TOKEN"), "expires_in": float64(7200)}
    }
    client, _ := mx.CreateByteDanceDouYinACClient(cfg)
    resp, err := client.GetContentVideoClient().List(context.Background(), &video.ListReq{Cursor: 0})
    fmt.Println(resp, err)
}
```

CLI：`go run ./cmd/accesstoken -provider byte_dance -provider-app douyin -action tokens.list -config config.yaml`

## 6. 与 `/debug` 调试台联动

1. 运行 `go run ./cmd/accesstoken/server -config config.yaml`。
2. 打开 `http://127.0.0.1:7071/debug` → 选择 “字节跳动 / 抖音 Douyin” → “同步模板”。
3. 点击 “发起授权” 完成 OAuth；Flow 写入内存/Redis（key 形如 `accesstoken:oauth:byte_dance_douyin:default:<mode>`）。
4. 点击 “解析 AccessToken” 或 “API 调用” 可直接使用最新 Flow Token，详见 `docs/develop/access-token/byteDance/debug.md`。
5. 如果 Flow 表中出现 `status=need reauth` 或 TTL 即将过期，刷新按钮会提示需要重新授权。

## 7. 常见问题

| 问题 | 处理方式 |
| --- | --- |
| 收到 `unsupported provider` | 检查 `provider_code` 是否写成 `byte_dance_douyin`，且 `auth_modes` 中存在当前 key。 |
| “OAuth scope 未配置” | 需要在 `config.yaml` 或 `DOUYIN_SCOPE` 中填入 DouYin 要求的 scope。 |
| 调试台提示 `storage_backend=memory` | Redis 未连接，授权记录仅存内存；设置 `ACCESSTOKEN_REDIS_ADDR` 后重启即可。 |
| `/accesstoken/call` 返回 `need reauth` | Flow 中的 refresh_token 已失效或无法刷新。请在 `/debug` 重新发起授权，新的 Flow 会自动覆盖旧记录。 |
| 429 或 5xx 重试 | DouYin action 被限流时，服务会记录 `retry_count` 与 `last_backoff_ms`，最多重试 3 次。可通过日志定位 action 并适当降低频率。 |
| CLI 调用 400/401 | 确认 `client_id/client_secret` 正确、应用开了对应权限。可在 `/debug` 的 AccessToken 面板查看 token 来源及 TTL。 |

> DouYin OAuth 流程可能要求提前登录抖音主体账号，首次授权会跳转到官方登录页，这是正常行为。完成登录后即可回到调试台并生成 Flow 记录。

## 8. Need Reauth 与自动刷新排障

1. `/accesstoken/call` 会在 Flow 的 AccessToken 剩余 ≤5 分钟时自动调用 DouYin Refresh API，并把新的 `access_token`、`refresh_token`、`expires_in` 写回 Flow。
2. 刷新成功：响应 JSON 中 `retry_count=0`，日志 `event=token.call provider=byte_dance_douyin status=refreshed`，Flow 记录的 `status` 会更新为 `refreshed` 并重置 `flow_ttl_seconds`。
3. 刷新失败：接口返回 `{"error":"need reauth"}`；日志包含 `invalidate flow reason=douyin_refresh_failed:*`；Flow 会被移除，需要重新在 `/debug` 授权。
4. 诊断步骤：
   - 查看 `logs/accesstoken-server-error.log` 是否有 DouYin refresh 错误。
   - 使用 `redis-cli keys accesstoken:oauth:byte_dance_douyin*` 确认 Flow 是否仍存在。
   - 重新授权后建议运行一次 `/accesstoken/call` 验证 `retry_count=0`。

## 9. CLI / Curl 脚本示例

```bash
# 解析 AccessToken
curl -H "Authorization: Bearer $ACCESSTOKEN_API_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"provider_code":"byte_dance_douyin","config_path":"config.yaml"}' \
     http://127.0.0.1:7071/accesstoken/token | jq

# 调用 DouYin API 并自动刷新
curl -H "Authorization: Bearer $ACCESSTOKEN_API_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"provider_code":"byte_dance_douyin","action":"douyin.video.list","payload":{"count":5}}' \
     http://127.0.0.1:7071/accesstoken/call | jq
```

所有脚本均可直接复制 `/debug` 中的 JSON；保持 `provider_app=douyin` 即可指向默认实例。

## 10. Success Criteria 验证工具

| 指标 | 工具/步骤 |
| --- | --- |
| SC-001 DouYin OAuth 成功率 ≥80% | `scripts/douyin_sc001_oauth_success.sh logs/accesstoken-server-info.log 5` 读取最近 5 次 `oauth.complete`/`oauth.start.error`，统计成功率并在低于 80% 时提示排障建议。 |
| SC-002 Flow 回填命中率 ≥95% | `EXPECTED_BACKEND=redis scripts/douyin_sc002_flow_backfill.sh http://127.0.0.1:7071 20`（Redis 模式） + `EXPECTED_BACKEND=memory ...`（内存模式），脚本会遍历 `/api/oauth/tokens` 并逐条调用 `/accesstoken/flow/replay` 计算命中率。 |
| SC-003 `/accesstoken/call` p95 ≤2s | `ITERATIONS=15 scripts/douyin_sc003_call_benchmark.sh http://127.0.0.1:7071` 对 DouYin action 连续调用，输出平均值与 p95，超过 2 秒会标黄。 |
| SC-004 文档走查无阻碍 | 依照下方模板安排 2 位未参与开发的伙伴逐步完成 Quickstart，记录耗时与阻塞项；若失败则在本文/Quickstart 中补充排障指引。 |

> 三个脚本默认读取 `ACCESSTOKEN_API_TOKEN` 和 `BASE_URL` 环境变量，可根据实际部署改写。建议将输出结果与日志一同归档，以便复验。

### SC-004 文档走查记录模板

| 日期 | 参与人 | 耗时 (分钟) | 结果 | 备注/阻塞 |
| --- | --- | --- | --- | --- |
|  |  |  | ✅/⚠️ |  |
|  |  |  | ✅/⚠️ |  |

执行流程：

1. 让被访者仅凭 `specs/006-bytedance-access-token/quickstart.md`、本文以及 `docs/develop/access-token/byteDance/debug.md` 操作。
2. 观察是否能在 60 分钟内完成配置→授权→`/accesstoken/call`→Redis 证据截图。
3. 将反馈记录到上表，并在相关文档中补充排障说明。若有阻塞，需复测通过后再关闭 SC-004。
