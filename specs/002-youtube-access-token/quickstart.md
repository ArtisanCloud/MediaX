# Quickstart

1. **准备配置**
   - 复制 `config.example.yaml` → `config.yaml`。
   - 在 `google_youtube_config` 中填写 `client_id/client_secret/scope`，若已有长期 token，可直接填入 `oauth.refresh_token` 或通过环境变量提供。
   - 设定 `oauth_key`，方便 CLI/Playground 记录 token 来源。
   - 本地调试推荐将 `oauth.redirect_url` 设置为 `http://localhost:7071/debug/callback`，这样 Google OAuth 完成后可直接在 Web 沙盒页面查看回调日志。

2. **注入 AccessToken**
   - 临时 token：运行 `export GOOGLE_YOUTUBE_ACCESS_TOKEN='<token>'`。
   - 长期 token：实现 `GetOAuthToken` 回调（见 `playground/google.go` 示例）或通过 CLI `-access-token` 指定。

3. **运行 CLI 调试**
   ```bash
   make accesstoken ARGS='-action videos.list -part snippet -ids dQw4w9WgXcQ'
   make accesstoken ARGS='-action search.list -part snippet -query "MediaX" -max-results 3'
   make accesstoken ARGS='-action playlists.list -part snippet -channel-id UC_x5XG1OV2P6uZZ5FSM9Ttw'
   ```
   - 若失败，检查错误信息：`invalid_grant`（刷新 token 失效）、`quotaExceeded`、`insufficientPermissions` 等，详细排查流程见 `docs/develop/access-token/google/debug.md`。
   - 想要和 SessionToken 一样“先 export 再 make”也可以：`export ACCESSTOKEN_ACTION=videos.list`、`export ACCESSTOKEN_PART=snippet` 等，然后直接 `make accesstoken`，flag>环境变量>配置文件。

4. **启动 AccessToken 调试服务（Web 沙盒）**
   - `make accesstoken-serve`（或 `go run ./cmd/accesstoken/server -config config.yaml`）默认监听 `http://127.0.0.1:7071`，页面路径 `/debug`（与 `/debug/accesstoken` 等效）。
   - AccessToken JSON：`{"config_path":"config.yaml","access_token":"","access_token_ttl":3600}` → “解析 AccessToken”按钮会告诉你 token 来源与脱敏值。
   - API JSON：`{"action":"videos.list","part":"snippet","ids":"dQw4w9WgXcQ"}` → “执行 API”后会把 Google 响应写到 `<pre>`，同时复用 CLI 的日志/代理/缓存设置。
   - Provider/Provider App 下拉会根据配置自动生成（默认 “Google YouTube (default)”），未来新增平台可直接复用该调试台。
   - OAuth 回调日志：任何命中 `/debug/callback` 的请求都会显示在表格里，可一键清空；通过 `ACCESSTOKEN_API_TOKEN`（默认为 `dev-accesstoken`）保护 API，支持 `ACCESSTOKEN_LISTEN_ADDR`/`ACCESSTOKEN_REDIS_*` 等环境变量。

5. **验证 Playground**
   - 设置开关：`export PLAYGROUND_GOOGLE_YOUTUBE=1`（可选 `PLAYGROUND_CACHE_MODE=memory` 切换缓存模式）。
   - 运行 `go run ./main.go`，观察控制台输出和 `logs/info.log`/`logs/error.log`。
   - `playground/google.go` 会自动复写 `GetOAuthToken` 并脱敏日志，如需接入 Vault/Redis，可直接修改该回调或覆盖 `GOOGLE_YOUTUBE_ACCESS_TOKEN`。完成 CLI → Playground 即可对照 `docs/develop/access-token/google/develop.md` 中的闭环表格检查“订阅→视频→发布→评论”各环节。

6. **故障排查**
   - 清理缓存：`redis-cli --raw keys 'mediax.access_token.*' | xargs -I{} redis-cli DEL {}`。
   - 检查日志：`tail -f logs/info.log | rg youtube`。
   - 代理/网络：若需要代理，在 `google_youtube_config.proxy_api_url` 设置并再次运行 CLI。
