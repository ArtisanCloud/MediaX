# Research Log - DouYin AccessToken Debug Integration

## Decision 1 - 自动刷新策略
- **Decision**: DouYin Flow 在 `/accesstoken/call` 前检测 access_token 剩余 ≤5 分钟或已过期时自动使用 refresh_token 刷新，并更新缓存/Redis。
- **Rationale**: 调试 CLI/脚本复用 Flow，需要在无人工干预下长时间保持可用；DouYin 官方返回 refresh_token，可直接借助 `ByteDanceDouYinACClient` 调用刷新端点，避免重复授权。
- **Alternatives considered**: 仅提示重新授权（会频繁中断调试）；后台定时刷新（无请求也占用配额，且无法把失败结果立刻反馈给调用者）。

## Decision 2 - Flow TTL & 失效处理
- **Decision**: Flow TTL 与 `oauthTokenRecord.ExpiresIn` 同步（等于 DouYin `expires_in`），刷新成功后重置 TTL；刷新失败立即删除 Flow 并返回 `need reauth`。
- **Rationale**: TTL 与 AccessToken 生命周期一致可避免 UI 展示过期数据；失败即删除能防止调用者继续拿失效 token。
- **Alternatives considered**: 固定 24h TTL（会让 UI 看到已过期 Flow）；失败时保留旧 Flow 并标记（仍需要手工清理，复杂度更高）。

## Decision 3 - 速率限制与重试
- **Decision**: 对 DouYin action 应用 per-action 1 QPS 速率限制，并在遇到 429/5xx 时进行最多 3 次指数退避重试。
- **Rationale**: 调试流量虽小，但 DouYin API 的风控敏感，加入轻量速率限制能防止操作失误造成封禁；有限次重试兼顾可用性与透明度。
- **Alternatives considered**: 无速率限制（风控风险高）；固定间隔重试（容易造成雪崩）；更低 QPS（会让多人联合调试时体验下降）。

## Decision 4 - 单应用配置
- **Decision**: `byte_dance_douyin_config` 仅支持单个 app，`<app>=default`，通过文档提示如需多 app 需启动多实例。
- **Rationale**: 现有 UI/Flow 命名只暴露单 app 选择；多 app 会引入大量配置分支且超出当前调试需求。
- **Alternatives considered**: 支持数组（需要改动 UI、Flow 键名、文档）；依赖 CLI 热切换（复杂度与可维护性不佳）。
