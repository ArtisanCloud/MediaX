package juGuang

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/redbook/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/redbook/juGuang/account"
	"github.com/ArtisanCloud/MediaX/pkg/client/redbook/juGuang/dataReport/offline"
	"github.com/ArtisanCloud/MediaX/pkg/client/redbook/juGuang/dataReport/realtime"
	"github.com/ArtisanCloud/MediaX/pkg/client/redbook/juGuang/note"
	"github.com/ArtisanCloud/MediaX/pkg/client/redbook/juGuang/promote/campaign"
	"github.com/ArtisanCloud/MediaX/pkg/client/redbook/juGuang/promote/creativity"
	"github.com/ArtisanCloud/MediaX/pkg/client/redbook/juGuang/promote/unit"
	"github.com/ArtisanCloud/MediaX/pkg/client/redbook/juGuang/tools"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// https://ad-market.xiaohongshu.com/docs-center?articleId=3180&bizType=943
type RedBookJuGuangACClient struct {
	RedBookClient      *core.RedBookClient
	JuGuangConfig      *config.RedBookJuGuangConfig
	AccessTokenHandler *core.RedBookAccessTokenHandler

	// clients
	account        *account.JuGuangAccountClient
	offlineReport  *offline.JuGuangDataReportOfflineClient
	realtimeReport *realtime.JuGuangDataReportRealtimeClient
	note           *note.JuGuangNoteClient
	campaign       *campaign.JuGuangPromoteCampaignClient
	creativity     *creativity.JuGuangPromoteCreativityClient
	unit           *unit.JuGuangPromoteUnitClient
	tools          *tools.JuGuangToolClient
}

func NewRedBookJuGuangACClient(cfg *config.RedBookJuGuangConfig, logger *logger.Logger, cache cache.ICache) (*RedBookJuGuangACClient, error) {
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = "https://adapi.xiaohongshu.com/"
	}
	c, err := core.NewRedBookClient(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	handler, err := core.NewRedBookAccessTokenHandler(cfg.ClientConfig, logger, cache)
	if err != nil {
		return nil, err
	}

	// bind token handler to client
	c.TokenHandler = handler.AccessTokenHandler

	// override get custom token
	c.TokenHandler.GetCustomToken = cfg.GetOAuthToken

	return &RedBookJuGuangACClient{
		RedBookClient:      c,
		JuGuangConfig:      cfg,
		AccessTokenHandler: handler,
	}, nil
}

func (client *RedBookJuGuangACClient) GetAccountClient() *account.JuGuangAccountClient {
	if client.account == nil {
		client.account = account.NewClient(client.RedBookClient.BaseClient)
	}
	return client.account
}

func (client *RedBookJuGuangACClient) GetOfflineReportClient() *offline.JuGuangDataReportOfflineClient {
	if client.offlineReport == nil {
		client.offlineReport = offline.NewClient(client.RedBookClient.BaseClient)
	}
	return client.offlineReport
}

func (client *RedBookJuGuangACClient) GetRealtimeReportClient() *realtime.JuGuangDataReportRealtimeClient {
	if client.realtimeReport == nil {
		client.realtimeReport = realtime.NewClient(client.RedBookClient.BaseClient)
	}
	return client.realtimeReport
}

func (client *RedBookJuGuangACClient) GetNoteClient() *note.JuGuangNoteClient {
	if client.note == nil {
		client.note = note.NewClient(client.RedBookClient.BaseClient)
	}
	return client.note
}

func (client *RedBookJuGuangACClient) GetCampaignClient() *campaign.JuGuangPromoteCampaignClient {
	if client.campaign == nil {
		client.campaign = campaign.NewClient(client.RedBookClient.BaseClient)
	}
	return client.campaign
}

func (client *RedBookJuGuangACClient) GetCreativityClient() *creativity.JuGuangPromoteCreativityClient {
	if client.creativity == nil {
		client.creativity = creativity.NewClient(client.RedBookClient.BaseClient)
	}
	return client.creativity
}

func (client *RedBookJuGuangACClient) GetUnitClient() *unit.JuGuangPromoteUnitClient {
	if client.unit == nil {
		client.unit = unit.NewClient(client.RedBookClient.BaseClient)
	}
	return client.unit
}

func (client *RedBookJuGuangACClient) GetToolsClient() *tools.JuGuangToolClient {
	if client.tools == nil {
		client.tools = tools.NewClient(client.RedBookClient.BaseClient)
	}
	return client.tools
}
