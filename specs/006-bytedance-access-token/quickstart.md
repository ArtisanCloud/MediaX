# Quickstart - DouYin AccessToken Debug Integration

## 1. 准备环境
1. 复制 `config.example.yaml` → `config.yaml`，确保 `access_token_providers` 中存在 `byte_dance/douyin` 分组且 `byte_dance_douyin_config` 只包含一个 `default` 模式。
2. 导出 DouYin OAuth 变量：
   ```bash
   export DOUYIN_CLIENT_ID="kxxxx"
   export DOUYIN_CLIENT_SECRET="sxxxx"
   export DOUYIN_SCOPE="video.list im.message.send"
   export DOUYIN_REDIRECT_URL="http://127.0.0.1:7071/debug/callback"
   export DOUYIN_OAUTH_URL="https://open.douyin.com/platform/oauth/connect"
   export DOUYIN_ACCESS_TOKEN_URL="https://open.douyin.com/oauth/access_token/"
   export ACCESSTOKEN_API_TOKEN="dev-accesstoken"
   export ACCESSTOKEN_CONFIG="./config.yaml"
   ```
3. 可选：配置 Redis（`ACCESSTOKEN_REDIS_ADDR=127.0.0.1:6379`），否则服务退回内存模式并在日志中提示 `storage_backend=memory`。

## 2. 启动 `cmd/accesstoken`
```bash
go run ./cmd/accesstoken/server -config $ACCESSTOKEN_CONFIG -listen :7071
```
日志应显示 `provider=byte_dance_douyin storage_backend=redis|memory`。

## 3. `/debug` 发起 DouYin OAuth
1. 浏览器打开 `http://127.0.0.1:7071/debug`，在 Provider 列表中选择 “字节跳动 / DouYin”。
2. 点击 “发起授权” → 浏览器跳至 DouYin 登录页，完成授权。
3. 回调后 `/debug` 的 Flow 列表出现新记录：`provider_code=byte_dance_douyin`、`flow_id=oauth-xxxx`、`expire_at`=当前时间+`expires_in` 秒。
4. 若缺少配置，UI toast 应提示 `OAuth scope 未配置` 等错误，同时日志输出 `provider=byte_dance_douyin event=oauth.start`。若 Flow 标记为 “需重新授权”，说明自动刷新失败，需要重新走 OAuth。

## 4. `/accesstoken/call` 调试 DouYin action
```bash
curl -H "Authorization: Bearer $ACCESSTOKEN_API_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
           "provider_code":"byte_dance_douyin",
           "provider_app":"douyin",
           "action":"douyin.video.list",
           "payload":{"cursor":0,"count":10}
         }' \
     http://127.0.0.1:7071/accesstoken/call | jq
```
- 响应应包含 `result`, `provider_code`, `flow_id`, `token_source`, `retry_count`。
- 若 AccessToken 即将过期，调用会先触发自动刷新；日志包含 `event=token.call ... status=refreshed` 且 Flow 的 `last_refresh_at` 更新时间。
- 若 refresh 失败，接口返回 `{"error":"need reauth"}`，同时 Flow 被删除。
- 若同一 action 在 1 秒内多次调用，第二次会返回 `429`。查看 `retry_count/last_backoff_ms` 可确认是否发生了指数退避。

## 5. Flow 回填与排障
1. `curl -H "Authorization: Bearer $ACCESSTOKEN_API_TOKEN" "http://127.0.0.1:7071/api/oauth/tokens?provider_code=byte_dance_douyin"` 查看 Flow 详情。
2. 点击 `/debug` Flow 列表中的 “填充” 将 JSON 写回调试表单。
3. 常见问题：
   - `provider code not found` → 检查 `access_token_providers` 配置是否包含 `byte_dance/douyin`。
   - `storage_backend=memory` 且 Flow 丢失 → Redis 未配置或故障，重新授权后提醒用户（Flow 会只保存在内存中）。
   - `need reauth` → Flow 的 refresh token 已失效，回到 `/debug` 重新发起授权即可。
   - `429 rate limit` → 调整脚本节奏，确保同一 action 至少间隔 1 秒。`retry_count` 会告诉你服务端已经重试了几次。

## 6. 文档与排障

- 更详细的配置说明：`docs/develop/access-token/byteDance/develop.md`
- `/debug` 操作细节与日志示例：`docs/develop/access-token/byteDance/debug.md`
- 若需将流程写入团队 wiki，可直接引用上述文档，确保强调“单实例 `<app>=default`”与 `need reauth` 排障步骤。

## 7. 日志验收示例

成功跑通上述 Quickstart 之后，`logs/accesstoken-server-info.log` 中应该出现类似记录：

```
accesstoken-server: event=oauth.start provider_code=byte_dance_douyin provider_app=douyin flow_id=oauth-XYZ token_source=pending
accesstoken-server: event=oauth.complete provider_code=byte_dance_douyin provider_app=douyin flow_id=oauth-XYZ token_source=authorization_code
accesstoken-server: event=token.call provider_code=byte_dance_douyin provider_app=douyin action=douyin.video.list flow_id=oauth-XYZ retry_count=0 storage_backend=redis
```

若自动刷新被触发，可看到 `status=refreshed` 与 `retry_count=0`；若 refresh 失败则会出现 `invalidate flow reason=douyin_refresh_failed:*`，与本文第 5 节的排障说明对应。

## 8. 文档走查（SC-004）建议流程

1. 将下表复制到 `docs/develop/access-token/byteDance/develop.md` 的“文档走查记录”，或直接在项目 Wiki 中登记：

   | 日期 | 参与人 | 耗时 (min) | 结果 | 阻塞/改进 |
   | --- | --- | --- | --- | --- |
   |  |  |  |  |  |

2. 仅向被访者提供 Quickstart/Develop/Debug 三份文档，引导他们从环境配置到 `/accesstoken/call` 全流程操作。
3. 运行过程中请保留 `/debug` 截图与 Redis 证据（`accesstoken:oauth:*`），并在表格中记录任何阻塞；若结果为 ⚠️，需要在相关文档新增排障说明后再复测。
