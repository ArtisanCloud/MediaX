# Google YouTube AccessTokenClient 开发指南

> 面向需要在外部应用中集成 YouTube 官方账号能力的研发，介绍如何基于 MediaX 的 `GoogleYouTubeACClient` 完成视频上传、播放列表管理、频道互动等操作。

## 1. 适用场景

- 统一管理多个频道的 OAuth 授权、AccessToken 刷新与缓存。
- 构建视频/直播/订阅等运营工具，复用 SDK 已封装的 20+ 子客户端。
- 希望通过 MediaX 与其他平台（抖音、知乎、小红书）共享日志、限流与配置体系。

完整能力矩阵见 `docs/plan/google/youtube_access_token_client.md`，涵盖 `video`、`playlists`、`comments`、`members`、`thumbnails`、`watermarks` 等模块。

## 2. 前置条件

1. 已完成 MediaX 的基础初始化（`client.MediaX` + 全局 `logger`/`cache`）。
2. Google Cloud 项目启用 **YouTube Data API v3**，持有 OAuth Client（`client_id/client_secret`）。
3. 每个频道至少完成一次 OAuth 授权并存储 `refresh_token` 或短期 `access_token`。
4. 业务方负责持久化 AccessToken/RefreshToken（可放在自有 DB 或 Secrets Manager），MediaX 负责在调用前注入。

## 3. 快速开始

```go
package youtube

import (
    "context"

    mediaxClient "github.com/ArtisanCloud/MediaX/pkg/client"
    "github.com/ArtisanCloud/MediaX/pkg/client/config"
    videoSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/video/schema"
    "github.com/ArtisanCloud/MediaXCore/pkg/cache"
    "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
    "github.com/ArtisanCloud/MediaXCore/utils/object"
)

func ListVideos(ctx context.Context, accountKey string, mxCfg *config.MediaXConfig, ytCfg *config.GoogleYouTubeConfig, c cache.ICache) (*videoSchema.YouTubeVideoListRes, error) {
    mediaX := mediaxClient.NewMediaX(mxCfg, c)

    // 自定义 token 提供者，可从数据库、Vault、外部服务读取
    ytCfg.GetOAuthToken = func(key string, refresh bool) object.HashMap {
        token, err := LoadGoogleToken(accountKey, refresh) // 业务自定义
        if err != nil {
            panic(err)
        }
        return token
    }

    client, err := mediaX.CreateGoogleYouTubeACClient(ytCfg)
    if err != nil {
        return nil, err
    }

    return client.GetVideoClient().List(ctx, &videoSchema.YouTubeVideoListReq{
        Part:       "snippet,contentDetails,statistics",
        Id:         []string{"abc123"},
        RegionCode: "US",
    })
}
```

### 3.1 配置要点

```yaml
access_token_providers:
  providers:
    - code: "google"
      name: "Google"
      apps:
        - code: "youtube"
          name: "Google YouTube"
          provider_code: "google_youtube"
          api_version: "v3"
          auth_modes:
            - key: "default"
              label: "默认 OAuth"
              google_youtube_config:
                api_url: "https://www.googleapis.com"
                timeout: 10
                http_debug: false
                oauth:
                  access_token_url: "https://oauth2.googleapis.com/token"
                  client_id: "xxx.apps.googleusercontent.com"
                  client_secret: "yyy"
                  redirect_url: "http://localhost:7071/debug/callback"
                  scope: "https://www.googleapis.com/auth/youtube.force-ssl"
                oauth_key: "tenantA.youtube"   # 对应 token 仓库的 key
        - code: "blogger"
          name: "Google Blogger"
          provider_code: "google_blogger"
          api_version: "v3"
          auth_modes:
            - key: "default"
              label: "默认 OAuth"
              google_blogger_config:
                api_url: "https://www.googleapis.com"
                timeout: 15
                http_debug: false
                oauth:
                  access_token_url: "https://oauth2.googleapis.com/token"
                  client_id: "${GOOGLE_BLOGGER_CLIENT_ID}"
                  client_secret: "${GOOGLE_BLOGGER_CLIENT_SECRET}"
                  redirect_url: "http://localhost:7071/debug/callback"
                  scope: "https://www.googleapis.com/auth/blogger"
                oauth_key: "tenantA.blogger"
```

> **授权回调**：若希望由本地调试服务接收 Google OAuth 回调，请在 Google Console 中允许 `http://localhost:7071/debug/callback`，并在对应模式（如 `google_youtube_config.oauth.redirect_url`）中填写同一地址。这样在浏览器授权后即可在调试页的“OAuth 回调日志”看到 `code/state`，不需要额外搭建回调服务器。

以上结构强调 **Provider → App → 授权模式**：同一 Provider（如 Google）可以同时包含 YouTube/Blogger，多套授权模式通过 `auth_modes` 区分不同 `oauth_key` 或租户。CLI、Web 调试台与 Playground 都会读取同一个层级，自动生成 `provider_code/provider_app/provider_auth_mode`，并把配置注入客户端。

建议按照以下顺序配置：

1. 复制 `config.example.yaml` → `config.yaml`，保留 `access_token_providers` 里各 Provider 的占位符（例如 `${GOOGLE_YOUTUBE_CLIENT_ID}`）。
2. 针对不同 App/授权模式导入对应的环境变量，或使用 Secret Manager 注入 `client_id/client_secret/refresh_token/access_token`。
3. 通过 `oauth_key` 区分不同租户的凭证仓库，后续 CLI/调试台的“授权模式”下拉会直接展示该 key。

`GetOAuthToken` 的返回值需要至少包含 `access_token`，可选 `expires_in/refresh_token/token_type`。若需让 SDK 自动发起刷新，可以在 `OAuthConfig.refresh_token` 中填入长期有效的值；`oauth_key` 用于在日志/缓存中区分不同账号。

### 3.2 CLI 调试与样例输出

`cmd/accesstoken` 为配置章节提供了“所见即所得”的校验方式：

1. **注入凭证**：可选择 `-access-token`、`GOOGLE_YOUTUBE_ACCESS_TOKEN`/`YOUTUBE_ACCESS_TOKEN` 环境变量、或 `config.yaml` 中的 `oauth.access_token`。优先级为 flag > env > config，CLI 日志会打印 `token_source`。
2. **运行命令**：`make accesstoken ARGS='-action videos.list -part snippet -ids <videoId>'`。命令会读取 `access_token_providers.google.youtube` 的默认授权模式（含代理、超时、`http_debug`）并在 stdout 输出 JSON。
3. **查看结果**：标准输出返回结构化响应，同时 `logs/info.log` 将记录 `provider=google action=<action> channel_id=...`，便于对照后台日志。

```bash
export MEDIA_X_CONFIG=$PWD/config.yaml
export GOOGLE_YOUTUBE_ACCESS_TOKEN="$(pass show youtube/dev-token)"

make accesstoken ARGS='-action search.list -part snippet -query "MediaX" -max-results 3'
make accesstoken ARGS='-action playlists.list -part snippet -mine -max-results 5'

# 或参照 SessionToken 的方式，先设置环境变量再直接 make：
export ACCESSTOKEN_ACTION=videos.list
export ACCESSTOKEN_PART=snippet
export ACCESSTOKEN_IDS=dQw4w9WgXcQ
make accesstoken
```

示例输出（节选）：

```json
{
  "kind": "youtube#videoListResponse",
  "pageInfo": {
    "resultsPerPage": 3
  },
  "items": [
    {
      "id": "dQw4w9WgXcQ",
      "snippet": {
        "title": "MediaX Demo Video",
        "channelTitle": "ArtisanCloud"
      }
    }
  ]
}
```

约束提示：

- `-action`、`-part` 必填，`-max-results` 仅允许 1~50；`videos.list` 需 `-ids` 或 `-chart`；`playlists.list` 中 `-mine` 与 `-ids` 互斥。
- CLI 支持 `-config` 覆盖配置路径，并自动尊重 `MEDIA_X_CONFIG`、`HTTPS_PROXY` 等环境变量。
- 无论 CLI 还是 Playground 都会复用 `MediaXCore` 日志器，输出脱敏后的 `token_source`、`provider` 与参数摘要，可直接粘贴到工单里定位问题。
- 遇到 `invalid_grant`/`quotaExceeded` 等错误时，请跳转到 `docs/develop/access-token/google/debug.md` 的 “常见问题排查” 小节获取具体指引。

### 3.3 Playground 调试步骤

`main.go` 已内置开关，可在不修改源码的情况下运行 Google Playground：

1. **准备配置/Token**：与 CLI 相同，确保 `MEDIA_X_CONFIG`（默认为仓库根目录的 `config.yaml`）可读取 `access_token_providers.google.youtube` 的授权模式，并通过 `GOOGLE_YOUTUBE_ACCESS_TOKEN` 或 `oauth.access_token` 提供可用 AccessToken。
2. **选择缓存**：默认使用 Redis（`PLAYGROUND_REDIS_ADDR=127.0.0.1:6379`），若希望快速验证可 `export PLAYGROUND_CACHE_MODE=memory` 切到内存缓存；日志会告知 cache 选择。
3. **启用 Playground**：`export PLAYGROUND_GOOGLE_YOUTUBE=1`，然后运行 `go run ./main.go`。关闭时移除该环境变量即可。
4. **日志落地**：`buildLoggerConfig` 会把日志写入 `logs/info.log`/`logs/error.log`，并默认开启控制台输出，方便观察 `provider=google action=videos.list token_source=env` 等字段；`playground/google.go` 内部会脱敏 AccessToken。
5. **Token 回调示例**：示例代码会自动复写 `GetOAuthToken`，注入一个来自环境变量/配置的 AccessToken。可根据业务需要替换成自定义读取逻辑，例如：

```go
googleCfg.GetOAuthToken = func(key string, refresh bool) object.HashMap {
    token, ttl := myVault.FetchToken(key) // 自行实现
    mediaX.Logger.InfoF("playground: inject token source=vault key=%s refresh=%t token=%s", key, refresh, maskToken(token))
    return object.HashMap{
        "access_token": token,
        "expires_in":   ttl.Seconds(),
    }
}
```

Playground 默认调用 `videos.list`（可通过 `PLAYGROUND_YOUTUBE_VIDEO_IDS/PLAYGROUND_YOUTUBE_REGION` 等环境变量覆盖），因此非常适合和 CLI 互为对照验证。若需要覆盖订阅/评论等模块，可在 `playground/google.go` 中复制 `PlayGoogleYouTube` 的结构并替换请求体。

> 快速复盘：先在 Quickstart（`specs/002-youtube-access-token/quickstart.md`）完成配置 → 依次执行 CLI 与 Playground → 若出现异常参考本页与 `debug.md`。这样文档与代码形成闭环，减少重复排查成本。

### 3.4 AccessToken 调试服务（Web 沙盒）

> 与 SessionToken 调试台一致：本地运行一个 HTTP 服务，浏览器即能完成 AccessToken 注入、API 调用与 OAuth 回调观测。底层仍然走 `MediaX.CreateGoogleYouTubeACClient`，因此 CLI/Playground/调试页三者共享配置与日志。

1. **启动方式**
   ```bash
   # 可通过 ARGS 透传 -config / -port
   make accesstoken-serve

   # 或自行指定
   go run ./cmd/accesstoken/server -config config.yaml -port 8080
   ```
   - 默认监听 `http://127.0.0.1:7071`，可使用 `-port` 或 `ACCESSTOKEN_LISTEN_ADDR=:8080` 覆盖。
   - API Token 从 `ACCESSTOKEN_API_TOKEN` 读取（默认 `dev-accesstoken`），网页脚本会自动携带 `Authorization: Bearer ...` 调用 JSON API。
   - 缓存：默认内存，若设置 `ACCESSTOKEN_REDIS_ADDR=127.0.0.1:6379` 会自动切换到 Redis（并支持 `ACCESSTOKEN_REDIS_DB/USERNAME/PASSWORD`）。

2. **页面入口**
   - 访问 `http://127.0.0.1:7071/debug`（或 `/debug/accesstoken`），顶部的“基础配置”区会列出当前 `config.yaml` 中启用的所有 AccessToken Provider（Google YouTube / Google Blogger / 字节抖音 / 小红书聚光 / 哔哩哔哩等），并展示 `config_path`、`oauth_key`、默认回调地址；API 版本输入框会随 Provider 自动填入（YouTube=`v3`、Blogger=`v3` 等），便于未来做版本切换。
   - **AccessToken 解析**：JSON 模板随 Provider 自动更新，点击“解析 AccessToken”即可调用 `POST /accesstoken/token`。后端会按“请求体 → Provider 对应的环境变量 → `oauth.access_token`”顺序解析，并返回 token 来源/脱敏值/TTL/http_debug。
   - **API 调用**：JSON 同样包含 `provider_code/provider_app/config_path`。目前仅对 `google_youtube` Provider 开放 `POST /accesstoken/call`（复用 CLI 的 `videos.list/search.list/playlists.list`），其余 Provider 会提示“暂未开放 API 调试”，为后续扩展预留空间。
   - **OAuth 回调日志**：所有命中 `/debug/callback` 的请求（默认 `http://127.0.0.1:7071/debug/callback`）都会记录最近 50 条，可随时刷新/清空，用来排查 `redirect_uri`。

3. **RESTful API**
   | Endpoint | Method | 说明 |
   | --- | --- | --- |
| `/accesstoken/token` | POST | `{"provider_code":"google_youtube","config_path":"config.yaml","access_token":""}` → 根据 Provider 自动回退到环境变量/配置文件，返回 token 来源、脱敏值、TTL 和 http_debug。 |
| `/accesstoken/call` | POST | `{"provider_code":"google_youtube","action":"videos.list","part":"snippet","ids":"dQw4w9WgXcQ"}` → 等价于 CLI 调用，只对 `google_youtube` Provider 生效。 |
   | `/debug/callback` | ALL | 提供给 OAuth/第三方回调的演练地址，记录 Query、Headers、Body，开放访问无需 API token。 |
   | `/api/callbacks` | GET | 需要携带 API token，返回最近的回调记录；`/api/callbacks/clear` 用于清空。 |

4. **环境变量一览**
   | 变量 | 作用 |
   | --- | --- |
   | `ACCESSTOKEN_API_TOKEN` | 保护 `/accesstoken/*` 与 `/api/*` 路由，默认 `dev-accesstoken`。 |
| `ACCESSTOKEN_LISTEN_ADDR` | 监听地址（默认 `:7071`）。 |
   | `ACCESSTOKEN_LOG_LEVEL` | 日志等级，写入 `logs/accesstoken-server-*.log`。 |
   | `ACCESSTOKEN_REDIS_ADDR/DB/USERNAME/PASSWORD` | 如需让调试服务也走 Redis，可在此配置。 |
   | `ACCESSTOKEN_CONFIG` | 针对服务端覆盖配置路径，优先级高于 `MEDIA_X_CONFIG`。 |

5. **常见用法**
   - 复制 CLI payload → 浏览器 JSON → 一键验证（可截图提交给 QA/产品）。
   - Provider/Provider App 下拉由配置自动生成：新增 `google_blogger_config`、`byte_dance_douyin_config`、`redbook_juguang_config`、`bilbili_config` 等段落后即可出现在页面内，并复用同一份 AccessToken/代理设置。
   - 在回调页面观察第三方返回的 `code`/`state`，同步给后端排查，必要时可直接 `curl -XPOST` 命中 `/debug/callback` 注入自定义 payload。
   - 使用 `curl` 直接调用 API：`curl -H "Authorization: Bearer $ACCESSTOKEN_API_TOKEN" -H "Content-Type: application/json" -d '{"provider_code":"google_youtube","action":"videos.list","part":"snippet","ids":"dQw4w9WgXcQ"}' http://127.0.0.1:7071/accesstoken/call`.

至此，AccessToken 的“CLI ↔ Playground ↔ Web 调试台”形成闭环：CLI 负责脚本化，Playground 贴近 SDK，调试服务则面向非工程同学，三者共享 `config.yaml`、环境变量与日志。

## 4. Token 管理策略

MediaX 的 `GoogleAccessTokenHandler` 默认会：

1. 先调用 `TokenHandler.GetCustomToken`（即上面自定义的函数）；如果返回值包含 `expires_in`，SDK 会写入缓存。
2. 若未提供自定义逻辑，则会根据 `oauth` 段配置，向 `access_token_url` 发送请求。
3. 缓存 key 为 `mediax.access_token.<md5(credentials)>`，缓存时间等于 `expires_in-60s`。

推荐做法：

- 对每个频道维护 `oauth_key` → `refresh_token` 的映射。
- `GetOAuthToken` 中，根据 `key` 加载最新 AccessToken；若 `refresh==true` 或缓存失效则主动刷新并写回。
- 对刷新失败添加报警，避免频繁触发 `invalid_grant`。

## 5. 常用子客户端

| 模块 | Client 入口 | 典型方法 | 说明 |
| --- | --- | --- | --- |
| 视频 | `GetVideoClient()` | `List/Insert/Update/Delete/Rate/ReportAbuse` | 上传与元数据管理。 |
| 播放列表 | `GetPlaylistsClient()`/`GetPlaylistItemsClient()` | `List/Insert/Update/Delete` | 维护频道首页内容。 |
| 搜索 | `GetSearchClient()` | `List` | 支持关键词、地理位置、直播筛选。 |
| 评论 | `GetCommentThreadsClient()`/`GetCommentsClient()` | `List/Insert/Update/SetModerationStatus/Delete` | 构建评论运营工具。 |
| 频道装饰 | `GetChannelSectionsClient()`、`GetChannelBannersClient()`、`GetThumbnailsClient()`、`GetWatermarksClient()` | - | 用于横幅/水印/封面上传。 |
| 会员 | `GetMembersClient()`、`GetMembershipsLevelsClient()` | `List` | 读取会员列表和等级。 |
| 国际化 | `GetI18nLanguagesClient()`、`GetI18nRegionsClient()` | `List` | 构建语言/地区下拉列表。 |

详细字段请查阅对应的 `schema/*.go`。

## 6. 错误与重试

- `kernel.BaseClient` 返回的 `error` 可能是：
  - Go 层网络错误（`context deadline exceeded`、`dial tcp` 等）。
  - Google API 的 `ResponseError`。可通过 `errors.As(err, *response.BasicError)` 获取 `Code`、`Message`、`Errors[]`。
- 默认不做自动重试；业务方可根据 Google 提供的 `reason`（如 `rateLimitExceeded`、`backendError`）自行退避。
- 上传类接口（`videos.insert`/`thumbnails.set`）属于 `multipart` 请求，失败时可能仍消耗配额，调试需谨慎。

## 7. 最佳实践

1. **使用 Context**：调用时传入含超时的 `ctx`，避免接口阻塞。`ctx, cancel := context.WithTimeout(...); defer cancel()`。
2. **精简 `part` 字段**：YouTube API 以 `part` 计费，尽量只请求必要字段。
3. **日志**：在对应授权模式设置 `google_youtube_config.http_debug=true` 时，`logs/info.log` 会记录完整请求/响应。生产环境建议关闭，仅保留 `warn/error`。
4. **幂等与配额**：写接口请记录 `request_id` 或使用自定义冪等键，避免重复提交导致配额浪费。
5. **文件上传**：`thumbnails.Set`、`watermarks.Set` 需要准备 `Files`（`map[string]string{"file": "/path/to/file"}`）。可配合 `os.Open` + `io.Reader`。
6. **多租户隔离**：将 `oauth_key`、`tenant_uuid` 写入自定义上下文，方便日志筛选与配额统计。

## 8. 结合外部服务的流程示例

1. 用户在后台发起“更新视频封面”的请求 → 业务服务从数据库读取 `channel_id`、`oauth_key`。
2. 调用上文的 `ListVideos` 或直接 `thumbnails.Set`，传入 `context.Context` 与文件内容。
3. 处理返回值：
   - 成功：写操作审计表，记录 `youtube#thumbnailSetResponse`。
   - 失败：根据 `error.Code` 判断是否需要重试或提示用户重新授权。
4. 监控：在服务层记录自定义 metric（`provider=google resource=thumbnails action=set`），并将 Google 返回的 `X-YouTube-Quota-Used` 写入日志，方便跟踪配额消耗。

## 9. 本地验证与演练

1. 复制 `config.example.yaml` 到 `config.yaml`，确保 `access_token_providers.google.youtube` 下的 `google_youtube_config` 填入真实的 `client_id/client_secret/refresh_token`（或临时 `access_token`）。
2. 运行 `make accesstoken ARGS='-action videos.list -ids <video_id> -part snippet -access-token <token>'`，或设置 `GOOGLE_YOUTUBE_ACCESS_TOKEN` 环境变量后仅指定 `-action/-part`。CLI 会读取 `config.yaml` 并输出格式化 JSON，方便对比 Google 控制台中的响应。
3. 若需要更完整的上下文，可启用 `main.go` 中的 `playground.PlayGoogleYouTube`，它复用了相同的配置与 `MediaX.CreateGoogleYouTubeACClient`，更贴近正式代码路径。
4. 调试完成后，将 `config.yaml` 中的敏感信息移入 Vault/Secret Manager，仅保留必要字段或通过环境变量传入，以免误提交。

## 9. 延伸阅读

- 能力细节：`docs/plan/google/youtube_access_token_client.md`
- Token 处理：`pkg/client/google/core/accessTokenHandler.go`
- 其它平台开发指南：`docs/develop/*`

## 10. CLI/Playground 闭环演练

> 对新人来说，按照“订阅 → 视频 → 发布 → 评论”顺序走通一次，可以验证配置、权限与 AccessToken 是否全部就绪。

| 阶段 | 工具入口 | 命令/示例 | 目标 |
| --- | --- | --- | --- |
| 订阅 | Playground `subscriptions.List`（`yt.GetSubscriptionsClient().List`） | `go run ./main.go` 并在 `playground/google.go` 中调用 `subscriptionsListExample(ctx, "snippet,contentDetails", true)` | 验证 AccessToken 能访问订阅资源，并记录 `provider=google action=subscriptions.list`。 |
| 视频/搜索 | CLI `videos.list` / `search.list` | `make accesstoken ARGS='-action videos.list -part snippet -ids <videoId>'` 或 `make accesstoken ARGS='-action search.list -part snippet -query "MediaX demo" -channel-id <channelId>'` | 读取具体视频或频道数据，确认 CLI 配置、代理与日志都正常。 |
| 发布 | Playground `video.Insert` 示例 | `go run ./main.go`，在 Playground 中启用 `videoInsertExample` 并提供文件流 | 真实演练上传流程（CLI 当前以读操作为主），并观察 `logs/info.log`。 |
| 评论 | Playground `commentThreads.List` + `comments.Insert` | `go run ./main.go`，执行 `commentFlowExample`（先 List 再 Insert） | 完成闭环最后一步：读取并回复评论，确保权限 Scope 足够。 |

完成表格中的 4 步，即可把 CLI 与 Playground 两条链路串联起来，形成可以复用的验收脚本。

按照以上步骤即可在自有项目中稳定复用 `GoogleYouTubeACClient`，并与 MediaX 的日志、缓存和观测体系保持一致。若需新增尚未封装的 API，可参考现有子客户端结构，在 `pkg/client/google/youtube/accessTokenClient/<module>` 目录内新增 `client.go` + `schema` 文件，再补充文档。
