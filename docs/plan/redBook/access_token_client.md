# RedBook JuGuang AccessTokenClient 规划

## 1. 背景与目标
小红书聚光开放平台提供了广告投放、笔记资产、数据报表等接口。MediaX 已在 `pkg/client/redBook/juGuang/accessTokenClient` 下实现了对应的 AccessToken 客户端。相比之下，Google（`google_youtube`）与 BiliBili（`bilbili`）已经完整接入 `cmd/accesstoken/server` 与 `/debug` 调试台，可以通过 `go run ./cmd/accesstoken/server -config config.yaml` 直接发起 OAuth、复用 Flow。本计划文档专门用于把聚光客户端“入住”同一套调试环境，目标不是另起炉灶，而是复用现有 cmd/accesstoken 流程，补齐配置/文档/接口路由后即可像 Google/BiliBili 那样在线调试。重点说明如下：

- 盘点现有代码结构与功能域，便于产品、研发快速对齐能力覆盖范围；
- 说明授权模式与 AccessToken 处理方式，指导如何通过 `docs/develop/access-token/*` 的调试服务复用这些接口；
- 梳理后续需要完成的文档、调试入口与自动化测试工作，确保聚光能力可以像 Google/BiliBili 一样对外开放。

## 2. 代码结构与依赖

| 目录 | 说明 |
| --- | --- |
| `pkg/client/redBook/core` | 复用 `kernel.BaseClient` 与 AccessTokenHandler，封装请求签名、重试、日志能力。 |
| `pkg/client/redBook/juGuang/accessTokenClient/client.go` | `RedBookJuGuangACClient` 聚合器，负责实例化各业务子客户端并绑定 `RedBookJuGuangConfig`。 |
| `pkg/client/redBook/juGuang/accessTokenClient/account` | 账户纬度接口（余额等）。 |
| `…/dataReport/{realtime,offline}` | 聚光实时/离线报表客户端。 |
| `…/note` | 笔记资产操作、落地页、否定词、直达链接等。 |
| `…/promote/{campaign,unit,creativity}` | 推广计划/单元/创意的 CRUD 与状态切换。 |
| `…/tools` | 关键词推荐、行业属性、定向信息、操作记录等工具接口。 |
| `pkg/client/config/redbook.go` | `RedBookJuGuangConfig` 定义，含 `api_url`、`oauth`、代理等字段，可直接映射到 `config.yaml` 的 `redbook_juguang_config`。 |

### 2.1 AccessToken 处理
- `RedBookJuGuangACClient` 通过 `MediaX.CreateRedBookJuGuangACClient` 创建，底层继承 `kernel.BaseClient`，并透传 `cache.ICache` 用于 Token 或报表缓存。
- `RedBookJuGuangConfig` 支持 `oauth` 块（`client_id/client_secret/redirect_url/scope`），可与 AccessToken Server 的授权模式对齐；若需要注入外部 Token，可赋值 `GetOAuthToken`。
- 统一日志：依赖 `MediaXCore/pkg/logger`，所有 `HttpPost/HttpGet` 均会记录请求，用于与 `docs/develop/access-token` 的 `/debug` 页面联动。

### 2.2 调试入口关联
- `config.yaml` 中示例（`access_token_providers -> redbook -> apps -> auth_modes -> redbook_juguang_config`）可直接被 `cmd/accesstoken/server` 读取，调试台会展示 Provider/App/Mode/OAuth Key，与 Google/BiliBili 一致。
- 调试流程复用 `docs/develop/access-token/google/debug.md` 中的机制：通过 `/accesstoken/oauth/start` 触发 OAuth、`/debug/callback` 接收 code，然后使用 Flow ID 重放 Token。新增文档将指导如何把聚光 API 接口接入 `/accesstoken/call` 或独立 HTTP Handler。

## 3. 能力矩阵

### 3.1 账号与余额（`account`）
| 方法 | API | 功能 |
| --- | --- | --- |
| `GetAccountBalance` | `POST /api/open/jg/account/balance/info` | 查询广告主账户余额（冻结/可用/今日花费等指标）。 |

### 3.2 笔记资产（`note`）
覆盖笔记发布、直达链接、否定词、门店/POI、模版等能力：
| 方法示例 | API 路径 | 备注 |
| --- | --- | --- |
| `PostVideoNote` | `/api/open/jg/note/video/publish` | 发布视频笔记，支持素材 ID、投放参数。 |
| `DeleteNote` | `/api/open/jg/note/delete` | 删除笔记。 |
| `CreateDirectLink`/`EditDirectLink`/`DeleteDirectLink` | `/api/open/jg/directLink/*` | 维护聚光直达链接。 |
| `BatchAddNegativeWord`/`GetNegativeWordList` | `/api/open/jg/keyword/negative/*` | 否定词管理。 |
| `ListLandingPage`/`GetAssetInfo`/`GetSpuList`/`GetPOIList` | 多条 `/api/open/jg/*` | 查询落地页、事件资产、商品、门店。 |

### 3.3 推广投放（`promote`）
| 子域 | 核心方法 | API |
| --- | --- | --- |
| 计划（`campaign`） | `List` `Create` `Update` `UpdateStatus` | `/api/open/jg/ad/campaign/*` |
| 单元（`unit`） | `List` `Create` `Update` `UpdateStatus` | `/api/open/jg/ad/unit/*` |
| 创意（`creativity`） | `CreateNote` `CreateLandingPage` `CreateProgrammaticPage` `Update` `UpdateStatus` `Search` | `/api/open/jg/ad/creativity/*` |

### 3.4 数据报表（`dataReport`）
| 维度 | 实时 (`realtime`) | 离线 (`offline`) |
| --- | --- | --- |
| 账户/计划/单元/创意/定向 | `AccountLevel`、`CampaignLevel`、`UnitLevel`、`CreativeLevel`、`TargetLevel` | `AccountLevel`、`CampaignLevel`、`UnitLevel`、`GroupReport` |
| 功能 | `GET/POST /api/open/jg/data-report/*` | 同步/异步任务，返回流量、消耗、曝光等指标，便于 BI 系统拉取。 |

### 3.5 工具能力（`tools`）
| 方法 | API | 用途 |
| --- | --- | --- |
| `GetOperationRecord` | `/api/open/jg/tools/operation/record` | 查询操作日志。 |
| `KeywordRecommend`/`GetKeywordMatch`/`GetKeywordIndustry` | `/api/open/jg/keyword/*` | 关键词规划与行业匹配。 |
| `GetIndustryAttribute` | `/api/open/jg/tools/industry/attribute` | 拉取行业属性。 |
| `CrowdEstimate` | `/api/open/jg/crowd/estimate` | 人群预估。 |
| `GetWordBagList`/`GetTargetInfo` | `/api/open/jg/wordbag/list`、`/api/open/jg/target/info` | 辅助定向。 |
| `CheckDuplicatedName` | `/api/open/jg/tools/checkName` | 校验计划/单元重名。 |

## 4. 与 AccessToken 调试模式的结合
1. **配置**：在 `config.yaml` 中新增：
   ```yaml
   access_token_providers:
     providers:
       - code: "redbook"
         name: "小红书"
         apps:
           - code: "juguang"
             name: "聚光"
             provider_code: "redbook_juguang"
             auth_modes:
               - key: "default"
                 label: "默认 OAuth"
                 redbook_juguang_config:
                   api_url: "https://ad.xiaohongshu.com"
                   timeout: 10
                   http_debug: false
                   oauth:
                     oauth_url: "<授权URL>"
                     access_token_url: "<token URL>"
                     client_id: "${JUGUANG_CLIENT_ID}"
                     client_secret: "${JUGUANG_CLIENT_SECRET}"
                     redirect_url: "http://127.0.0.1:7071/debug/callback"
                     scope: "notes.read,notes.write,..."
   ```
2. **服务启动**：`go run ./cmd/accesstoken/server -config config.yaml` 后，调试台会出现 “小红书 / 聚光” Provider。点击“发起授权”→ 聚光 OAuth → 回调写入 Redis → Flow ID 可用于 CLI/脚本。
3. **API 调试**：`/accesstoken/call` 现在兼容 YouTube 与聚光。通过 `action: "redbook.account.balance"` + `payload.advertiser_id` 即可复用既有调试台，不需要额外的路由或页面。

## 5. 后续工作
1. **文档补齐**：新增 `docs/develop/access-token/redbook/develop.md`、`debug.md`，直接照搬 Google/BiliBili 的结构，重点写明如何在 `cmd/accesstoken/server` 中选择“小红书/聚光”并获取 Flow，避免独立脚本。
2. **调试接口**：在 `cmd/accesstoken/server`（与 YouTube/BiliBili 共享）里为 `redbook_juguang` 增设 `POST /accesstoken/call` 或 `/accesstoken/redbook/*` endpoints，让调试页表单可以直接调用聚光客户端。
3. **自动化测试**：为关键客户端（Note、Campaign、DataReport）添加 mock server 或 contract tests，保证 schema 与官方一致。
4. **配置示例与 Quickstart**：更新 `config.example.yaml`、`docs/develop/access-token/google/bilibili/...` 旁边的 quickstart，确保复制模板并执行 `go run ./cmd/accesstoken/server -config config.yaml` 即能在同一个调试台里切换 Google/BiliBili/RedBook。
5. **观测与限流**：对齐 `docs/plan/google/...` 的做法，在未来的服务端中输出 `redbook_metric`（provider_app/action/status/latency）并预留配额告警。

> 该文档会随 `pkg/client/redBook/juGuang/accessTokenClient` 演进而更新。新增接口、能力或配置项时，请同步维护本计划与开发/调试文档，确保外部团队始终能找到最新的接入指引。
