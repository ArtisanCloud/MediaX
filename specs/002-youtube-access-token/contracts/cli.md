# CLI Contract: cmd/accesstoken

## Command
```
make accesstoken ARGS='[flags]'
# or
go run ./cmd/accesstoken [flags]
```

## Flags
| Flag | Type | Required | Description |
| --- | --- | --- | --- |
| `-config` | string | No | 配置文件路径（默认 `config.yaml` 或 `MEDIA_X_CONFIG`） |
| `-action` | string | Yes | `videos.list` / `search.list` / `playlists.list` |
| `-part` | string | Yes | YouTube API `part` 参数（例如 `snippet,contentDetails`） |
| `-ids` | string | Cond. | 逗号分隔的资源 ID，适用于 `videos.list`/`playlists.list` |
| `-query` | string | Cond. | `search.list` 关键字 |
| `-channel-id` | string | Cond. | 频道 ID（`search.list`/`playlists.list`） |
| `-search-type` | string | No | `search.list` 的 `type` 参数（video/channel/playlist） |
| `-mine` | bool | No | `playlists.list` 是否拉取当前账号（与 `-ids` 互斥） |
| `-search-mine` | bool | No | `search.list` 是否限制为当前账号 |
| `-max-results` | int | No | 1-50，默认 5，越界会直接报错 |
| `-page-token` | string | No | 分页 token |
| `-region` | string | No | ISO 区域码（videos/search） |
| `-category` | string | No | `videoCategoryId` |
| `-chart` | string | No | `mostPopular` 等 |
| `-access-token` | string | Cond. | OAuth AccessToken，若未在配置/环境变量提供则必填 |
| `-access-token-ttl` | int | No | AccessToken 缓存 TTL，默认 3600 秒 |

## Input Sources
1. 配置文件 `google_youtube_config`
2. 通用环境变量：`GOOGLE_YOUTUBE_ACCESS_TOKEN` / `YOUTUBE_ACCESS_TOKEN`
3. AccessToken 专用环境变量：`ACCESSTOKEN_ACTION`、`ACCESSTOKEN_PART`、`ACCESSTOKEN_IDS`、`ACCESSTOKEN_QUERY`、`ACCESSTOKEN_CHANNEL_ID`、`ACCESSTOKEN_MINE`、`ACCESSTOKEN_SEARCH_MINE`、`ACCESSTOKEN_SEARCH_TYPE`、`ACCESSTOKEN_MAX_RESULTS`、`ACCESSTOKEN_PAGE_TOKEN`、`ACCESSTOKEN_REGION`、`ACCESSTOKEN_CATEGORY`、`ACCESSTOKEN_CHART`、`ACCESSTOKEN_ACCESS_TOKEN`、`ACCESSTOKEN_ACCESS_TOKEN_TTL`、`ACCESSTOKEN_CONFIG`
4. CLI flag（最终可覆盖所有参数）

优先级：flag > `ACCESSTOKEN_*` 环境变量 > 通用环境变量 > 配置文件。

## Validation Rules
- `-action`、`-part` 必填；`videos.list` 需 `-ids` 或 `-chart`；`playlists.list` 需 `-ids`、`-channel-id`、`-mine` 三选一且 `-mine` 与 `-ids` 互斥；`search.list` 需至少提供 `-query`、`-channel-id` 或 `-search-mine`。
- `-max-results` 强制在 1~50 范围；违反时 CLI 直接退出。
- 参数校验覆盖在 `cmd/accesstoken/main_test.go`，执行 `go test ./cmd/accesstoken` 可回归。

## Output
- 成功：格式化 JSON（与对应 schema 一致）。
- 失败：标准错误文本 + 退出码 1，包含 Google API error message，并在日志中输出 `provider=google action=<action> token_source=<flag/env/config>`。

## Error Handling
- 缺少 `action/part/access-token` → 输出 usage + error。
- Google API 错误 → 透传 `reason`/`message`，同时在日志中记录 `provider=google action=<action>`。
- 配置解析错误 → 返回 `load config <file>: ...`。

## Sample Commands（闭环摘录）
1. 查看视频：`make accesstoken ARGS='-action videos.list -part snippet -ids dQw4w9WgXcQ'`
2. 搜索频道：`make accesstoken ARGS='-action search.list -part snippet -query "MediaX" -channel-id UC_x5XG1OV2P6uZZ5FSM9Ttw'`
3. 拉取播放列表：`make accesstoken ARGS='-action playlists.list -part snippet -mine -max-results 5'`

更多 CLI/Playground 闭环示例见 `docs/develop/access-token/google/debug.md` / `develop.md`。
