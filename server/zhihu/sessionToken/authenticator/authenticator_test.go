package authenticator

import (
	"context"
	"testing"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
)

func TestAuthenticatorSelectsEntryByMetadata(t *testing.T) {
	auth, err := NewAuthenticator(&config.ZhihuSessionTokenAuthenticatorConfig{
		Entries: []config.ZhihuSessionTokenEntryConfig{
			{Type: "pc", URL: "https://www.zhihu.com/signin"},
			{Type: "mobile", URL: "https://m.zhihu.com/signin"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	flow := &sessiontoken.Flow{
		FlowID: "stf_1",
		State:  "state",
		Metadata: map[string]string{
			"login_entry": "mobile",
		},
	}
	url, err := auth.BuildAuthorizeURL(context.Background(), flow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url == "" || url[:19] != "https://m.zhihu.com" {
		t.Fatalf("expected mobile entry, got %s", url)
	}
}

func TestAuthenticatorDefaultsToFirstEntry(t *testing.T) {
	auth, _ := NewAuthenticator(&config.ZhihuSessionTokenAuthenticatorConfig{
		Entries: []config.ZhihuSessionTokenEntryConfig{
			{Type: "pc", URL: "https://www.zhihu.com/signin"},
		},
	})
	flow := &sessiontoken.Flow{FlowID: "stf_1", State: "state"}
	url, err := auth.BuildAuthorizeURL(context.Background(), flow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url == "" {
		t.Fatalf("expected default url")
	}
}
