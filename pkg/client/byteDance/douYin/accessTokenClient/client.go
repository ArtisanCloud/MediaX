package accessTokenClient

import (
	"context"
	"fmt"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/connection/data"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/connection/fan"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/connection/fanData"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/content/activity"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/content/task"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/content/video"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/group"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/message"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/im/tool"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/market/service"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/oauth"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/search"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/tools/micApp"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/tools/sandbox"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient/tools/ticket"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/schemas"
	core2 "github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core/response"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// https://developer.open-douyin.com/docs/resource/zh-CN/dop/overview/usage-guide
type ByteDanceDouYinACClient struct {
	ByteDanceClient    *core2.ByteDanceClient
	DouYinConfig       *config.ByteDanceDouYinConfig
	AccessTokenHandler *core2.ByteDanceTokenHandler

	// clients
	video             *video.DouYinContentVideoClient
	task              *task.DouYinContentTaskClient
	schemas           *schemas.DouYinContentSchemasClient
	activity          *activity.DouYinContentActivityClient
	search            *search.DouYinSearchClient
	oauth             *oauth.DouYinOAuthClient
	connectionFan     *fan.DouYinConnectionFanClient
	connectionFanData *fanData.DouYinConnectionFanDataClient
	connectionData    *data.DouYinConnectionDataClient
	imMessage         *message.DouYinIMMessageClient
	imTool            *tool.DouYinIMToolClient
	imGroup           *group.DouYinIMGroupClient
	sandbox           *sandbox.DouYinSandboxClient
	micApp            *micApp.DouYinMicAppClient
	ticket            *ticket.DouYinTicketClient
	marketService     *service.DouYinMarketServiceClient
}

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

func (c *ByteDanceDouYinACClient) GetContentVideoClient() *video.DouYinContentVideoClient {
	if c.video == nil {
		c.video = video.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.video
}

func (c *ByteDanceDouYinACClient) GetOAuthClient() *oauth.DouYinOAuthClient {
	if c.oauth == nil {
		c.oauth = oauth.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.oauth
}

func (c *ByteDanceDouYinACClient) GetSearchClient() *search.DouYinSearchClient {
	if c.search == nil {
		c.search = search.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.search
}

func (c *ByteDanceDouYinACClient) GetConnectionFanClient() *fan.DouYinConnectionFanClient {
	if c.connectionFan == nil {
		c.connectionFan = fan.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionFan
}

func (c *ByteDanceDouYinACClient) GetConnectionFanDataClient() *fanData.DouYinConnectionFanDataClient {
	if c.connectionFanData == nil {
		c.connectionFanData = fanData.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionFanData
}

func (c *ByteDanceDouYinACClient) GetConnectionDataClient() *data.DouYinConnectionDataClient {
	if c.connectionData == nil {
		c.connectionData = data.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.connectionData
}

func (c *ByteDanceDouYinACClient) GetIMMessageClient() *message.DouYinIMMessageClient {
	if c.imMessage == nil {
		c.imMessage = message.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imMessage
}

func (c *ByteDanceDouYinACClient) GetIMToolClient() *tool.DouYinIMToolClient {
	if c.imTool == nil {
		c.imTool = tool.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imTool
}

func (c *ByteDanceDouYinACClient) GetIMGroupClient() *group.DouYinIMGroupClient {
	if c.imGroup == nil {
		c.imGroup = group.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.imGroup
}

func (c *ByteDanceDouYinACClient) GetContentTaskClient() *task.DouYinContentTaskClient {
	if c.task == nil {
		c.task = task.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.task
}

func (c *ByteDanceDouYinACClient) GetContentActivityClient() *activity.DouYinContentActivityClient {
	if c.activity == nil {
		c.activity = activity.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.activity
}

func (c *ByteDanceDouYinACClient) GetContentSchemaClient() *schemas.DouYinContentSchemasClient {
	if c.schemas == nil {
		c.schemas = schemas.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.schemas
}

func (c *ByteDanceDouYinACClient) GetSandboxClient() *sandbox.DouYinSandboxClient {
	if c.sandbox == nil {
		c.sandbox = sandbox.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.sandbox
}

func (c *ByteDanceDouYinACClient) GetMicAppClient() *micApp.DouYinMicAppClient {
	if c.micApp == nil {
		c.micApp = micApp.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.micApp
}

func (c *ByteDanceDouYinACClient) GetTicketClient() *ticket.DouYinTicketClient {
	if c.ticket == nil {
		c.ticket = ticket.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.ticket
}

func (c *ByteDanceDouYinACClient) GetMarketServiceClient() *service.DouYinMarketServiceClient {
	if c.marketService == nil {
		c.marketService = service.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.marketService
}
