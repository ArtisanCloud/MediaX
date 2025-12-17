package sessiontoken

import (
	"context"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken/callback"
)

// SessionTokenClient 定义了暴露给外部的 Flow 管理入口。
type SessionTokenClient interface {
	CreateFlow(ctx context.Context, opts *CreateFlowOptions) (*Flow, error)
	GetFlow(ctx context.Context, flowID string) (*Flow, error)
	CompleteFlowSuccess(ctx context.Context, flowID string, credentials *callback.CredentialPayload) (*Flow, error)
	CompleteFlowFailed(ctx context.Context, flowID string, failure *FlowFailure) (*Flow, error)
	MarkTokenInvalid(ctx context.Context, flowID string, code string, message string, api string) (*Flow, error)
	FindFlowBySessionToken(ctx context.Context, token string) (*Flow, error)
}

// ReusableSessionFetcher 定义了可重用 Session 的读取能力。
type ReusableSessionFetcher interface {
	FetchReusableSession(ctx context.Context, flow *Flow) (*callback.CredentialPayload, error)
}

// Authenticator 负责生成 authorize URL 或指令。
type Authenticator interface {
	BuildAuthorizeURL(ctx context.Context, flow *Flow) (string, error)
}

// CredentialHarvester 负责监听并抓取凭证。
type CredentialHarvester interface {
	Watch(ctx context.Context, flow *Flow) (*callback.CredentialPayload, error)
}

// CallbackDispatcher 负责向业务侧回调凭证结果。
type CallbackDispatcher interface {
	Dispatch(ctx context.Context, req *callback.Request) (int, error)
}

// CreateFlowOptions 用于描述 Flow 创建参数。
type CreateFlowOptions struct {
	ProviderCode    string
	ProviderAppCode string
	TenantUUID      string
	AccountID       string
	State           string
	CallbackURL     string
	Metadata        map[string]string
	TTL             time.Duration
}
