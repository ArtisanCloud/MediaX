# Data Model

## 1. GoogleYouTubeConfig
- **Location**: `pkg/client/config/google.go`
- **Fields**:
  - `ApiUrl` (string, default `https://www.googleapis.com`)
  - `ProxyApiUrl` (string, optional)
  - `Timeout` (float64 seconds)
  - `HttpDebug` (bool)
  - `OAuth` (struct)
    - `AccessTokenUrl`
    - `ProxyAccessTokenUrl`
    - `ClientID`
    - `ClientSecret`
    - `RedirectUrl`
    - `Scope`
    - `RefreshToken`
    - `AccessToken`
    - `TokenExpiry`
  - `OauthKey` (string,映射外部 token 仓库)
  - `GetOAuthToken` (func,运行时注入)
- **Relationships**: 传入 `MediaX.CreateGoogleYouTubeACClient` → `core.GoogleAccessTokenHandler`。
- **Validation**: `ClientID`/`ClientSecret` 必填；AccessToken/RefreshToken 至少提供一种；`Timeout>0`；敏感字段仅可通过 env/secret store 注入。

## 2. AccessToken CLI Invocation
- **Location**: `cmd/accesstoken/main.go`
- **Attributes**:
  - `Action` (enum: videos.list/search.list/playlists.list)
  - `Part` (string,必填)
  - `IDs` (string,可选)
  - `Query` (string,search 用)
  - `ChannelID` (string)
  - `Mine` (bool,playlists)
  - `SearchMine` (bool,search)
  - `MaxResults` (int,默认5,1-50)
  - `PageToken` (string)
  - `Region`/`VideoCategory`/`Chart`（videos-specific）
  - `AccessToken` (string,来自 flag/env/config)
  - `AccessTokenTTL` (int,秒)
- **Relationships**: CLI 将参数映射至 `video/search/playlists` schema 请求 → `GoogleYouTubeACClient`。
- **Validation**: Action + Part 必填；`MaxResults` ≤50；互斥组合（如 `ids` vs `mine`）需在 CLI 层校验。

## 3. Quickstart Flow
- **Steps**:
  1. 复制 `config.example.yaml` → `config.yaml`。
  2. 填入 `google_youtube_config` 凭证或设置环境变量。
  3. 运行 `make accesstoken ...` 验证；可选 `go run ./main.go` 启动 Playground。
- **State transitions**: `未配置` → `已配置` → `CLI 成功`/`CLI 失败`（触发故障排查）。
