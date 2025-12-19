# BiliBili AccessTokenClient 开发指南

> 说明 MediaX 如何封装 B 站开放平台的 AccessToken 客户端、何时适合接入、如何初始化配置以及和调试台联动。接口矩阵请参考 `docs/plan/bilibili/access_token_client.md`。

## 1. 适用场景
- 已在哔哩哔哩开放平台创建应用，并拿到 `client_id/client_secret` 与授权回调域名。
- 需要批量管理视频、专栏、数据、直播、第三方开播等 API，统一走 MediaX 的 Logger/Cache/HTTP 基座。
- 希望沿用 AccessToken Server、CLI 与 `/debug` 页面校验 `access_token`、存档授权记录。

## 2. 前置条件
1. Go 1.18+，并已拉取 `github.com/ArtisanCloud/MediaX`。
2. 准备好 Redis 或其他外部存储，用于保存 OAuth 结果（`access_token/refresh_token/expires_in`）。
3. 在 config.yaml 的 `access_token_providers` 节点注册 BiliBili Provider，或在业务代码中直接构造 `config.BiliBiliConfig`。
4. 获取授权：B 站 OAuth 依赖 `authorization_code`，需要引导创作者登录一次，拿到 `code` 后换取长效 `refresh_token` 并妥善存储。

## 3. 快速开始（Go SDK）
```go
package bilibili

import (
    "context"

    mediax "github.com/ArtisanCloud/MediaX/pkg/client"
    accessTokenClient "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient"
    "github.com/ArtisanCloud/MediaX/pkg/client/bilibili/b/accessTokenClient/video/schema"
    "github.com/ArtisanCloud/MediaX/pkg/client/config"
    "github.com/ArtisanCloud/MediaXCore/pkg/cache"
    "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
    "github.com/ArtisanCloud/MediaXCore/utils/object"
)

type Store interface {
    LoadToken(account string) (TokenRecord, error)
}

type TokenRecord struct {
    AccessToken  string
    RefreshToken string
    ExpiresIn    int
}

func SubmitDraft(ctx context.Context, account string, mxCfg *config.MediaXConfig, store Store, redis cache.ICache) error {
    mx := mediax.NewMediaX(mxCfg, redis)
    biliCfg := &config.BiliBiliConfig{
        ClientConfig: &config.ClientConfig{
            BaseConfig: &config.BaseConfig{
                ApiUrl:    config.BiliBiliAPIUrl,
                Timeout:   10,
                HttpDebug: true,
            },
            OAuthConfig: &config.OAuthConfig{
                ClientID:     getenv("BILIBILI_CLIENT_ID"),
                ClientSecret: getenv("BILIBILI_CLIENT_SECRET"),
                RedirectUrl:  "https://example.com/oauth/bilibili/callback",
                Scope:        "ARC_BASE,ARTICLE_BASE,LIVE_ROOM_DATA",
            },
        },
        GetOAuthToken: func(key string, refresh bool) object.HashMap {
            token, err := store.LoadToken(account)
            if err != nil {
                panic(err)
            }
            return object.HashMap{
                "access_token":  token.AccessToken,
                "refresh_token": token.RefreshToken,
                "expires_in":    float64(token.ExpiresIn),
            }
        },
    }

    client, err := accessTokenClient.NewBiliBiliACClient(biliCfg, mx.Logger, redis)
    if err != nil {
        return err
    }

    videoClient := client.GetVideoClient()
    _, err = videoClient.SubmitVideo(ctx, &schema.BiliBiliVideoSubmitVideoReq{
        UploadToken: "upload-token",
        Title:       "MediaX 自动投稿",
        Tid:         171, // 分区 ID
        Tag:         "MediaX,自动化",
        Desc:        "通过 MediaX SDK 提交",
        Copyright:   1,
    })
    return err
}
```
要点：
- **Token 注入**：`GetOAuthToken` 是 BiliBili 访问令牌的唯一入口，SDK 会在每次 HTTP 请求前调用一次，开发者可自由决定从 DB/Redis/自建服务读取。
- **客户端复用**：同一个 `BiliBiliACClient` 持有 `video/user/article/data/live/liveWS/liveThirdParty` 等子客户端，不必重复创建。
- **HTTP Debug**：当 `BaseConfig.HttpDebug=true` 时，`kernel.BaseClient` 会打印完整 URL/Headers，便于排查 401/签名问题。

## 4. 配置模板（access_token_providers）
```yaml
access_token_providers:
  - code: bilibili
    name: BiliBili
    apps:
      - code: content_center
        name: 内容运营
        provider_code: bilibili_arcopen
        auth_modes:
          - key: default
            label: 默认 OAuth
            bilbili_config:
              api_url: "https://member.bilibili.com"
              timeout: 10
              http_debug: false
              oauth:
                access_token_url: "https://member.bilibili.com/x/credential/token"
                oauth_url: "https://member.bilibili.com/x/credential/oauth2/authorize"
                client_id: "${BILIBILI_CLIENT_ID}"
                client_secret: "${BILIBILI_CLIENT_SECRET}"
                redirect_url: "https://dev.example.com/oauth/bilibili/callback"
                scope: "ARC_BASE,ARTICLE_BASE,LIVE_ROOM_DATA"
              meta:
                oauth_key: "tenantA.bilibili"
```
建议：
- 把所有敏感字段放入环境变量或 Secret Manager，避免直接写在仓库。
- `oauth_key`（可写入 `meta` 或自定义字段）用于识别租户/账号，调试页会回显该值。
- 若要在 AccessToken Server 的 `/debug` 页面里使用，只需把 `config_path` 指向此文件，前端就会读取 Provider/App/Mode 列表。

## 5. AccessToken 生命周期
- **授权**：调试台或业务系统可引导创作者访问 `/accesstoken/oauth/start`，得到 `authorize_url` → 登录后回调 `/debug/callback`，AccessToken Server 会把 `code` 交换成 `access_token/refresh_token` 并写入 Redis：`accesstoken:oauth:bilibili:content_center:default`。
- **刷新**：`GetOAuthToken` 可根据 `refresh` 参数决定是否强制刷新。如果返回的 `access_token` 已过期，`kernel.TokenHandler` 会触发 `refresh_token` 流程；若希望完全由业务方控制，可忽略 `refresh` 并直接返回最新 token。
- **缓存与日志**：所有 HTTP 调用会打印 `bilibili: request path=...`、`bilibili: response code=...`，并自动脱敏 token。可结合 `logs/accesstoken-server-info.log` 定位问题。

## 6. 功能范围速览
- 视频：上传、分片合并、封面、投稿、查询列表/详情、生成客户端分享 URL。
- 专栏：增删改查、分类、卡片、图片上传、文集合并。
- 数据中心：用户/稿件/专栏维度的实时与增量数据。
- 用户：OpenID 权限范围、基础信息、UnionID。
- 直播：房间信��、长链心跳、第三方开播授权/推流/收流。

详细参数、请求体结构与注意事项可在 `pkg/client/bilibili/b/accessTokenClient/*` 以及 `docs/plan/bilibili/access_token_client.md` 中查阅，新增接口时务必要同步更新两处文档。

## 7. 调试与 CLI
- Web：`go run ./cmd/accesstoken/server -config config.yaml` 后访问 <http://127.0.0.1:7071/debug>，即可在浏览器中完成 OAuth、查看 AccessToken、复制 Flow ID。针对 BiliBili 目前仅开放“AccessToken 解析 + OAuth 回调记录”，API 调试面板仍只支持 Google YouTube（后续会扩展）。见 `debug.md`。
- CLI：`go run ./cmd/accesstoken -provider bilibili -provider-app content_center -auth-mode default -action token` 可用于校验配置是否生效；真正的业务调用建议直接使用 Go SDK 或自建 HTTP 层包装。

## 8. 常见集成注意事项
| 问题 | 处理方法 |
| --- | --- |
| `invalid client` | 检查 `client_id/client_secret` 是否与后台一致、OAuth 回调域名是否实名匹配。|
| SDK 提示 "missing access token" | `GetOAuthToken` 返回的对象缺少 `access_token` 字段，或 Redis 中读取失败。请打印日志确认。|
| 频繁 401 | 确认 `expires_in` 是否正确写入，或 B 站后台是否吊销了 `refresh_token`。建议在 `GetOAuthToken` 中捕获 `refresh==true` 的场景主动刷新。|
| 上传接口 413/400 | `HttpUpload` 默认通过 `BaseClient` 直传到 B 站域名，注意文件大小、Content-Type、`name`/`part_number` 参数需与官方要求一致。|
| 调试台 API 模块不可选 BiliBili | 当前版本 `POST /accesstoken/call` 仅支持 Google YouTube。BiliBili 相关 API 请使用 Go SDK 或未来的定制 CLI。|

> 当你需要把 AccessToken 能力以 HTTP 的形式暴露给外部团队时，可复用 AccessToken Server 的 `/accesstoken/token` 接口，让他们直接获取当前 token 及来源，再结合本指南提供的 SDK 调用具体业务接口。
