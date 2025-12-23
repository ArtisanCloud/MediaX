package clientTokenClient

import (
	"context"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/activity"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/schemas"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/task"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/video"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/search"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/tools/micApp"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/tools/sandbox"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/tools/ticket"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core"
	douyinresponse "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// ByteDanceDouYinCTClient 抖音客户端Token客户端
// 提供抖音开放平台的各种功能接口
// 文档：https://developer.open-douyin.com/docs/resource/zh-CN/dop/overview/usage-guide
type ByteDanceDouYinCTClient struct {
	ByteDanceClient    *core.ByteDanceClient         // 字节跳动基础客户端
	DouYinConfig       *config.ByteDanceDouYinConfig // 抖音配置
	ClientTokenHandler *core.ByteDanceTokenHandler   // 客户端Token处理器
	tokenCache         *clientTokenCache

	// clients
	video       *video.DouYinContentVideoClient       // 视频管理客户端
	task        *task.DouYinContentTaskClient         // 任务管理客户端
	schemas     *schemas.DouYinContentSchemasClient   // 内容模板管理客户端
	activity    *activity.DouYinContentActivityClient // 活动管理客户端
	search      *search.DouYinSearchClient            // 搜索管理客户端
	toolMicApp  *micApp.DouYinToolMicAppClient        // 小程序工具客户端
	toolSandbox *sandbox.DouYinToolSandboxClient      // 沙箱工具客户端
	toolTicket  *ticket.DouYinToolTicketClient        // 票据工具客户端
}

// NewByteDanceDouYinCTClient 创建新的抖音客户端Token客户端实例
// cfg: 抖音配置
// logger: 日志记录器
// cache: 缓存接口
func NewByteDanceDouYinCTClient(cfg *config.ByteDanceDouYinConfig, logger *logger.Logger, cache cache.ICache) (*ByteDanceDouYinCTClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = config.ByteDanceDouYinAPIUrl
	}
	c, err := core.NewByteDanceClient(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	handler, err := core.NewByteDanceTokenHandler(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	// bind token handler to client
	c.TokenHandler = handler.TokenHandler

	client := &ByteDanceDouYinCTClient{
		ByteDanceClient:    c,
		DouYinConfig:       cfg,
		ClientTokenHandler: handler,
	}
	client.applyTokenCache(cache)
	return client, nil
}

// applyTokenCache 将 DouYin client_token 缓存策略绑定到 kernel.TokenHandler：
// 统一 redis key、refresh_before 判定，并在 TokenHandler.GetTokenQuery 中复用 EnsureClientToken。
func (c *ByteDanceDouYinCTClient) applyTokenCache(cacheBackend cache.ICache) {
	store := newClientTokenCache(c.DouYinConfig, cacheBackend)
	if store == nil {
		return
	}
	c.tokenCache = store
	tokenKey := c.ClientTokenHandler.TokenHandler.TokenKey
	if tokenKey == "" {
		tokenKey = "access_token"
	}
	c.ClientTokenHandler.TokenHandler.GetTokenQuery = func(ctx context.Context) (arrayQuery *object.StringMap, arrayHeader *object.StringMap, err error) {
		record, err := c.EnsureClientToken(ctx, false)
		if err != nil || record == nil {
			return nil, nil, err
		}
		query := object.StringMap{
			tokenKey: record.AccessToken,
		}
		return &query, nil, nil
	}
}

// EnsureClientToken 读取缓存的 DouYin client_token，若剩余 TTL 低于阈值则自动刷新并回写。
func (c *ByteDanceDouYinCTClient) EnsureClientToken(ctx context.Context, forceRefresh bool) (*clientTokenCacheRecord, error) {
	if c.tokenCache == nil {
		return c.refreshClientToken(ctx)
	}
	if forceRefresh {
		return c.tokenCache.forceRefresh(ctx, c.refreshClientTokenFetch)
	}
	return c.tokenCache.ensure(ctx, c.refreshClientTokenFetch)
}

// refreshClientTokenFetch 调用 DouYin API 刷新 Token。
func (c *ByteDanceDouYinCTClient) refreshClientTokenFetch(ctx context.Context) (*douyinresponse.ByteDanceAccessTokenRes, error) {
	tokenRes := &douyinresponse.ByteDanceAccessTokenRes{}
	if err := c.ClientTokenHandler.TokenHandler.GetToken(ctx, true, tokenRes); err != nil {
		return nil, err
	}
	return tokenRes, nil
}

func (c *ByteDanceDouYinCTClient) refreshClientToken(ctx context.Context) (*clientTokenCacheRecord, error) {
	tokenRes, err := c.refreshClientTokenFetch(ctx)
	if err != nil {
		return nil, err
	}
	now := nowFunc()
	ttl := c.DouYinConfig.EffectiveTTLSeconds()
	if tokenRes != nil && tokenRes.ExpiresIn > 0 {
		ttl = int(tokenRes.ExpiresIn)
	}
	if ttl <= 0 {
		ttl = 7000
	}
	record := &clientTokenCacheRecord{
		AppID:       c.DouYinConfig.ClientTokenCredential().ClientKeyValue(),
		AccessToken: tokenRes.AccessToken,
		StoredAt:    now,
		ExpireAt:    now.Add(time.Duration(ttl) * time.Second),
		Source:      "douyin.client",
	}
	if c.tokenCache != nil {
		_ = c.tokenCache.save(ctx, record, time.Duration(ttl)*time.Second)
	}
	return record, nil
}

// GetContentVideoClient 获取抖音内容管理-视频管理客户端
func (c *ByteDanceDouYinCTClient) GetContentVideoClient() *video.DouYinContentVideoClient {
	if c.video == nil {
		c.video = video.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.video
}

// GetContentTaskClient 获取抖音内容管理-任务管理客户端
func (c *ByteDanceDouYinCTClient) GetContentTaskClient() *task.DouYinContentTaskClient {
	if c.task == nil {
		c.task = task.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.task
}

// GetContentSchemasClient 获取抖音内容管理-模板管理客户端
func (c *ByteDanceDouYinCTClient) GetContentSchemasClient() *schemas.DouYinContentSchemasClient {
	if c.schemas == nil {
		c.schemas = schemas.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.schemas
}

// GetContentActivityClient 获取抖音内容管理-活动管理客户端
func (c *ByteDanceDouYinCTClient) GetContentActivityClient() *activity.DouYinContentActivityClient {
	if c.activity == nil {
		c.activity = activity.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.activity
}

// GetSearchClient 获取抖音搜索管理客户端
func (c *ByteDanceDouYinCTClient) GetSearchClient() *search.DouYinSearchClient {
	if c.search == nil {
		c.search = search.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.search
}

// GetContentSchemaClient 获取抖音内容管理-模板管理客户端（与GetContentSchemasClient功能相同）
func (c *ByteDanceDouYinCTClient) GetContentSchemaClient() *schemas.DouYinContentSchemasClient {
	if c.schemas == nil {
		c.schemas = schemas.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.schemas
}

// GetMicAppClient 获取抖音小程序工具客户端
func (c *ByteDanceDouYinCTClient) GetMicAppClient() *micApp.DouYinToolMicAppClient {
	if c.toolMicApp == nil {
		c.toolMicApp = micApp.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.toolMicApp
}

// GetSandboxClient 获取抖音沙箱工具客户端
func (c *ByteDanceDouYinCTClient) GetSandboxClient() *sandbox.DouYinToolSandboxClient {
	if c.toolSandbox == nil {
		c.toolSandbox = sandbox.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.toolSandbox
}

// GetTicketClient 获取抖音票据工具客户端
func (c *ByteDanceDouYinCTClient) GetTicketClient() *ticket.DouYinToolTicketClient {
	if c.toolTicket == nil {
		c.toolTicket = ticket.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.toolTicket
}
