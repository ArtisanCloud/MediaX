# Data Model: SessionTokenClient Provider Architecture

## Flow
| Field | Type | Description | Rules |
|-------|------|-------------|-------|
| flow_id | string (`stf_<uuid>`) | 全局唯一 Flow 标识 | 由 SDK 生成；Redis key = `sessiontoken:flow:<flow_id>` |
| provider_code | string | 平台代码（如 `zhihu`） | 与 MediaX provider registry 对齐 |
| tenant_uuid | string | 调用方租户 | 参与日志与缓存 key |
| account_id | string/int | 业务账号 | 可选；同租户内用于定位 Flow |
| state | string | 插件自定义值 | 必填，回调时原样返回 |
| status | enum(`pending`,`authorizing`,`succeeded`,`failed`) | Flow 当前状态 | 状态机严格递进，不可回退 |
| authorize_url | string | 登录入口 URL | Authenticator 生成，可能包含一次性参数 |
| expires_at | RFC3339 timestamp | Flow 过期时间 | TTL 到期即清理；查询返回 `expired` |
| metadata | map[string]string | 登录模式等附加信息 | 透传插件；用于 UI 展示 |
| result | CredentialPayload | 凭证结果 | 仅在 succeeded 时存在；存储脱敏版本 |
| last_error | string | 最近一次失败描述 | 包含错误码/摘要，脱敏 |
| callback_url | string | 插件回调地址 | 来自创建请求；用于签名派发 |
| retry_attempts | int | 回调已重试次数 | 与 3 次指数退避策略对齐 |

## CredentialPayload
| Field | Type | Description | Rules |
|-------|------|-------------|-------|
| session_token | string | 模拟登录产出的 token | 存储加密或脱敏版本，仅日志打印前后缀 |
| cookies_json | []Cookie | 需要同步给插件的 cookies | 可序列化为 JSON；字段包含 name/value/domain |
| headers_json | map[string]string | 需要注入的 header | 例如 `User-Agent`、`X-XSRF-TOKEN` |
| expires_at | RFC3339 timestamp | 凭证过期时间 | 由插件用于续期提醒 |
| captured_at | RFC3339 timestamp | 捕获时间 | 供审计/排查使用 |
| note | string | 登录方式/提示 | 例如 “手机号验证码登录” |

## SessionTokenConfig (per provider)
| Section | Fields | Purpose |
|---------|--------|---------|
| Service | `base_url`, `api_token`, `timeout` | 指向 SessionToken 服务与 API 鉴权 |
| Authenticator | `entries[]`(pc/h5/qrcode)、`default_user_agent`, `script_ids`, `captcha_strategy` | 配置登录入口及脚本 |
| Harvester | `watch_cookies[]`, `watch_headers[]`, `harvest_script_id` | 指示需要抓取的字段与脚本 |
| Callback | `callback_secret`, `max_retry` (可覆盖默认 3)、`retry_backoff` | 回调签名与重试策略 |
| Network | `proxy_pool`, `ip_strategy`, `request_timeout` | 代理/网络层控制 |

## FlowStateHistory (optional audit stream)
| Field | Type | Description |
|-------|------|-------------|
| flow_id | string | 关联 Flow |
| from_status | enum | 原状态 |
| to_status | enum | 新状态 |
| changed_at | RFC3339 timestamp | 变更时间 |
| reason | string | 触发原因（超时/用户操作） |

> FlowStateHistory 可存入 Redis Stream 或外部日志，用于合规审计；不是主接口的强制字段，但建议在实现阶段补充。
