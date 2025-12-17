package sanitizer

import (
	"testing"

	"github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken/callback"
)

func TestMaskString(t *testing.T) {
	if masked := MaskString("abcdefg"); masked != "abc*efg" {
		t.Fatalf("unexpected masked string: %s", masked)
	}
	if masked := MaskString("abc"); masked != "***" {
		t.Fatalf("expected length preserved, got %s", masked)
	}
}

func TestMaskCredentialPayload(t *testing.T) {
	payload := &callback.CredentialPayload{
		SessionToken: "token-123456",
		Cookies:      []callback.Cookie{{Name: "z_c0", Value: "cookie-value"}},
		Headers:      map[string]string{"X-XSRF-TOKEN": "header-value"},
	}
	masked := MaskCredentialPayload(payload)
	if masked.SessionToken == payload.SessionToken {
		t.Fatal("token should be masked")
	}
	if masked.Cookies[0].Value == payload.Cookies[0].Value {
		t.Fatal("cookie should be masked")
	}
	if masked.Headers["X-XSRF-TOKEN"] == "header-value" {
		t.Fatal("header should be masked")
	}
	if payload.SessionToken != "token-123456" {
		t.Fatal("original payload must stay untouched")
	}
}
