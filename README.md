# MediaX

## 简介

MediaX 是一个支持多平台内容发布的轻量级 SDK，旨在帮助开发者快速集成到主流自媒体平台，实现统一的内容发布接口。

ArisanCloud 团队已经成功维护了 [PowerWechat](https://powerwechat.artisan-cloud.com) 开源项目，基本上 SDK 的框架接口趋于稳定，所以我们有延伸出这个 MediaX，同时也做了模块优化， 按需加载模块，支持更多平台的接口使用。

## MediaX 系列产品介绍

- MediaX：一个 开源的 Golang SDK，可以直接被其他项目引用。
- MediaX Studio：基于 MediaX 构建的商业化服务，包含 API（gRPC/HTTP）、高级功能和闭源插件，面向企业用户。

## 功能特点

1. **多平台支持**：支持主流的自媒体平台如抖音、小红书，youtube 等。
2. **统一接口**：通过标准化的接口，实现对多个平台的内容发布。
3. **开源与闭源结合**：Studio 支持插件化架构，开源插件可自由扩展，闭源插件提供增强功能。
4. **跨语言支持**：Studio 通过 gRPC 实现跨语言调用，支持 Go 和其他语言集成。
5. **可扩展性**：灵活的 Provider 机制，轻松接入新平台。

## MediaX 快速开始

### 环境要求

- Go 1.18 或更高版本

### 使用知识

- OAuth2.0 授权流程，可以参考[理解 OAuth 2.0](https://www.ruanyifeng.com/blog/2014/05/oauth_2_0.html)文章

### 安装

1. 初始化项目并下载 MediaX：

   ```bash
   go get github.com/ArtisanCloud/MediaX
   ```

2. 创建一个简单 客户端凭证模式（Client Credentials Grant）（非用户授权）的示例，本项目作者正在自己系统中使用，陆续会迭代版本：

   ```go

   import (
      "github.com/ArtisanCloud/MediaX/pkg/client"
      config2 "github.com/ArtisanCloud/MediaX/pkg/client/config"
      "github.com/ArtisanCloud/MediaX/pkg/utils"
      "github.com/ArtisanCloud/MediaXCore/pkg/cache"
      "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
      "github.com/ArtisanCloud/MediaXCore/utils/fmt"
      "github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient/publish/schema"
      "github.com/redis/go-redis/v9"
   )

   // 配置Media Client实例的信息
   mediaXClient := client.NewMediaX(&config2.MediaXConfig{
     &config.LogConfig{
         Level:   "debug",
         Console: true,
         File: config.FileConfig{
             Enable: true,
         },
     },
   }, c)

   // 从MediaXClient实例中获取到微信平台中公众号的实例，该实例是Client Token模式，不需要用户授权
   wechatOAClient, err := mediaXClient.MediaXClient.NewWeChatOfficialAccountCTClient(&config2.WeChatOfficialAccountConfig{
     ClientConfig: &ClientConfig{
         BaseConfig: &BaseConfig{
            Timeout: 30,
            HttpDebug: true,
         },
         OAuthConfig: &OAuthConfig{
            ClientID: "your client/app id"
            ClientSecret: "your client/app secret"
         },
     },
   })
   if err != nil {
     panic(err)
   }

   // 调用 wechatOAClient 的方法
   ctx := context.Background()
   var reqData = &schema.DraftAddReq{}
   resData, err := oaClient.GetPublishClient().DraftAdd(ctx, reqData)
   if err != nil {
      return nil, err
   }

   ```

3. OAuth 2.0 授权码模式（Authorization Code Grant）（用户授权）的示例，本项目作者正在自己系统中使用，陆续会迭代版本：

   ```go
   import (
      "github.com/ArtisanCloud/MediaX/pkg/client"
      config2 "github.com/ArtisanCloud/MediaX/pkg/client/config"
      "github.com/ArtisanCloud/MediaX/pkg/utils"
      "github.com/ArtisanCloud/MediaXCore/pkg/cache"
      "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
      "github.com/redis/go-redis/v9"
      "github.com/ArtisanCloud/MediaX/pkg/client/config"
      "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/video/schema"
      "github.com/ArtisanCloud/MediaXCore/utils/fmt"
      "github.com/ArtisanCloud/MediaXCore/utils/object"
   )

   googleYouTubeClient, err := mediaXClient.CreateGoogleYouTubeACClient(localConfig.GoogleYouTubeConfig)
   if err != nil {
   	panic(err)
   }

   // 设置AccessToken
   googleYouTubeClient.GoogleClient.TokenHandler.GetCustomToken = func(key string, refresh bool) object.HashMap {
   	fmt.Dump("GetCustomToken", key, refresh)
  	return object.HashMap{
   		// 这个acess token需要开发这来维护，或者可以通过MediaX Studio的UI界面来维护
   		"access_token": "72_ggzUdSgH99StJ2EhmuaIbHHUP9_3rDvdnQVQ9eoX5gwmNfuLpJgBUb5uPgdoh4aoVv9jYz3EKglRT73ppWqgRwzirNQM-bHaToDQ83ux1sFdCr5GK7jxYQfAESoCOEaAHAKWM",
   		"expires_in":   float64(7200),
  	}

### Zhihu SessionToken 配置片段

在 `config.yaml` 中启用 `zhihu_config.sessionToken`，为服务端、Authenticator、Harvester、回调与网络策略提供必要的字段：

```yaml
zhihu_config:
  sessionToken:
    service:
      base_url: ${ZH_SESSIONTOKEN_BASE_URL}
      api_token: ${ZH_SESSIONTOKEN_API_TOKEN}
      timeout: 30
    authenticator:
      entries:
        - type: pc
          url: ${ZH_SESSIONTOKEN_PC_ENTRY_URL}
      default_user_agent: ${ZH_SESSIONTOKEN_DEFAULT_UA}
      script_ids: [${ZH_SESSIONTOKEN_SCRIPT_ID}]
      captcha_strategy: ${ZH_SESSIONTOKEN_CAPTCHA_STRATEGY}
    harvester:
      watch_cookies: [z_c0]
      watch_headers: [X-XSRF-TOKEN]
      harvest_script_id: ${ZH_SESSIONTOKEN_HARVEST_SCRIPT_ID}
    callback:
      callback_secret: ${ZH_SESSIONTOKEN_CALLBACK_SECRET}
      max_retry: 3
      retry_backoff: [2, 4, 8]
    network:
      proxy_pool: ${ZH_SESSIONTOKEN_PROXY_POOL}
      ip_strategy: ${ZH_SESSIONTOKEN_IP_STRATEGY}
      request_timeout: 60
```

### 启动 SessionToken HTTP 服务

1. **统一环境变量**：确保 MediaX 与插件共用以下值（可写入 `.env`）：
   ```bash
   export POWERX_SESSION_TOKEN_BASE_URL="http://127.0.0.1:7070"
   export POWERX_SESSION_TOKEN_API_TOKEN="dev-session-token"
   export POWERX_SESSION_TOKEN_CALLBACK_URL="https://plugin.local/api/v1/admin/platforms/session-token/callback"
   ```
2. **启动服务**：在 MediaX 仓库根目录运行 `make sessiontoken`（内部执行 `go run ./cmd/sessiontoken -config config.yaml`，默认监听 `:7070`）。
3. **启动插件**：在插件或 MediaX Studio 侧，将 `POWERX_SESSION_TOKEN_*` 指向上一步的 BaseURL/API Token/Callback URL，再启动 `/publish/platforms` 等入口触发模拟登录。
4. **顺序要求**：必须先启动 SessionToken 服务，待日志出现 `sessiontoken: server listening...` 后再启动插件/浏览器容器，否则插件会因无法访问 `/session-token/flows` 返回 5xx。

### SessionToken 监控与排查

- 关注 `sessiontoken_metric`（Flow 创建/查询/完成）与 `sessiontoken_callback`（回调重试/成功）日志，字段中包含 `provider/provider_app/tenant_uuid/flow_id/retry/latency_ms`，便于将其采集到日志或指标系统实现 SLA 追踪。
- 快速排查可以直接 tail + ripgrep：
  ```bash
  tail -f logs/mediax.log | rg 'sessiontoken_(metric|callback)'
  ```
  典型日志：
  ```
  sessiontoken_metric: action=create_flow provider=zhihu provider_app=zhihu_article tenant_uuid=tenant_x account_id=acct_demo state=ui-flow flow_id=stf_xxx status=pending latency_ms=12 retry=0
  sessiontoken_callback: success provider=zhihu provider_app=zhihu_article tenant_uuid=tenant_x flow_id=stf_xxx state=ui-flow flow_status=succeeded retry=0 http_status=200 latency_ms=5
  ```
- 若 `latency_ms` 长期高于目标值，优先检查 Redis/第三方登录入口；若 `retry` >= 3，可结合 `last_error` 与插件回调响应定位网络问题。
   }


   // 调用 Youtube的VideoClient 的方法
   ctx := context.Background()
   video := googleYouTubeClient.GetVideoClient()
   res, err := video.List(ctx, &schema.YouTubeVideoListReq{})
   if err != nil {
   	panic(err)
   }
   fmt.Dump(res)

   ```

## 文档与接口说明

详细的功能接口文档和使用指南，请访问我们的[开发者中心](https://mediax.artisan-cloud.com)。在这里您可以找到:

- 完整的 API 参考文档
- 快速入门指南
- 最佳实践示例
- 常见问题解答

## 功能矩阵

| 平台     | 应用    | 授权类型          | 图文 | 视频 | 素材管理 | 评论管理 | 数据管理 |
| -------- | ------- | ----------------- | ---- | ---- | -------- | -------- | -------- |
| WeChat   | 公众号  | client_credential | ✔    | ✘    | ✔        | ✘        | ✘        |
| 小红书   | 聚光    | auth_code         | ✔    | ✔    | ✔        | ✘        | ✔        |
| 字节     | 抖音    | auth_code         | ✔    | ✔    | ✔        | ✔        | ✔        |
| 字节     | 抖音    | client_credential | ✔    | ✔    | ✔        | ✔        | ✔        |
| Bilibili | B 站    | auth_code         | ✔    | ✔    | ✔        | ✘        | ✔        |
| Google   | YouTube | auth_code         | ✔    | ✔    | ✔        | ✔        | ✔        |
| Google   | Blogger | auth_code         | ✔    | ✘    | ✘        | ✔        | ✔        |
| ......   |         |                   |      |      |          |          |          |

> 注：✔ 表示支持该功能，✘ 表示暂不支持该功能，🚧 表示功能开发中。

<!-- ## 功能介绍

\*\*\* [项目功能的开发安排](https://github.com/orgs/ArtisanCloud/projects/5/views/2) -->

## 产品主要维护者

Michael Hu

申请添加好友时，请备注产品名称，比如：“关注 MediaX”

<img src="https://mediax.artisan-cloud.com/assets/contact-qr-matrix-x.CMxSV8Gs.jpg" alt="请扫我" style="display:inline; width: 150px;"/>

## 许可证

MediaX SDK 项目 采用 [MIT License](./LICENSE) 开源。
