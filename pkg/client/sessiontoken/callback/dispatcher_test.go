package callback

import (
	"testing"
)

func TestSignerGeneratesDeterministicSignature(t *testing.T) {
	signer := NewSigner("secret")
	body := []byte(`{"flow_id":"stf_1"}`)
	sig, err := signer.Sign(1700000000, "nonce", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "1b2c2a2e56ef8cf2097648be3355b7f622161836fe3ef34fc3f519513661c8f1"
	if sig != expected {
		t.Fatalf("expected %s, got %s", expected, sig)
	}
}

func TestPayloadBody(t *testing.T) {
	payload := &Payload{FlowID: "stf_1", State: "demo", Status: "pending", Timestamp: 1, Nonce: "nonce"}
	body, err := payload.Body()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("expected body to be encoded")
	}
}
