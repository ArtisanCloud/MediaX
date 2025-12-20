# RedBook JuGuang AccessToken 调试沙盒

> 与 Google/BiliBili 相同，所有调试能力都来自 `go run ./cmd/accesstoken/server -config config.yaml`。本文仅补充小红书聚光在该调试台中的差异点，强调“同一套服务、同一份配置、同一个 `/debug` 页面”。

## 1. 前置准备

- Go 1.21、`redis-server`（可选，用于 Flow 持久化；未配置时默认尝试 `127.0.0.1:6379`，再回退内存）。
- `config.yaml` 已按 `access_token_providers -> redbook -> apps -> auth_modes -> redbook_juguang_config` 填写（参考 `docs/develop/access-token/redbook/develop.md`）。
- 环境变量：
  ```bash
  export REDBOOK_JUGUANG_CLIENT_ID=xxx
  export REDBOOK_JUGUANG_CLIENT_SECRET=yyy
  export REDBOOK_JUGUANG_SCOPE='notes.read,notes.write'
  export REDBOOK_JUGUANG_OAUTH_URL='https://ad.xiaohongshu.com/oauth2/authorize'
  export REDBOOK_JUGUANG_ACCESS_TOKEN_URL='https://ad.xiaohongshu.com/oauth2/access_token'
  export REDBOOK_JUGUANG_REDIRECT_URL='http://127.0.0.1:7071/debug/callback'
  export REDBOOK_JUGUANG_OAUTH_KEY='tenantA.juguang'
  export ACCESSTOKEN_API_TOKEN='dev-accesstoken'
  # Redis 可选：
  export ACCESSTOKEN_REDIS_ADDR='127.0.0.1:6379'
  export ACCESSTOKEN_REDIS_DB='6'
  ```

## 2. 启动命令

```bash
cd /private/var/www/html/ArtisanCloud/X/MediaX/core/MediaX
# 与 Google/BiliBili 相同的命令
go run ./cmd/accesstoken/server -config config.yaml
```

预期日志：
```
{"msg":"accesstoken-server: event=listening ... provider=redbook ... storage_backend=redis"}
```
- 如果未连上 Redis，会看到 `redis disabled, Flow 数据将存储在内存中...`，页面顶部也会出现黄色警告。
- `flow_ttl_seconds` 默认为 86400，可通过 `ACCESSTOKEN_FLOW_TTL_SECONDS` 覆盖。

## 3. 访问 `/debug`

打开 http://127.0.0.1:7071/debug，右上角输入 API Token（或 URL 附带 `?api_token=dev-accesstoken`）。

### 3.1 Provider 下拉 & 新卡片

1. 选择 **Provider → “小红书”**，应用列表会出现 “小红书聚光”。
2. `Provider Card`（新增组件）会立即展示：
   - 标题：`小红书 / 小红书聚光`
   - API 版本：`API v1`
   - 标签：`模式：默认 OAuth`、`OAuth Key：tenantA.juguang`
   - 元信息：`Provider Code redbook_juguang`、`App Code juguang`、`Config Path config.yaml`
   - **默认模板**：当前 Provider/App/Mode 的 JSON 摘要，可直接复制给 CLI 或脚本。
3. 编辑 `Config Path` 时卡片会实时更新，确保不同配置文件的调试动作一目了然。

### 3.2 同步模板 & OAuth

- 点击 “同步模板” 会把 Provider/Card 信息写入 “AccessToken 解析” 与 “API 调用” 两个 JSON 表单，模板格式示例：
  ```json
  {
    "provider_code": "redbook_juguang",
    "provider_app": "juguang",
    "provider_auth_mode": "default",
    "config_path": "config.yaml"
  }
  ```
- 点击 “发起授权” → 浏览器跳至聚光登录页 → 成功后回调 `http://127.0.0.1:7071/debug/callback`。
- 页面第 5 节会出现回调日志；第 3 节的 “授权记录” 列表会出现新的 Flow，字段含 `storage_backend`、`expire_at`、`oauth_key`。

### 3.3 Flow 回填 & Redis 检查

- 使用 “刷新授权记录” 或者在 “Flow ID 回填” 输入 `oauth-xxxxx`，调试台会通过 `/api/oauth/tokens` 查询内存/Redis。
- 对应的 Redis 键：
  - 主记录：`accesstoken:oauth:redbook_juguang:juguang:default`
  - 索引：`accesstoken:oauth:flow:<flow_id>`
- 如果 Redis 处于禁用状态，顶部会持续提示 “Flow 数据当前存储在内存中...”，请尽快填好 `ACCESSTOKEN_REDIS_ADDR`，否则重启即丢失。

## 4. API 调用示例

`/accesstoken/call` 现在同时支持 YouTube 与聚光调试。若要调用聚光账户余额接口，只需在 “API 调用” JSON 中写入：

```json
{
  "provider_code": "redbook_juguang",
  "provider_app": "juguang",
  "provider_auth_mode": "default",
  "config_path": "config.yaml",
  "action": "redbook.account.balance",
  "payload": {
    "advertiser_id": 123456789
  }
}
```

点击 “执行 API” 后，服务端会复用当前 Flow/环境变量中的 AccessToken，通过聚光客户端转发到 `/api/open/jg/account/balance/info`，响应 JSON 将直接展示在调试台。无需额外的 `/accesstoken/redbook/*` 路由，所有 Provider 共用同一入口。

AccessToken 的优先级仍为：请求体 → `REDBOOK_ACCESS_TOKEN/XHS_ACCESS_TOKEN` → `config.yaml->oauth.access_token` → 最近 Flow；因此可以与 Google/BiliBili 保持完全一致的操作方式。

## 5. 常见提示

| 场景 | 提示 |
| --- | --- |
| 调试页提示 `unsupported provider` | 说明请求体里 `provider_code` 填成 `redbook` 等；请使用 `redbook_juguang`。|
| “OAuth scope 未配置” | 新增字段，必须在 `config.yaml` 或 `REDBOOK_JUGUANG_SCOPE` 中给出，逗号分隔多个能力。|
| 页面显示 “Flow 数据当前存储在内存中…” | Redis 地址为空或连接失败；请检查 `ACCESSTOKEN_REDIS_ADDR/DB/PASSWORD`。|
| 页面只显示 Google/BiliBili | 需要重启服务以重新加载配置，或检查 `auth_modes` 是否为空。|

> 与其它 Provider 一样：复制 `config.example.yaml` → 注入环境变量 → `go run ./cmd/accesstoken/server -config config.yaml` → 打开 `/debug`，即可完成聚光 OAuth & Flow 管理，无需新增 CLI 或网页。
