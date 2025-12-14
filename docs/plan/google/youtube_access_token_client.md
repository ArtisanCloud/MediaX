# Google YouTube AccessTokenClient 功能介绍

## 1. 背景与目标
MediaX 的 Google YouTube AccessTokenClient（`pkg/client/google/youtube/accessTokenClient`）封装了 YouTube Data API v3 的主要资源操作，统一处理 OAuth 访问令牌、日志与缓存策略，便于业务方在 MediaX SDK 内按需组合视频上传、播放列表维护、互动管理与频道装饰能力。本篇文档用于梳理现有子客户端的功能范围，方便产品/研发快速对齐可复用的接口以及潜在的补强点。

## 2. 客户端整体结构
1. `GoogleYouTubeACClient`（`client.go`）持有 `core.GoogleClient`、`GoogleYouTubeConfig` 与 `GoogleAccessTokenHandler`，实例化时会绑定 AccessTokenHandler 并允许传入自定义的 `cfg.GetOAuthToken`。
2. 通过 `GetXXXClient` 方法按需懒加载各功能域子客户端，所有子客户端共享 `kernel.BaseClient`，具备统一的 `HttpGet/HttpPost/HttpUpload` 等能力，以及 `MediaXCore/utils/object.StructToStringMap` 的参数构建方式。
3. 子客户端聚焦在业务语义封装，不重新处理认证，共享日志、重试、缓存策略，保持与官方 [YouTube Data API v3](https://developers.google.com/youtube/v3/docs?hl=zh-cn) 一致的 REST 路由与请求体。

## 3. 认证、依赖与通用特性
- **认证链路**：所有请求经过 `core.GoogleAccessTokenHandler` 自动注入 Access Token，可选自定义 token 获取逻辑，满足多租户或外部授权缓存的场景。
- **异常与日志**：依赖 `MediaXCore/pkg/logger`，在 `kernel.BaseClient` 层记录请求响应、配额使用情况，便于与 MediaX 统一链路观测接入。
- **缓存与限流**：留空 `cache.ICache` 接口，便于按业务需要注入 Redis/MemCache，用于 token、配额策略或接口数据预热。
- **参数转换**：统一使用 `object.StructToStringMap` 将请求结构体转换为 query string，减少重复代码。
- **配置模板**：`config.example.yaml`/`config.yaml` 已内置 `google_youtube_config` 示例（含 `oauth_key` 与环境变量占位符）；调试指南详见 `docs/develop/access-token/google/develop.md`、`debug.md`。

## 4. 能力矩阵
| 功能域 | 子客户端 / 路径 | 已实现方法 | 说明 |
| --- | --- | --- | --- |
| 视频管理 | `video/client.go` | `List` `Insert` `Update` `Delete` `Rate` `GetRating` `ReportAbuse` | 覆盖上传、元数据维护、评分与举报。 |
| 视频分类 | `videoCategory/client.go` | `List` | 按地区/语言读取分类。 |
| 举报原因 | `videoAbuseReportReasons/client.go` | `List` | 获取官方举报原因枚举。 |
| 字幕 | `captions/client.go` | `List` `Insert` `Update` `Download` `Delete` | 处理字幕轨道生命周期，含文件上传/下载。 |
| 频道横幅 | `channelBanners/client.go` | `Insert` | 上传频道横幅图片。 |
| 频道版块 | `channelSections/client.go` | `List` `Insert` `Update` `Delete` | 频道首页分区管理。 |
| 频道信息 | `channels/client.go` | `List` `Update` | 查询频道详情并维护 branding/localization。 |
| 评论线程 | `commentThreads/client.go` | `List` `Insert` | 列表 + 创建顶级评论。 |
| 评论 | `comments/client.go` | `List` `Insert` `Update` `SetModerationStatus` `Delete` | 直接操作评论资源与审核状态。 |
| 国际化语言 | `i18nLanguages/client.go` | `List` | 读取平台支持的语言。 |
| 国际化地区 | `i18nRegions/client.go` | `List` | 读取平台支持的地区。 |
| 频道会员 | `members/client.go` | `List` | 拉取频道会员及筛选。 |
| 会员等级 | `membershipsLevels/client.go` | `List` | 列举频道可售卖的会员等级。 |
| 播放列表图片 | `playlistImages/client.go` | `List` `Insert` `Update` `Delete` | 实验性 API，用于管理 playlist artwork。 |
| 播放列表项 | `playlistItems/client.go` | `List` `Insert` `Update` `Delete` | 控制列表内容顺序与内容。 |
| 播放列表 | `playlists/client.go` | `List` `Insert` `Update` `Delete` | 维护播放列表基本信息及隐私。 |
| 搜索 | `search/client.go` | `List` | 支持多维过滤的内容搜索。 |
| 订阅 | `subscriptions/client.go` | `List` `Insert` `Delete` | 查询/创建/取消订阅关系。 |
| 缩略图 | `thumbnails/client.go` | `Set` | 上传视频封面，支持文件流与表单。 |
| 频道水印 | `watermarks/client.go` | `Set` `Unset` | 管理频道右下角水印。 |
| 活动（占位） | `activities/client.go` | — | 仅存在空客户端，尚未封装具体 API。 |

## 5. 能力详解
### 5.1 视频与合规能力
- **视频 CRUD（`video/client.go`）**：`List` 通过 `/youtube/v3/videos` 支持按 `part/id/chart/regionCode` 检索；`Insert`、`Update`、`Delete` 覆盖视频上传与元数据维护；`Rate`/`GetRating` 用于点赞、获取评分；`ReportAbuse` 暴露官方举报通道，配合 `videoAbuseReportReasons.List` 输出的原因枚举实现完整合规模块。
- **视频分类（`videoCategory/client.go`）**：`List` 允许传入 `regionCode`、`hl` 获取多语言分类，适合生成下拉选项或构建本地缓存。

### 5.2 字幕全流程
- `captions.List/Insert/Update/Delete` 完整覆盖字幕轨道管理，`Download` 借助 `HttpGet` 返回 `Content-Type` + 二进制内容，可配合 `tfmt/tlang` 构建多格式、多语言字幕分发。
- 插入、更新均支持 `onBehalfOfContentOwner`，满足 MCN 场景；上传接口默认使用 JSON 体，可扩展为 multipart。

### 5.3 播放列表与媒体容器
- `playlists` + `playlistItems` 是常规列表维护；`playlistImages` 针对播放列表封面，提供额外的 `Insert/Update` 路由。
- 所有方法都允许 `onBehalfOfContentOwner`，便于内容合作伙伴批量操作子频道的播放列表。

### 5.4 评论与互动
- `commentThreads.List/Insert` 支持拉取与创建顶级评论，`comments` 客户端进一步支持回复、修改、审核状态（`SetModerationStatus`）以及删除，方便搭建风控托管流程。

### 5.5 频道装饰与品牌化
- `channels.Update` 可更新 `brandingSettings`、`invideoPromotion`、`localizations`；
- `channelBanners.Insert`、`thumbnails.Set`、`watermarks.Set/Unset` 提供素材上传；
- `channelSections` 管理频道首页布局。

### 5.6 会员体系
- `members.List` 支持 `mode`, `filterByMemberChannelId`, `hasAccessToLevel` 三种筛选方式，覆盖常见 CRM 场景；
- `membershipsLevels.List` 输出所有上架等级的上下文，可用于渲染支付页或进行权限校验。

### 5.7 搜索、订阅与发现
- `search.List` 把官方 query 条件（地理位置、字幕、视频类型等）全部显式暴露，便于扩展内容洞察；
- `subscriptions.List/Insert/Delete` 覆盖自有频道订阅管理以及粉丝列表查询，`Insert` 只需传入 `snippet.resourceId.channelId` 即可。

### 5.8 国际化辅助
- `i18nLanguages.List` 与 `i18nRegions.List` 将 YouTube 支持的语言/地区按 `hl` 参数本地化，方便构建前端多语言选择器。

### 5.9 尚未落地的 `activities`
- `activities` 目录当前仅提供客户端骨架，尚未封装 `List` 或 `Insert` API，如需使用 YouTube 活动 Feed，可参考现有模式新增方法和 schema。

## 6. 使用示例
```go
package demo

import (
    "context"

    mediaxClient "github.com/ArtisanCloud/MediaX/pkg/client"
    "github.com/ArtisanCloud/MediaX/pkg/client/config"
    videoSchema "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient/video/schema"
    "github.com/ArtisanCloud/MediaXCore/pkg/cache"
)

func FetchChannelVideos(ctx context.Context, mxCfg *config.MediaXConfig, ytCfg *config.GoogleYouTubeConfig, c cache.ICache) error {
    mx := mediaxClient.NewMediaX(mxCfg, c)
    yt, err := mx.CreateGoogleYouTubeACClient(ytCfg)
    if err != nil {
        return err
    }

    res, err := yt.GetVideoClient().List(ctx, &videoSchema.YouTubeVideoListReq{
        Part:       "snippet,contentDetails,statistics",
        Id:         []string{"abc123", "xyz789"},
        RegionCode: "US",
    })
    if err != nil {
        return err
    }

    for _, item := range res.Items {
        // TODO: 自定义业务处理
    }
    return nil
}
```
> 说明：示例依赖 `MediaX.CreateGoogleYouTubeACClient` 完成实例化，子客户端与 schema 可直接复用 `pkg/client/google/youtube/accessTokenClient/.../schema` 中的结构体。

## 7. 扩展建议
1. **补齐 activities**：按照 `List` 或 `Insert` 的常见参数增加 schema 与方法，完善频道动态能力。
2. **完善示例**：可在 `playground/google.go` 中补充更多认证、上传示例，帮助第三方快速上手。
3. **文档链接校验**：部分如 `playlistImages` 属于实验性接口，建议标注可用性或提供 fallback。
4. **配额与错误映射**：后续可在 `kernel.BaseClient` 记录配额消耗与错误码映射，方便治理。
