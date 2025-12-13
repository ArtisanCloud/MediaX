package sessiontoken

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/callback"
)

type orchestratorManagerStub struct {
	flow         *Flow
	successCh    chan struct{}
	failureCh    chan *FlowFailure
	mu           sync.Mutex
	calls        int
	reusePayload *callback.CredentialPayload
}

func (m *orchestratorManagerStub) CreateFlow(ctx context.Context, opts *CreateFlowOptions) (*Flow, error) {
	m.mu.Lock()
	m.calls++
	defer m.mu.Unlock()
	return m.cloneFlow(), nil
}

func (m *orchestratorManagerStub) GetFlow(ctx context.Context, flowID string) (*Flow, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.flow == nil {
		return nil, ErrFlowNotFound
	}
	return m.cloneFlow(), nil
}

func (m *orchestratorManagerStub) CompleteFlowSuccess(ctx context.Context, flowID string, credentials *callback.CredentialPayload) (*Flow, error) {
	if m.successCh != nil {
		select {
		case <-m.successCh:
		default:
			close(m.successCh)
		}
	}
	return m.flow, nil
}

func (m *orchestratorManagerStub) CompleteFlowFailed(ctx context.Context, flowID string, failure *FlowFailure) (*Flow, error) {
	if m.failureCh != nil {
		m.failureCh <- failure
	}
	return m.flow, nil
}

func (m *orchestratorManagerStub) MarkTokenInvalid(ctx context.Context, flowID string, code string, message string, api string) (*Flow, error) {
	if m.failureCh != nil {
		m.failureCh <- &FlowFailure{Reason: message, Code: code, LastFailedAPI: api}
	}
	return m.flow, nil
}

func (m *orchestratorManagerStub) FindFlowBySessionToken(ctx context.Context, token string) (*Flow, error) {
	return m.flow, nil
}

func (m *orchestratorManagerStub) FetchReusableSession(ctx context.Context, flow *Flow) (*callback.CredentialPayload, error) {
	if m.reusePayload == nil {
		return nil, ErrFlowNotFound
	}
	return m.reusePayload, nil
}

func (m *orchestratorManagerStub) cloneFlow() *Flow {
	if m.flow == nil {
		return nil
	}
	cloned := *m.flow
	if len(m.flow.Metadata) > 0 {
		meta := make(map[string]string, len(m.flow.Metadata))
		for k, v := range m.flow.Metadata {
			meta[k] = v
		}
		cloned.Metadata = meta
	}
	return &cloned
}

func (m *orchestratorManagerStub) setMetadata(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.flow.Metadata == nil {
		m.flow.Metadata = make(map[string]string)
	}
	if strings.TrimSpace(value) == "" {
		delete(m.flow.Metadata, key)
		return
	}
	m.flow.Metadata[key] = value
}

type stubHarvester struct {
	payload *callback.CredentialPayload
	err     error
	mu      sync.Mutex
	calls   int
}

func (h *stubHarvester) Watch(ctx context.Context, flow *Flow) (*callback.CredentialPayload, error) {
	h.mu.Lock()
	h.calls++
	h.mu.Unlock()
	if h.err != nil {
		return nil, h.err
	}
	return h.payload, nil
}

func TestFlowOrchestratorTriggersHarvesterSuccess(t *testing.T) {
	flow := &Flow{FlowID: "stf_success", Status: FlowStatusPending}
	successCh := make(chan struct{})
	manager := &orchestratorManagerStub{flow: flow, successCh: successCh}
	h := &stubHarvester{payload: &callback.CredentialPayload{SessionToken: "token"}}
	orch := NewFlowOrchestrator(context.Background(), manager, h, nil)
	orch.metadataWaitTimeout = 0

	if _, err := orch.CreateFlow(context.Background(), &CreateFlowOptions{}); err != nil {
		t.Fatalf("CreateFlow returned error: %v", err)
	}
	select {
	case <-successCh:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for CompleteFlowSuccess")
	}
	if h.calls != 1 {
		t.Fatalf("expected harvester to be called once, got %d", h.calls)
	}
}

func TestFlowOrchestratorMarksFailureWhenHarvesterErrors(t *testing.T) {
	flow := &Flow{FlowID: "stf_failed", Status: FlowStatusPending}
	failureCh := make(chan *FlowFailure, 1)
	manager := &orchestratorManagerStub{flow: flow, failureCh: failureCh}
	h := &stubHarvester{err: fmt.Errorf("metadata.session_token is empty")}
	orch := NewFlowOrchestrator(context.Background(), manager, h, nil)
	orch.metadataWaitTimeout = 0

	if _, err := orch.CreateFlow(context.Background(), &CreateFlowOptions{}); err != nil {
		t.Fatalf("CreateFlow returned error: %v", err)
	}
	select {
	case failure := <-failureCh:
		if failure == nil || !strings.Contains(failure.Reason, "metadata.session_token") {
			t.Fatalf("unexpected failure reason: %+v", failure)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for CompleteFlowFailed")
	}
}

func TestFlowOrchestratorAvoidsDuplicateRuns(t *testing.T) {
	flow := &Flow{FlowID: "stf_dup", Status: FlowStatusPending}
	manager := &orchestratorManagerStub{flow: flow}
	h := &stubHarvester{payload: &callback.CredentialPayload{SessionToken: "token"}}
	orch := NewFlowOrchestrator(context.Background(), manager, h, nil)
	orch.metadataWaitTimeout = 0

	if _, err := orch.CreateFlow(context.Background(), &CreateFlowOptions{}); err != nil {
		t.Fatalf("CreateFlow returned error: %v", err)
	}
	if _, err := orch.CreateFlow(context.Background(), &CreateFlowOptions{}); err != nil {
		t.Fatalf("second CreateFlow returned error: %v", err)
	}
	time.Sleep(50 * time.Millisecond)
	if h.calls != 1 {
		t.Fatalf("expected harvester to be called once, got %d", h.calls)
	}
}

func TestFlowOrchestratorReusesSessionWhenAvailable(t *testing.T) {
	flow := &Flow{
		FlowID: "stf_reuse",
		Status: FlowStatusPending,
		Metadata: map[string]string{
			metadataReuseSessionKey: "true",
		},
	}
	successCh := make(chan struct{})
	payload := &callback.CredentialPayload{SessionToken: "cached_token"}
	manager := &orchestratorManagerStub{flow: flow, successCh: successCh, reusePayload: payload}
	h := &stubHarvester{payload: &callback.CredentialPayload{SessionToken: "new"}}
	orch := NewFlowOrchestrator(context.Background(), manager, h, nil)

	if _, err := orch.CreateFlow(context.Background(), &CreateFlowOptions{}); err != nil {
		t.Fatalf("CreateFlow returned error: %v", err)
	}
	select {
	case <-successCh:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for reuse CompleteFlowSuccess")
	}
	if h.calls != 0 {
		t.Fatalf("expected harvester not to run when reuse succeeds, got %d", h.calls)
	}
}

func TestFlowOrchestratorWaitsForMetadataBeforeHarvest(t *testing.T) {
	flow := &Flow{
		FlowID: "stf_wait",
		Status: FlowStatusPending,
		Metadata: map[string]string{
			metadataReuseSessionKey: "false",
		},
	}
	successCh := make(chan struct{})
	manager := &orchestratorManagerStub{flow: flow, successCh: successCh}
	h := &stubHarvester{payload: &callback.CredentialPayload{SessionToken: "token"}}
	orch := NewFlowOrchestrator(context.Background(), manager, h, nil)
	orch.metadataPollInterval = 10 * time.Millisecond
	orch.metadataWaitTimeout = 200 * time.Millisecond

	go func() {
		time.Sleep(50 * time.Millisecond)
		manager.setMetadata(metadataSessionTokenKey, "cookie=abc")
	}()

	if _, err := orch.CreateFlow(context.Background(), &CreateFlowOptions{}); err != nil {
		t.Fatalf("CreateFlow returned error: %v", err)
	}
	select {
	case <-successCh:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for CompleteFlowSuccess after metadata ready")
	}
	if h.calls != 1 {
		t.Fatalf("expected harvester to run once after metadata ready, got %d", h.calls)
	}
}

func TestFlowOrchestratorFailsWhenMetadataTimeout(t *testing.T) {
	flow := &Flow{FlowID: "stf_timeout", Status: FlowStatusPending}
	failureCh := make(chan *FlowFailure, 1)
	manager := &orchestratorManagerStub{flow: flow, failureCh: failureCh}
	h := &stubHarvester{payload: &callback.CredentialPayload{SessionToken: "unused"}}
	orch := NewFlowOrchestrator(context.Background(), manager, h, nil)
	orch.metadataPollInterval = 5 * time.Millisecond
	orch.metadataWaitTimeout = 30 * time.Millisecond

	if _, err := orch.CreateFlow(context.Background(), &CreateFlowOptions{}); err != nil {
		t.Fatalf("CreateFlow returned error: %v", err)
	}
	select {
	case failure := <-failureCh:
		if failure == nil || !strings.Contains(failure.Reason, "metadata session_token") {
			t.Fatalf("unexpected failure: %+v", failure)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for metadata wait failure")
	}
	if h.calls != 0 {
		t.Fatalf("expected harvester not to run when metadata missing, got %d", h.calls)
	}
}
