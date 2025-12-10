package sessiontoken

import (
	"context"
	"testing"
	"time"
)

type fakeFlowStore struct {
	savedFlow      *Flow
	saveTTL        time.Duration
	saveCount      int
	getByStateFlow *Flow
	getByStateErr  error
}

func (s *fakeFlowStore) Save(ctx context.Context, flow *Flow, ttl time.Duration) error {
	s.savedFlow = flow
	s.saveTTL = ttl
	s.saveCount++
	return nil
}

func (s *fakeFlowStore) Update(ctx context.Context, flow *Flow, ttl time.Duration) error { return nil }

func (s *fakeFlowStore) Get(ctx context.Context, flowID string) (*Flow, error) {
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
