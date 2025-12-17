package callback

import (
	"encoding/json"
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
	payload := &Payload{
		FlowID:        "stf_1",
		State:         "demo",
		Status:        "pending",
		Timestamp:     1,
		Nonce:         "nonce",
		Code:          "ZH_COOKIE_EXPIRED",
		Message:       "cookie expired",
		LastFailedAPI: "GET /zhihu/api",
		Metadata: map[string]string{
			"last_failed_api": "GET /zhihu/api",
		},
	}
	body, err := payload.Body()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("expected body to be encoded")
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}
	if decoded["code"] != "ZH_COOKIE_EXPIRED" {
		t.Fatalf("expected code to be serialized, got %v", decoded["code"])
	}
	meta, ok := decoded["metadata"].(map[string]interface{})
	if !ok || meta["last_failed_api"] != "GET /zhihu/api" {
		t.Fatalf("expected metadata.last_failed_api to be serialized, got %v", decoded["metadata"])
	}
}
