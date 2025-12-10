package session_token

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
	sessionapi "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/api"
	callback "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/callback"
)

type fakeManager struct {
	flow *sessiontoken.Flow
	err  error
}

func (f *fakeManager) CreateFlow(ctx context.Context, opts *sessiontoken.CreateFlowOptions) (*sessiontoken.Flow, error) {
	return f.flow, f.err
}

func (f *fakeManager) GetFlow(ctx context.Context, flowID string) (*sessiontoken.Flow, error) {
	return nil, sessiontoken.ErrFlowNotFound
}

func (f *fakeManager) CompleteFlowSuccess(ctx context.Context, flowID string, credentials *callback.CredentialPayload) (*sessiontoken.Flow, error) {
	return nil, nil
}

func (f *fakeManager) CompleteFlowFailed(ctx context.Context, flowID string, reason string) (*sessiontoken.Flow, error) {
	return nil, nil
}

func TestSessionTokenFlowCreateHandlerSuccess(t *testing.T) {
	flow := &sessiontoken.Flow{FlowID: "stf_1", AuthorizeURL: "https://auth", ExpiresAt: time.Now().Add(time.Minute)}
	handler := NewSessionTokenFlowCreateHandler(&fakeManager{flow: flow}, nil)

	body, _ := json.Marshal(sessionapi.CreateFlowRequest{
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		State:           "state",
		CallbackURL:     "https://callback",
	})
	req := httptest.NewRequest(http.MethodPost, SessionTokenFlowCreatePath, bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestSessionTokenFlowCreateHandlerValidation(t *testing.T) {
	handler := NewSessionTokenFlowCreateHandler(&fakeManager{}, nil)
	req := httptest.NewRequest(http.MethodPost, SessionTokenFlowCreatePath, bytes.NewReader([]byte(`{"provider_code":""}`)))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestSessionTokenFlowCreateHandlerMethodNotAllowed(t *testing.T) {
	handler := NewSessionTokenFlowCreateHandler(&fakeManager{}, nil)
	req := httptest.NewRequest(http.MethodGet, SessionTokenFlowCreatePath, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestRegisterSessionTokenFlowCreateRouteRequiresToken(t *testing.T) {
	mux := http.NewServeMux()
	if err := RegisterSessionTokenFlowCreateRoute(mux, &fakeManager{}, "", nil); err == nil {
		t.Fatalf("expected error for missing token")
	}
}

func TestRegisterSessionTokenFlowCreateRouteWithMiddleware(t *testing.T) {
	flow := &sessiontoken.Flow{FlowID: "stf_2", AuthorizeURL: "https://auth", ExpiresAt: time.Now().Add(time.Minute)}
	mux := http.NewServeMux()
	if err := RegisterSessionTokenFlowCreateRoute(mux, &fakeManager{flow: flow}, "secret", nil); err != nil {
		t.Fatalf("register route failed: %v", err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, SessionTokenFlowCreatePath, bytes.NewReader([]byte(`{}`)))
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	validBody, _ := json.Marshal(sessionapi.CreateFlowRequest{
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		State:           "state",
		CallbackURL:     "https://callback",
	})
	req = httptest.NewRequest(http.MethodPost, SessionTokenFlowCreatePath, bytes.NewReader(validBody))
	req.Header.Set("Authorization", "Bearer secret")
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 with valid token, got %d", rec.Code)
	}
}
