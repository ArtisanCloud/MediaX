package sessiontoken

import (
	"testing"
	"time"
)

func TestStateMachineAllowsValidTransitions(t *testing.T) {
	sm := NewStateMachine()
	fixed := time.Unix(1700000000, 0)
	sm.WithClock(func() time.Time { return fixed })

	flow := &Flow{Status: FlowStatusPending}
	if err := sm.Transition(flow, FlowStatusAuthorizing); err != nil {
		t.Fatalf("expected pending->authorizing allowed, got %v", err)
	}
	if !flow.UpdatedAt.Equal(fixed.UTC()) {
		t.Fatalf("expected updated timestamp, got %v", flow.UpdatedAt)
	}

	if err := sm.Transition(flow, FlowStatusSucceeded); err != nil {
		t.Fatalf("expected authorizing->succeeded allowed, got %v", err)
	}
}

func TestStateMachineRejectsInvalidTransitions(t *testing.T) {
	sm := NewStateMachine()
	flow := &Flow{Status: FlowStatusPending}

	if err := sm.Transition(flow, FlowStatusSucceeded); err == nil {
		t.Fatal("expected pending->succeeded to be rejected")
	}
}

func TestStateMachineAllowsIdempotentTransitions(t *testing.T) {
	sm := NewStateMachine()
	flow := &Flow{Status: FlowStatusPending}

	if err := sm.Transition(flow, FlowStatusPending); err != nil {
		t.Fatalf("expected noop transition, got %v", err)
	}
}
