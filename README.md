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

在 `config.yaml` 中通过 `session_token_providers` 配置知乎 SessionToken，支持多 Provider/App/Mode：

```yaml
session_token_providers:
  providers:
    - code: zhihu
      name: 知乎
      apps:
        - code: web
          name: 知乎 Web
          provider_code: zhihu_sessiontoken
          auth_modes:
            - key: default
              label: 默认
              zhihu_session_token_config:
                strategy: zhihu_pc_v1
                service:
                  base_url: http://127.0.0.1:7070
                  api_token: dev-session-token
                  timeout: 30
                  http_debug: false
                  api_version: v4
                authenticator:
                  entries:
                    - type: pc
                      url: https://www.zhihu.com/signin?next=%2F
                  default_user_agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.58 Safari/537.36
                  script_ids:
                    - sessiontoken.zhihu.auth.pc.v1
                  captcha_strategy: auto
                harvester:
                  watch_cookies:
                    - SESSIONID
                    - JOID
                    - osd
                    - q_c1
                    - d_c0
                    - unlock_ticket
                    - z_c0
                  watch_headers:
                    - X-XSRF-TOKEN
                  harvest_script_id: sessiontoken.zhihu.harvest.pc.v1
                callback:
                  secret: mediax-sessiontoken-callback
                  max_retry: 3
                  retry_backoff: [2, 4, 8]
                network:
                  proxy: ""
                  proxy_pool: zhihu-default
                  ip_strategy: china_rotating
                  request_timeout: 60
                api:
                  retry_backoff: [2, 4, 8]
```
> 若仓库首次拉取尚无 `config.yaml`，运行 `make sessiontoken-bootstrap`（或 `./scripts/sessiontoken-bootstrap.sh`）即可从 `config.example.yaml` 复制上述默认模板，再按环境覆盖 `callback secret`/Redis 等字段。可配合 `.env.example` 统一维护调试所需的环境变量。

### SessionToken 环境变量

本地/插件应保持以下环境变量一致，可直接复制 `.env.example`：

| 变量 | 说明 | 默认 |
| --- | --- | --- |
| `SESSIONTOKEN_LISTEN_ADDR` | SessionToken HTTP 服务监听地址 | `:7070` |
| `SESSIONTOKEN_API_TOKEN` / `POWERX_SESSION_TOKEN_API_TOKEN` | Flow 接口鉴权 | `dev-session-token` |
| `POWERX_SESSION_TOKEN_BASE_URL` | 插件访问的 SessionToken Base URL | `http://127.0.0.1:7070` |
| `POWERX_SESSION_TOKEN_CALLBACK_URL` | 插件自身的回调接口 | `https://plugin.local/api/v1/...` |
| `POWERX_SESSION_TOKEN_CALLBACK_SECRET` | 回调签名 secret，对应 `callback.secret` | `mediax-sessiontoken-callback` |
| `SESSIONTOKEN_DEBUG_CALLBACK_URL` | `/debug` 页面默认回调地址 | `http://127.0.0.1:7070/debug/callback` |
| `SESSIONTOKEN_REDIS_ADDR` | Flow 存储 Redis | `127.0.0.1:6379` |
| `SESSIONTOKEN_ZHIHU_PROXY` | 代理地址，覆盖 `network.proxy`（可填 `none`/`direct` 表示禁用代理；未设置时默认 `none`） | 空 |
| `SESSIONTOKEN_ZHIHU_API_BASE_URL` | 可覆盖默认 `https://www.zhihu.com` | 空 |
| `SESSIONTOKEN_ZHIHU_API_VERSION` | 强制指定 Zhihu API 版本（如 `v4`） | `v4` |

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

### 本地调试与 Zhihu API 验证

- 浏览器访问 `http://127.0.0.1:7070/debug` 可使用调试台：Provider→App 级联、API Token/Callback 默认填入 `dev-session-token` + `/debug/callback`，所有输入会记忆在 `localStorage`。点击“创建 Flow”后会显示 Flow ID、authorize_url，并在 Mock 区域实时展示 `/debug/callback` 的最近 20 条记录。
- **手动写回 Cookie**：授权页是普通浏览器标签页，需要手动在 DevTools Console 抓取 Cookie 并回传服务。典型流程：
  1. 在 `/debug` 创建 Flow 并复制 `authorize_url`，使用新的浏览器 profile 打开并完成知乎登录。
  2. 登录成功后在该标签页执行：
     ```js
     (() => {
       const cookie = document.cookie;
       const pick = name => (new RegExp(`${name}=([^;]+)`)).exec(cookie)?.[1] || '';
       return {
         session_token: cookie,
         cookie_sessionid: pick('SESSIONID'),
         cookie_joid: pick('JOID'),
         cookie_osd: pick('osd'),
         cookie_q_c1: pick('q_c1'),
         cookie_d_c0: pick('d_c0'),
         cookie_unlock_ticket: pick('unlock_ticket'),
         cookie_z_c0: pick('z_c0'),
         credentials_note: navigator.userAgent
       };
     })();
     ```
  3. 将上一条的 JSON 通过 `POST http://127.0.0.1:7070/debug/flows/<flow_id>/metadata` 写回（可直接在 Console 中 `fetch` 发送）。服务端收到后会立即触发 harvester 并向 `/debug/callback` 推送成功记录。
- **Playwright 自动调试器**（可选）：项目内新增 `tools/sessiontoken-browser/` CLI。首次运行时请在仓库根目录执行 `pnpm install`（或 `npm install`）以及 `npx playwright install chromium` 安装依赖/浏览器。随后即可使用 `pnpm sessiontoken:browser --flow stf_xxx [--base http://127.0.0.1:7070] [--api-token dev-session-token]`。若未显式提供 `--api-token`，会读取 `SESSIONTOKEN_API_TOKEN`/`POWERX_SESSION_TOKEN_API_TOKEN` 环境变量。脚本会：
  1. 调用 `/session-token/flows/{flow_id}` 解析 `authorize_url`。
  2. 通过 Playwright 拉起无痕 Chromium，自动打开上述 URL，提示你在弹出的窗口完成登录。
  3. 登录完成后脚本利用 CDP 读取 `document.cookie` 并调用 `POST /debug/flows/<flow_id>/metadata` 写回；终端中会显示 `/debug/callback` 返回的状态。
  4. 关闭窗口后即可在 `/debug` 页面刷新“Mock 回调”看到成功记录，再继续调用 `zhihu/v1/*` API。该 CLI 便于 QA 将整套流程脚本化，也可以在调试页里通过按钮提示用户运行。
- 如浏览器本身已处于登录状态，可勾选“复用 Cookie”选项，MediaX 会读取同租户/账号最近一次成功 Flow 的凭证并直接 `succeeded`，无需再次手动登录；若缓存缺失则自动回退到手动登录流程。
- **API 调试面板（Beta）**：`/debug` 页新增“API 调试”区块，会根据所选 Provider/版本（如 Zhihu v4）动态展示可调试的 Zhihu API（method/path/说明），并可一键复制示例 query/body。点击“从 Flow 填充 SessionToken”将调用 `GET /session-token/flows/<flow_id>` 自动填充 `X-SessionToken`，然后填写 query/path 参数并点“发送请求”即可直接调用 `/zhihu/v1/*`，响应会在页面下方实时展示。该面板默认携带 `Authorization: Bearer <API Token>`，也支持自行输入 `session_token` 以模拟失效流程。
- 终端快速调试可使用 `scripts/sessiontoken-debug.sh`：
```bash
scripts/sessiontoken-debug.sh flow stf_xxx
scripts/sessiontoken-debug.sh followings "$SESSION_TOKEN"
scripts/sessiontoken-debug.sh channels "$SESSION_TOKEN" zhihu_column_id 10 0
scripts/sessiontoken-debug.sh sanity "$SESSION_TOKEN"
```
其中 `SESSION_TOKEN` 取自 Flow metadata 或插件回调的 `metadata.session_token` 字段，脚本会自动读取 `POWERX_SESSION_TOKEN_BASE_URL` 及 API Token，并在安装 `jq` 时输出格式化 JSON。需要复用已有 Flow 时，可在本地执行 `redis-cli --raw keys 'sessionToken:flow:*'` 列出 key，再用 `redis-cli --raw GET "sessionToken:flow:<id>" | jq '.'` 查看 metadata，然后把 Flow ID 粘回 `/debug` 调试。
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
