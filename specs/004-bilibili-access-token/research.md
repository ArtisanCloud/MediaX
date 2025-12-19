# Research Log

## MediaXCore BaseClient & TokenHandler
- **Decision**: 调试服务中的 OAuth 回调、`/accesstoken/token` 解析仍通过 `github.com/ArtisanCloud/MediaXCore/kernel.BaseClient` 与 `TokenHandler` 提供的 HTTP/缓存逻辑。
- **Rationale**: BaseClient 已内建日志、重试、脱敏字段，且与生产 AccessToken/SessionToken 流程保持一致，可复用刷新检测与 `GetOAuthToken` hook；减少编写重复 HTTP/签名代码。
- **Alternatives considered**: 直接在调试服务里使用 `net/http` 手写 OAuth 请求（缺乏统一日志与构建器，违背宪章 I/III）以及新建轻量客户端（需要重新实现缓存/Hook, 成本更高），因此拒绝。

## go-redis/v9 Flow 持久化
- **Decision**: 若设置 `ACCESSTOKEN_REDIS_ADDR`，使用共享 `redis.Client`（db、密码、TLS 均来源 env）写入 `accesstoken:oauth:*` key，并为 Flow 与最新记录统一设置 24h TTL；Redis 不可用时回退到进程内 map，并在日志/UI 明示数据易失。
- **Rationale**: Redis 可跨进程协作并允许 CLI/外部工具直接读取；统一 TTL 便于支持回放与自动清理；go-redis v9 提供 context-aware API 与自动连接池，符合现有模块。
- **Alternatives considered**: 新增 SQLite/文件持久化（增加部署复杂度，不满足“调试即可启动”目标）或无限 TTL（违反安全/隐私要求，导致敏感 token 长期存在），因此拒绝。

## YAML Provider 配置与监听地址
- **Decision**: 继续依赖 `config.yaml` 中的 `access_token_providers`（仅包含 `${ENV}` 占位）驱动 UI 模板；调试服务默认监听 `127.0.0.1:7071`，仅在显式设置 `ACCESSTOKEN_LISTEN_ADDR` 时暴露其它地址。
- **Rationale**: 现有 YAML schema 已被其它 Provider 复用，保持结构即可保证 parity；默认本地监听配合 `ACCESSTOKEN_API_TOKEN` 防止意外暴露，在需要共享时 operator 可自主 override。
- **Alternatives considered**: 为 BiliBili 单独定义 JSON 模板或默认监听 `0.0.0.0`。前者造成配置重复与维护成本；后者在未设防火墙时可能泄露调试页面，不符合 Config-Layered Security。

## Flow TTL 与可观测性
- **Decision**: Flow 记录（Redis 与内存）统一 TTL = 24 小时 (86,400 秒)，日志中追加 `flow_ttl_seconds`、`storage_backend`、`listen_addr` 字段以帮助支持团队判断持久化状态。
- **Rationale**: 24h 覆盖常见排障周期又能限制敏感 token 暴露时间；额外日志字段确保易于追踪“数据为何丢失”并满足宪章 Observability 要求。
- **Alternatives considered**: 与 BiliBili `expires_in` 完全对齐的动态 TTL（过短，不利于回放）或无限 TTL（安全风险），因此拒绝；日志保持旧字段（缺少监听/存储信息，难诊断）亦被拒绝。
