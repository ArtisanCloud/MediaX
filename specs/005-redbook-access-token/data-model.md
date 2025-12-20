# Data Model

## ProviderMetadata
- **provider_code** (`string`): 例如 `redbook_juguang`，由 config.yaml 的 `provider_code` 派生。
- **app_code** (`string`): 聚光应用标识（默认 `juguang`）。
- **auth_mode** (`string`): OAuth 模式键（默认 `default`）。
- **config_path** (`string`): 配置文件路径，写入 `/debug` 模板与日志。
- **oauth_key** (`string`): 区分租户/账号的键，写入前端模板与 Flow 记录。
- **api_version** (`string`): 展示在调试 UI，便于后续多版本。

## OAuthFlowRecord
- **flow_id** (`string`): `oauth-<state>`，同时作为 Redis flow 索引。
- **provider_code/provider_app/provider_auth_mode** (`string`): 过滤条件，必须大小写不敏感匹配。
- **storage_backend** (`enum: redis|memory`): 指示 Flow 存储位置，便于调试提示。
- **access_token / masked_token** (`string`): 实际 token（仅内部保存）与脱敏展示值。
- **token_source / token_source_detail** (`string`): 记录 token 来源（request/env/config/flow）。
- **flow_ttl_seconds / flow_expire_at** (`int` / `timestamp`): TTL 与过期时间，用于 `/accesstoken/flows` 排序。
- **callback** (`object`): 保存最近一次 `/debug/callback` 请求/响应摘要，用于审计。

## RedBookJuGuangConfig (扩展字段)
- **api_url/proxy_api_url/timeout/http_debug**: 继承 `ClientConfig` 行为。
- **oauth** (`OAuthConfig`): 需补齐 `oauth_url`, `access_token_url`, `scope`, `redirect_url`。
- **meta.oauth_key** (`string`): 供 UI 回显及缓存分组。
- **GetOAuthToken func**: 允许 CLI/Playground 注入自定义 token，调试台透传 Flow ID 时调用。

## Flow API Endpoint Payload（账户余额示例）
- **provider_code/provider_app/provider_auth_mode**: 请求中必须填写；可由前端 JSON 模板生成。
- **config_path**: 限制服务读取指定的聚光配置。
- **payload** (`object`): 包含 `advertiser_id`，由 `JuGuangAccountGetAccountBalanceReq` 定义。
