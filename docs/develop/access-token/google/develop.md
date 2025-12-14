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
google_youtube_config:
  api_url: "https://www.googleapis.com"
  timeout: 10
  http_debug: false
  oauth:
    access_token_url: "https://oauth2.googleapis.com/token"
    client_id: "xxx.apps.googleusercontent.com"
    client_secret: "yyy"
    scope: "https://www.googleapis.com/auth/youtube.force-ssl"
  oauth_key: "tenantA.youtube"   # 对应 token 仓库的 key
```

建议按照以下顺序配置：

1. 复制 `config.example.yaml` → `config.yaml`，保留 `${GOOGLE_YOUTUBE_CLIENT_ID}` 等占位符。
2. 在 shell 中导出对应环境变量，或使用 Secret Manager 注入 `client_id/client_secret/refresh_token/access_token`。
3. 可选：通过 `oauth_key` 区分不同租户的凭证仓库。

`GetOAuthToken` 的返回值需要至少包含 `access_token`，可选 `expires_in/refresh_token/token_type`。若需让 SDK 自动发起刷新，可以在 `OAuthConfig.refresh_token` 中填入长期有效的值；`oauth_key` 用于在日志/缓存中区分不同账号。

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
3. **日志**：在 `google_youtube_config.http_debug=true` 时，`logs/info.log` 会记录完整请求/响应。生产环境建议关闭，仅保留 `warn/error`。
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

1. 复制 `config.example.yaml` 到 `config.yaml`，确保 `google_youtube_config` 填入真实的 `client_id/client_secret/refresh_token`（或临时 `access_token`）。
2. 运行 `make accesstoken ARGS='-action videos.list -ids <video_id> -part snippet -access-token <token>'`，或设置 `GOOGLE_YOUTUBE_ACCESS_TOKEN` 环境变量后仅指定 `-action/-part`。CLI 会读取 `config.yaml` 并输出格式化 JSON，方便对比 Google 控制台中的响应。
3. 若需要更完整的上下文，可启用 `main.go` 中的 `playground.PlayGoogleYouTube`，它复用了相同的配置与 `MediaX.CreateGoogleYouTubeACClient`，更贴近正式代码路径。
4. 调试完成后，将 `config.yaml` 中的敏感信息移入 Vault/Secret Manager，仅保留必要字段或通过环境变量传入，以免误提交。

## 9. 延伸阅读

- 能力细节：`docs/plan/google/youtube_access_token_client.md`
- Token 处理：`pkg/client/google/core/accessTokenHandler.go`
- 其它平台开发指南：`docs/develop/*`

按照以上步骤即可在自有项目中稳定复用 `GoogleYouTubeACClient`，并与 MediaX 的日志、缓存和观测体系保持一致。若需新增尚未封装的 API，可参考现有子客户端结构，在 `pkg/client/google/youtube/accessTokenClient/<module>` 目录内新增 `client.go` + `schema` 文件，再补充文档。
