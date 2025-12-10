// Package client 提供了多平台媒体内容发布和管理的统一客户端接口
//
// 本包实现了对以下平台的支持：
// - 微信公众号：内容发布、用户管理等功能
// - YouTube：视频上传、频道管理、数据分析等
// - 抖音：视频发布、账号管理、数据同步等
// - 小红书：内容发布、数据分析等
//
// 使用说明：
// 1. 首先需要初始化 MediaX 实例，提供配置和缓存实现
// 2. 然后可以根据需要创建对应平台的客户端
// 3. 每个平台客户端都提供了特定的 API 调用方法
// 4. 支持 AccessToken 和 ClientToken 两种认证方式
//
// 详细文档请参考各平台的开发者文档：
// - 微信：https://developers.weixin.qq.com/doc/
// - YouTube：https://developers.google.com/youtube/v3/docs
// - 抖音：https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/
// - 小红书：https://open.xiaohongshu.com/document
package client

import (
	"errors"

	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/accessTokenClient"
	"github.com/ArtisanCloud/MediaX/pkg/client/byteDance/douYin/clientTokenClient"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	accessTokenClient3 "github.com/ArtisanCloud/MediaX/pkg/client/google/youtube/accessTokenClient"
	accessTokenClient2 "github.com/ArtisanCloud/MediaX/pkg/client/redBook/juGuang/accessTokenClient"
	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
	officialAccount "github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient"
	"github.com/ArtisanCloud/MediaX/pkg/client/zhihu/core"
	"github.com/ArtisanCloud/MediaX/pkg/client/zhihu/sessionToken/authenticator"
	zhCallback "github.com/ArtisanCloud/MediaX/pkg/client/zhihu/sessionToken/callback"
	zhHarvester "github.com/ArtisanCloud/MediaX/pkg/client/zhihu/sessionToken/harvester"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// MediaX 是统一的媒体平台管理客户端
// 提供了对各个平台的客户端创建和管理功能
type MediaX struct {
	Logger *logger.Logger // 全局日志记录器，用于记录所有平台的日志信息
	Cache  cache.ICache   // 全局缓存接口，用于存储 Token、配置等信息
}

// NewMediaX 初始化 MediaX 实例
//
// 该方法用于创建一个新的 MediaX 实例，用于统一管理各平台的客户端。
// Logger 和 Cache 在所有平台客户端中共享，用于日志记录和数据缓存。
//
// 参数：
//   - config: MediaX 全局配置，包含日志配置等
//   - cache: 缓存实现，用于存储 Token、配置等信息
//
// 返回值：
//   - *MediaX: 初始化完成的 MediaX 实例
func NewMediaX(config *config.MediaXConfig, cache cache.ICache) *MediaX {
	l := logger.NewLogger(config.Logger)
	return &MediaX{
		Logger: l,
		Cache:  cache,
	}
}

// CreateWechatOfficialAccount 创建微信公众号客户端
//
// 该方法使用 ClientToken 认证方式创建微信公众号客户端。
// 支持消息管理、用户管理、素材管理等功能。
//
// 参数：
//   - cfg: 微信公众号配置，包含 AppID、AppSecret 等
//
// 返回值：
//   - *WeChatOfficialAccountCTClient: 公众号客户端实例
//   - error: 创建过程中的错误信息
func (m *MediaX) CreateWechatOfficialAccount(cfg *config.WeChatOfficialAccountConfig) (*officialAccount.WeChatOfficialAccountCTClient, error) {
	return officialAccount.NewWeChatOfficialAccountCTClient(cfg, m.Logger, m.Cache)
}

// CreateGoogleYouTubeACClient 创建 YouTube 客户端
//
// 该方法使用 AccessToken 认证方式创建 YouTube 客户端。
// 支持视频上传、播放列表管理、频道管理等功能。
//
// 参数：
//   - cfg: YouTube 配置，包含 API Key、OAuth 认证信息等
//
// 返回值：
//   - *GoogleYouTubeACClient: YouTube 客户端实例
//   - error: 创建过程中的错误信息
func (m *MediaX) CreateGoogleYouTubeACClient(cfg *config.GoogleYouTubeConfig) (*accessTokenClient3.GoogleYouTubeACClient, error) {
	return accessTokenClient3.NewGoogleYouTubeACClient(cfg, m.Logger, m.Cache)
}

// CreateByteDanceDouYinACClient 创建抖音 AccessToken 客户端
//
// 该方法使用 AccessToken 认证方式创建抖音客户端。
// 支持视频发布、账号管理、数据分析等功能。
//
// 参数：
//   - cfg: 抖音开放平台配置，包含 ClientKey、ClientSecret 等
//
// 返回值：
//   - *ByteDanceDouYinACClient: 抖音 AccessToken 客户端实例
//   - error: 创建过程中的错误信息
func (m *MediaX) CreateByteDanceDouYinACClient(cfg *config.ByteDanceDouYinConfig) (*accessTokenClient.ByteDanceDouYinACClient, error) {
	return accessTokenClient.NewByteDanceDouYinACClient(cfg, m.Logger, m.Cache)
}

// CreateByteDanceDouYinCTClient 创建抖音 ClientToken 客户端
//
// 该方法使用 ClientToken 认证方式创建抖音客户端。
// 主要用于调用不需要用户授权的接口。
//
// 参数：
//   - cfg: 抖音开放平台配置，包含 ClientKey、ClientSecret 等
//
// 返回值：
//   - *ByteDanceDouYinCTClient: 抖音 ClientToken 客户端实例
//   - error: 创建过程中的错误信息
func (m *MediaX) CreateByteDanceDouYinCTClient(cfg *config.ByteDanceDouYinConfig) (*clientTokenClient.ByteDanceDouYinCTClient, error) {
	return clientTokenClient.NewByteDanceDouYinCTClient(cfg, m.Logger, m.Cache)
}

// CreateRedBookJuGuangACClient 创建小红书聚光平台客户端
//
// 该方法使用 AccessToken 认证方式创建小红书客户端。
// 支持笔记发布、数据分析、账号管理等功能。
//
// 参数：
//   - cfg: 小红书聚光平台配置，包含 AppKey、AppSecret 等
//
// 返回值：
//   - *RedBookJuGuangACClient: 小红书客户端实例
//   - error: 创建过程中的错误信息
func (m *MediaX) CreateRedBookJuGuangACClient(cfg *config.RedBookJuGuangConfig) (*accessTokenClient2.RedBookJuGuangACClient, error) {
	return accessTokenClient2.NewRedBookJuGuangACClient(cfg, m.Logger, m.Cache)
}

// CreateZhihuSessionTokenClient 创建知乎 SessionToken 管理器（返回 Manager + Harvester，方便调用方驱动凭证抓取）。
func (m *MediaX) CreateZhihuSessionTokenClient(cfg *config.ZhihuSessionTokenConfig, store sessiontoken.FlowStore) (*sessiontoken.Manager, sessiontoken.CredentialHarvester, error) {
	if cfg == nil {
		return nil, nil, errors.New("zhihu.sessiontoken: config is nil")
	}
	if store == nil {
		return nil, nil, errors.New("zhihu.sessiontoken: flow store is nil")
	}
	coreClient, err := core.NewClient(cfg, m.Logger, m.Cache)
	if err != nil {
		return nil, nil, err
	}
	auth, err := authenticator.NewAuthenticator(&cfg.Authenticator)
	if err != nil {
		return nil, nil, err
	}
	harvester, err := zhHarvester.NewHarvester(&cfg.Harvester)
	if err != nil {
		return nil, nil, err
	}
	dispatcher, err := zhCallback.NewDispatcher(&cfg.Callback, m.Logger)
	if err != nil {
		return nil, nil, err
	}

	manager := sessiontoken.NewManager(
		coreClient.BaseClient,
		m.Logger,
		m.Cache,
		store,
		sessiontoken.WithAuthenticator(auth),
		sessiontoken.WithCallbackDispatcher(dispatcher),
	)
	return manager, harvester, nil
}
