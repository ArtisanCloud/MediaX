package douYin

import (
	"context"
	"fmt"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/connection/data"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/connection/fan"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/connection/fanData"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/content/activity"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/content/schema"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/content/task"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/content/video"
	core2 "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
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
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// https://developer.open-douyin.com/docs/resource/zh-CN/dop/overview/usage-guide
type ByteDanceDouYinAccessTokenClient struct {
	ByteDanceClient    *core2.ByteDanceClient
	DouYinConfig       *config.ByteDanceDouYinConfig
	AccessTokenHandler *core2.ByteDanceTokenHandler

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

func NewByteDanceDouYinAccessTokenClient(cfg *config.ByteDanceDouYinConfig, logger *logger.Logger, cache cache.ICache) (*ByteDanceDouYinAccessTokenClient, error) {
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

	return &ByteDanceDouYinAccessTokenClient{
		ByteDanceClient:    c,
		DouYinConfig:       cfg,
		AccessTokenHandler: handler,
	}, nil
}

func (c *ByteDanceDouYinAccessTokenClient) OverrideGetQuery() {
	tHandler := c.AccessTokenHandler.TokenHandler
	tHandler.GetTokenQuery = func(ctx context.Context) (*object.StringMap, error) {
		// set the current token key
		var key string
		if tHandler.QueryName != "" {
			key = tHandler.QueryName
		} else {
			key = tHandler.TokenKey
		}

		// get token string power
		resToken := &response.ByteDanceAccessTokenRes{}
		err := tHandler.GetToken(ctx, false, resToken)
		if err != nil {
			return nil, err
		}
		if resToken.AccessToken == "" {
			return nil, fmt.Errorf("get access token error")
		}

		arrayReturn := &object.StringMap{
			key:       resToken.AccessToken,
			"open_id": resToken.OpenId,
		}

		return arrayReturn, err
	}
}

func (c *ByteDanceDouYinAccessTokenClient) GetContentVideoClient() *video.DouYinContentVideoClient {
	if c.video == nil {
		c.video = video.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.video
}

func (c *ByteDanceDouYinAccessTokenClient) GetOAuthClient() *oauth.DouYinOAuthClient {
	if c.oauth == nil {
		c.oauth = oauth.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.oauth
}

func (c *ByteDanceDouYinAccessTokenClient) GetSearchClient() *search.DouYinSearchClient {
	if c.search == nil {
		c.search = search.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.search
}

func (c *ByteDanceDouYinAccessTokenClient) GetConnectionFanClient() *fan.DouYinConnectionFanClient {
	if c.connectionFan == nil {
		c.connectionFan = fan.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionFan
}

func (c *ByteDanceDouYinAccessTokenClient) GetConnectionFanDataClient() *fanData.DouYinConnectionFanDataClient {
	if c.connectionFanData == nil {
		c.connectionFanData = fanData.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionFanData
}

func (c *ByteDanceDouYinAccessTokenClient) GetConnectionDataClient() *data.DouYinConnectionDataClient {
	if c.connectionData == nil {
		c.connectionData = data.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionData
}

func (c *ByteDanceDouYinAccessTokenClient) GetIMMessageClient() *message.DouYinIMMessageClient {
	if c.imMessage == nil {
		c.imMessage = message.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imMessage
}

func (c *ByteDanceDouYinAccessTokenClient) GetIMToolClient() *tool.DouYinIMToolClient {
	if c.imTool == nil {
		c.imTool = tool.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imTool
}

func (c *ByteDanceDouYinAccessTokenClient) GetIMGroupClient() *group.DouYinIMGroupClient {
	if c.imGroup == nil {
		c.imGroup = group.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imGroup
}

func (c *ByteDanceDouYinAccessTokenClient) GetContentTaskClient() *task.DouYinContentTaskClient {
	if c.task == nil {
		c.task = task.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.task
}

func (c *ByteDanceDouYinAccessTokenClient) GetContentActivityClient() *activity.DouYinContentActivityClient {
	if c.activity == nil {
		c.activity = activity.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.activity
}

func (c *ByteDanceDouYinAccessTokenClient) GetContentSchemaClient() *schema.DouYinContentSchemaClient {
	if c.contentSchema == nil {
		c.contentSchema = schema.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.contentSchema
}

func (c *ByteDanceDouYinAccessTokenClient) GetSandboxClient() *sandbox.DouYinSandboxClient {
	if c.sandbox == nil {
		c.sandbox = sandbox.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.sandbox
}

func (c *ByteDanceDouYinAccessTokenClient) GetMicAppClient() *micApp.DouYinMicAppClient {
	if c.micApp == nil {
		c.micApp = micApp.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.micApp
}

func (c *ByteDanceDouYinAccessTokenClient) GetTicketClient() *ticket.DouYinTicketClient {
	if c.ticket == nil {
		c.ticket = ticket.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.ticket
}

func (c *ByteDanceDouYinAccessTokenClient) GetMarketServiceClient() *service.DouYinMarketServiceClient {
	if c.marketService == nil {
		c.marketService = service.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.marketService
}
