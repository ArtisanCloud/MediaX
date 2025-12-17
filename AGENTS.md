# MediaX Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-12-09

## Active Technologies
- Go 1.18（仓库 `go.mod` 已指定） + `github.com/ArtisanCloud/MediaXCore`（BaseClient/Logger）、`github.com/redis/go-redis/v9`（可选 Token 缓存）、`gopkg.in/yaml.v3`（配置解析） (002-youtube-access-token)
- Redis（用于 AccessToken cache，可选；本特性主要读取本地 YAML + 环境变量） (002-youtube-access-token)

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
- 002-youtube-access-token: Added Go 1.18（仓库 `go.mod` 已指定） + `github.com/ArtisanCloud/MediaXCore`（BaseClient/Logger）、`github.com/redis/go-redis/v9`（可选 Token 缓存）、`gopkg.in/yaml.v3`（配置解析）

- 001-session-token-client: Added Go 1.18（仓库 `go.mod`） + `github.com/ArtisanCloud/MediaXCore`（日志/HTTP/BaseClient）、`github.com/redis/go-redis/v9`（Flow 持久化）、`gopkg.in/yaml.v3`（配置解析）

<!-- MANUAL ADDITIONS START -->
Always respond in Chinese-simplified
<!-- MANUAL ADDITIONS END -->
