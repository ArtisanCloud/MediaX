# Data Model - DouYin ClientToken Mode Onboarding

## Entity: DouYinClientTokenProvider
- **Description**: 配置项描述字节跳动/抖音 ClientToken Provider → App → Mode 组合。
- **Fields**:
  - `provider_code` (string, required, pattern `byte_dance_douyin_clienttoken`)
  - `app_code` (string, required, default `douyin_service`)
  - `auth_mode` (string, required, default `default`)
  - `byte_dance_douyin_config` (object)
    - `api_url` (string, default `https://open.douyin.com`)
    - `timeout` (int, seconds, default 5)
    - `http_debug` (bool)
    - `client_token.client_key` (string, env)
    - `client_token.client_secret` (string, env)
    - `cache.redis_key` (string, format `clientToken:douyin:<client_key>`)
    - `cache.ttl_seconds` (int)
    - `cache.refresh_before_seconds` (int)
    - `redis.addr` (string, env)
    - `redis.db` (int, env)
    - `redis.pass` (string, env, optional)
    - `api_token` (string, env, protects HTTP API)
- **Relationships**: Injected到 `cmd/clienttoken/server` Provider 列表，与缓存记录 1:n（一个 Provider 多个 token 记录）。

## Entity: ClientTokenCacheRecord
- **Description**: 存储 DouYin client_token 当前值与元数据的缓存结构。
- **Fields**:
  - `redis_key` (string, key)
  - `token_value_masked` (string, masked token)
  - `ttl_seconds` (int, DouYin 实际 TTL)
  - `refresh_before_seconds` (int)
  - `source` (enum: redis/memory)
  - `updated_at` (timestamp)
  - `provider_code` (string)
  - `client_key` (string)
  - `last_refresh_reason` (enum: manual/threshold/api_call)
  - `error` (string, optional for失败)
- **Relationships**: 与 Provider 关联 1:1，每个 client_key 仅一个有效记录；调试台/CLI 读取该实体。

## Entity: ClientTokenDebugAction
- **Description**: 代表一次刷新、查看或 API 调试操作。
- **Fields**:
  - `action_id` (uuid)
  - `action_type` (enum: refresh/cache_view/api_call/cache_clear)
  - `provider_code` (string)
  - `action` (string, DouYin API path, 可空)
  - `request_payload` (json, masked)
  - `result_status` (enum: success/error)
  - `error_hint` (string, optional)
  - `token_source` (enum: redis/memory)
  - `ttl_remaining` (int)
  - `created_at` (timestamp)
- **Relationships**: 归属于操作人会话/tenant，用于日志与调试历史。
