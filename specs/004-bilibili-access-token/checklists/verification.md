# Verification Evidence

| 项目 | 说明 & 证据 |
| --- | --- |
| ✅ `go test ./cmd/accesstoken/server` | 2025-12-19 在仓库根目录执行，输出 `ok github.com/ArtisanCloud/MediaX/cmd/accesstoken/server 0.73s`，覆盖 OAuth/Flow/Token/Replay 合约测试。 |
| ✅ Flow 列表 API | `curl -H "Authorization: Bearer $ACCESSTOKEN_API_TOKEN" "http://127.0.0.1:7071/accesstoken/flows?provider_code=bilibili&provider_app=content_center&provider_auth_mode=default"` 返回 `flows[0].status=authorized`、`masked_account=bilibili/content_center/default`。 |
| ✅ Flow Replay API | `curl -sS -X POST http://127.0.0.1:7071/accesstoken/flow/replay -H "Authorization: Bearer $ACCESSTOKEN_API_TOKEN" -H "Content-Type: application/json" -d '{"flow_id":"oauth-demo-flow"}'` → JSON 中 `payload.access_token` 与 `/accesstoken/token` 响应一致，`status=authorized`。 |
| ✅ 调试页交互 | 浏览器访问 `/debug?api_token=dev-accesstoken`：可点击“发起授权”完成流程，`授权记录` 表显示最新 Flow，点击“填充”后 JSON 表单自动注入 AccessToken。 |
