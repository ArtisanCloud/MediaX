# ByteDance DouYin ClientToken Client 规划

> 目标：沉淀字节跳动抖音 ClientToken 能力的入驻路径，使 `cmd/clienttoken/server`、CLI 与三方插件能复用统一配置、Redis 缓存与调试文档，与 WeChat ClientToken 体验保持一致。

## 1. 关键目录

| 路径 | 说明 |
| --- | --- |
| `pkg/client/byteDance/douYin/clientTokenClient` | 抖音 ClientToken SDK，包含视频、任务、模板、活动、MicApp、沙箱、票据等子模块。 |
| `pkg/client/config/douyin.go` | `ByteDanceDouYinConfig`，用于指定 API URL、OAuth client、缓存策略。 |
| `cmd/clienttoken/server` | ClientToken 调试台与 HTTP API，Provider/App/Mode 来自 `client_token_providers` 配置。 |
| `docs/develop/client-token/byteDance/*.md` | 本文档要求新增的开发/调试指南。 |
| `config.example.yaml` | ClientToken Provider 示例，需要追加 DouYin 节点与环境变量说明。 |

## 2. 功能矩阵（ClientToken 客户端）

| 域 | 模块 | 说明 |
| --- | --- | --- |
| 内容管理 | `content/video`, `content/task`, `content/activity`, `content/schemas` | 适用于不涉及用户授权的后台任务，如批量上传、模板配置。 |
| 搜索/开放工具 | `search`, `tools/micApp`, `tools/sandbox`, `tools/ticket` | 包含小程序、沙箱、票据等后台工具能力。 |
| Token 管理 | `core.ByteDanceTokenHandler` | 负责 client_token 的刷新、缓存、TTL 管控。 |

## 3. 配置约定

```yaml
client_token_providers:
  providers:
    - code: "byte_dance"
      name: "字节跳动"
      apps:
        - code: "douyin_service"
          name: "抖音服务端"
          provider_code: "byte_dance_douyin_clienttoken"
          auth_modes:
            - key: "default"
              label: "默认"
              byte_dance_douyin_config:
                api_url: "https://open.douyin.com"
                timeout: 5
                http_debug: false
                client_token:
                  client_key: "${DOUYIN_CLIENT_KEY}"
                  client_secret: "${DOUYIN_CLIENT_SECRET}"
                # Risk 字段仅在业务需要时启用，否则调试台会默认禁用相关 API 模板
                device_id: "${DOUYIN_DEVICE_ID:-}"
                risk_info: "${DOUYIN_RISK_INFO:-}"
                cache:
                  redis_key: "clientToken:douyin:${DOUYIN_CLIENT_KEY}"
                  ttl_seconds: 7000
                redis:
                  addr: "${CLIENTTOKEN_REDIS_ADDR:-127.0.0.1:6379}"
                  db: "${CLIENTTOKEN_REDIS_DB:-0}"
                api_token: "${CLIENTTOKEN_API_TOKEN:-dev-clienttoken}"
```

- `client_key/client_secret` 必须来自环境变量，禁止硬编码。
- `redis_key` 与 `ttl_seconds` 用于控制缓存策略，需在文档中说明如何在调试台查看，并确保当 TTL 低于 `refresh_before_seconds` 时触发刷新。
- 当 API 需要 `device_id`/`risk_info` 等额外字段时，默认通过调试台检测并提示；只有两个字段同时存在时才启用 Douyin API 模板，避免误调风控接口。

## 4. 集成步骤

1. **Provider wiring**：更新 `config.example.yaml` 与 `cmd/clienttoken/server` 的 Provider 加载逻辑，保证 “字节跳动 / 抖音服务端” 能出现在 `/client-token/debug` 下拉中。
2. **文档**：新增 `docs/develop/client-token/byteDance/develop.md`（环境、配置、CLI/Go 用例）与 `docs/develop/client-token/byteDance/debug.md`（调试台操作、常见问题），格式参考微信 ClientToken 文档。
3. **Redis/缓存**：文档需解释缓存键命名（如 `clientToken:douyin:<client_key>`）、TTL、刷新策略；调试台需支持“刷新 token / 查看缓存 / 调用 API”三段流程。
4. **CLI/Scripts**：提供 `scripts/clienttoken-douyin.sh`，覆盖 refresh/cache/call 场景，复用 `CLIENTTOKEN_API_TOKEN`，便于 CI 或命令行诊断。
5. **外部依赖**：若 DouYin ClientToken 接口需要 `device_id`、`risk_info` 等额外字段，应在计划中列出并于文档中解释。

## 5. 下一步

1. 在 `specs/` 下建立对应的 feature plan（例如 `006-douyin-client-token`），确保后续修改受控。
2. 在调试台 UI 中增加 DouYin ClientToken 的 API 模板（如 `content/video/upload`），方便 QA 与集成商验证。
3. 保障 `cmd/clienttoken/server` 日志输出 `clienttoken_event`，字段覆盖 provider_code/redis_key/ttl_remaining/storage_backend/api_token_subject，方便排障。

## 6. 验证记录

- **2025-02-14**：按照 `docs/develop/client-token/byteDance/develop.md` / Quickstart 步骤演练，`go test ./cmd/clienttoken/... ./pkg/client/byteDance/douYin/...` 均通过；`scripts/clienttoken-douyin.sh refresh|cache|call` 正常运行，日志输出包含 `clienttoken_event` 字段。
- 缺省未提供 `DOUYIN_DEVICE_ID/DOUYIN_RISK_INFO` 时，调试台会提示禁用 Douyin API 模板，确保风险字段显式配置后再开启。

> 有了本计划与配套文档，运维与合作方可以直接复制 `config.example.yaml` 段落、启动 `cmd/clienttoken/server` 并在 `/client-token/debug` 完成 Token 刷新/接口联调，避免重复造轮子。
