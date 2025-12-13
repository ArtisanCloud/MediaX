package harvester

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
	callbackpkg "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/callback"
)

// Harvester 根据 Flow.Metadata 构建凭证 payload。
type Harvester struct {
	cfg   *config.ZhihuSessionTokenHarvesterConfig
	clock func() time.Time
}

// NewHarvester 创建知乎 CredentialHarvester。
func NewHarvester(cfg *config.ZhihuSessionTokenHarvesterConfig) (*Harvester, error) {
	if cfg == nil {
		return nil, errors.New("zhihu.sessiontoken: harvester config is nil")
	}
	return &Harvester{
		cfg:   cfg,
		clock: time.Now,
	}, nil
}

// Watch 尝试从 Flow metadata 中构建凭证，该方法不会主动与浏览器交互，主要用于示例/测试。
func (h *Harvester) Watch(_ context.Context, flow *sessiontoken.Flow) (*callbackpkg.CredentialPayload, error) {
	if flow == nil {
		return nil, errors.New("zhihu.sessiontoken: flow is nil")
	}
	token := metadataValue(flow, "session_token")
	if token == "" {
		return nil, errors.New("zhihu.sessiontoken: metadata.session_token is empty")
	}
	payload := &callbackpkg.CredentialPayload{
		SessionToken: token,
		Note:         metadataValue(flow, "credentials_note"),
		ExpiresAt:    metadataValue(flow, "credentials_expires_at"),
		CapturedAt:   h.clock().UTC().Format(time.RFC3339),
		Cookies:      h.collectCookies(flow),
		Headers:      h.collectHeaders(flow),
	}
	return payload, nil
}

func (h *Harvester) collectCookies(flow *sessiontoken.Flow) []callbackpkg.Cookie {
	var cookies []callbackpkg.Cookie
	for _, name := range h.cfg.WatchCookies {
		value := metadataValue(flow, cookieMetadataKey(name))
		if value == "" {
			continue
		}
		cookies = append(cookies, callbackpkg.Cookie{
			Name:  name,
			Value: value,
		})
	}
	return cookies
}

func (h *Harvester) collectHeaders(flow *sessiontoken.Flow) map[string]string {
	headers := make(map[string]string, len(h.cfg.WatchHeaders))
	for _, name := range h.cfg.WatchHeaders {
		value := metadataValue(flow, headerMetadataKey(name))
		if value == "" {
			continue
		}
		headers[name] = value
	}
	if len(headers) == 0 {
		return nil
	}
	return headers
}

func metadataValue(flow *sessiontoken.Flow, key string) string {
	if flow == nil || len(flow.Metadata) == 0 {
		return ""
	}
	return strings.TrimSpace(flow.Metadata[key])
}

func cookieMetadataKey(name string) string {
	return fmt.Sprintf("cookie_%s", strings.ToLower(name))
}

func headerMetadataKey(name string) string {
	return fmt.Sprintf("header_%s", strings.ToLower(name))
}
