# Research Notes: Google YouTube AccessToken Enablement

## Decision 1: AccessToken 提供方式
- **Decision**: 默认仍由 `GoogleYouTubeConfig.GetOAuthToken` 回调提供 AccessToken，可选通过 CLI `--access-token` 或环境变量 `GOOGLE_YOUTUBE_ACCESS_TOKEN` 覆盖。
- **Rationale**: 回调与 `GoogleAccessTokenHandler` 已在 MediaX 体系内验证；CLI 仅作为临时注入层，不必重写刷新逻辑。
- **Alternatives considered**:
  - 直接在 CLI 内实现 refresh_token 流程：增加代码与凭证存储风险，且与现有 handler 重复。
  - 强制所有用户依赖 Redis 缓存：会阻碍快速验证场景。

## Decision 2: CLI Action 范围
- **Decision**: 首批 action 覆盖 `videos.list`, `search.list`, `playlists.list`，并支持参数映射（part/ids/query/channel-id/max-results/page-token/mine）。
- **Rationale**: 这些 API 覆盖最常见的调试需求（查看视频/搜索结果/播放列表），且对应 `pkg/client/...` 已实现的方法，易于维护。
- **Alternatives considered**:
  - 覆盖全部 20+ 子客户端：初期维护成本过高，且调试门槛上升。
  - 只实现 `videos.list`: 无法满足搜索/播放列表调试。

## Decision 3: 文档/Quickstart 输出
- **Decision**: 在 specs 中产出 `data-model.md`（配置/CLI 参数映射）、`quickstart.md`（步骤）、`contracts/cli.md`（命令、输入/输出示例），并同步 `docs/develop/access-token/google`。
- **Rationale**: 规范化文档结构，方便后续自动生成 README 或 onboarding 内容。
- **Alternatives considered**:
  - 仅更新 README：难以复用且与 Speckit 流程不对齐。
  - 没有契约文档：CLI 参数与输出可能在团队内口口相传，不利于维护。
