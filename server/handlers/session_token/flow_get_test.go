package session_token

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken"
	callback "github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken/callback"
)

type flowGetFakeManager struct {
	flow *sessiontoken.Flow
	err  error
}

func (f *flowGetFakeManager) CreateFlow(ctx context.Context, opts *sessiontoken.CreateFlowOptions) (*sessiontoken.Flow, error) {
	return nil, nil
}

func (f *flowGetFakeManager) GetFlow(ctx context.Context, flowID string) (*sessiontoken.Flow, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.flow, nil
}

func (f *flowGetFakeManager) CompleteFlowSuccess(ctx context.Context, flowID string, credentials *callback.CredentialPayload) (*sessiontoken.Flow, error) {
	return nil, nil
}

func (f *flowGetFakeManager) CompleteFlowFailed(ctx context.Context, flowID string, failure *sessiontoken.FlowFailure) (*sessiontoken.Flow, error) {
	return nil, nil
}

func (f *flowGetFakeManager) MarkTokenInvalid(ctx context.Context, flowID string, code string, message string, api string) (*sessiontoken.Flow, error) {
	return nil, nil
}

func (f *flowGetFakeManager) FindFlowBySessionToken(ctx context.Context, token string) (*sessiontoken.Flow, error) {
	return nil, sessiontoken.ErrFlowNotFound
}

func TestSessionTokenFlowGetHandlerSuccess(t *testing.T) {
	fixedNow := time.Unix(1800000000, 0).UTC()
	flow := &sessiontoken.Flow{
		FlowID:          "stf_get",
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		Status:          sessiontoken.FlowStatusPending,
		AuthorizeURL:    "https://auth",
		ExpiresAt:       fixedNow.Add(5 * time.Minute),
	}
	handler := NewSessionTokenFlowGetHandler(&flowGetFakeManager{flow: flow}, nil)
	handler.clock = func() time.Time { return fixedNow }

	req := httptest.NewRequest(http.MethodGet, SessionTokenFlowGetPrefix+flow.FlowID, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if _, ok := body["flow"]; !ok {
		t.Fatalf("expected flow field in response")
	}
}

func TestSessionTokenFlowGetHandlerNotFound(t *testing.T) {
	handler := NewSessionTokenFlowGetHandler(&flowGetFakeManager{err: sessiontoken.ErrFlowNotFound}, nil)
	req := httptest.NewRequest(http.MethodGet, SessionTokenFlowGetPrefix+"invalid", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestSessionTokenFlowGetHandlerMethodNotAllowed(t *testing.T) {
	handler := NewSessionTokenFlowGetHandler(&flowGetFakeManager{}, nil)
	req := httptest.NewRequest(http.MethodPost, SessionTokenFlowGetPrefix+"stf", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestRegisterSessionTokenFlowGetRouteRequiresToken(t *testing.T) {
	mux := http.NewServeMux()
	if err := RegisterSessionTokenFlowGetRoute(mux, &flowGetFakeManager{}, "", nil); err == nil {
		t.Fatalf("expected error when token missing")
	}
}

func TestRegisterSessionTokenFlowGetRouteWithAuth(t *testing.T) {
	flow := &sessiontoken.Flow{FlowID: "stf_auth", ExpiresAt: time.Now().Add(time.Minute)}
	mux := http.NewServeMux()
	if err := RegisterSessionTokenFlowGetRoute(mux, &flowGetFakeManager{flow: flow}, "secret", nil); err != nil {
		t.Fatalf("register route failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, SessionTokenFlowGetPrefix+flow.FlowID, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without header, got %d", rec.Code)
	}
	req = httptest.NewRequest(http.MethodGet, SessionTokenFlowGetPrefix+flow.FlowID, nil)
	req.Header.Set("Authorization", "Bearer secret")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with auth header, got %d", rec.Code)
	}
}
