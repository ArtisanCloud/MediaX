package clientTokenClient

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/activity"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/schemas"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/task"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/content/video"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/search"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/tools/micApp"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient/tools/sandbox"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// https://developer.open-douyin.com/docs/resource/zh-CN/dop/overview/usage-guide
type ByteDanceDouYinCTClient struct {
	ByteDanceClient    *core.ByteDanceClient
	DouYinConfig       *config.ByteDanceDouYinConfig
	ClientTokenHandler *core.ByteDanceTokenHandler

	// clients
	video       *video.DouYinContentVideoClient
	task        *task.DouYinContentTaskClient
	schemas     *schemas.DouYinContentSchemasClient
	activity    *activity.DouYinContentActivityClient
	search      *search.DouYinSearchClient
	toolMicApp  *micApp.DouYinToolMicAppClient
	toolSandbox *sandbox.DouYinSandboxClient
}

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

	return &ByteDanceDouYinCTClient{
		ByteDanceClient:    c,
		DouYinConfig:       cfg,
		ClientTokenHandler: handler,
	}, nil
}

func (c *ByteDanceDouYinCTClient) GetContentVideoClient() *video.DouYinContentVideoClient {
	if c.video == nil {
		c.video = video.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.video
}

func (c *ByteDanceDouYinCTClient) GetContentTaskClient() *task.DouYinContentTaskClient {
	if c.task == nil {
		c.task = task.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.task
}

func (c *ByteDanceDouYinCTClient) GetContentSchemasClient() *schemas.DouYinContentSchemasClient {
	if c.schemas == nil {
		c.schemas = schemas.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.schemas
}

func (c *ByteDanceDouYinCTClient) GetContentActivityClient() *activity.DouYinContentActivityClient {
	if c.activity == nil {
		c.activity = activity.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.activity
}

func (c *ByteDanceDouYinCTClient) GetSearchClient() *search.DouYinSearchClient {
	if c.search == nil {
		c.search = search.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.search
}

func (c *ByteDanceDouYinCTClient) GetContentSchemaClient() *schemas.DouYinContentSchemasClient {
	if c.schemas == nil {
		c.schemas = schemas.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.schemas
}
func (c *ByteDanceDouYinCTClient) GetMicAppClient() *micApp.DouYinToolMicAppClient {
	if c.toolMicApp == nil {
		c.toolMicApp = micApp.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.toolMicApp
}

func (c *ByteDanceDouYinCTClient) GetSandboxClient() *sandbox.DouYinSandboxClient {
	if c.toolSandbox == nil {
		c.toolSandbox = sandbox.NewClient(c.ByteDanceClient.BaseClient)
	}
	return c.toolSandbox
}
