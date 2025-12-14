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
| `-mine` | bool | No | `playlists.list` 是否拉取当前账号 |
| `-search-mine` | bool | No | `search.list` 是否限制为当前账号 |
| `-max-results` | int | No | 1-50，默认 5 |
| `-page-token` | string | No | 分页 token |
| `-region` | string | No | ISO 区域码（videos/search） |
| `-category` | string | No | `videoCategoryId` |
| `-chart` | string | No | `mostPopular` 等 |
| `-access-token` | string | Cond. | OAuth AccessToken，若未在配置/环境变量提供则必填 |
| `-access-token-ttl` | int | No | AccessToken 缓存 TTL，默认 3600 秒 |

## Input Sources
1. 配置文件 `google_youtube_config`
2. 环境变量 `GOOGLE_YOUTUBE_ACCESS_TOKEN`
3. CLI flag `-access-token`

优先级：flag > 环境变量 > 配置文件。

## Output
- 成功：格式化 JSON（与对应 schema 一致）。
- 失败：标准错误文本 + 退出码 1，包含 Google API error message。

## Error Handling
- 缺少 `action/part/access-token` → 输出 usage + error。
- Google API 错误 → 透传 `reason`/`message`，同时在日志中记录 `provider=google action=<action>`。
- 配置解析错误 → 返回 `load config <file>: ...`。
