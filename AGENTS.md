# MediaX Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-12-09

## Active Technologies
- Go 1.18（仓库 `go.mod` 已指定） + `github.com/ArtisanCloud/MediaXCore`（BaseClient/Logger）、`github.com/redis/go-redis/v9`（可选 Token 缓存）、`gopkg.in/yaml.v3`（配置解析） (002-youtube-access-token)
- Redis（用于 AccessToken cache，可选；本特性主要读取本地 YAML + 环境变量） (002-youtube-access-token)
- Go 1.18（仓库 go.mod） + `github.com/ArtisanCloud/MediaXCore`（BaseClient/日志/配置绑定）、`github.com/redis/go-redis/v9`（可选 Flow 缓存）、`gopkg.in/yaml.v3`（解析 access_token_providers） (004-bilibili-access-token)
- 进程内内存 map + 可选 Redis (key 前缀 `accesstoken:oauth:*`, TTL=24h) (004-bilibili-access-token)
- Go 1.18（`go.mod`） + `github.com/ArtisanCloud/MediaXCore`（BaseClient/Logger/Cache）、`github.com/redis/go-redis/v9`（可选 Flow 缓存）、`pkg/client/byteDance/douYin/accessTokenClient`（DouYin SDK）、`internal/accesstoken/handler/mask` (006-bytedance-access-token)
- 进程内 map + 可选 Redis（key `accesstoken:oauth:byte_dance_douyin:default:<mode>`，flow 索引 `accesstoken:oauth:flow:<flow_id>`） (006-bytedance-access-token)
- Go 1.18（仓库 go.mod） + `github.com/ArtisanCloud/MediaXCore`、`github.com/redis/go-redis/v9`、既有 `pkg/client/byteDance/douYin/clientTokenClient` (007-douyin-client-token)
- Redis（主缓存，key `clientToken:douyin:<client_key>`）+ 进程内内存降级 (007-douyin-client-token)

- Go 1.18（仓库 `go.mod`） + `github.com/ArtisanCloud/MediaXCore`（日志/HTTP/BaseClient）、`github.com/redis/go-redis/v9`（Flow 持久化）、`gopkg.in/yaml.v3`（配置解析） (001-session-token-client)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.18（仓库 `go.mod`）

## Code Style

Go 1.18（仓库 `go.mod`）: Follow standard conventions

## Recent Changes
- 007-douyin-client-token: Added Go 1.18（仓库 go.mod） + `github.com/ArtisanCloud/MediaXCore`、`github.com/redis/go-redis/v9`、既有 `pkg/client/byteDance/douYin/clientTokenClient`
- 007-douyin-client-token: Added [if applicable, e.g., PostgreSQL, CoreData, files or N/A]
- 006-bytedance-access-token: Added Go 1.18（`go.mod`） + `github.com/ArtisanCloud/MediaXCore`（BaseClient/Logger/Cache）、`github.com/redis/go-redis/v9`（可选 Flow 缓存）、`pkg/client/byteDance/douYin/accessTokenClient`（DouYin SDK）、`internal/accesstoken/handler/mask`


<!-- MANUAL ADDITIONS START -->
Always respond in Chinese-simplified
<!-- MANUAL ADDITIONS END -->
