# BiliBili AccessToken Client 现状与计划

## 1. 背景
- MediaX 在 `pkg/client/bilibili` 下已经实现面向 B 站开放平台（OAuth2 AccessToken）的业务客户端，覆盖视频投稿、专栏、数据、用户、直播等 API。
- 目前这些能力只暴露为 Go SDK，尚未像 WeChat ClientToken/SessionToken 那样接入统一的服务端、CLI、调试文档，因此需要整理接口覆盖范围，评估后续如何对接 HTTP Server/调试台。

## 2. 代码结构
| 目录 | 说明 |
| --- | --- |
| `pkg/client/bilibili/core` | 封装 `kernel.BaseClient`、`BiliBiliAccessTokenHandler`，负责 AccessToken 获取与刷新，默认 `grant_type=authorization_code`。|
| `pkg/client/bilibili/b/accessTokenClient` | 业务聚合层，`BiliBiliACClient` 负责实例化视频/文章/数据/直播等子客户端。|
| `pkg/client/config` | 含 `BiliBiliConfig`、`ClientConfig`，提供 API/授权地址、ClientID/Secret、Proxy、缓存等配置。|
| `docs/develop/...`（缺失） | 尚无 B 站开发/调试文档，后续需要与 ClientToken/SessionToken 系列保持一致。|

### 2.1 AccessToken 处理
- `core.NewBiliBiliClient`：复用 `kernel.NewBaseClient`，启动时 `OverrideCheckTokenNeedRefresh()`，保证所有 `HttpGet/Post/Upload` 会自动带 Token。
- `core.BiliBiliAccessTokenHandler`：封装 `kernel.TokenHandler`，覆盖 `GetCredentials` -> `{grant_type=authorization_code, client_id, client_secret}`，并预留 `GetCustomToken` 供外部注入 OAuth 结果。
- `BiliBiliACClient`：
  - 暴露 `GetVideoClient`/`GetUserClient`/`GetArticleClient`/`GetDataClient`/`GetLiveClient`/`GetLiveWSClient`/`GetLiveThirdPartyClient` 工厂。
  - 每个子客户端都共享 `BaseClient`，因此具备统一的鉴权、签名、重试、日志能力。

## 3. 已实现接口
> 以下内容来自 `pkg/client/bilibili/b/accessTokenClient/*`，罗列所有已封装的方法与对应的开放接口，便于后续挂载 HTTP/CLI。

### 3.1 视频稿件（`video`）
| 方法 | API 路径 | 功能 |
| --- | --- | --- |
| `Edit` | `POST /arcopen/fn/archive/edit` | 编辑稿件标题/分区/描述等。|
| `Delete` | `POST /arcopen/fn/archive/delete` | 删除稿件。|
| `PreUpload` | `POST /arcopen/fn/archive/video/init` | 获取上传 token，准备分片或小视频流程。|
| `UploadPart` | `POST /arcopen/fn/archive/video/part` | 分片上传视频。|
| `UploadComplete` | `POST /arcopen/fn/archive/video/complete` | 合并分片。|
| `CoverUpload` | `POST /arcopen/fn/archive/cover/upload` | 上传视频封面（multipart）。|
| `SubmitVideo` | `POST /arcopen/fn/archive/add-by-utoken` | 使用 `upload_token` 投稿。|
| `UploadSmallVideo` | `POST https://openupos.bilivideo.com/video/v2/part/upload` | 单文件上传（≤100MB）。|
| `AddShareUrl` | `POST /arcopen/fn/resource/add_share` | 生成客户端投稿唤起链接。|
| `GetVideoInfo` | `GET /arcopen/fn/archive/view` | 查询单条稿件详情。|
| `GetVideoList` | `GET /arcopen/fn/archive/viewlist` | 拉取稿件列表（分页）。|

### 3.2 专栏稿件（`article`）
| 方法 | API 路径 | 功能 |
| --- | --- | --- |
| `Add` | `POST /arcopen/fn/article/add` | 新增专栏稿件。|
| `Edit` | `POST /arcopen/fn/article/edit` | 编辑专栏。|
| `Delete` | `POST /arcopen/fn/article/delete` | 删除文章（支持批量）。|
| `Detail` | `POST /arcopen/fn/article/detail` | 拉取文章详情。|
| `List` | `GET /arcopen/fn/article/list` | 列出当前账号文章。|
| `Categories` | `GET /arcopen/fn/article/categories` | 查询分类。|
| `Card` | `GET /arcopen/fn/article/cards` | 获取视频/文章卡片片段。|
| `UploadImage` | `POST /arcopen/fn/article/upload/image` | 上传专栏图片（支持水印）。|

#### 文集管理（`article/collection`）
| 方法 | API 路径 | 功能 |
| --- | --- | --- |
| `Edit` | `POST /arcopen/fn/article/anthology/edit` | 编辑文集信息。|
| `EditList` | `POST /arcopen/fn/article/belong` | 调整文集所含文章。|
| `Delete` | `POST /arcopen/fn/article/anthology/delete` | 删除文集。|
| `List` | `GET /arcopen/fn/article/anthology/list` | 查询文集列表。|
| `Detail` | `GET /arcopen/fn/article/anthology/detail` | 查询文集详情。|

### 3.3 数据中心（`data`）
| 方法 | API 路径 | 功能 |
| --- | --- | --- |
| `GetUserStat` | `GET /arcopen/fn/data/user/stat` | 关注/粉丝/投稿等用户指标。|
| `GetVideoStat` | `GET /arcopen/fn/data/arc/stat` | 单稿件播放/弹幕/互动。|
| `GetVideoStatIncrement` | `GET /arcopen/fn/data/arc/inc-stats` | 投后整体增量指标。|
| `GetColumnStat` | `GET /arcopen/fn/data/art/stat` | 指定专栏数据。|
| `GetColumnStatIncrement` | `GET /arcopen/fn/data/art/inc-stats` | 专栏整体增量。|

### 3.4 用户（`user`）
| 方法 | API 路径 | 功能 |
| --- | --- | --- |
| `AccountScopes` | `GET /arcopen/fn/user/account/scopes` | 查询已授权的权限点。|
| `GetUserInfo` | `GET /arcopen/fn/user/account/info` | 获取用户头像/昵称/OpenID。|
| `GetUserUnionId` | `GET /arcopen/fn/user/account/union_id` | 查询 union_id。|

### 3.5 直播（`live`、`live/ws`、`live/thirdParty`）
| 子模块 | 方法 | API 路径 | 功能 |
| --- | --- | --- | --- |
| `live` | `GetRoomInfo` | `GET /arcopen/fn/live/room/info` | 查询主播房间 ID、标题、是否开播。|
| `live/ws` | `Start` | `POST /arcopen/fn/live/room/ws-start` | 获取长链/ConnID。|
| `live/ws` | `HeartBeat` | `POST /arcopen/fn/live/room/ws-heartbeat` | 单链接心跳。|
| `live/ws` | `BatchHeartBeat` | `POST /arcopen/fn/live/room/ws-batch-heartbeat` | 批量心跳。|
| `live/thirdParty` | `GetLiveGrantUrl` | `GET /liveopen/fn/live/thirdPartyLive/grantUrl` | 生成第三方开播授权链接。|
| `live/thirdParty` | `StartLive` | `POST /liveopen/fn/live/thirdPartyLive/startLive` | 开始第三方直播，获取推流地址。|
| `live/thirdParty` | `EndLive` | `POST /liveopen/fn/live/thirdPartyLive/endLive` | 结束第三方直播。|

## 4. 现状评估
- ✅ **接口覆盖**：B 站内容/数据/直播三大板块均具备基础 CRUD/查询能力；上传流程（视频/图）均以 `HttpUpload` 实现，可直接复用。
- ✅ **Token 管理**：依赖 `kernel.TokenHandler`，默认自动刷新，但当前 `grant_type` 固定为 `authorization_code`，需要结合业务 Flow 注入 `authorization_code` 或自定义 Token。
- ⚠️ **调试入口缺失**：没有 `cmd/bilibili` 服务、脚本或 `/debug` 页面；外部系统无法直接复用这些 SDK。
- ⚠️ **测试覆盖缺乏**：`pkg/client/bilibili` 下无单元/集成测试，也无 mock server，后续扩展需补充。
- ⚠️ **开放接口清单维护成本高**：目前所有 API 均散落在 GoDoc 注释中，建议把上述清单同步到 docs/develop/、规格/任务中，便于查看。

## 5. 后续建议
1. **输出服务端**：比照 `cmd/clienttoken/server`，新增 `cmd/accesstoken-bilibili`（或纳入现有 AccessToken Server）来调用上述 SDK，满足外部 HTTP 请求、Redis 缓存、报错透传。
2. **补齐文档**：
   - `docs/develop/access-token/bilibili/develop.md`：面向外部的调用说明。
   - `docs/develop/access-token/bilibili/debug.md`：本地调试与脚本指南。
   - 在本文件基础上维护接口矩阵，新增接口时同步更新。
3. **自动化测试**：为关键上传/投稿流程引入 mock server 或录制数据，至少覆盖序列化、签名、错误透传。
4. **配置统一**：在 `config.example.yaml` 中补充 `bilibili` provider 示例，指定 `client_id/secret`, `redirect_uri`, `scopes`，方便接入。
5. **观测指标**：对齐 `clienttoken_metric` 规范，在未来的 HTTP 服务里输出 `bilibili_metric`，记录 action/provider_app/status/latency，辅助排障。

> 本文档会随 `pkg/client/bilibili` 更新同步维护，后续新增接口时需在「已实现接口」章节追加条目。
