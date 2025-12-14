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
   - 若失败，检查错误信息：`invalid_grant`（刷新 token 失效）、`quotaExceeded`、`insufficientPermissions` 等。

4. **验证 Playground**
   - 在 `main.go` 中取消注释 `playground.PlayGoogleYouTube(localConfig, mediaX)`。
   - 运行 `go run ./main.go`，观察控制台输出和 `logs/info.log`。
   - 修改 `GetOAuthToken` 回调，从 Vault/Redis 读取 token 以模拟真实环境。

5. **故障排查**
   - 清理缓存：`redis-cli --raw keys 'mediax.access_token.*' | xargs -I{} redis-cli DEL {}`。
   - 检查日志：`tail -f logs/info.log | rg youtube`。
   - 代理/网络：若需要代理，在 `google_youtube_config.proxy_api_url` 设置并再次运行 CLI。
