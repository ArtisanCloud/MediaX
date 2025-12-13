package sessionTokenClient

import (
	"context"
	"os"
	"testing"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
	"github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/callback"
)

func TestNewClientDefaultsToV4(t *testing.T) {
	cfg := &config.ZhihuSessionTokenConfig{
		Service: config.ZhihuSessionTokenServiceConfig{},
	}
	client, err := NewClient(cfg, &stubManager{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := client.Version(); got != "v4" {
		t.Fatalf("expected version v4, got %s", got)
	}
}

func TestNewClientEnvOverridesVersion(t *testing.T) {
	t.Setenv("SESSIONTOKEN_ZHIHU_API_VERSION", "  v4 ")
	cfg := &config.ZhihuSessionTokenConfig{}
	client, err := NewClient(cfg, &stubManager{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Version() != "v4" {
		t.Fatalf("expected version v4, got %s", client.Version())
	}
}

func TestNewClientUnsupportedVersion(t *testing.T) {
	t.Setenv("SESSIONTOKEN_ZHIHU_API_VERSION", "")
	cfg := &config.ZhihuSessionTokenConfig{
		Service: config.ZhihuSessionTokenServiceConfig{
			APIVersion: "v5",
		},
	}
	_, err := NewClient(cfg, &stubManager{}, nil)
	if err == nil {
		t.Fatalf("expected error for unsupported version")
	}
}

type stubManager struct{}

func (s *stubManager) CreateFlow(ctx context.Context, opts *sessiontoken.CreateFlowOptions) (*sessiontoken.Flow, error) {
	return nil, nil
}

func (s *stubManager) GetFlow(ctx context.Context, flowID string) (*sessiontoken.Flow, error) {
	return nil, nil
}

func (s *stubManager) CompleteFlowSuccess(ctx context.Context, flowID string, credentials *callback.CredentialPayload) (*sessiontoken.Flow, error) {
	return nil, nil
}

func (s *stubManager) CompleteFlowFailed(ctx context.Context, flowID string, failure *sessiontoken.FlowFailure) (*sessiontoken.Flow, error) {
	return nil, nil
}

func (s *stubManager) MarkTokenInvalid(ctx context.Context, flowID string, code string, message string, api string) (*sessiontoken.Flow, error) {
	return nil, nil
}

func (s *stubManager) FindFlowBySessionToken(ctx context.Context, token string) (*sessiontoken.Flow, error) {
	return nil, nil
}

// Ensure stub satisfies interface at compile-time.
var _ sessiontoken.SessionTokenClient = (*stubManager)(nil)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
