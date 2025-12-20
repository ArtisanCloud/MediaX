# Quickstart: RedBook 聚光接入调试台

1. **准备配置**
   - 复制 `config.example.yaml` 为 `config.yaml`。
   - 在 `access_token_providers.redbook.apps[].auth_modes[].redbook_juguang_config.oauth` 中填入真实 `client_id`/`client_secret`/`scope`，`redirect_url` 设为 `http://127.0.0.1:7071/debug/callback`。
   - 设置环境变量：
     ```bash
     export REDBOOK_JUGUANG_CLIENT_ID=xxx
     export REDBOOK_JUGUANG_CLIENT_SECRET=xxx
     export ACCESSTOKEN_API_TOKEN=dev-accesstoken
     export ACCESSTOKEN_REDIS_ADDR=127.0.0.1:6379
     ```

2. **启动调试服务**
   ```bash
   go run ./cmd/accesstoken/server -config config.yaml
   ```
   - 成功日志应包含 `provider_code=redbook_juguang` 的加载提示。

3. **OAuth 授权**
   - 打开 `http://127.0.0.1:7071/debug`，选择“小红书 / 聚光” Provider → App → 授权模式。
   - 点击“发起授权”，在浏览器中完成聚光登录授权。
   - 回调后，页面底部 “授权记录” 会显示新 Flow，可点击“填充”将 token 写入 JSON。

4. **调用聚光 API**
   - 在调试页或终端调用账户余额接口：
     ```bash
     curl -H "Authorization: Bearer dev-accesstoken" \
          -H "Content-Type: application/json" \
          -d '{
                "provider_code": "redbook_juguang",
                "provider_app": "juguang",
                "provider_auth_mode": "default",
                "config_path": "config.yaml",
                "action": "redbook.account.balance",
                "payload": {"advertiser_id": 12345}
              }' \
          http://127.0.0.1:7071/accesstoken/call | jq .
     ```
   - 服务端日志会显示 `event=token.call provider=redbook_juguang action=redbook.account.balance`，调试页「API 调用」面板同样会输出 JSON 结果与错误信息。

5. **Flow 回填 / Redis 验证**
   - 使用 `/accesstoken/flows?provider_code=redbook_juguang` 查看 Flow 列表。
   - 通过 `GET /api/oauth/tokens?flow_id=...` 可将 Flow ID 回写到 CLI。

6. **常见问题**
   - 缺少 `scope` → `/accesstoken/oauth/start` 返回 `OAuth scope 未配置`。
   - Redis 未连通 → 日志提示 `storage_backend=memory`，重启后需重新授权。
   - 调试 API 返回 `redbook action ... 暂未开放` → 确认 `action` 是否填写为 `redbook.account.balance`，或等待更多 Action 接入。
   - 聚光 API 报签名/401 → 检查 `advertiser_id` 与授权账号是否匹配，并确认 `REDBOOK_ACCESS_TOKEN`/Flow 中的 token 仍在有效期内。
