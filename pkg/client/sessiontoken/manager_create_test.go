package sessiontoken

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken/callback"
	corecache "github.com/ArtisanCloud/MediaXCore/pkg/cache"
)

type fakeFlowStore struct {
	savedFlow      *Flow
	saveTTL        time.Duration
	saveCount      int
	getByStateFlow *Flow
	getByStateErr  error
	getFlow        *Flow
	getErr         error
	updateFlow     *Flow
	updateTTL      time.Duration
	updateCount    int
	updateErr      error
}

func (s *fakeFlowStore) Save(ctx context.Context, flow *Flow, ttl time.Duration) error {
	s.savedFlow = flow
	s.saveTTL = ttl
	s.saveCount++
	return nil
}

func (s *fakeFlowStore) Update(ctx context.Context, flow *Flow, ttl time.Duration) error {
	s.updateFlow = flow
	s.updateTTL = ttl
	s.updateCount++
	return s.updateErr
}

func (s *fakeFlowStore) Get(ctx context.Context, flowID string) (*Flow, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	if s.getFlow != nil {
		return s.getFlow, nil
	}
	return nil, ErrFlowNotFound
}

func (s *fakeFlowStore) GetByState(ctx context.Context, tenantUUID, state string) (*Flow, error) {
	return s.getByStateFlow, s.getByStateErr
}

func (s *fakeFlowStore) Delete(ctx context.Context, flowID string) error { return nil }

type fakeAuthenticator struct {
	url string
	err error
}

func (f *fakeAuthenticator) BuildAuthorizeURL(ctx context.Context, flow *Flow) (string, error) {
	return f.url, f.err
}

type fakeDispatcher struct {
	attempts   int
	err        error
	lastReq    *callback.Request
	callCount  int
	sleepDelay time.Duration
}

func (d *fakeDispatcher) Dispatch(ctx context.Context, req *callback.Request) (int, error) {
	d.callCount++
	d.lastReq = req
	if d.sleepDelay > 0 {
		<-time.After(d.sleepDelay)
	}
	return d.attempts, d.err
}

func TestManagerCreateFlowSuccess(t *testing.T) {
	store := &fakeFlowStore{getByStateErr: ErrFlowNotFound}
	fixedNow := time.Unix(1700000000, 0)
	mgr := NewManager(nil, nil, nil, store, WithAuthenticator(&fakeAuthenticator{url: "https://auth"}), WithClock(func() time.Time { return fixedNow }), WithFlowTTL(10*time.Minute))

	opts := &CreateFlowOptions{
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		State:           "state",
		CallbackURL:     "https://callback",
		AccountID:       "account",
		TTL:             5 * time.Minute,
	}

	flow, err := mgr.CreateFlow(context.Background(), opts)
	if err != nil {
		t.Fatalf("CreateFlow returned error: %v", err)
	}
	if flow.AuthorizeURL != "https://auth" {
		t.Fatalf("unexpected authorize url: %s", flow.AuthorizeURL)
	}
	if flow.State != "state" || flow.TenantUUID != "tenant" {
		t.Fatalf("flow fields not copied from options")
	}
	expectedExpiry := fixedNow.Add(5 * time.Minute)
	if !flow.ExpiresAt.Equal(expectedExpiry) {
		t.Fatalf("expected expires_at %v, got %v", expectedExpiry, flow.ExpiresAt)
	}
	if store.saveCount != 1 {
		t.Fatalf("expected store save to be called once, got %d", store.saveCount)
	}
	if store.saveTTL != 5*time.Minute {
		t.Fatalf("expected ttl 5m, got %v", store.saveTTL)
	}
}

func TestManagerCreateFlowRequiresAuthenticator(t *testing.T) {
	store := &fakeFlowStore{getByStateErr: ErrFlowNotFound}
	mgr := NewManager(nil, nil, nil, store)
	opts := &CreateFlowOptions{
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		State:           "state",
		CallbackURL:     "https://callback",
	}
	if _, err := mgr.CreateFlow(context.Background(), opts); err == nil {
		t.Fatalf("expected error when authenticator missing")
	}
}

func TestManagerCreateFlowDedupByState(t *testing.T) {
	existing := &Flow{FlowID: "stf_existing", ExpiresAt: time.Now().Add(time.Minute)}
	store := &fakeFlowStore{getByStateFlow: existing}
	mgr := NewManager(nil, nil, nil, store, WithAuthenticator(&fakeAuthenticator{url: "https://auth"}))

	opts := &CreateFlowOptions{
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		State:           "state",
		CallbackURL:     "https://callback",
	}

	flow, err := mgr.CreateFlow(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flow != existing {
		t.Fatalf("expected existing flow to be returned")
	}
	if store.saveCount != 0 {
		t.Fatalf("expected no save when returning existing flow")
	}
}

func TestManagerCreateFlowExpiredStateCreatesNew(t *testing.T) {
	fixedNow := time.Unix(1700000500, 0)
	existing := &Flow{FlowID: "stf_old", ExpiresAt: fixedNow.Add(-time.Minute)}
	store := &fakeFlowStore{getByStateFlow: existing}
	mgr := NewManager(nil, nil, nil, store, WithAuthenticator(&fakeAuthenticator{url: "https://auth"}), WithClock(func() time.Time { return fixedNow }))

	opts := &CreateFlowOptions{
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		State:           "state",
		CallbackURL:     "https://callback",
	}

	flow, err := mgr.CreateFlow(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flow == existing {
		t.Fatalf("expected new flow when existing expired")
	}
	if store.saveCount != 1 {
		t.Fatalf("expected Save to be called when creating new flow")
	}
	if flow.ExpiresAt.Before(fixedNow) {
		t.Fatalf("expected expires_at to be after current time")
	}
}

func TestManagerGetFlowSuccess(t *testing.T) {
	store := &fakeFlowStore{}
	now := time.Unix(1800000000, 0)
	expected := &Flow{
		FlowID:       "stf_1",
		ExpiresAt:    now.Add(5 * time.Minute),
		ProviderCode: "zhihu",
		TenantUUID:   "tenant",
	}
	store.getFlow = expected
	mgr := NewManager(nil, nil, nil, store, WithClock(func() time.Time { return now }))

	flow, err := mgr.GetFlow(context.Background(), "stf_1")
	if err != nil {
		t.Fatalf("expected success, got err=%v", err)
	}
	if flow != expected {
		t.Fatalf("expected returned flow to match store")
	}
}

func TestManagerGetFlowExpiredStillAccessible(t *testing.T) {
	store := &fakeFlowStore{}
	now := time.Unix(1800001000, 0)
	expired := &Flow{
		FlowID:    "stf_expired",
		ExpiresAt: now.Add(-time.Minute),
	}
	store.getFlow = expired
	mgr := NewManager(nil, nil, nil, store, WithClock(func() time.Time { return now }))
	flow, err := mgr.GetFlow(context.Background(), "stf_expired")
	if err != nil {
		t.Fatalf("expected expired flow to be returned, got err=%v", err)
	}
	if flow != expired {
		t.Fatalf("expected underlying store flow to be returned")
	}
}

func TestManagerCompleteFlowSuccess(t *testing.T) {
	store := &fakeFlowStore{}
	now := time.Unix(1900000000, 0)
	flow := &Flow{
		FlowID:          "stf_success",
		State:           "state",
		Status:          FlowStatusAuthorizing,
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		CallbackURL:     "https://callback",
		Metadata:        map[string]string{"mode": "password"},
		ExpiresAt:       now.Add(15 * time.Minute),
		ProviderCode:    "zhihu",
	}
	store.getFlow = flow
	dispatcher := &fakeDispatcher{}
	mgr := NewManager(nil, nil, nil, store,
		WithClock(func() time.Time { return now }),
		WithCallbackDispatcher(dispatcher),
	)
	creds := &callback.CredentialPayload{
		SessionToken: "abcdef1234567890",
		Cookies:      []callback.Cookie{{Name: "z_c0", Value: "cookie-value"}},
		Note:         "Chrome 123",
		ExpiresAt:    "2026-01-01T00:00:00Z",
	}

	if _, err := mgr.CompleteFlowSuccess(context.Background(), "stf_success", creds); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flow.Status != FlowStatusSucceeded {
		t.Fatalf("expected status succeeded, got %s", flow.Status)
	}
	if flow.Result == nil || flow.Result.SessionToken == creds.SessionToken {
		t.Fatalf("expected stored result to be masked")
	}
	if dispatcher.lastReq == nil || dispatcher.lastReq.Payload.Credentials.SessionToken != "abcdef1234567890" {
		t.Fatalf("expected dispatcher to receive raw credentials")
	}
	if dispatcher.lastReq.Payload.ProviderCode != "zhihu" || dispatcher.lastReq.Payload.ProviderAppCode != "app" {
		t.Fatalf("expected callback payload to carry provider metadata")
	}
	if dispatcher.lastReq.Payload.TenantUUID != "tenant" {
		t.Fatalf("expected callback payload to carry tenant uuid")
	}
	if flow.CredentialsNote != "Chrome 123" {
		t.Fatalf("expected credentials note to be recorded")
	}
	if flow.CredentialsExpiresHint != "2026-01-01T00:00:00Z" {
		t.Fatalf("expected credentials expires hint to be recorded")
	}
	if v := flow.Metadata[metadataCredentialsNote]; v != "Chrome 123" {
		t.Fatalf("expected metadata to include credentials note")
	}
	if v := flow.Metadata[metadataCredentialsExpiresHint]; v != "2026-01-01T00:00:00Z" {
		t.Fatalf("expected metadata to include expires hint")
	}
	if v := flow.Metadata[metadataSessionTokenKey]; v != "abcdef1234567890" {
		t.Fatalf("expected metadata to include session token copy")
	}
	if store.updateCount != 2 {
		t.Fatalf("expected 2 updates (before and after callback), got %d", store.updateCount)
	}
}

func TestManagerCompleteFlowSuccessCallbackFailure(t *testing.T) {
	store := &fakeFlowStore{}
	now := time.Unix(1900000500, 0)
	flow := &Flow{
		FlowID:          "stf_callback_fail",
		State:           "state",
		Status:          FlowStatusAuthorizing,
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		CallbackURL:     "https://callback",
		ExpiresAt:       now.Add(10 * time.Minute),
	}
	store.getFlow = flow
	dispatcher := &fakeDispatcher{attempts: 2, err: errors.New("http 500")}
	mgr := NewManager(nil, nil, nil, store,
		WithClock(func() time.Time { return now }),
		WithCallbackDispatcher(dispatcher),
	)
	creds := &callback.CredentialPayload{SessionToken: "token"}
	if _, err := mgr.CompleteFlowSuccess(context.Background(), "stf_callback_fail", creds); err == nil {
		t.Fatalf("expected error when dispatcher fails")
	}
	if flow.LastError == "" || !strings.Contains(flow.LastError, "http 500") {
		t.Fatalf("expected flow last error to include dispatcher error, got %s", flow.LastError)
	}
	if flow.RetryAttempts != 2 {
		t.Fatalf("expected retry attempts to be recorded, got %d", flow.RetryAttempts)
	}
}

func TestManagerCompleteFlowFailed(t *testing.T) {
	store := &fakeFlowStore{}
	now := time.Unix(1900000600, 0)
	flow := &Flow{
		FlowID:          "stf_failed",
		State:           "state",
		Status:          FlowStatusAuthorizing,
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		CallbackURL:     "https://callback",
		ExpiresAt:       now.Add(5 * time.Minute),
	}
	store.getFlow = flow
	dispatcher := &fakeDispatcher{}
	mgr := NewManager(nil, nil, nil, store,
		WithClock(func() time.Time { return now }),
		WithCallbackDispatcher(dispatcher),
	)
	failure := &FlowFailure{Reason: "timeout", Code: "ZH_COOKIE_EMPTY", Message: "cookie empty", LastFailedAPI: "GET /me"}
	if _, err := mgr.CompleteFlowFailed(context.Background(), "stf_failed", failure); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flow.Status != FlowStatusFailed {
		t.Fatalf("expected flow status failed")
	}
	if flow.LastError != "timeout" {
		t.Fatalf("expected last_error=timeout, got %s", flow.LastError)
	}
	if flow.Code != "ZH_COOKIE_EMPTY" || flow.Message != "cookie empty" {
		t.Fatalf("expected failure code/message recorded")
	}
	if flow.LastFailedAPI != "GET /me" {
		t.Fatalf("expected last failed api recorded")
	}
	if v := flow.Metadata[metadataLastFailedAPI]; v != "GET /me" {
		t.Fatalf("expected metadata last_failed_api to be recorded")
	}
	if dispatcher.lastReq == nil || dispatcher.lastReq.Payload.Credentials != nil {
		t.Fatalf("expected failure callback without credentials")
	}
	if dispatcher.lastReq.Payload.ProviderCode != "zhihu" || dispatcher.lastReq.Payload.TenantUUID != "tenant" {
		t.Fatalf("expected provider/tenant metadata in failure callback")
	}
}

func TestManagerMarkTokenInvalid(t *testing.T) {
	store := &fakeFlowStore{}
	now := time.Unix(1900000700, 0)
	flow := &Flow{
		FlowID:          "stf_invalid",
		State:           "state",
		Status:          FlowStatusAuthorizing,
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		CallbackURL:     "https://callback",
		ExpiresAt:       now.Add(5 * time.Minute),
	}
	store.getFlow = flow
	dispatcher := &fakeDispatcher{}
	mgr := NewManager(nil, nil, nil, store,
		WithClock(func() time.Time { return now }),
		WithCallbackDispatcher(dispatcher),
	)
	if _, err := mgr.MarkTokenInvalid(context.Background(), "stf_invalid", "ZH_COOKIE_EXPIRED", "cookie expired", "GET https://www.zhihu.com/api/v4/me"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flow.Status != FlowStatusFailed {
		t.Fatalf("expected flow status failed after MarkTokenInvalid")
	}
	if flow.Code != "ZH_COOKIE_EXPIRED" {
		t.Fatalf("expected flow code to be set, got %s", flow.Code)
	}
	if flow.LastFailedAPI != "GET https://www.zhihu.com/api/v4/me" {
		t.Fatalf("expected last_failed_api to be recorded")
	}
	if dispatcher.lastReq == nil || dispatcher.lastReq.Payload.Code != "ZH_COOKIE_EXPIRED" {
		t.Fatalf("expected callback payload to contain failure code")
	}
}

func TestManagerFindFlowBySessionToken(t *testing.T) {
	store := &fakeFlowStore{}
	cacheStore := corecache.NewMemoryCache()
	now := time.Unix(1900000800, 0)
	flow := &Flow{
		FlowID:          "stf_lookup",
		State:           "state",
		Status:          FlowStatusAuthorizing,
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		CallbackURL:     "https://callback",
		ExpiresAt:       now.Add(30 * time.Minute),
	}
	store.getFlow = flow
	dispatcher := &fakeDispatcher{}
	mgr := NewManager(nil, nil, cacheStore, store,
		WithClock(func() time.Time { return now }),
		WithCallbackDispatcher(dispatcher),
	)
	token := "SESSIONID=abc123; JOID=xyz"
	creds := &callback.CredentialPayload{SessionToken: token}
	if _, err := mgr.CompleteFlowSuccess(context.Background(), flow.FlowID, creds); err != nil {
		t.Fatalf("CompleteFlowSuccess returned error: %v", err)
	}
	found, err := mgr.FindFlowBySessionToken(context.Background(), token)
	if err != nil {
		t.Fatalf("FindFlowBySessionToken returned error: %v", err)
	}
	if found == nil || found.FlowID != flow.FlowID {
		t.Fatalf("expected to retrieve flow %s, got %#v", flow.FlowID, found)
	}
}

func TestManagerFindFlowBySessionTokenNotFound(t *testing.T) {
	cacheStore := corecache.NewMemoryCache()
	store := &fakeFlowStore{}
	mgr := NewManager(nil, nil, cacheStore, store)
	if _, err := mgr.FindFlowBySessionToken(context.Background(), ""); !errors.Is(err, ErrFlowNotFound) {
		t.Fatalf("expected ErrFlowNotFound for empty token, got %v", err)
	}
	if _, err := mgr.FindFlowBySessionToken(context.Background(), "SESSIONID=missing"); !errors.Is(err, ErrFlowNotFound) {
		t.Fatalf("expected ErrFlowNotFound for missing mapping, got %v", err)
	}
}

func TestManagerReusableSessionCache(t *testing.T) {
	store := &fakeFlowStore{}
	cacheStore := corecache.NewMemoryCache()
	now := time.Unix(1900000900, 0)
	flow := &Flow{
		FlowID:          "stf_reuse",
		State:           "state",
		Status:          FlowStatusAuthorizing,
		ProviderCode:    "zhihu",
		ProviderAppCode: "app",
		TenantUUID:      "tenant",
		AccountID:       "acct",
		CallbackURL:     "https://callback",
		ExpiresAt:       now.Add(30 * time.Minute),
	}
	store.getFlow = flow
	dispatcher := &fakeDispatcher{}
	mgr := NewManager(nil, nil, cacheStore, store,
		WithClock(func() time.Time { return now }),
		WithCallbackDispatcher(dispatcher),
	)
	creds := &callback.CredentialPayload{
		SessionToken: "SESSIONID=token",
		Cookies: []callback.Cookie{
			{Name: "SESSIONID", Value: "token"},
		},
		Headers: map[string]string{"X-XSRF-TOKEN": "abc"},
	}
	if _, err := mgr.CompleteFlowSuccess(context.Background(), flow.FlowID, creds); err != nil {
		t.Fatalf("CompleteFlowSuccess returned error: %v", err)
	}
	cached, err := mgr.FetchReusableSession(context.Background(), flow)
	if err != nil {
		t.Fatalf("FetchReusableSession returned error: %v", err)
	}
	if cached == nil || cached.SessionToken != creds.SessionToken {
		t.Fatalf("expected cached session token, got %+v", cached)
	}
	if len(cached.Cookies) != 1 || cached.Cookies[0].Name != "SESSIONID" {
		t.Fatalf("expected cookies to be preserved, got %+v", cached.Cookies)
	}
}
