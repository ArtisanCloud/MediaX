# Google YouTube AccessTokenClient Debug 指南

> 目标：在本地仓库内快速验证 `pkg/client/google/youtube/accessTokenClient` 的 API 封装，覆盖 `AccessToken` 注入、示例调用与日志排查，方便定位 401/403、Quota 与参数错误。

## 1. 环境与依赖

1. Go 1.18+。
2. Redis（用作 AccessToken 缓存，可复用 SessionToken 环境的 `127.0.0.1:6379`）。
3. Google Cloud 项目已启用 **YouTube Data API v3**，且创建了 OAuth Client（推荐「桌面应用」）。
4. 拥有至少一个长期可用的 `refresh_token`（获得方式：在浏览器完成 OAuth 授权后拦截 `code` 并调用 token 接口，或使用 Google API Explorer）。
5. 仓库根目录：`/private/var/www/html/ArtisanCloud/X/MediaX/core/MediaX`。
6. `config.yaml` 里已经从 `config.example.yaml` 拷贝了 `google_youtube_config` 段落。

## 2. 准备 OAuth 凭证

1. 在 <https://console.cloud.google.com/apis/credentials> 创建 OAuth Client，Scope 至少包含：
   - `https://www.googleapis.com/auth/youtube.force-ssl`（读写绝大多数频道能力）
   - 仅读场景可换成 `youtube.readonly`。
2. 通过 OAuth 授权获取 `refresh_token`，之后可用以下命令置换 AccessToken：

   ```bash
   curl -s https://oauth2.googleapis.com/token \
     -d client_id=$GOOGLE_CLIENT_ID \
     -d client_secret=$GOOGLE_CLIENT_SECRET \
     -d refresh_token=$GOOGLE_REFRESH_TOKEN \
     -d grant_type=refresh_token | jq .
   ```

3. 为避免频繁刷新，可将 `refresh_token` 与最新 AccessToken 存入自有的 KV/Secrets，再由 `GetCustomToken` 回调读取。

## 3. 更新 `config.yaml`

1. 复制模板：`cp config.example.yaml config.yaml`（如已有配置，确保包含 `google_youtube_config` 块）。
2. 在 `google_youtube_config` 内填充 OAuth 凭证；模板支持 `${GOOGLE_YOUTUBE_CLIENT_ID}` 等环境变量占位符，可先在 shell 中 `export` 对应值以防止明文写入仓库。
3. 如需代理，将 `proxy_api_url` 指向公司内网代理（CLI 会继承该设置）；若使用 Redis 缓存 AccessToken，请确认 `SESSIONTOKEN_REDIS_ADDR` 或自定义地址已经启动。

```yaml
google_youtube_config:
  api_url: "https://www.googleapis.com"
  proxy_api_url: ""                 # 如需科学上网，可指向公司内网代理
  timeout: 15
  http_debug: true                  # 打开后 BaseClient 会打印所有 HTTP 请求/响应
  oauth:
    access_token_url: "https://oauth2.googleapis.com/token"
    client_id: "YOUR_CLIENT_ID.apps.googleusercontent.com"
    client_secret: "YOUR_CLIENT_SECRET"
    refresh_token: "1//0gExampleRefreshToken"
    scope: "https://www.googleapis.com/auth/youtube.force-ssl"
    grant_type: "refresh_token"
  oauth_key: "youtube.default"      # 可选，用于在回调中区分多租户
```

> **提示**：`config.example.yaml` 中已包含完整模板，若后续需要不同环境的配置，可用 `MEDIA_X_CONFIG=/path/to/config` 覆盖 CLI 读取的路径。`ClientConfig` 会自动缓存 AccessToken，缓存 key 形如 `mediax.access_token.<md5>`。如需强制刷新，可手动删除 Redis 中的对应 key。

### 3.1 AccessToken 缺失/过期排查

- **缺失**：CLI 会提示 `missing access token`，请确认是否设置 `GOOGLE_YOUTUBE_ACCESS_TOKEN` 或在 `config.yaml` 中写入占位符并导出环境变量。
- **过期**：若出现 `invalid_grant`/`401`，执行 `redis-cli --raw keys 'mediax.access_token.*' | xargs -I{} redis-cli DEL {}` 清除缓存，再刷新 token。
- **代理导致连接失败**：检查 `proxy_api_url` 是否可访问 Google API，必要时临时禁用。

## 4. 启动本地调试脚本

### 4.1 使用 `cmd/accesstoken` CLI（推荐）

1. 准备 AccessToken：可通过 `GOOGLE_YOUTUBE_ACCESS_TOKEN` 环境变量或 `--access-token` 参数提供，CLI 会把它注入到 `GetOAuthToken`。如需长期使用 `refresh_token`，可在 `config.yaml` 中维护 `oauth.refresh_token` 并让 CLI 返回预估 TTL。
2. 在项目根目录执行：

   ```bash
   # 示例：查询视频详情
   make accesstoken ARGS='-action videos.list -ids dQw4w9WgXcQ -part snippet,statistics'

   # 搜索
   make accesstoken ARGS='-action search.list -query "MediaX demo" -max-results 3'

   # 拉取频道播放列表
   make accesstoken ARGS='-action playlists.list -channel-id UC_x5XG1OV2P6uZZ5FSM9Ttw -part snippet -max-results 10'
   ```

   支持的核心参数：
   - `-config`：配置文件路径（默认为 `config.yaml` 或 `MEDIA_X_CONFIG`）。
   - `-action`：`videos.list` / `search.list` / `playlists.list`。
   - `-part`：必填，例 `snippet,contentDetails`。
   - `-ids`：逗号分隔的视频/播放列表 ID。
   - `-query`、`-channel-id`、`-region`、`-page-token`、`-max-results` 等常见查询参数。
   - `-access-token`、`-access-token-ttl`：覆盖 AccessToken 及其 TTL。

   CLI 会把响应格式化为 JSON 输出，失败时直接返回 Go error，方便脚本化联调。`cmd/accesstoken` 也会读取 `google_youtube_config` 中的 `proxy_api_url/timeout/http_debug`，因此无需额外写死网络配置。

#### 4.1.1 订阅 → 视频 → 评论闭环

> 典型外部应用希望：读取订阅频道 → 拉取频道视频 → 发布自有视频 → 获取/回复评论。以下检查清单可在本地逐项验证。

| 场景 | 本地入口 | 最小参数 | 说明 |
| --- | --- | --- | --- |
| 获取订阅频道列表 | `pkg/client/google/youtube/accessTokenClient/subscriptions`（`yt.GetSubscriptionsClient().List`） | `part=snippet,contentDetails` + `mine=true` | 建议在 Playground 或自定义脚本中调用，日志记录 `provider=google api=subscriptions.list`。 |
| 获取订阅频道最新视频 | CLI `videos.list`（`-action videos.list -part snippet -channel-id <channelId>`）或 `pkg/client/.../playlistItems` | `channelId` / `playlistId` | 对应 “获取每个频道的 videos 列表” 用例。 |
| 发布自有视频 | `pkg/client/google/youtube/accessTokenClient/video`. `Insert` | `part=snippet,status` + 上传媒资 | 通过 Playground 或自定义工具执行，可复用 CLI 的 AccessToken 与代理配置。 |
| 获取自己发布视频的评论 | `pkg/client/google/youtube/accessTokenClient/commentThreads`. `List` | `videoId` + `part=snippet,replies` | 建议在 Playground 中演练，输出 JSON 并记录日志。 |
| 回复视频评论 | `pkg/client/google/youtube/accessTokenClient/comments`. `Insert` 或 `commentThreads.Insert` | `part=snippet` + `parentId`/`videoId` | 验证完 `List` 后即可调用，确保 CLI/Playground 使用相同 AccessToken。 |

调试提示：

1. **订阅 & 视频列表**：可通过 CLI（videos.list）或在 Playground 中调用 `subscriptions.List` + `playlistItems.List` 组合来完成链路。
2. **发布与评论**：推荐在 Playground 中添加辅助函数，调用 `video.Insert`、`commentThreads.List/Insert` 等 API；这些调用与 CLI 共用 `google_youtube_config`，无需额外配置。
3. **验证流程**：完成上述五个动作后，记录 CLI/Playground 响应与日志，作为“订阅→视频→评论”闭环的验收依据。

### 4.2 使用 `playground/google.go`（保持与旧教程兼容）

1. 启动 Redis：`docker run --rm -p 6379:6379 redis:7-alpine`（若已经在用 SessionToken，可跳过）。
2. 打开 `main.go`，取消注释 `playground.PlayGoogleYouTube(localConfig, mediaX)`。
3. 运行：

   ```bash
   go run ./main.go
   ```

4. `playground/google.go` 会：
   - 读取 `config.yaml` → 初始化 `MediaX`。
   - 执行 `mediaX.CreateGoogleYouTubeACClient`。
   - 通过 `GoogleClient.TokenHandler.GetCustomToken` 注入 AccessToken（默认示例直接返回写死的 Token；请替换为从 Secret/Redis 读取的逻辑，或返回刚在第 2 步获取的 AccessToken）。
   - 调用 `video.List`，并将响应输出到控制台，同时在 `logs/info.log` 记录完整 HTTP 请求。

若看到形如 `youtube#videoListResponse` 的 JSON，说明签名与网络均正常。

## 5. 常用调试场景

### 5.1 验证不同模块

可直接修改 `playground/google.go` 中的示例调用，例如：

```go
playlist := googleYouTubeClient.GetPlaylistsClient()
res, err := playlist.List(ctx, &playlistSchema.YouTubePlaylistsListReq{
    Part:     "snippet,status",
    Mine:     true,
    MaxResults: 10,
})
```

保存后再次 `go run ./main.go` 即可。`GoogleYouTubeACClient` 暴露的 `GetXXXClient` 与 `docs/plan/google/youtube_access_token_client.md` 中的能力矩阵一致，常见用途：

- `GetSearchClient().List`：验证关键词/频道搜索。
- `GetThumbnailsClient().Set`：测试封面上传，需提供 `Files` 或 `Form`。
- `GetPlaylistsClient().Insert/Update/Delete`：调试播放列表写入。

### 5.2 切换 AccessToken 来源

如果要动态刷新 Token，可将 `GetCustomToken` 修改为：

```go
googleYouTubeClient.GoogleClient.TokenHandler.GetCustomToken = func(key string, refresh bool) object.HashMap {
    token, err := myVault.FetchAccessToken(key, refresh) // 从 Vault/Redis 读取
    if err != nil {
        panic(err)
    }
    return token
}
```

`key` 默认等于 `oauth_key`。当返回值包含 `expires_in` 时，MediaX 会自动写入缓存并按剩余时长刷新。

### 5.3 清理缓存/查看日志

```bash
# 清理 AccessToken 缓存
redis-cli --raw keys 'mediax.access_token.*' | xargs -I{} redis-cli DEL {}

# 持续查看 HTTP 日志
tail -f logs/info.log | rg youtube
```

## 6. 常见问题排查

| 现象 | 日志/错误 | 排查建议 |
| --- | --- | --- |
| `invalid_grant` | 获取 AccessToken 时报 400 | 检查 refresh_token 是否被吊销，若使用桌面应用需重新授权。 |
| `insufficientPermissions` | YouTube API 返回 403 | Scope 不包含目标接口所需权限，或频道未开通会员等功能。 |
| `quotaExceeded` | 403 + `rateLimitExceeded`/`quotaExceeded` | 查看 Google Cloud→Quota，调低 `part` 字段、避免 `videos.insert` 等高配额操作在调试阶段频繁触发。 |
| `403 forbidden` 且 `reason=channelSuspended` | 频道/视频被限制 | 与 AccessToken 无关，需要业务方处理。 |
| `401 Unauthorized` | BaseClient 自动重试后仍失败 | 多半是 AccessToken 过期且 `GetCustomToken` 返回的仍是旧值；删除 Redis 缓存并重新置换 token。 |

## 7. 建议的调试流程回顾

1. 准备 OAuth Client + refresh token。
2. 在 `config.yaml` 配置 `google_youtube_config`。
3. 启动 Redis，运行 `go run ./main.go`。
4. 根据需要在 `playground/google.go` 调整要调用的子客户端。
5. 通过日志观察请求体、响应体；必要时抓包或查看 Google API Dashboard。
6. 调试完成后记得注释 `PlayGoogleYouTube`，避免在默认 `go run` 中误打 YouTube API。
