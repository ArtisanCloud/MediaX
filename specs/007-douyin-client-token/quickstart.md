# Quickstart - DouYin ClientToken Mode

1. **准备配置**
   - 复制 `config.example.yaml` 中的 `client_token_providers.providers[byte_dance]` 节点到本地 `config.yaml`。
   - 设置环境变量：
     ```bash
     export DOUYIN_CLIENT_KEY=xxx
     export DOUYIN_CLIENT_SECRET=yyy
     # 若业务涉及安全/风控字段，请同时设置，缺省则调试台禁用对应 API 模板
      export DOUYIN_DEVICE_ID=${DOUYIN_DEVICE_ID:-""}
      export DOUYIN_RISK_INFO=${DOUYIN_RISK_INFO:-""}
     export CLIENTTOKEN_REDIS_ADDR=127.0.0.1:6379
     export CLIENTTOKEN_REDIS_DB=0
     export CLIENTTOKEN_API_TOKEN=dev-clienttoken
     ```

2. **启动调试服务**
   ```bash
   CLIENTTOKEN_CONFIG=config.yaml \
   CLIENTTOKEN_LISTEN_ADDR=:7072 \
   go run ./cmd/clienttoken/server -config "$CLIENTTOKEN_CONFIG"
   ```
   - 日志应包含 `provider=byte_dance_douyin_clienttoken`。

3. **刷新 Token**
   - 浏览器访问 `http://127.0.0.1:7072/debug?api_token=dev-clienttoken`。
   - 选择 Provider→App→Mode，点击“同步模板”→“刷新 Token”。
   - 确认“查看缓存”显示 `source=redis`（或 `memory`）。

4. **CLI/脚本复用**
   - `./scripts/clienttoken-douyin.sh refresh`：刷新 token。
   - `./scripts/clienttoken-douyin.sh cache show|clear`：查看或清空缓存。
   - `./scripts/clienttoken-douyin.sh call open_api/1/content/video/list/ POST "" '{"page":1,"size":10}'`：发起 Douyin API 调试。
   - 使用 `redis-cli --raw GET clientToken:douyin:$DOUYIN_CLIENT_KEY | jq '.'` 验证 TTL 与值。

5. **API 调试**
   - 在调试页面填写 `action=open_api/1/content/video/list/`、请求体。
   - 若页面提示 “Douyin API 模板已禁用”，请确认已同时配置 `DOUYIN_DEVICE_ID` 与 `DOUYIN_RISK_INFO`。
   - 点击“调用 API”，检查响应及日志字段 `token_source`, `ttl_remaining`、`storage_backend`。

6. **排障**
   - `provider code not found`：确认配置已加载并重启服务。
   - `storage_backend=memory`：说明 Redis 连不上，检查 `CLIENTTOKEN_REDIS_*`。
   - Token 将到期：刷新日志出现 `ttl_remaining<refresh_before_seconds`，系统会自动刷新；如失败，手动点击“刷新 Token”，若仍失败则终端/页面会提示 `token refresh failed`。
   - Douyin API 模板禁用：确保两个风控变量同时存在，否则只能使用“自定义”模式自行输入接口请求。

## 验证记录

- **2025-02-14**：按上述步骤完成一次端到端验证，`go test ./cmd/clienttoken/... ./pkg/client/byteDance/douYin/...` 通过；`scripts/clienttoken-douyin.sh refresh/cache/call` 输出与调试台一致，Redis 记录 `clientToken:douyin:$DOUYIN_CLIENT_KEY` TTL 正常。
