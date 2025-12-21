# Data Model - DouYin AccessToken Debug Integration

## Entity: DouYinProviderConfig
| Field | Type | Source | Rules |
|-------|------|--------|-------|
| provider_code | string | config.yaml (`byte_dance_douyin_config`) | 固定 `byte_dance_douyin`；与 UI Provider 卡片唯一对应 |
| app_code | string | config.yaml | 固定 `douyin`（映射 `<app>=default`） |
| auth_mode | string | config.yaml | 默认为 `default`；只有一个模式 |
| oauth_key | string | config.yaml | 必填；供 `/debug` UI 展示 |
| api_version | string | config.yaml | 可选，显示在 Provider 卡片 |
| client_id | string | env (`DOUYIN_CLIENT_ID`) | 必填；非空，<=128 字符 |
| client_secret | string | env (`DOUYIN_CLIENT_SECRET`) | 必填；仅通过 env 注入 |
| scope | string | env (`DOUYIN_SCOPE`) | 必填；空格分隔；序列化前需去重排序 |
| redirect_url | string | env (`DOUYIN_REDIRECT_URL`) | 可覆盖默认 `/debug/callback`；需合法 URL |
| oauth_url | string | env (`DOUYIN_OAUTH_URL`) | 必填；HTTPS |
| access_token_url | string | env (`DOUYIN_ACCESS_TOKEN_URL`) | 必填；HTTPS |
| http_debug | bool | config.yaml | 可选；默认 false |
| storage_backend | enum(`redis`,`memory`) | server runtime | 由 `buildCacheStore` 判定并显示 |

## Entity: DouYinOAuthFlowRecord
| Field | Type | Source | Rules |
|-------|------|--------|-------|
| flow_id | string | `/debug` state | 唯一，格式 `oauth-<state>`；索引键 `accesstoken:oauth:flow:<flow_id>` |
| provider_code | string | resolver | 固定 `byte_dance_douyin` |
| provider_app | string | resolver | 固定 `douyin` |
| auth_mode | string | resolver | 固定 `default` |
| config_path | string | request | 用户指定的配置路径，原样回显 |
| access_token | string | DouYin token API | 必填，保存全量字符串但 UI 只展示 mask |
| refresh_token | string | DouYin token API | 必填，用于自动刷新；在 UI/日志中脱敏 |
| expires_in | int | DouYin token API | 必填，>0；同样写入 `flow_ttl_seconds` |
| token_expire_at | timestamp | server computed | `stored_at + expires_in` |
| flow_ttl_seconds | int | server computed | 直接取 `expires_in`，刷新后同步更新 |
| flow_expire_at | timestamp | server computed | `stored_at + flow_ttl_seconds` |
| storage_backend | enum | runtime | `redis` 或 `memory`，来自 cache 层 |
| source | enum | server | `authorization_code`、`refresh_token`、`flow_replay` |
| retry_count | int | runtime | `/accesstoken/call` 调用 DouYin API 时记录的重试次数 |
| last_refresh_at | timestamp | runtime | 自动刷新成功的时间；失败则置空并删除 Flow |
| status | enum(`fresh`,`refreshing`,`need_reauth`) | derived | `fresh`=授权完成；`refreshing`=refresh 中（debug UI 提示）；`need_reauth`=刷新失败/Flow 被删除 |

## Relationships
- **DouYinProviderConfig → DouYinOAuthFlowRecord**：一对多。单 app (`douyin`，映射 `<app>=default`) 可生成若干 Flow（不同 `flow_id`、不同 `mode`）。Flow key `accesstoken:oauth:byte_dance_douyin:default:<mode>` 指回最新记录。
- **Flow ↔ `/accesstoken/call`**：每次调用读取最新 Flow 记录并可能触发自动刷新，刷新成功写回同一条记录。

## Validation Rules & Constraints
- `client_id`, `client_secret`, `scope`, `oauth_url`, `access_token_url` 缺失时 `/accesstoken/oauth/start` 必须报错并记录日志（FR-002）。
- Flow TTL 必须与 `expires_in` 一致，不允许长于 DouYin token；刷新成功后立刻更新 Redis TTL。
- 自动刷新失败后 Flow 必须删除，`/accesstoken/call` 返回 `need reauth` 并在 `/debug` Flow 列表中提示（FR-012）。
- 速率限制在 `action` 维度 1 QPS，超限返回 429（带 `retry_after` 提示）。
- `retry_count`、`token_source`、`flow_id` 必须写入日志（宪章 Observability）。

## State Transitions
1. `pending`（发起 `/accesstoken/oauth/start`）
2. `authorizing`（用户在 DouYin 页面完成登录）
3. `succeeded`（`/debug/callback` 交换成功 → Flow 写入）
4. `refreshing`（调用 `/accesstoken/call` 检测 nearing expiry → 调用 refresh）
5. `succeeded`（刷新成功，字段更新）或 `need_reauth`（刷新失败 → Flow 删除，用户需重新授权）

状态变化需记录 `flow_id`、`provider_code`、`tenant_uuid`（若可用）以及时间戳，方便审计。
