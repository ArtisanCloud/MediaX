# Research: SessionTokenClient Provider Architecture

## Decision 1: Flow 持久化方式
- **Decision**: 以 Redis (go-redis v9) 作为唯一主存储，依赖 TTL/索引实现 pending→authorizing→succeeded/failed 生命周期，并支持多实例共享。
- **Rationale**: Redis 已在 MediaXCore 中作为默认缓存依赖；支持毫秒级 TTL、原子操作与 Lua/事务，便于实现幂等写入与状态恢复。无需额外 DB schema 即可满足 Flow 数量级需求。
- **Alternatives considered**: 关系型数据库（需要维护额外 schema/迁移）；进程内存（无法跨实例共享，丢失重启状态）。

## Decision 2: 回调签名与重试策略
- **Decision**: 所有 Callback payload 使用 HMAC-SHA256(state+payload+timestamp+nonce)；失败时固定 3 次指数退避（2s→4s→8s），超出即标记 Flow failed 并记录 last_error。
- **Rationale**: 符合宪章的安全要求，保障插件可验证来源；有限次数的指数退避能兼顾可靠性与避免压垮插件 API。
- **Alternatives considered**: 无限重试（风险放大，难以回收）；线性延迟或人工触发（SLA 不稳定）。

## Decision 3: Provider 配置与工厂注入
- **Decision**: 在 `pkg/client/config/<provider>.go` 增加 `SessionTokenConfig`（Service/Auth/Harvester/Callback/Network），并由 `MediaX.Create<Provider>SessionTokenClient` 注入到 adapter 内；配置值全部来源于环境变量或 secret 管理。
- **Rationale**: 与宪章的 Config-Layered Security 一致，方便未来接入其他 provider，且保持工厂注入缓存/Logger 的模式。
- **Alternatives considered**: 全局 SessionTokenConfig（无法表达 provider 差异）；在 adapter 中硬编码脚本/代理常量（违反配置规范）。
