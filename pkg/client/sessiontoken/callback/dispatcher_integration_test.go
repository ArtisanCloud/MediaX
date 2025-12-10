package callback

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestHTTPDispatcherRetriesAndSucceeds(t *testing.T) {
	var attempts int32
	var receivedSignature string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		if atomic.LoadInt32(&attempts) < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		receivedSignature = r.Header.Get("X-MediaX-Signature")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	payload := &Payload{
		FlowID:    "stf_retry",
		State:     "state",
		Status:    "succeeded",
		Timestamp: 1700000000,
		Nonce:     "nonce",
	}
	req := &Request{
		URL:     ts.URL,
		Payload: payload,
	}
	dispatcher := NewHTTPDispatcher("secret", nil, WithRetrySchedule([]time.Duration{time.Millisecond, time.Millisecond}))
	retries, err := dispatcher.Dispatch(context.Background(), req)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if retries != 2 {
		t.Fatalf("expected 2 retries, got %d", retries)
	}
	expectedSig, _ := BuildSignature("secret", payload.Timestamp, payload.Nonce, mustPayloadBody(payload))
	if receivedSignature != expectedSig {
		t.Fatalf("expected signature %s, got %s", expectedSig, receivedSignature)
	}
}

func TestHTTPDispatcherFailsAfterRetries(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	dispatcher := NewHTTPDispatcher("secret", nil, WithRetrySchedule([]time.Duration{time.Millisecond, time.Millisecond}))
	req := &Request{
		URL: ts.URL,
		Payload: &Payload{
			FlowID: "stf_fail",
			State:  "state",
			Status: "failed",
		},
	}
	if _, err := dispatcher.Dispatch(context.Background(), req); err == nil {
		t.Fatalf("expected dispatcher to return error on repeated failure")
	}
}

func mustPayloadBody(p *Payload) []byte {
	body, err := p.Body()
	if err != nil {
		panic(err)
	}
	return body
}
