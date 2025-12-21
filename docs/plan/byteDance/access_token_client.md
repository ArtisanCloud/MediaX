# ByteDance DouYin AccessToken Client 规划

> 目标：在 MediaX 中以“与 Google/BiliBili 相同的方式”交付抖音 AccessToken 能力，让 `cmd/accesstoken`、CLI 与外部服务能够统一加载 DouYin Provider、完成 OAuth / Flow 管理并调用 SDK 中的各类接口。

## 1. 关键目录

| 路径 | 说明 |
| --- | --- |
| `pkg/client/byteDance/douYin/accessTokenClient` | DouYin AccessToken SDK，涵盖内容发布、粉丝数据、IM、活动等模块。 |
| `pkg/client/config/douyin.go` | `ByteDanceDouYinConfig` 定义，继承 `ClientConfig`，配置 `client_id`/`client_secret`、代理、超时等。 |
| `cmd/accesstoken/server` | Web/HTTP 调试台，Provider 列表来自 `config.yaml`，需要在此注册 DouYin Provider 元数据。 |
| `docs/develop/access-token/byteDance/*.md` | 面向 SDK/运营的开发与调试指南（本文档要求的输出）。 |
| `config.example.yaml` | AccessToken Provider 示范配置，需要补充 DouYin 片段与环境变量说明。 |

## 2. 功能矩阵（AccessToken 客户端）

| 域 | 模块 | 说明 |
| --- | --- | --- |
| 内容管理 | `content/video`, `content/task`, `content/activity` | 处理视频上传、任务、活动及素材任务盒子。 |
| 搜索/互动 | `search`, `connection/fan`, `connection/fanProfile`, `connection/data` | 提供粉丝画像、互动统计、搜索投放等接口。 |
| 即时通讯 | `im/message`, `im/group`, `im/tool/*` | 支持企业私信、群聊、小程序模板、留存卡片。 |
| 市场/票据 | `market/service`, `tools/ticket` | 申请票据、调用市场服务。 |
| OAuth | `oauth` | 负责授权 URL、换取 AccessToken、刷新逻辑。 |

## 3. 配置约定

```yaml
access_token_providers:
  providers:
    - code: "byte_dance"
      name: "字节跳动"
      apps:
        - code: "douyin"
          name: "抖音 Douyin"
          provider_code: "byte_dance_douyin"
          api_version: "v1"
          auth_modes:
            - key: "default"
              label: "默认"
              byte_dance_douyin_config:
                api_url: "https://open.douyin.com"
                timeout: 5
                http_debug: false
                oauth:
                  client_id: "${DOUYIN_CLIENT_ID}"
                  client_secret: "${DOUYIN_CLIENT_SECRET}"
                  redirect_url: "${DOUYIN_REDIRECT_URL:-http://127.0.0.1:7071/debug/callback}"
                  scope: "user_info,video.create"
                oauth_key: "${DOUYIN_OAUTH_KEY:-sandbox}"
```

- 所有敏感字段必须经由环境变量注入；文档需说明如何在 `.env` 中准备 `DOUYIN_CLIENT_ID/SECRET/SCOPE`。
- 若未来开放多租户，只需追加新的 `auth_modes`。

## 4. 集成步骤

1. **Provider wiring**：在 `cmd/accesstoken/server/server.go` 加入 DouYin Provider 元数据（若尚未存在），确保 `/debug` 能读取 `byte_dance` 分组与 `douyin` App。完成后在 `docs/develop/access-token/byteDance/develop.md` 记录使用方法。
2. **OAuth 流程**：复用现有 `handleOAuthStart` / `completeOAuthFlow`，检查 DouYin 是否需要额外参数（state、scope、device_id 等）。若需要自定义 `authorize_url`，在 `buildOAuthAuthorizeURL` 中根据 `provider_code` 定制。
3. **Flow/Redis**：沿用 `accesstoken:oauth:<provider_code>:<app>:<mode>` 命名；文档中示例 Provider Code 应为 `byte_dance_douyin`，方便排障。
4. **调试文档**：新增 `develop.md`（环境、配置、Go/CLI 示例）与 `debug.md`（`/debug` 页面操作、常见问题），格式参考 google/redbook 文档。
5. **后续扩展**：等 AccessToken API 调试面板开放 DouYin Action 时，再补充 `/accesstoken/call` 的 action presets；本阶段重点在 OAuth/Token 解析。

## 5. 下一步

1. 补齐 `specs/00x` 流程：创建相应 spec/tasks 追踪 DouYin Provider 入驻。
2. 根据 DouYin OAuth 要求补充回调校验逻辑（例如 `state` 需包含 `device_id`）。
3. 评估是否需要将 DouYin access token 回写到外部缓存/Studio，以便多端复用。

## 6. 交付物 & 参考文档

- **Quickstart**：`specs/006-bytedance-access-token/quickstart.md`，涵盖端到端的配置→OAuth→`/accesstoken/call` 调试→排障 FAQ。
- **开发指南**：`docs/develop/access-token/byteDance/develop.md`，列出所有环境变量、`config.yaml` 示例、Need Reauth & 自动刷新排障步骤、Go/CLI 调用示例。
- **调试指南**：`docs/develop/access-token/byteDance/debug.md`，说明如何在 `/debug` 中同步模板、发起授权、解析 Flow、查看 `retry_count`/`last_backoff_ms`、处理 Redis 降级/速率限制。
- **README 引用**：顶层 README 的 AccessToken 章节增加了 DouYin Quickstart/文档链接，方便团队统一入口。
- **日志/授权契约**：`cmd/accesstoken/server/logging.go` 保证 `token.call` 日志包含 `provider_code/provider_app/action/flow_id/token_source/retry_count/storage_backend`，满足 Observability 的守则。
- **回归验证**：`go test ./cmd/accesstoken/...` 会覆盖 Google、BiliBili、RedBook 现有合约（`oauth_contract_test.go`, `token_contract_test.go`, `flow_contract_test.go`），持续集成时确保这些 Provider 行为未被 DouYin 修改影响。

> 交付上述计划后，接入方可以复用 MediaX 的 AccessToken 基建（UI、Redis、日志、Flow 管理、自动刷新/限流），仅需在 `/accesstoken/call` 中扩展对应 action 即可完成 DouYin API 联调。
