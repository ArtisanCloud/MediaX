package accessTokenClient

import (
	"context"
	"fmt"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/connection/data"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/connection/fan"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/connection/fanProfile"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/content/activity"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/content/task"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/content/video"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/group"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/message"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/tool/appletTemplate"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/tool/retainCard"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/market/service"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/oauth"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/search"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/tools/ticket"
	core2 "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// ByteDanceDouYinACClient 抖音访问Token客户端
// 提供抖音开放平台的各种功能接口，使用访问Token进行认证
// 文档：https://developer.open-douyin.com/docs/resource/zh-CN/dop/overview/usage-guide
type ByteDanceDouYinACClient struct {
	ByteDanceClient    *core2.ByteDanceClient        // 字节跳动基础客户端
	DouYinConfig       *config.ByteDanceDouYinConfig // 抖音配置
	AccessTokenHandler *core2.ByteDanceTokenHandler  // 访问Token处理器

	// clients
	video *video.DouYinContentVideoClient // 视频管理客户端
	task  *task.DouYinContentTaskClient   // 任务管理客户端

	activity             *activity.DouYinContentActivityClient            // 活动管理客户端
	search               *search.DouYinSearchClient                       // 搜索管理客户端
	oauth                *oauth.DouYinOAuthClient                         // OAuth认证客户端
	connectionFan        *fan.DouYinConnectionFanClient                   // 粉丝连接管理客户端
	connectionFanProfile *fanProfile.DouYinConnectionFanProfileClient     // 粉丝资料管理客户端
	connectionData       *data.DouYinConnectionDataClient                 // 数据连接管理客户端
	imMessage            *message.DouYinIMMessageClient                   // IM消息管理客户端
	imGroup              *group.DouYinIMGroupClient                       // IM群组管理客户端
	imToolAppletTemplate *appletTemplate.DouYinIMToolAppletTemplateClient // IM小程序模板工具客户端
	imToolRetainCard     *retainCard.DouYinIMToolRetainCardClient         // IM留存卡片工具客户端
	ticket               *ticket.DouYinTicketClient                       // 票据管理客户端
	marketService        *service.DouYinMarketServiceClient               // 市场服务管理客户端
}

// NewByteDanceDouYinACClient 创建新的抖音访问Token客户端实例
// cfg: 抖音配置
// logger: 日志记录器
// cache: 缓存接口
func NewByteDanceDouYinACClient(cfg *config.ByteDanceDouYinConfig, logger *logger.Logger, cache cache.ICache) (*ByteDanceDouYinACClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = config.ByteDanceDouYinAPIUrl
	}
	c, err := core2.NewByteDanceClient(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	handler, err := core2.NewByteDanceTokenHandler(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	// bind token handler to client
	c.TokenHandler = handler.TokenHandler

	return &ByteDanceDouYinACClient{
		ByteDanceClient:    c,
		DouYinConfig:       cfg,
		AccessTokenHandler: handler,
	}, nil
}

// OverrideGetQuery 重写获取查询参数的方法
// 用于自定义获取访问Token时的查询参数和请求头
func (c *ByteDanceDouYinACClient) OverrideGetQuery() {
	tHandler := c.AccessTokenHandler.TokenHandler
	tHandler.GetTokenQuery = func(ctx context.Context) (arrayQuery *object.StringMap, arrayHeader *object.StringMap, err error) {
		// set the current token key
		var key string
		if tHandler.QueryName != "" {
			key = tHandler.QueryName
		} else {
			key = tHandler.TokenKey
		}

		// get token string power
		resToken := &response.ByteDanceAccessTokenRes{}
		err = tHandler.GetToken(ctx, false, resToken)
		if err != nil {
			return nil, nil, err
		}
		if resToken.AccessToken == "" {
			return nil, nil, fmt.Errorf("get access token error")
		}

		arrayQuery = &object.StringMap{
			"open_id": resToken.OpenId,
		}
		arrayHeader = &object.StringMap{
			key: resToken.AccessToken,
		}

		return arrayQuery, arrayHeader, err
	}
}

// GetContentVideoClient 获取抖音内容管理-视频管理客户端
func (c *ByteDanceDouYinACClient) GetContentVideoClient() *video.DouYinContentVideoClient {
	if c.video == nil {
		c.video = video.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.video
}

// GetOAuthClient 获取抖音OAuth认证客户端
func (c *ByteDanceDouYinACClient) GetOAuthClient() *oauth.DouYinOAuthClient {
	if c.oauth == nil {
		c.oauth = oauth.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.oauth
}

// GetSearchClient 获取抖音搜索管理客户端
func (c *ByteDanceDouYinACClient) GetSearchClient() *search.DouYinSearchClient {
	if c.search == nil {
		c.search = search.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.search
}

// GetConnectionFanClient 获取抖音粉丝连接管理客户端
func (c *ByteDanceDouYinACClient) GetConnectionFanClient() *fan.DouYinConnectionFanClient {
	if c.connectionFan == nil {
		c.connectionFan = fan.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionFan
}

// GetConnectionFanProfileClient 获取抖音粉丝资料管理客户端
func (c *ByteDanceDouYinACClient) GetConnectionFanProfileClient() *fanProfile.DouYinConnectionFanProfileClient {
	if c.connectionFanProfile == nil {
		c.connectionFanProfile = fanProfile.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionFanProfile
}

// GetConnectionDataClient 获取抖音数据连接管理客户端
func (c *ByteDanceDouYinACClient) GetConnectionDataClient() *data.DouYinConnectionDataClient {
	if c.connectionData == nil {
		c.connectionData = data.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionData
}

// GetIMMessageClient 获取抖音IM消息管理客户端
func (c *ByteDanceDouYinACClient) GetIMMessageClient() *message.DouYinIMMessageClient {
	if c.imMessage == nil {
		c.imMessage = message.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imMessage
}

// GetIMToolAppletTemplateClient 获取抖音IM小程序模板工具客户端
func (c *ByteDanceDouYinACClient) GetIMToolAppletTemplateClient() *appletTemplate.DouYinIMToolAppletTemplateClient {
	if c.imToolAppletTemplate == nil {
		c.imToolAppletTemplate = appletTemplate.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imToolAppletTemplate
}

// GetIMToolRetainCardClient 获取抖音IM留存卡片工具客户端
func (c *ByteDanceDouYinACClient) GetIMToolRetainCardClient() *retainCard.DouYinIMToolRetainCardClient {
	if c.imToolRetainCard == nil {
		c.imToolRetainCard = retainCard.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imToolRetainCard
}

// GetIMGroupClient 获取抖音IM群组管理客户端
func (c *ByteDanceDouYinACClient) GetIMGroupClient() *group.DouYinIMGroupClient {
	if c.imGroup == nil {
		c.imGroup = group.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imGroup
}

// GetContentTaskClient 获取抖音内容管理-任务管理客户端
func (c *ByteDanceDouYinACClient) GetContentTaskClient() *task.DouYinContentTaskClient {
	if c.task == nil {
		c.task = task.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.task
}

// GetMarketServiceClient 获取抖音市场服务管理客户端
func (c *ByteDanceDouYinACClient) GetContentActivityClient() *activity.DouYinContentActivityClient {
	if c.activity == nil {
		c.activity = activity.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.activity
}

// GetTicketClient 获取抖音票据管理客户端
func (c *ByteDanceDouYinACClient) GetTicketClient() *ticket.DouYinTicketClient {
	if c.ticket == nil {
		c.ticket = ticket.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.ticket
}

// GetContentActivityClient 获取抖音内容管理-活动管理客户端
func (c *ByteDanceDouYinACClient) GetMarketServiceClient() *service.DouYinMarketServiceClient {
	if c.marketService == nil {
		c.marketService = service.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.marketService
}
