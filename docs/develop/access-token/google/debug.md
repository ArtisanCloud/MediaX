# AccessToken 调试沙盒使用指南

> 目标：像 `cmd/sessiontoken` 一样提供“所见即所得”的 Web 沙盒。本文只关注 `go run ./cmd/accesstoken/server` 启动的调试页面；CLI/Playground 细节请参考 `docs/develop/access-token/google/develop.md`。

## 1. 前置准备

- Go 1.18+、Redis（可选，用于缓存授权记录，默认连接 `127.0.0.1:6379`）。
- 仓库路径：`/private/var/www/html/ArtisanCloud/X/MediaX/core/MediaX`。
- `config.yaml` 已按照 `access_token_providers -> provider -> apps -> auth_modes` 的层级维护配置。示例：

  ```yaml
  access_token_providers:
    - code: google
      name: Google
      apps:
        - code: youtube
          name: YouTube
          api_version: v4
          auth_modes:
            - code: sandbox
              name: 沙盒默认
              google_youtube_config:
                oauth:
                  client_id: "${GOOGLE_YOUTUBE_CLIENT_ID}"
                  client_secret: "${GOOGLE_YOUTUBE_CLIENT_SECRET}"
                  redirect_url: "http://localhost:7071/debug/callback"
                  scope: "https://www.googleapis.com/auth/youtube.force-ssl"
  ```

- `redirect_url` 必须指向本地服务端口（默认 `7071`），否则 Google OAuth 回调不到 `/debug/callback`。
- `.env`/shell 中准备好 `GOOGLE_YOUTUBE_*`、`ACCESSTOKEN_*` 变量即可，系统会在请求体缺失时自动读取。

## 2. 启动调试服务

### 2.1 直接运行 Go

```bash
cd /private/var/www/html/ArtisanCloud/X/MediaX/core/MediaX
go run ./cmd/accesstoken/server -config config.yaml
```

启动成功会输出：

```
accesstoken-server: listening addr=:7071 api_token=de***en config=config.yaml
```

### 2.2 使用 Makefile 包装

```bash
make accesstoken-serve                  # 等价于 go run ./cmd/accesstoken/server -config config.yaml -port 7071
make accesstoken-serve ARGS='-port 8080 -config ./configs/google.yaml'
```

常用参数：

| 参数/环境变量 | 说明 |
| --- | --- |
| `-config` / `ACCESSTOKEN_CONFIG` | 配置文件路径，默认为 `MEDIA_X_CONFIG` 或 `config.yaml` |
| `-port` / `ACCESSTOKEN_LISTEN_ADDR` | 监听端口，默认 `:7071` |
| `-api-token` / `ACCESSTOKEN_API_TOKEN` | 保护调试接口的 Bearer Token，默认 `dev-accesstoken` |
| `ACCESSTOKEN_REDIS_ADDR/DB/PASS` | 配置 Redis；未显式设置地址时默认尝试 `127.0.0.1:6379`（连接失败才回退内存） |

> ⚠️ 记得把 `config.yaml` 里的 `redirect_url` 与实际监听端口保持一致，例如 `http://localhost:7071/debug/callback`。

## 3. 调试台操作流程（http://127.0.0.1:7071/debug）

1. **选择 Provider**：左上角“Provider”下拉会列出所有在配置中启用的提供商。例如 `Google` 下含 `YouTube`、`Blogger`，`字节跳动` 下含 `抖音` 等。
2. **选择 Provider App**：第二个下拉列出该 Provider 下的所有 App。默认选中 `google/youtube`，并自动展示 API 版本（如 `v4`）、配置路径与 OAuth Key。
3. **选择授权模式 + 同步模板**：第三个下拉对应 `auth_modes`，用于区分不同租户/环境。切换完成后点击“同步模板”即可把最新 Provider/App/Mode 写入“AccessToken 解析 / API 调用”两个 JSON，保持和 CLI 入参一致。
4. **AccessToken 解析**：在真正调用 API 之前先点“解析 AccessToken”（`POST /accesstoken/token`），系统会按“payload → 环境变量（`GOOGLE_YOUTUBE_ACCESS_TOKEN` / `ACCESSTOKEN_ACCESS_TOKEN` 等）→ config.yaml → 最近授权记录（Flow ID）”的优先级自动注入 token，并输出来源、脱敏值、TTL、`oauth_key`、`http_debug`，用来排查 401/配置差异。
5. **发起授权**：
   - 点击 “发起授权” → 浏览器打开 Google 授权页。
   - 正常登录后，Google 会回跳至 `http://127.0.0.1:7071/debug/callback?code=xxx&state=yyy`。
   - 页面下方 “授权回调记录” 会新增一条记录，包含 Provider、App、State、Code 以及完成时间。
6. **刷新 / Flow ID 回填**：点击 “刷新授权记录” 可从内存/Redis 拉取列表；若服务已重启，也可在“Flow ID 回填”输入框中直接填写 `oauth-xxxx` 并点击“加载 Flow ID”，背后会调用 `GET /api/oauth/tokens?flow_id=` 从 Redis 读取。页面提示了 Redis 主 key `accesstoken:oauth:<provider_code>:<app_code>:<mode>` 与索引 `accesstoken:oauth:flow:<flow_id>`，可用 `redis-cli GET ...` 直查。
7. **复用 AccessToken**：无论是刚刷新的记录还是 Flow ID 查询，都可以点击 “填充” 把 token 写回“AccessToken 解析 / API 调用”两个 JSON，之后继续解析或发起 API 调试，体验与 SessionToken 调试台一致。

> 调试台的所有数据都来自 `config.yaml`，因此增加 Provider/App 只需新增配置并重启服务即可。

## 4. OAuth 记录与缓存

- 成功回调后，服务器会调用 `completeOAuthFlow`：
  1. 使用 `client_id + client_secret + code` 换取 AccessToken/RefreshToken。
  2. 将完整结果写入内存索引，并（若配置）同步到 Redis，主 key 形如 `accesstoken:oauth:<provider_code>:<app_code>:<auth_mode>`。
  3. 额外写入 Flow ID 索引 `accesstoken:oauth:flow:<flow_id>`，用于在服务重启后根据 flow id 直接加载。
  4. 页面会显示 flow ID，可用于接口调用或线下排查。
- 缓存配置示例：

```bash
# 未设置 ACCESSTOKEN_REDIS_ADDR 时会自动连接 127.0.0.1:6379
export ACCESSTOKEN_REDIS_ADDR=127.0.0.1:6379
export ACCESSTOKEN_REDIS_DB=6
export ACCESSTOKEN_REDIS_PASS=''
```

- 若不希望落地 Redis，可省略上述环境变量，此时所有授权仅保存在进程内存，重启后需要重新授权，同时 Flow ID 回填功能将失效。

## 5. API 调试

调试台内置两个主要操作：

1. **解析 AccessToken**（`POST /accesstoken/token`）  
   - 用途：无需真正访问 Google API，就能核实当前 JSON 将使用哪一份 AccessToken、TTL、`oauth_key`、`http_debug`，是排查“为什么 CLI 能调通而 Web 沙盒 401/403”的首选途径。
   - 优先级：payload → 环境变量（`GOOGLE_YOUTUBE_ACCESS_TOKEN`、`ACCESSTOKEN_ACCESS_TOKEN` 等）→ `config.yaml` 中的 `oauth.access_token` → 最新授权记录（Flow ID/Redis）。
   - 返回值包含 token 来源、脱敏 token、TTL、HTTP 调试配置等，可直接截图给外部团队做证据。

2. **调试 API**（`POST /accesstoken/call`）  
   调试台的“API 调用”区域新增了和 SessionToken `/debug` 类似的表单，可直接在网页上选择 Action（`videos.list/search.list/playlists.list`）、填写 `part/ids/query/channel_id/max_results` 或勾选 `mine/search_mine`，页面会自动把这些值同步到 JSON，无需手动编辑。示例 JSON：

   ```json
   {
     "provider_code": "google",
     "provider_app": "youtube",
     "provider_auth_mode": "sandbox",
     "action": "videos.list",
     "payload": {
       "part": "snippet,statistics",
       "ids": "dQw4w9WgXcQ",
       "max_results": 3
     },
     "oauth": {
       "access_token": ""
     }
   }
   ```

   - 如果 `oauth.access_token` 为空，服务会尝试读取刚刚保存的授权记录并自动注入。
   - “调用结果” 面板会展示格式化 JSON，并附带请求日志（URL、响应时间、HTTP 状态）。
   - 若要调试 Blogger/其它 Provider，可先完成配置，前端会展示对应模板，后端接口将在后续迭代实现。

### 5.1 Flow ID / Redis 回填 API

- 浏览器端的“Flow ID 回填”按钮底层调用 `GET /api/oauth/tokens?flow_id=<id>`，服务会先查内存，再查 Redis 的索引 key（`accesstoken:oauth:flow:<id>`），最后返回标准的授权记录数组。
- 若需要脚本化复用，可以直接调用该接口，然后把返回的 `access_token` 手工填入 CLI/Playground。例如：

  ```bash
  curl -H 'Authorization: Bearer dev-accesstoken' \
       'http://127.0.0.1:7071/api/oauth/tokens?flow_id=oauth-abc123' | jq .
  ```

- 若接口返回空数组，请确认该 flow ID 是否仍在有效期内（`expire_at` 字段），或 Redis 是否仍保留对应的 key。

## 6. CLI / Playground 快速引用

虽然 Web 沙盒覆盖了绝大多数场景，但仍可通过下列命令快速重用配置：

- CLI：`go run ./cmd/accesstoken -config config.yaml -provider google -provider-app youtube -auth-mode sandbox -action videos.list -part snippet -ids dQw4w9WgXcQ`
- Makefile：`make accesstoken ARGS='-provider google -provider-app youtube -action search.list -query "MediaX demo" -part snippet'`
- Playground：`export PLAYGROUND_GOOGLE_YOUTUBE=1 && go run ./main.go`（会复用同一份配置、AccessToken）。

这些入口读取同一套 `access_token_providers` 配置，因此文档不再重复参数说明，详情见 `develop.md`。

## 7. 常见问题

| 问题 | 处理方式 |
| --- | --- |
| `http://127.0.0.1:7071/debug` 404 | 服务未启动或端口占用，重新执行 `go run ./cmd/accesstoken/server -config config.yaml -port 7071` |
| 页面只有 Google/YouTube，没有 Blogger | `config.yaml` 中缺少 `apps: [{code: blogger, ...}]` 或该 App 未启用授权模式，补齐后重启即可 |
| “发起授权” 跳回 404 | `oauth.redirect_url` 未更新为 `http://localhost:7071/debug/callback` |
| 授权记录刷新后消失 | 未配置 Redis，进程重启即丢失；或 `ACCESSTOKEN_API_TOKEN` 不匹配导致接口 401 |
| “调试 API” 返回 `暂未开放` | 该 Provider/App 还没实现后端调用逻辑，目前只有 `google/youtube v4` 支持 |
| 仍想使用旧入口 `go run ./main.go` | 仍然支持，但推荐统一在 Web 沙盒调好参数后，再复制 JSON 到 CLI/Playground，避免重复配置 |

> 结论：访问 `http://127.0.0.1:7071/debug`，选择 Provider → App → 授权模式 → 发起授权 → 回调 → 复用 token → 调试 API，即可完整模拟外部应用接入 AccessToken 客户端，操作路径与 SessionToken 调试台保持一致。
