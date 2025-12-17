package harvester

import (
	"context"
	"testing"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken"
)

func TestHarvesterBuildsPayloadFromMetadata(t *testing.T) {
	h, err := NewHarvester(&config.ZhihuSessionTokenHarvesterConfig{
		WatchCookies: []string{"z_c0"},
		WatchHeaders: []string{"X-XSRF-TOKEN"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h.clock = func() time.Time { return time.Unix(1700000000, 0) }
	flow := &sessiontoken.Flow{
		Metadata: map[string]string{
			"session_token":          "raw_token",
			"cookie_z_c0":            "cookie-value",
			"header_x-xsrf-token":    "header-value",
			"credentials_expires_at": "2025-01-01T00:00:00Z",
		},
	}
	payload, err := h.Watch(context.Background(), flow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payload.SessionToken != "raw_token" {
		t.Fatalf("expected raw token, got %s", payload.SessionToken)
	}
	if len(payload.Cookies) != 1 || payload.Cookies[0].Value != "cookie-value" {
		t.Fatalf("unexpected cookies %#v", payload.Cookies)
	}
	if payload.Headers["X-XSRF-TOKEN"] != "header-value" {
		t.Fatalf("unexpected headers %#v", payload.Headers)
	}
	if payload.CapturedAt == "" {
		t.Fatalf("expected captured_at to be set")
	}
}
