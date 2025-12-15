# Google YouTube AccessTokenClient Debug 指南

> 目标：在本地仓库内快速验证 `pkg/client/google/youtube/accessTokenClient` 的 API 封装，覆盖 `AccessToken` 注入、示例调用与日志排查，方便定位 401/403、Quota 与参数错误。CLI/Playground 配置方法可结合 `docs/develop/access-token/google/develop.md` 阅读，两篇文章互相引用，保证指令与代码一致。

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

> AccessToken 优先级：`-access-token` flag > `GOOGLE_YOUTUBE_ACCESS_TOKEN`/`YOUTUBE_ACCESS_TOKEN` 环境变量 > `config.yaml` 中的 `oauth.access_token`。CLI 日志会打印 `token_source=<flag/env/config>`，方便确认来源。
2. 在项目根目录执行：

   ```bash
   # 注入环境变量，便于在命令之间复用
   export GOOGLE_YOUTUBE_ACCESS_TOKEN="$(pass show youtube/dev-token)"

   # 查询指定视频
   make accesstoken ARGS='-action videos.list -ids dQw4w9WgXcQ -part snippet,statistics'

   # 拉取 trending 榜单，需要 chart + region
   make accesstoken ARGS='-action videos.list -chart mostPopular -region US -part snippet -max-results 5'

   # 关键词/频道搜索
   make accesstoken ARGS='-action search.list -query "MediaX demo" -channel-id UC_x5XG1OV2P6uZZ5FSM9Ttw -part snippet -max-results 3'

   # 拉取频道播放列表
   make accesstoken ARGS='-action playlists.list -channel-id UC_x5XG1OV2P6uZZ5FSM9Ttw -part snippet -max-results 10'

   # mine 示例（与 -ids 互斥），查看当前账号的播放列表
   make accesstoken ARGS='-action playlists.list -mine -part snippet -max-results 5'
   ```

   支持的核心参数：
   - `-config`：配置文件路径（默认为 `config.yaml` 或 `MEDIA_X_CONFIG`）。
   - `-action`：`videos.list` / `search.list` / `playlists.list`。
   - `-part`：必填，例 `snippet,contentDetails`。
   - `-ids`：逗号分隔的视频/播放列表 ID。
   - `-query`、`-channel-id`、`-region`、`-page-token`、`-max-results` 等常见查询参数。
   - `-chart`（`videos.list`）、`-category`（`videoCategoryId`）、`-search-type`、`-search-mine`、`-mine`（`playlists.list`）。
   - `-access-token`、`-access-token-ttl`：覆盖 AccessToken 及其 TTL，TTL 仅影响日志与缓存 metadata。
   - 校验提示：`-max-results` 必须在 1~50 内；`videos.list` 需要 `-ids` 或 `-chart`；`playlists.list` 中 `-mine` 与 `-ids` 互斥，CLI 会给出错误提示。

   CLI 会把响应格式化为 JSON 输出，失败时直接返回 Go error，方便脚本化联调。`cmd/accesstoken` 也会读取 `google_youtube_config` 中的 `proxy_api_url/timeout/http_debug`，并尊重 `HTTPS_PROXY` 等环境变量，因此无需额外写死网络配置。所有调用都会输出结构化日志（`provider=google action=videos.list token_source=flag channel_id=...`），排查时可结合 `logs/info.log`；如需更详细的配置说明可回看 `develop.md` 的 CLI 小节。

> 与 `cmd/sessiontoken` 一样，也可以通过环境变量准备参数，再直接运行 `make accesstoken` 而无需 `ARGS`，例如：
> ```bash
> export ACCESSTOKEN_ACTION=videos.list
> export ACCESSTOKEN_PART=snippet
> export ACCESSTOKEN_IDS=dQw4w9WgXcQ
> export ACCESSTOKEN_MAX_RESULTS=5
> make accesstoken
> ```
> 环境变量优先级低于 CLI flag，高于配置文件；布尔值支持 `1/0/true/false`。

| 环境变量 | 含义 |
| --- | --- |
| `ACCESSTOKEN_CONFIG` | 覆盖 `-config`（默认仍会 fallback 到 `MEDIA_X_CONFIG`/`config.yaml`） |
| `ACCESSTOKEN_ACTION` / `ACCESSTOKEN_PART` | 对应 `-action` / `-part` |
| `ACCESSTOKEN_IDS` / `ACCESSTOKEN_QUERY` / `ACCESSTOKEN_CHANNEL_ID` | 对应 `-ids` / `-query` / `-channel-id` |
| `ACCESSTOKEN_MINE` / `ACCESSTOKEN_SEARCH_MINE` | 布尔型，控制 `-mine`、`-search-mine` |
| `ACCESSTOKEN_SEARCH_TYPE` / `ACCESSTOKEN_CHART` / `ACCESSTOKEN_CATEGORY` / `ACCESSTOKEN_REGION` / `ACCESSTOKEN_PAGE_TOKEN` | 对应同名 flag |
| `ACCESSTOKEN_MAX_RESULTS` | 覆盖 `-max-results`（1~50） |
| `ACCESSTOKEN_ACCESS_TOKEN` / `ACCESSTOKEN_ACCESS_TOKEN_TTL` | 覆盖 `-access-token` 与 `-access-token-ttl` |

其余未设置的参数依旧从 `config.yaml` 与通用环境变量（`GOOGLE_YOUTUBE_ACCESS_TOKEN` 等）中获取，这样就能像 SessionToken 调试一样，通过“export env + make accesstoken”固定一套命令。

   常见 CLI 报错：
   - `missing access token`：补充 `-access-token` 或设置环境变量。
   - `videos.list: provide at least -ids or -chart`、`playlists.list: -mine and -ids are mutually exclusive`：按提示调整参数组合。
   - `invalid_grant`/`insufficientPermissions`/`quotaExceeded`/`401`/`403`：参考下方“常见问题排查”章节，通常与 OAuth 凭证、Scope 或配额相关。

#### 4.1.1 订阅 → 视频 → 评论闭环

> 典型外部应用希望：读取订阅频道 → 拉取频道视频 → 发布自有视频 → 获取/回复评论。以下检查清单可在本地逐项验证。

| 环节（闭环） | 使用入口 | CLI/Playground 命令示例 | 说明 |
| --- | --- | --- | --- |
| 订阅 | Playground `subscriptions.List`（`yt.GetSubscriptionsClient().List`） | `go run ./main.go # 在 playground/google.go 中启用 subscriptionsListExample(ctx, "snippet,contentDetails", true)` | CLI 暂未涵盖订阅接口，可在 Playground 中按需添加 `subscriptions.List` 调用，日志同样会输出 `provider=google action=subscriptions.list`。 |
| 视频列表 | CLI `videos.list`/`search.list` | `make accesstoken ARGS='-action videos.list -part snippet -ids <videoId>,<videoIdB>'` 或 `make accesstoken ARGS='-action search.list -part snippet -query "MediaX" -channel-id <channelId>'` | CLI 负责调试视频/搜索接口，借助 `-ids/-chart/-channel-id` 快速定位素材。 |
| 发布 | Playground `video.Insert` | `go run ./main.go # 在 playground/google.go 的 videoInsertExample 中设置 Files/Form 后运行` | 需要上传文件流，因此仍建议在 Playground 或业务代码里调用 `video.Insert`，复用同一份配置与 AccessToken。 |
| 评论（读取/回复） | Playground `commentThreads.List` + `comments.Insert` | `go run ./main.go # 添加 commentFlowExample，先 List 再 Insert` | 调试评论链路时，可先使用 CLI 查询视频 ID，再在 Playground 中读取/回复，确保日志持续输出 `provider=google action=commentThreads.list/comments.insert`。 |

调试提示：

1. **订阅 & 视频列表**：先用 Playground 获取订阅列表，再使用 CLI 的 `videos.list`/`search.list` 跟踪具体视频，或切换到 `playlists.list` 验证频道首页。
2. **发布与评论**：推荐在 Playground 中添加辅助函数调用 `video.Insert`、`commentThreads.List/Insert`，这些调用与 CLI 共用 `google_youtube_config`，无需重复配置 AccessToken。
3. **验证流程**：完成“订阅→视频→发布→评论”4 步后，记录 CLI/Playground 响应与日志，作为闭环验收依据。

### 4.2 使用 `playground/google.go`

1. （可选）启动 Redis：`docker run --rm -p 6379:6379 redis:7-alpine`。若希望零依赖，可 `export PLAYGROUND_CACHE_MODE=memory`，示例将使用内存缓存。
2. 设置配置与 Token：`export MEDIA_X_CONFIG=$PWD/config.yaml`、`export GOOGLE_YOUTUBE_ACCESS_TOKEN='<token>'`。
3. 启动 Playground：`export PLAYGROUND_GOOGLE_YOUTUBE=1`，随后执行：

   ```bash
   go run ./main.go
   ```

4. `main.go` 会加载配置、初始化 `MediaX`（日志写入 `logs/info.log`/`logs/error.log`）并根据开关执行 `playground.PlayGoogleYouTube`。完成后可 `unset PLAYGROUND_GOOGLE_YOUTUBE` 关闭示例。
5. `playground/google.go` 会：
   - 读取 `config.yaml` → 初始化 `MediaX`。
   - 执行 `mediaX.CreateGoogleYouTubeACClient`。
   - 通过 `GoogleClient.TokenHandler.GetCustomToken` 注入 AccessToken（优先读取 `GOOGLE_YOUTUBE_ACCESS_TOKEN`，其余 fallback 为 `oauth.access_token`），日志会输出 `token_source=<env/config>` 且自动脱敏。
   - 调用 `video.List`，并将响应输出到控制台，同时在 `logs/info.log` 记录 `provider=google action=videos.list ids=...` 等摘要，方便与 CLI 对比。

若看到形如 `youtube#videoListResponse` 的 JSON，说明签名与网络均正常。常见自定义项包括 `PLAYGROUND_YOUTUBE_VIDEO_IDS`（覆盖默认视频 ID）与 `PLAYGROUND_YOUTUBE_REGION`（配合 `chart=mostPopular` 调整地区）。

### 4.3 自动化校验（TDD）

每次改动 CLI 后建议运行：

```bash
go test ./cmd/accesstoken
```

测试覆盖参数校验、AccessToken 优先级与错误提示。若要在同一脚本中快速验证 CLI 行为，可将该命令与一条示例调用组合，例如：

```bash
export GOOGLE_YOUTUBE_ACCESS_TOKEN="$(pass show youtube/dev-token)"
make accesstoken ARGS='-action videos.list -ids dQw4w9WgXcQ -part snippet'
```

运行结果记录在 `docs/develop/access-token/google/debug.md` 即为调试凭证的一部分，提交 PR 前请截图或粘贴输出摘要以便 reviewer 回放。

### 4.4 AccessToken 调试服务（规划中）

> 目标：像 `cmd/sessiontoken` 一样提供浏览器调试台，方便非开发同学模拟授权 + API 调用；需求详情参见 `specs/002-youtube-access-token/spec.md#user-story-4`。

**规划概要**：

1. **启动命令**：新增 `make accesstoken-serve`（或 `go run ./cmd/accesstoken/server`），默认监听 `http://127.0.0.1:7070`，可用 `-port`/`ACCESSTOKEN_LISTEN_ADDR` 指定其他端口。
2. **认证**：服务会读取 `ACCESSTOKEN_API_TOKEN`（默认 `dev-accesstoken`），以 `Authorization: Bearer <token>` 的方式保护 API，逻辑复用 SessionToken。
3. **页面入口**：`/debug/accesstoken` 使用与 `cmd/sessiontoken/debug_page.go` 相同的模板结构，提供：
   - Provider / App / API 版本选择；App 默认读取 `google_youtube_config.oauth_key`。
   - AccessToken 操作区：支持刷新 token、临时覆盖 AccessToken、查看 `token_source`。
   - API 调试区：直接在页面上填写 action、part、ids/query 等参数，底层调用 `GoogleYouTubeACClient` 并展示 JSON。
   - 回调日志区：复用 `callbackLogStore`，记录 `/debug/callback` 收到的请求（时间、FlowID、Query、Headers、Body），支持清空。
   - 快捷命令：输出 `make accesstoken ...`、`curl` 或 Playwright 命令，帮助把浏览器操作迁移到 CLI。
4. **REST 接口**：页面将调用 `POST /accesstoken/token`（刷新/注入 token）、`POST /accesstoken/call`（执行 API）、`POST /debug/callback`（模拟 OAuth 回调），返回结构体内会脱敏 token 并包含 `provider/action/token_source` 字段。
5. **共享配置**：同 CLI/Playground 一样，服务会加载 `MEDIA_X_CONFIG`/`config.yaml` 与 `GOOGLE_YOUTUBE_*` 环境变量，支持 `HTTPS_PROXY`、`SESSIONTOKEN_REDIS_*` 等代理/缓存设置。

目前服务仍在规划阶段，短期内调试仍以 CLI + Playground 为主；若需要浏览器沙盒，可对照 `cmd/sessiontoken` 目录提前评估需求或参与贡献。

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
