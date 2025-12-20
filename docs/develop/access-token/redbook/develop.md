# RedBook JuGuang AccessTokenClient 开发指南

> 目标：在 **同一套** `cmd/accesstoken` 调试环境中复用小红书聚光 SDK，不再另起脚本；本文与已有的 Google/BiliBili 文档并列，所有命令、配置、Redis/Flow 管控完全一致。

## 1. 适用场景

- 需要用 MediaX SDK 访问聚光开放平台（账号余额、投放、笔记、报表等接口）。
- 希望像 Google/BiliBili 一样，统一在 `go run ./cmd/accesstoken/server -config config.yaml` 的 `/debug` 页面完成授权、Flow 回填与 API 调试。
- 已经拥有聚光开放平台应用（`client_id/client_secret`），并且愿意把授权回调指向 `http://127.0.0.1:7071/debug/callback` 以便复用本地调试台。

## 2. 目录与依赖

| 目录 | 说明 |
| --- | --- |
| `pkg/client/redBook/juGuang/accessTokenClient` | 聚光 AccessToken 客户端，包含 `account`/`note`/`promote`/`dataReport`/`tools` 等子域。 |
| `pkg/client/config/redbook.go` | `RedBookJuGuangConfig` 定义，内嵌 `ClientConfig + OAuthConfig`，支持 `oauth_key`/`GetOAuthToken`。 |
| `cmd/accesstoken/server` | Web 调试台（`/debug`）与 HTTP 接口；Provider 列表自动读取 `config.yaml`。 |
| `docs/develop/access-token/redbook/debug.md` | 本地调试台使用指南（与 Google/BiliBili 结构一致）。 |

核心依赖：Go 1.21、`github.com/ArtisanCloud/MediaXCore`（日志/Cache/BaseClient）、Redis（可选 Flow 持久化）。

## 3. 环境变量一览

| 变量 | 用途 |
| --- | --- |
| `REDBOOK_JUGUANG_CLIENT_ID` / `REDBOOK_JUGUANG_CLIENT_SECRET` | OAuth Client 凭证。 |
| `REDBOOK_JUGUANG_OAUTH_URL` | 授权地址（官方默认 `https://ad.xiaohongshu.com/oauth2/authorize`，可根据代理覆盖）。 |
| `REDBOOK_JUGUANG_ACCESS_TOKEN_URL` | Token 交换地址（官方默认 `https://ad.xiaohongshu.com/oauth2/access_token`）。 |
| `REDBOOK_JUGUANG_REDIRECT_URL` | 回调地址，建议保持 `http://127.0.0.1:7071/debug/callback` 以便调试台接收授权。 |
| `REDBOOK_JUGUANG_SCOPE` | 授权范围，例如 `notes.read,notes.write,ads.read`。 |
| `REDBOOK_JUGUANG_OAUTH_KEY` | 用于区分租户/账号的 key，会显示在 `/debug` Provider 卡片与 Flow 记录里。 |
| `REDBOOK_ACCESS_TOKEN` / `XHS_ACCESS_TOKEN` | （可选）绕过 OAuth 的临时 AccessToken，`/accesstoken/token` 会按 env→config→Flow 的优先级加载。 |
| `ACCESSTOKEN_REDIS_*` | 配置调试服务的 Redis（addr/db/username/password），未设置时默认尝试 `127.0.0.1:6379`，失败才回退内存。 |

## 4. `config.yaml`/`config.example.yaml` 片段

```yaml
access_token_providers:
  providers:
    - code: "redbook"
      name: "小红书"
      apps:
        - code: "juguang"
          name: "小红书聚光"
          provider_code: "redbook_juguang"
          api_version: "v1"
          auth_modes:
            - key: "default"
              label: "默认 OAuth"
              redbook_juguang_config:
                api_url: "https://adapi.xiaohongshu.com"
                proxy_api_url: ""
                timeout: 5
                http_debug: false
                oauth:
                  oauth_url: "${REDBOOK_JUGUANG_OAUTH_URL:-https://ad.xiaohongshu.com/oauth2/authorize}"
                  access_token_url: "${REDBOOK_JUGUANG_ACCESS_TOKEN_URL:-https://ad.xiaohongshu.com/oauth2/access_token}"
                  client_id: "${REDBOOK_JUGUANG_CLIENT_ID}"
                  client_secret: "${REDBOOK_JUGUANG_CLIENT_SECRET}"
                  redirect_url: "${REDBOOK_JUGUANG_REDIRECT_URL:-http://127.0.0.1:7071/debug/callback}"
                  scope: "${REDBOOK_JUGUANG_SCOPE:-notes.read,notes.write}"
                oauth_key: "${REDBOOK_JUGUANG_OAUTH_KEY:-default}"
```

> ⚠️ 所有敏感字段请通过环境变量注入；配置文件只保留占位符，避免误提交。`oauth_key` 会直接出现在 `/debug` 页面，便于识别 Flow、Redis key 与日志。

## 5. 快速上手

### 5.1 调试台操作

1. 复制 `config.example.yaml`，填入聚光 `client_id/client_secret/scope` 与 `redirect_url=http://127.0.0.1:7071/debug/callback`。
2. 运行 `go run ./cmd/accesstoken/server -config config.yaml`，确认日志包含 `provider=redbook_juguang`。
3. 打开 `http://127.0.0.1:7071/debug` → 选择“小红书 / 聚光” → 点击“发起授权”完成 OAuth。
4. 在 “授权记录” 中点击“填充”，再到 “API 调用” JSON 写入：
   ```json
   {
     "provider_code": "redbook_juguang",
     "provider_app": "juguang",
     "provider_auth_mode": "default",
     "config_path": "config.yaml",
     "action": "redbook.account.balance",
     "payload": {"advertiser_id": 123456789}
   }
   ```
5. 点击 “执行 API”，即可看到 `/api/open/jg/account/balance/info` 的响应；若需命令行调用，参见 quickstart 示例（同样走 `/accesstoken/call`）。

### 5.2 Go 代码示例

```go
package redbook

import (
    "context"

    mediax "github.com/ArtisanCloud/MediaX/pkg/client"
    "github.com/ArtisanCloud/MediaX/pkg/client/config"
    accountSchema "github.com/ArtisanCloud/MediaX/pkg/client/redBook/juGuang/accessTokenClient/account/schema"
    "github.com/ArtisanCloud/MediaXCore/pkg/cache"
    logcfg "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
    "github.com/ArtisanCloud/MediaXCore/utils/object"
)

func FetchBalance(ctx context.Context, mxCfg *config.MediaXConfig, rbCfg *config.RedBookJuGuangConfig, c cache.ICache, advertiserID string) (*accountSchema.GetAccountBalanceRes, error) {
    client := mediax.NewMediaX(mxCfg, c)
    rbCfg.GetOAuthToken = func(key string, refresh bool) object.HashMap {
        // 业务方自定义，示例：从 Redis/DB 读取 AccessToken
        return object.HashMap{"access_token": LoadTokenFromVault(key)}
    }

    cli, err := client.CreateRedBookJuGuangACClient(rbCfg)
    if err != nil {
        return nil, err
    }

    return cli.GetAccountClient().GetAccountBalance(ctx, &accountSchema.GetAccountBalanceReq{
        AdvertiserID: advertiserID,
    })
}
```

- `MediaXConfig` 中的 `Logger`/`Cache` 可直接复用项目里已有的配置。
- 如果期望 SDK 自动刷新 Token，可在 `rbCfg.ClientConfig.OAuthConfig.RefreshToken` 中写入长期凭证，并在回调后更新。

> 详细的网页交互说明、Redis Key 结构与截图，请见 `docs/develop/access-token/redbook/debug.md`，与 Google/BiliBili 文档保持一致。

## 7. Redis / Flow 与日志

- 未显式设置 `ACCESSTOKEN_REDIS_ADDR` 时，服务会优先尝试 `127.0.0.1:6379`。连接失败才退回内存，并在页面顶部显示 “Flow 数据当前存储在内存中…” 警告。
- `flow_ttl_seconds` 默认为 86400，可通过环境变量 `ACCESSTOKEN_FLOW_TTL_SECONDS` 调整。Redis key TTL 亦保持一致。
- 所有 OAuth/API 日志均包含 `provider=redbook_juguang`、`provider_app=juguang`、`flow_id`、`token_source`，并使用 `mask.Token` 工具自动脱敏。

## 8. 常见问题

| 问题 | 处理方式 |
| --- | --- |
| “unsupported provider” | 检查 `config.yaml` 中 `provider_code` 是否为 `redbook_juguang`，以及 `auth_modes` 是否包含当前 `provider_auth_mode`。 |
| “OAuth scope 未配置” | 新增的字段 `oauth.scope` 必填；可通过 `REDBOOK_JUGUANG_SCOPE` 注入。 |
| Redis 报错后 Flow 丢失 | 启动日志会提示 `storage_backend=memory`，说明 Redis 不通；设置正确的 `ACCESSTOKEN_REDIS_*` 后重启即可。 |
| 不想每次都发起 OAuth | 可将已有 AccessToken 写入 `REDBOOK_ACCESS_TOKEN` 或 `config.yaml -> oauth.access_token`，调试台会自动复用。 |

> 只要严格按照本文与 `google/bilibili` 文档的结构维护配置，即可同时在 `/debug` 中切换多个 Provider，RedBook 不需要特殊的调试命令。
