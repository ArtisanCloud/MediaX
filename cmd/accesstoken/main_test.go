package main

import (
	"testing"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
)

func TestValidateOptions(t *testing.T) {
	t.Run("videos list requires ids or chart", func(t *testing.T) {
		opts := baseOptions()
		opts.ids = ""
		opts.chart = ""
		if err := validateOptions(opts); err == nil || err.Error() != "videos.list: provide at least -ids or -chart" {
			t.Fatalf("expected ids/chart validation error, got %v", err)
		}
	})

	t.Run("missing action", func(t *testing.T) {
		opts := baseOptions()
		opts.action = ""
		if err := validateOptions(opts); err == nil || err.Error() != "missing required flag: -action" {
			t.Fatalf("expected action error, got %v", err)
		}
	})

	t.Run("missing part", func(t *testing.T) {
		opts := baseOptions()
		opts.part = ""
		if err := validateOptions(opts); err == nil || err.Error() != "missing required flag: -part" {
			t.Fatalf("expected part error, got %v", err)
		}
	})

	t.Run("invalid action", func(t *testing.T) {
		opts := baseOptions()
		opts.action = "videos.get"
		if err := validateOptions(opts); err == nil || err.Error() != `unsupported action "videos.get"` {
			t.Fatalf("expected unsupported action error, got %v", err)
		}
	})

	t.Run("invalid max results", func(t *testing.T) {
		opts := baseOptions()
		opts.maxResults = 100
		if err := validateOptions(opts); err == nil || err.Error() != "invalid -max-results 100: must be between 1 and 50" {
			t.Fatalf("expected max results error, got %v", err)
		}
	})

	t.Run("search must have query or channel or mine", func(t *testing.T) {
		opts := baseOptions()
		opts.action = "search.list"
		opts.ids = ""
		opts.query = ""
		opts.channelID = ""
		opts.searchMine = false
		if err := validateOptions(opts); err == nil || err.Error() != "search.list: provide -query, -channel-id, or -search-mine" {
			t.Fatalf("expected search validation error, got %v", err)
		}
	})

	t.Run("playlists requires ids channel or mine", func(t *testing.T) {
		opts := baseOptions()
		opts.action = "playlists.list"
		opts.ids = ""
		opts.channelID = ""
		opts.mine = false
		if err := validateOptions(opts); err == nil || err.Error() != "playlists.list: provide at least -ids, -channel-id, or -mine" {
			t.Fatalf("expected playlist validation error, got %v", err)
		}
	})

	t.Run("playlists mine and ids mutually exclusive", func(t *testing.T) {
		opts := baseOptions()
		opts.action = "playlists.list"
		opts.ids = "PL123"
		opts.mine = true
		if err := validateOptions(opts); err == nil || err.Error() != "playlists.list: -mine and -ids are mutually exclusive" {
			t.Fatalf("expected mutual exclusivity error, got %v", err)
		}
	})

	t.Run("valid videos options", func(t *testing.T) {
		if err := validateOptions(baseOptions()); err != nil {
			t.Fatalf("expected valid options, got %v", err)
		}
	})

	t.Run("valid search options", func(t *testing.T) {
		opts := baseOptions()
		opts.action = "search.list"
		opts.query = "MediaX"
		opts.ids = ""
		if err := validateOptions(opts); err != nil {
			t.Fatalf("expected valid search options, got %v", err)
		}
	})

	t.Run("valid playlist options with mine", func(t *testing.T) {
		opts := baseOptions()
		opts.action = "playlists.list"
		opts.ids = ""
		opts.mine = true
		if err := validateOptions(opts); err != nil {
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
		token, source := resolveAccessToken("flag-token", cfg)
		if token != "flag-token" || source != "flag" {
			t.Fatalf("expected flag token, got token=%q source=%q", token, source)
		}
	})

	t.Run("env overrides config", func(t *testing.T) {
		t.Setenv("GOOGLE_YOUTUBE_ACCESS_TOKEN", "env-token")
		token, source := resolveAccessToken("", cfg)
		if token != "env-token" || source != "env" {
			t.Fatalf("expected env token, got token=%q source=%q", token, source)
		}
	})

	t.Run("config fallback", func(t *testing.T) {
		t.Setenv("GOOGLE_YOUTUBE_ACCESS_TOKEN", "")
		token, source := resolveAccessToken("", cfg)
		if token != "config-token" || source != "config" {
			t.Fatalf("expected config token, got token=%q source=%q", token, source)
		}
	})

	t.Run("empty when no sources", func(t *testing.T) {
		t.Setenv("GOOGLE_YOUTUBE_ACCESS_TOKEN", "")
		token, source := resolveAccessToken("", &config.GoogleYouTubeConfig{})
		if token != "" || source != "" {
			t.Fatalf("expected empty token, got token=%q source=%q", token, source)
		}
	})
}

func baseOptions() *options {
	return &options{
		configPath: "/tmp/config.yaml",
		action:     "videos.list",
		part:       "snippet",
		ids:        "abc123",
		maxResults: 5,
	}
}
