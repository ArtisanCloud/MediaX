# Data Model

## ProviderSelection
| Field | Type | Source | Validation / Notes |
|-------|------|--------|--------------------|
| `provider_code` | string | YAML `access_token_providers[].code` | 必须等于 `bilibili`; 参与日志/metrics 标签 |
| `provider_app` | string | YAML `apps[].code` | 与 Auth 模式联合唯一；驱动 UI 模板与 Redis key |
| `provider_auth_mode` | string | YAML `auth_modes[].key` | 必填；决定 OAuth client/redirect 配置 |
| `config_path` | string | CLI/env (`ACCESSTOKEN_CONFIG`) | 有效文件路径；用于 `/accesstoken/token` 响应回显 |
| `listen_addr` | string | env `ACCESSTOKEN_LISTEN_ADDR` | 默认 `127.0.0.1:7071`；用于日志/安全检查 |

**Relationships**: ProviderSelection→AccessTokenRecord (1:N)；同一 selection 的最新记录覆盖 `accesstoken:oauth:{provider_app}:{mode}`。

## AccessTokenRecord
| Field | Type | Source | Validation / Notes |
|-------|------|--------|--------------------|
| `flow_id` | string | FlowRecord reference | 形如 `oauth-<uuid>`；必须存在 FlowRecord 才可写入 |
| `masked_token` | string | Derived | 只保留头尾 4 字符；任何 API 响应均为脱敏值 |
| `token_source` | enum(`flow`,`env`,`config`,`request`) | Resolver | 决定展示优先级与日志字段 |
| `oauth_key` | string | Provider config | 映射回 OAuth app id；用于排障 |
| `expire_at` | timestamp | TTL calc | `issued_at + expires_in`；随响应返回 TTL |
| `storage_backend` | enum(`memory`,`redis`) | runtime | 指示当前记录来源，便于 UI 消费 |
| `created_at` | timestamp | server clock | 记录写入时间 |
| `updated_at` | timestamp | server clock | Flow 重放时刷新 |

**State**: `pending`（等待 OAuth）→`issued`（Flow 回调成功并写 masked_token）→`expired`（TTL 过期，UI/API 提示重新授权）。

## FlowRecord
| Field | Type | Source | Validation / Notes |
|-------|------|--------|--------------------|
| `flow_id` | string | Generated UUID | key=`accesstoken:oauth:flow:<flow_id>`；24h TTL |
| `provider_key` | string | ProviderSelection composite | `<provider_code>:<app>:<auth_mode>`；索引主记录 |
| `authorize_url` | string | OAuth start response | 记录发起时 URL，便于重放 |
| `callback_payload` | object | `/debug/callback` body | 完整保留 `code`, `state`, `cookies` 脱敏字段 |
| `redis_pointer` | string | Redis storage | 指向主 key；内存模式该字段可为空 |
| `tenant_uuid` | string | env/context | 必须写入日志/record；默认 `default` 但可扩展 |
| `expires_in` | int | OAuth provider | 与 BiliBili 响应对齐；用于 TTL check |
| `ttl_seconds` | int | constant 86400 | Flow 与主记录都采用 24h TTL；与 `expires_in` 取 min |
| `status` | enum(`pending`,`authorized`,`invalid`,`expired`) | derived | 根据回调进度/TTL 更新 |

**Lifecycle**: 由 `/accesstoken/oauth/start` 创建 `pending`；回调后置为 `authorized` 和 `status_reason=success`；回填失败或 TTL 过后标记 `invalid/expired` 并在 UI 提示。

## DebugServiceBinding
| Field | Type | Source | Validation / Notes |
|-------|------|--------|--------------------|
| `listen_addr` | string | env flag/default | 默认 `127.0.0.1:7071`; 仅允许 operator 明确覆盖 |
| `api_token` | string | env `ACCESSTOKEN_API_TOKEN` | 必填；用于 Bearer 校验；日志仅显示前两位 |
| `redis_available` | bool | env presence | 控制 UI 提示与存储路径 |
| `memory_store` | map | process scoped | 保存 Flow/Token 结构；需互斥保护 |

**Interactions**: Binding→Handler wiring: CLI 启动 -> validates config -> registers HTTP routes -> wires Flow store (memory or Redis) based on `redis_available`.
