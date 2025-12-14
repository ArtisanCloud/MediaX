package v4

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
	"github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/callback"
	zhmiddleware "github.com/ArtisanCloud/MediaX/server/zhihu/sessionToken/middleware"
)

const testBaseURL = "https://example.com"

func TestMeFollowingsSuccess(t *testing.T) {
	stub := &stubSessionTokenManager{
		flow: &sessiontoken.Flow{
			FlowID:          "stf_test",
			TenantUUID:      "tenant_demo",
			ProviderAppCode: "zhihu_article",
			Status:          sessiontoken.FlowStatusSucceeded,
			AccountID:       "acct_user",
		},
	}
	httpClient := newHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if got := req.Header.Get("Cookie"); !strings.Contains(got, "SESSION=abc") {
			t.Fatalf("expected cookie header to propagate session token, got %s", got)
		}
		if req.URL.Path != "/api/v4/members/acct_user/following-columns" {
			t.Fatalf("unexpected upstream path %s", req.URL.Path)
		}
		q := req.URL.Query()
		if q.Get("limit") != "5" || q.Get("offset") != "0" {
			t.Fatalf("unexpected pagination params: %s", req.URL.RawQuery)
		}
		if q.Get("include") != defaultFollowingsInclude {
			t.Fatalf("expected include=%s, got %s", defaultFollowingsInclude, q.Get("include"))
		}
		return jsonResponse(t, http.StatusOK, map[string]any{"data": []string{"foo"}}), nil
	})
	client, err := NewClient(&config.ZhihuSessionTokenConfig{}, stub, nil, WithBaseURL(testBaseURL), WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	mux := http.NewServeMux()
	if err := client.RegisterRoutes(mux, "dev-session-token"); err != nil {
		t.Fatalf("RegisterRoutes failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/zhihu/v1/me/followings?limit=5", nil)
	req.Header.Set("Authorization", "Bearer dev-session-token")
	req.Header.Set(zhmiddleware.HeaderSessionToken, "SESSION=abc")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	var resp apiSuccessResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Meta.FlowID != stub.flow.FlowID {
		t.Fatalf("expected flow id %s, got %s", stub.flow.FlowID, resp.Meta.FlowID)
	}
	if resp.Meta.TenantUUID != stub.flow.TenantUUID {
		t.Fatalf("expected tenant uuid %s, got %s", stub.flow.TenantUUID, resp.Meta.TenantUUID)
	}
	if resp.Meta.ProviderApp != stub.flow.ProviderAppCode {
		t.Fatalf("expected provider app code %s, got %s", stub.flow.ProviderAppCode, resp.Meta.ProviderApp)
	}
	if stub.lastToken != "SESSION=abc" {
		t.Fatalf("expected token lookup to use raw header, got %s", stub.lastToken)
	}
}

func TestMeFollowingsResolveAccountIDFromProfile(t *testing.T) {
	hits := 0
	stub := &stubSessionTokenManager{
		flow: &sessiontoken.Flow{
			FlowID:          "stf_missing",
			TenantUUID:      "tenant_demo",
			ProviderAppCode: "zhihu_article",
			Status:          sessiontoken.FlowStatusSucceeded,
		},
	}
	httpClient := newHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		hits++
		switch req.URL.Path {
		case "/api/v4/me":
			return jsonResponse(t, http.StatusOK, map[string]string{"url_token": "artisancloud"}), nil
		case "/api/v4/members/artisancloud/following-columns":
			return jsonResponse(t, http.StatusOK, map[string]any{"data": []string{"ok"}}), nil
		default:
			t.Fatalf("unexpected upstream path %s", req.URL.Path)
			return nil, nil
		}
	})
	client, err := NewClient(&config.ZhihuSessionTokenConfig{}, stub, nil, WithBaseURL(testBaseURL), WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	mux := http.NewServeMux()
	if err := client.RegisterRoutes(mux, "dev-session-token"); err != nil {
		t.Fatalf("RegisterRoutes failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/zhihu/v1/me/followings", nil)
	req.Header.Set("Authorization", "Bearer dev-session-token")
	req.Header.Set(zhmiddleware.HeaderSessionToken, "SESSION=abc")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	if hits != 2 {
		t.Fatalf("expected two upstream calls (profile + followings), got %d", hits)
	}
}

func TestMeFollowingsAccountIDOverride(t *testing.T) {
	hits := 0
	stub := &stubSessionTokenManager{
		findErr: sessiontoken.ErrFlowNotFound,
	}
	httpClient := newHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		hits++
		if req.URL.Path != "/api/v4/members/custom_user/following-columns" {
			t.Fatalf("unexpected upstream path %s", req.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, map[string]any{"data": []string{"ok"}}), nil
	})
	client, err := NewClient(&config.ZhihuSessionTokenConfig{}, stub, nil, WithBaseURL(testBaseURL), WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	mux := http.NewServeMux()
	if err := client.RegisterRoutes(mux, "dev-session-token"); err != nil {
		t.Fatalf("RegisterRoutes failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/zhihu/v1/me/followings?account_id=custom_user&limit=2", nil)
	req.Header.Set("Authorization", "Bearer dev-session-token")
	req.Header.Set(zhmiddleware.HeaderSessionToken, "SESSION=cookie")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	if hits != 1 {
		t.Fatalf("expected only followings call, got %d hits", hits)
	}
	var resp apiSuccessResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Meta.FlowID != "" {
		t.Fatalf("expected empty flow id when manager has no flow")
	}
}

func TestMeFollowingsProfileUnauthorized(t *testing.T) {
	stub := &stubSessionTokenManager{
		flow: &sessiontoken.Flow{
			FlowID:          "stf_missing",
			TenantUUID:      "tenant_demo",
			ProviderAppCode: "zhihu_article",
			Status:          sessiontoken.FlowStatusSucceeded,
		},
	}
	httpClient := newHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/api/v4/me" {
			return jsonResponse(t, http.StatusUnauthorized, map[string]string{"message": "invalid"}), nil
		}
		t.Fatalf("unexpected upstream path %s", req.URL.Path)
		return nil, nil
	})
	client, err := NewClient(&config.ZhihuSessionTokenConfig{}, stub, nil, WithBaseURL(testBaseURL), WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	mux := http.NewServeMux()
	if err := client.RegisterRoutes(mux, "dev-session-token"); err != nil {
		t.Fatalf("RegisterRoutes failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/zhihu/v1/me/followings", nil)
	req.Header.Set("Authorization", "Bearer dev-session-token")
	req.Header.Set(zhmiddleware.HeaderSessionToken, "SESSION=bad")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
	}
	var resp apiError
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Code != codeCookieExpired {
		t.Fatalf("expected %s, got %s", codeCookieExpired, resp.Code)
	}
}

func TestChannelsArticlesForbiddenMarksInvalid(t *testing.T) {
	stub := &stubSessionTokenManager{
		flow: &sessiontoken.Flow{
			FlowID:          "stf_401",
			TenantUUID:      "tenant_a",
			ProviderAppCode: "zhihu_article",
			Status:          sessiontoken.FlowStatusSucceeded,
		},
	}
	httpClient := newHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		return jsonResponse(t, http.StatusForbidden, map[string]string{"message": "risk detected"}), nil
	})
	client, err := NewClient(&config.ZhihuSessionTokenConfig{}, stub, nil, WithBaseURL(testBaseURL), WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	mux := http.NewServeMux()
	if err := client.RegisterRoutes(mux, "dev-session-token"); err != nil {
		t.Fatalf("RegisterRoutes failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/zhihu/v1/channels/abc/articles", nil)
	req.Header.Set("Authorization", "Bearer dev-session-token")
	req.Header.Set(zhmiddleware.HeaderSessionToken, "SESSION=def")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", rec.Code, rec.Body.String())
	}
	var resp apiError
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Code != codeRiskBlock {
		t.Fatalf("expected code %s, got %s", codeRiskBlock, resp.Code)
	}
	if len(stub.markCalls) != 1 {
		t.Fatalf("expected MarkTokenInvalid to be called once, got %d", len(stub.markCalls))
	}
	if stub.markCalls[0].Code != codeRiskBlock {
		t.Fatalf("expected failure code recorded, got %+v", stub.markCalls[0])
	}
}

func TestArticlePostValidationError(t *testing.T) {
	stub := &stubSessionTokenManager{
		flow: &sessiontoken.Flow{
			FlowID:          "stf_post",
			TenantUUID:      "tenant",
			ProviderAppCode: "zhihu_article",
			Status:          sessiontoken.FlowStatusSucceeded,
		},
	}
	client, err := NewClient(&config.ZhihuSessionTokenConfig{}, stub, nil)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	mux := http.NewServeMux()
	if err := client.RegisterRoutes(mux, "dev-session-token"); err != nil {
		t.Fatalf("RegisterRoutes failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/zhihu/v1/articles", strings.NewReader(`{"title":""}`))
	req.Header.Set("Authorization", "Bearer dev-session-token")
	req.Header.Set(zhmiddleware.HeaderSessionToken, "SESSION=xyz")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
	if len(stub.markCalls) != 0 {
		t.Fatalf("expected mark invalid not to be called on validation error")
	}
}

func TestSanityCheckFlowMissing(t *testing.T) {
	stub := &stubSessionTokenManager{
		findErr: sessiontoken.ErrFlowNotFound,
	}
	httpClient := newHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/api/v4/me" {
			t.Fatalf("unexpected upstream path %s", req.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, map[string]string{"ok": "true"}), nil
	})
	client, err := NewClient(&config.ZhihuSessionTokenConfig{}, stub, nil, WithBaseURL(testBaseURL), WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	mux := http.NewServeMux()
	if err := client.RegisterRoutes(mux, "dev-session-token"); err != nil {
		t.Fatalf("RegisterRoutes failed: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/zhihu/v1/sanity/check", nil)
	req.Header.Set("Authorization", "Bearer dev-session-token")
	req.Header.Set(zhmiddleware.HeaderSessionToken, "SESSION=xyz")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestCallZhihuRetriesWithAPIConfig(t *testing.T) {
	hits := 0
	stub := &stubSessionTokenManager{}
	cfg := &config.ZhihuSessionTokenConfig{
		API: config.ZhihuSessionTokenAPIConfig{
			RetryBackoff: []int{1},
		},
	}
	httpClient := newHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		hits++
		if hits == 1 {
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			}, nil
		}
		return jsonResponse(t, http.StatusOK, map[string]string{"ok": "true"}), nil
	})
	client, err := NewClient(cfg, stub, nil,
		WithBaseURL(testBaseURL),
		WithHTTPClient(httpClient),
		WithRetrySchedule([]time.Duration{0}),
	)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	status, _, _, err := client.callZhihu(context.Background(), http.MethodGet, "/foo", url.Values{}, nil, "SESSION=abc")
	if err != nil {
		t.Fatalf("callZhihu returned error: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}
	if hits != 2 {
		t.Fatalf("expected upstream to be hit twice, got %d", hits)
	}
}

type stubSessionTokenManager struct {
	flow      *sessiontoken.Flow
	findErr   error
	lastToken string
	markCalls []markInvalidCall
}

type markInvalidCall struct {
	FlowID  string
	Code    string
	Message string
	API     string
}

func (s *stubSessionTokenManager) CreateFlow(ctx context.Context, opts *sessiontoken.CreateFlowOptions) (*sessiontoken.Flow, error) {
	return nil, errors.New("not implemented")
}

func (s *stubSessionTokenManager) GetFlow(ctx context.Context, flowID string) (*sessiontoken.Flow, error) {
	return nil, errors.New("not implemented")
}

func (s *stubSessionTokenManager) CompleteFlowSuccess(ctx context.Context, flowID string, credentials *callback.CredentialPayload) (*sessiontoken.Flow, error) {
	return nil, errors.New("not implemented")
}

func (s *stubSessionTokenManager) CompleteFlowFailed(ctx context.Context, flowID string, failure *sessiontoken.FlowFailure) (*sessiontoken.Flow, error) {
	return nil, errors.New("not implemented")
}

func (s *stubSessionTokenManager) MarkTokenInvalid(ctx context.Context, flowID string, code string, message string, api string) (*sessiontoken.Flow, error) {
	s.markCalls = append(s.markCalls, markInvalidCall{FlowID: flowID, Code: code, Message: message, API: api})
	return s.flow, nil
}

func (s *stubSessionTokenManager) FindFlowBySessionToken(ctx context.Context, token string) (*sessiontoken.Flow, error) {
	s.lastToken = token
	if s.findErr != nil {
		return nil, s.findErr
	}
	if s.flow == nil {
		return nil, sessiontoken.ErrFlowNotFound
	}
	return s.flow, nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newHTTPClient(t *testing.T, fn roundTripFunc) *http.Client {
	t.Helper()
	if fn == nil {
		t.Fatalf("round trip func is nil")
	}
	return &http.Client{
		Transport: roundTripFunc(fn),
		Timeout:   5 * time.Second,
	}
}

func jsonResponse(t *testing.T, status int, body interface{}) *http.Response {
	t.Helper()
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal json: %v", err)
		}
	}
	resp := &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(payload)),
	}
	resp.Header.Set("Content-Type", "application/json")
	return resp
}
