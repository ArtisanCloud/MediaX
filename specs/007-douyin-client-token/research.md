# Research Notes - DouYin ClientToken Mode Onboarding

## Decision 1: Token 自动刷新触发时机
- **Decision**: 读取或使用 client_token 前检测 TTL，低于 `refresh_before_seconds` 立即刷新。
- **Rationale**: 与当前 AccessToken 流程一致，可避免后台守护进程及多副本竞争，并确保调试台/CLI 均获取最新凭证。
- **Alternatives considered**:
  - 定时任务批量刷新：需要额外调度器，且难以兼顾多租户使用率。
  - 纯手动刷新：容易遗漏，无法满足无人值守任务。

## Decision 2: 缓存后端优先级
- **Decision**: 默认 Redis（key `clientToken:douyin:<client_key>`）并记录 TTL，Redis 不可达时自动降级进程内缓存，UI/日志必须高亮 `storage_backend`。
- **Rationale**: Redis 能支持多实例共享和 TTL 续期；降级确保本地调试不中断，同时可提示风险。
- **Alternatives considered**:
  - 仅支持 Redis：阻碍本地无 Redis 环境。
  - 自建文件缓存：牵涉权限同步与清理复杂度。

## Decision 3: 调试日志与脱敏策略
- **Decision**: 所有刷新/调用日志通过 `client.Logger.WithContext` 输出结构化字段（`event`,`provider`,`action`,`redis_key`,`ttl_remaining`），token 仅保留前后 4 位并使用 `mask` 包。
- **Rationale**: 满足宪章的 Observability 要求并与 AccessToken 模块一致，便于排障与审计。
- **Alternatives considered**:
  - 只记录原始 JSON：存在敏感信息泄露风险。
  - 仅在调试台显示：缺乏统一日志链路，难以远程排障。
