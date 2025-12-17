package main

import (
	"testing"

	app "github.com/ArtisanCloud/MediaX/cmd/accesstoken/internal/app"
	"github.com/ArtisanCloud/MediaX/pkg/client/config"
)

func TestValidateOptions(t *testing.T) {
	t.Run("videos list requires ids or chart", func(t *testing.T) {
		opts := baseOptions()
		opts.IDs = ""
		opts.Chart = ""
		opts.Normalize()
		if err := opts.Validate(); err == nil || err.Error() != "videos.list: provide at least ids or chart" {
			t.Fatalf("expected ids/chart validation error, got %v", err)
		}
	})

	t.Run("missing action", func(t *testing.T) {
		opts := baseOptions()
		opts.Action = ""
		opts.Normalize()
		if err := opts.Validate(); err == nil || err.Error() != "missing required field: action" {
			t.Fatalf("expected action error, got %v", err)
		}
	})

	t.Run("missing part", func(t *testing.T) {
		opts := baseOptions()
		opts.Part = ""
		opts.Normalize()
		if err := opts.Validate(); err == nil || err.Error() != "missing required field: part" {
			t.Fatalf("expected part error, got %v", err)
		}
	})

	t.Run("invalid action", func(t *testing.T) {
		opts := baseOptions()
		opts.Action = "videos.get"
		opts.Normalize()
		if err := opts.Validate(); err == nil || err.Error() != `unsupported action "videos.get"` {
			t.Fatalf("expected unsupported action error, got %v", err)
		}
	})

	t.Run("invalid max results", func(t *testing.T) {
		opts := baseOptions()
		opts.MaxResults = 100
		opts.Normalize()
		if err := opts.Validate(); err == nil || err.Error() != "invalid max_results 100: must be between 1 and 50" {
			t.Fatalf("expected max results error, got %v", err)
		}
	})

	t.Run("search must have query or channel or mine", func(t *testing.T) {
		opts := baseOptions()
		opts.Action = "search.list"
		opts.IDs = ""
		opts.Query = ""
		opts.ChannelID = ""
		opts.SearchMine = false
		opts.Normalize()
		if err := opts.Validate(); err == nil || err.Error() != "search.list: provide query, channel_id, or search_mine" {
			t.Fatalf("expected search validation error, got %v", err)
		}
	})

	t.Run("playlists requires ids channel or mine", func(t *testing.T) {
		opts := baseOptions()
		opts.Action = "playlists.list"
		opts.IDs = ""
		opts.ChannelID = ""
		opts.Mine = false
		opts.Normalize()
		if err := opts.Validate(); err == nil || err.Error() != "playlists.list: provide at least ids, channel_id, or mine" {
			t.Fatalf("expected playlist validation error, got %v", err)
		}
	})

	t.Run("playlists mine and ids mutually exclusive", func(t *testing.T) {
		opts := baseOptions()
		opts.Action = "playlists.list"
		opts.IDs = "PL123"
		opts.Mine = true
		opts.Normalize()
		if err := opts.Validate(); err == nil || err.Error() != "playlists.list: mine and ids are mutually exclusive" {
			t.Fatalf("expected mutual exclusivity error, got %v", err)
		}
	})

	t.Run("valid videos options", func(t *testing.T) {
		opts := baseOptions()
		opts.Normalize()
		if err := opts.Validate(); err != nil {
			t.Fatalf("expected valid options, got %v", err)
		}
	})

	t.Run("valid search options", func(t *testing.T) {
		opts := baseOptions()
		opts.Action = "search.list"
		opts.Query = "MediaX"
		opts.IDs = ""
		opts.Normalize()
		if err := opts.Validate(); err != nil {
			t.Fatalf("expected valid search options, got %v", err)
		}
	})

	t.Run("valid playlist options with mine", func(t *testing.T) {
		opts := baseOptions()
		opts.Action = "playlists.list"
		opts.IDs = ""
		opts.Mine = true
		opts.Normalize()
		if err := opts.Validate(); err != nil {
			t.Fatalf("expected valid playlist options, got %v", err)
		}
	})
}

func TestResolveAccessTokenPriority(t *testing.T) {
	cfg := &config.GoogleYouTubeConfig{
		ClientConfig: &config.ClientConfig{
			OAuthConfig: &config.OAuthConfig{
				AccessToken: "config-token",
			},
		},
	}

	t.Run("flag overrides env and config", func(t *testing.T) {
		t.Setenv("GOOGLE_YOUTUBE_ACCESS_TOKEN", "env-token")
		token, source := app.ResolveAccessToken("flag-token", "google_youtube", cfg.ClientConfig)
		if token != "flag-token" || source != "flag" {
			t.Fatalf("expected flag token, got token=%q source=%q", token, source)
		}
	})

	t.Run("env overrides config", func(t *testing.T) {
		t.Setenv("GOOGLE_YOUTUBE_ACCESS_TOKEN", "env-token")
		token, source := app.ResolveAccessToken("", "google_youtube", cfg.ClientConfig)
		if token != "env-token" || source != "env:GOOGLE_YOUTUBE_ACCESS_TOKEN" {
			t.Fatalf("expected env token, got token=%q source=%q", token, source)
		}
	})

	t.Run("config fallback", func(t *testing.T) {
		t.Setenv("GOOGLE_YOUTUBE_ACCESS_TOKEN", "")
		token, source := app.ResolveAccessToken("", "google_youtube", cfg.ClientConfig)
		if token != "config-token" || source != "config" {
			t.Fatalf("expected config token, got token=%q source=%q", token, source)
		}
	})

	t.Run("empty when no sources", func(t *testing.T) {
		t.Setenv("GOOGLE_YOUTUBE_ACCESS_TOKEN", "")
		token, source := app.ResolveAccessToken("", "google_youtube", &config.ClientConfig{})
		if token != "" || source != "" {
			t.Fatalf("expected empty token, got token=%q source=%q", token, source)
		}
	})
}

func baseOptions() *app.Options {
	return &app.Options{
		ConfigPath: "/tmp/config.yaml",
		Action:     "videos.list",
		Part:       "snippet",
		IDs:        "abc123",
		MaxResults: 5,
	}
}
