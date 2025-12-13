package sessionTokenClient

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
	v4 "github.com/ArtisanCloud/MediaX/pkg/client/zhihu/web/sessionTokenClient/v4"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// Client 负责根据配置选择对应版本的 Zhihu API handler。
type Client struct {
	version string
	router  versionRouter
}

type versionRouter interface {
	RegisterRoutes(mux *http.ServeMux, apiToken string) error
	Version() string
}

// Option 透传给版本实现，当前等价于 v4.Option。
type Option = v4.Option

// WithHTTPClient 覆盖默认 HTTP Client。
var WithHTTPClient = v4.WithHTTPClient

// WithBaseURL 覆盖默认上游 base url。
var WithBaseURL = v4.WithBaseURL

// WithRetrySchedule 自定义上游重试间隔。
var WithRetrySchedule = v4.WithRetrySchedule

// NewClient 根据配置创建对应版本的 Zhihu API handler。
func NewClient(cfg *config.ZhihuSessionTokenConfig, manager sessiontoken.SessionTokenClient, log *logger.Logger, opts ...Option) (*Client, error) {
	version := resolveAPIVersion(cfg)
	router, err := buildRouter(version, cfg, manager, log, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		version: router.Version(),
		router:  router,
	}, nil
}

// RegisterRoutes 委托给具体版本实现。
func (c *Client) RegisterRoutes(mux *http.ServeMux, apiToken string) error {
	if c == nil {
		return fmt.Errorf("zhihu.sessiontoken: client is nil")
	}
	if c.router == nil {
		return fmt.Errorf("zhihu.sessiontoken: router is nil")
	}
	return c.router.RegisterRoutes(mux, apiToken)
}

// Version 返回当前使用的 API 版本。
func (c *Client) Version() string {
	if c == nil {
		return ""
	}
	return c.version
}

func buildRouter(version string, cfg *config.ZhihuSessionTokenConfig, manager sessiontoken.SessionTokenClient, log *logger.Logger, opts ...Option) (versionRouter, error) {
	switch strings.ToLower(strings.TrimSpace(version)) {
	case "", "v4":
		return v4.NewClient(cfg, manager, log, opts...)
	default:
		return nil, fmt.Errorf("zhihu.sessiontoken: unsupported api version %s", version)
	}
}

func resolveAPIVersion(cfg *config.ZhihuSessionTokenConfig) string {
	if env := strings.TrimSpace(os.Getenv("SESSIONTOKEN_ZHIHU_API_VERSION")); env != "" {
		return env
	}
	if cfg != nil {
		if v := strings.TrimSpace(cfg.Service.APIVersion); v != "" {
			return v
		}
	}
	return "v4"
}
