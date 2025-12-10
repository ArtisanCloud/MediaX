package sessiontoken

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/callback"
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

func TestManagerGetFlowExpired(t *testing.T) {
	store := &fakeFlowStore{}
	now := time.Unix(1800001000, 0)
	store.getFlow = &Flow{
		FlowID:    "stf_expired",
		ExpiresAt: now.Add(-time.Minute),
	}
	mgr := NewManager(nil, nil, nil, store, WithClock(func() time.Time { return now }))
	if _, err := mgr.GetFlow(context.Background(), "stf_expired"); err == nil {
		t.Fatalf("expected error for expired flow")
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
	if _, err := mgr.CompleteFlowFailed(context.Background(), "stf_failed", "timeout"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flow.Status != FlowStatusFailed {
		t.Fatalf("expected flow status failed")
	}
	if flow.LastError != "timeout" {
		t.Fatalf("expected last_error=timeout, got %s", flow.LastError)
	}
	if dispatcher.lastReq == nil || dispatcher.lastReq.Payload.Credentials != nil {
		t.Fatalf("expected failure callback without credentials")
	}
	if dispatcher.lastReq.Payload.ProviderCode != "zhihu" || dispatcher.lastReq.Payload.TenantUUID != "tenant" {
		t.Fatalf("expected provider/tenant metadata in failure callback")
	}
}
