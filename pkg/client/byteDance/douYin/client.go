package douYin

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/connection/data"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/connection/fan"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/connection/fanData"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/content/activity"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/content/schema"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/content/task"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/content/video"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/im/group"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/im/message"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/im/tool"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/market/service"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/oauth"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/search"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/tools/micApp"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/tools/sandbox"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/tools/ticket"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// https://developer.open-douyin.com/docs/resource/zh-CN/dop/overview/usage-guide
type ByteDanceDouYinClient struct {
	ByteDanceClient    *core.ByteDanceClient
	DouYinConfig       *config.ByteDanceDouYinConfig
	AccessTokenHandler *core.ByteDanceAccessTokenHandler

	// clients
	video             *video.DouYinContentVideoClient
	oauth             *oauth.DouYinOAuthClient
	search            *search.DouYinSearchClient
	connectionFan     *fan.DouYinConnectionFanClient
	connectionFanData *fanData.DouYinConnectionFanDataClient
	connectionData    *data.DouYinConnectionDataClient
	imMessage         *message.DouYinIMMessageClient
	imTool            *tool.DouYinIMToolClient
	imGroup           *group.DouYinIMGroupClient
	task              *task.DouYinContentTaskClient
	activity          *activity.DouYinContentActivityClient
	contentSchema     *schema.DouYinContentSchemaClient
	sandbox           *sandbox.DouYinSandboxClient
	micApp            *micApp.DouYinMicAppClient
	ticket            *ticket.DouYinTicketClient
	marketService     *service.DouYinMarketServiceClient
}

func NewByteDanceDouYinClient(cfg *config.ByteDanceDouYinConfig, logger *logger.Logger, cache cache.ICache) (*ByteDanceDouYinClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = config.ByteDanceDouYinAPIUrl
	}
	c, err := core.NewByteDanceClient(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	handler, err := core.NewByteDanceAccessTokenHandler(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	// bind token handler to client
	c.TokenHandler = handler.AccessTokenHandler

	return &ByteDanceDouYinClient{
		ByteDanceClient:    c,
		DouYinConfig:       cfg,
		AccessTokenHandler: handler,
	}, nil
}

func (c *ByteDanceDouYinClient) GetContentVideoClient() *video.DouYinContentVideoClient {
	if c.video == nil {
		c.video = video.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.video
}

func (c *ByteDanceDouYinClient) GetOAuthClient() *oauth.DouYinOAuthClient {
	if c.oauth == nil {
		c.oauth = oauth.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.oauth
}

func (c *ByteDanceDouYinClient) GetSearchClient() *search.DouYinSearchClient {
	if c.search == nil {
		c.search = search.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.search
}

func (c *ByteDanceDouYinClient) GetConnectionFanClient() *fan.DouYinConnectionFanClient {
	if c.connectionFan == nil {
		c.connectionFan = fan.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionFan
}

func (c *ByteDanceDouYinClient) GetConnectionFanDataClient() *fanData.DouYinConnectionFanDataClient {
	if c.connectionFanData == nil {
		c.connectionFanData = fanData.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionFanData
}

func (c *ByteDanceDouYinClient) GetConnectionDataClient() *data.DouYinConnectionDataClient {
	if c.connectionData == nil {
		c.connectionData = data.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionData
}

func (c *ByteDanceDouYinClient) GetIMMessageClient() *message.DouYinIMMessageClient {
	if c.imMessage == nil {
		c.imMessage = message.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imMessage
}

func (c *ByteDanceDouYinClient) GetIMToolClient() *tool.DouYinIMToolClient {
	if c.imTool == nil {
		c.imTool = tool.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imTool
}

func (c *ByteDanceDouYinClient) GetIMGroupClient() *group.DouYinIMGroupClient {
	if c.imGroup == nil {
		c.imGroup = group.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imGroup
}

func (c *ByteDanceDouYinClient) GetContentTaskClient() *task.DouYinContentTaskClient {
	if c.task == nil {
		c.task = task.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.task
}

func (c *ByteDanceDouYinClient) GetContentActivityClient() *activity.DouYinContentActivityClient {
	if c.activity == nil {
		c.activity = activity.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.activity
}

func (c *ByteDanceDouYinClient) GetContentSchemaClient() *schema.DouYinContentSchemaClient {
	if c.contentSchema == nil {
		c.contentSchema = schema.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.contentSchema
}

func (c *ByteDanceDouYinClient) GetSandboxClient() *sandbox.DouYinSandboxClient {
	if c.sandbox == nil {
		c.sandbox = sandbox.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.sandbox
}

func (c *ByteDanceDouYinClient) GetMicAppClient() *micApp.DouYinMicAppClient {
	if c.micApp == nil {
		c.micApp = micApp.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.micApp
}

func (c *ByteDanceDouYinClient) GetTicketClient() *ticket.DouYinTicketClient {
	if c.ticket == nil {
		c.ticket = ticket.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.ticket
}

func (c *ByteDanceDouYinClient) GetMarketServiceClient() *service.DouYinMarketServiceClient {
	if c.marketService == nil {
		c.marketService = service.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.marketService
}
