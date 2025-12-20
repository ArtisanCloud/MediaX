# Verification Evidence

- [x] `curl` 调用 `/accesstoken/call`（action=redbook.account.balance）截图/JSON：`curl -H "Authorization: Bearer dev-accesstoken" -d '{...}' http://127.0.0.1:7071/accesstoken/call | jq .`。
- [x] Redis 中 `accesstoken:oauth:redbook_juguang:juguang:default` 与 `accesstoken:oauth:flow:<flow_id>` 的键值截图，证明 Flow 持久化成功。
- [x] `/debug` 页面授权记录 + Flow 回填截图，展示 token 填充后的 API 调用结果。
