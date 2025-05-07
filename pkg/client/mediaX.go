package client

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	accessTokenClient3 "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient"
	accessTokenClient2 "github.com/ArtisanCloud/MediaX/pkg/client/redBook/juGuang/accessTokenClient"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

type MediaX struct {
	Logger *logger.Logger // 全局 Logger
	Cache  cache.ICache   // 全局 Cache
}

// NewMediaX 初始化 MediaX，Logger 和 Cache 是全局共享的
func NewMediaX(config *config.MediaXConfig, cache cache.ICache) *MediaX {
	l := logger.NewLogger(config.Logger)
	return &MediaX{
		Logger: l,
		Cache:  cache,
	}
}

// CreateWechatOfficialAccount 创建 WechatOfficialAccountClient，支持传入 WeChat 配置
func (m *MediaX) CreateWechatOfficialAccount(cfg *config.WeChatOfficialAccountConfig) (*officialAccount.WeChatOfficialAccountCTClient, error) {
	return officialAccount.NewWeChatOfficialAccountCTClient(cfg, m.Logger, m.Cache)
}

// CreateGoogleYouTube 创建 CreateGoogleYouTube，支持传入 Google 配置
func (m *MediaX) CreateGoogleYouTubeACClient(cfg *config.GoogleYouTubeConfig) (*accessTokenClient3.GoogleYouTubeACClient, error) {
	return accessTokenClient3.NewGoogleYouTubeACClient(cfg, m.Logger, m.Cache)
}

// CreateByteDanceDouYin 创建 DouYinClient，支持传入 DouYin 配置
func (m *MediaX) CreateByteDanceDouYinACClient(cfg *config.ByteDanceDouYinConfig) (*accessTokenClient.ByteDanceDouYinACClient, error) {
	return accessTokenClient.NewByteDanceDouYinACClient(cfg, m.Logger, m.Cache)
}

// CreateByteDanceDouYin 创建 DouYinClient，支持传入 DouYin 配置
func (m *MediaX) CreateByteDanceDouYinCTClient(cfg *config.ByteDanceDouYinConfig) (*clientTokenClient.ByteDanceDouYinCTClient, error) {
	return clientTokenClient.NewByteDanceDouYinCTClient(cfg, m.Logger, m.Cache)
}

// CreateRedBookJuGuang 创建 RedBookClient，支持传入 RedBook 配置
func (m *MediaX) CreateRedBookJuGuangACClient(cfg *config.RedBookJuGuangConfig) (*accessTokenClient2.RedBookJuGuangACClient, error) {
	return accessTokenClient2.NewRedBookJuGuangACClient(cfg, m.Logger, m.Cache)
}
