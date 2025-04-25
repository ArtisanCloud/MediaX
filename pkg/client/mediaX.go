package client

import (
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	"github.com/ArtisanCloud/MediaX/pkg/client/google/youtube"
	"github.com/ArtisanCloud/MediaX/pkg/client/redbook/juGuang"
	"github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount"
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
func (m *MediaX) CreateWechatOfficialAccount(cfg *config.WeChatOfficialAccountConfig) (*officialAccount.WeChatOfficialAccountClient, error) {
	return officialAccount.NewWeChatOfficialAccountClient(cfg, m.Logger, m.Cache)
}

// CreateGoogleYouTube 创建 CreateGoogleYouTube，支持传入 Google 配置
func (m *MediaX) CreateGoogleYouTube(cfg *config.GoogleYouTubeConfig) (*youtube.GoogleYouTubeClient, error) {
	return youtube.NewGoogleYouTubeClient(cfg, m.Logger, m.Cache)
}

// CreateByteDanceDouYin 创建 DouYinClient，支持传入 DouYin 配置
func (m *MediaX) CreateByteDanceDouYin(cfg *config.ByteDanceDouYinConfig) (*douYin.ByteDanceDouYinClient, error) {
	return douYin.NewByteDanceDouYinClient(cfg, m.Logger, m.Cache)
}

// CreateRedBookJuGuang 创建 RedBookClient，支持传入 RedBook 配置
func (m *MediaX) CreateRedBookJuGuang(cfg *config.RedBookJuGuangConfig) (*juGuang.RedBookJuGuangClient, error) {
	return juGuang.NewRedBookJuGuangClient(cfg, m.Logger, m.Cache)
}
