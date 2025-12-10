package authenticator

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
)

// Authenticator 负责根据配置构建知乎登录入口 URL。
type Authenticator struct {
	cfg *config.ZhihuSessionTokenAuthenticatorConfig
}

// NewAuthenticator 创建新的知乎 SessionToken Authenticator。
func NewAuthenticator(cfg *config.ZhihuSessionTokenAuthenticatorConfig) (*Authenticator, error) {
	if cfg == nil {
		return nil, errors.New("zhihu.sessiontoken: authenticator config is nil")
	}
	if len(cfg.Entries) == 0 {
		return nil, errors.New("zhihu.sessiontoken: authenticator entries is empty")
	}
	return &Authenticator{cfg: cfg}, nil
}

// BuildAuthorizeURL 依据 metadata 或默认入口返回拉起浏览器的 URL。
func (a *Authenticator) BuildAuthorizeURL(_ context.Context, flow *sessiontoken.Flow) (string, error) {
	if flow == nil {
		return "", errors.New("zhihu.sessiontoken: flow is nil")
	}
	entryType := ""
	if flow.Metadata != nil {
		entryType = flow.Metadata["login_entry"]
		if entryType == "" {
			entryType = flow.Metadata["login_mode"]
		}
	}
	entry := a.findEntry(entryType)
	if entry == nil {
		return "", fmt.Errorf("zhihu.sessiontoken: entry not found for type=%s", entryType)
	}
	u, err := url.Parse(entry.URL)
	if err != nil {
		return "", fmt.Errorf("zhihu.sessiontoken: invalid entry url: %w", err)
	}
	q := u.Query()
	q.Set("flow_id", flow.FlowID)
	q.Set("state", flow.State)
	if flow.AccountID != "" {
		q.Set("account_id", flow.AccountID)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// DefaultUserAgent 返回配置中的默认 UA。
func (a *Authenticator) DefaultUserAgent() string {
	return a.cfg.DefaultUserAgent
}

func (a *Authenticator) findEntry(entryType string) *config.ZhihuSessionTokenEntryConfig {
	if entryType != "" {
		normalized := strings.ToLower(entryType)
		for _, e := range a.cfg.Entries {
			if strings.ToLower(e.Type) == normalized {
				return &e
			}
		}
	}
	if len(a.cfg.Entries) == 0 {
		return nil
	}
	return &a.cfg.Entries[0]
}
