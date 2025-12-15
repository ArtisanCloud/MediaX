# Quickstart

1. **准备配置**
   - 复制 `config.example.yaml` → `config.yaml`。
   - 在 `google_youtube_config` 中填写 `client_id/client_secret/scope`，若已有长期 token，可直接填入 `oauth.refresh_token` 或通过环境变量提供。
   - 设定 `oauth_key`，方便 CLI/Playground 记录 token 来源。

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

4. **启动 AccessToken 调试服务（开发中）**
   - 目标是与 SessionToken 调试台一致：`make accesstoken-serve`（或 `go run ./cmd/accesstoken/server`）后默认监听 `http://127.0.0.1:7070`，页面路径为 `/debug/accesstoken`。
   - 服务会读取 `config.yaml` 与 `GOOGLE_YOUTUBE_*` 环境变量，支持在 UI 中切换 Provider/App/API 版本、刷新 AccessToken、调用 `videos.list/search.list/playlists.list`，并展示 `/debug/callback` 收到的回调日志。
   - 可通过 `ACCESSTOKEN_API_TOKEN`（默认 `dev-accesstoken`）限制访问；若需要代理/Redis，可沿用 CLI 的 `SESSIONTOKEN_REDIS_*` 与 `HTTPS_PROXY` 变量。该服务目前处于规划阶段，请关注 `specs/002-youtube-access-token/tasks.md` 中的 US4 进度。

5. **验证 Playground**
   - 设置开关：`export PLAYGROUND_GOOGLE_YOUTUBE=1`（可选 `PLAYGROUND_CACHE_MODE=memory` 切换缓存模式）。
   - 运行 `go run ./main.go`，观察控制台输出和 `logs/info.log`/`logs/error.log`。
   - `playground/google.go` 会自动复写 `GetOAuthToken` 并脱敏日志，如需接入 Vault/Redis，可直接修改该回调或覆盖 `GOOGLE_YOUTUBE_ACCESS_TOKEN`。完成 CLI → Playground 即可对照 `docs/develop/access-token/google/develop.md` 中的闭环表格检查“订阅→视频→发布→评论”各环节。

6. **故障排查**
   - 清理缓存：`redis-cli --raw keys 'mediax.access_token.*' | xargs -I{} redis-cli DEL {}`。
   - 检查日志：`tail -f logs/info.log | rg youtube`。
   - 代理/网络：若需要代理，在 `google_youtube_config.proxy_api_url` 设置并再次运行 CLI。
