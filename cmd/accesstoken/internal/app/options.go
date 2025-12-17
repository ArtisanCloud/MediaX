package app

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// DefaultConfigPath is the fallback path for the CLI/service when no config is provided.
	DefaultConfigPath = "config.yaml"

	maxResultsMin = 1
	maxResultsMax = 50

	// DefaultAccessTokenTTLSeconds is used when no TTL is supplied.
	DefaultAccessTokenTTLSeconds = 3600
)

// Options describes a single AccessToken invocation, shared by CLI & HTTP service.
type Options struct {
	ConfigPath     string `json:"config_path,omitempty"`
	Provider       string `json:"provider,omitempty"`
	ProviderApp    string `json:"provider_app,omitempty"`
	AuthMode       string `json:"provider_auth_mode,omitempty"`
	ProviderCode   string `json:"provider_code,omitempty"`
	Action         string `json:"action"`
	Part           string `json:"part"`
	IDs            string `json:"ids,omitempty"`
	Chart          string `json:"chart,omitempty"`
	VideoCategory  string `json:"video_category,omitempty"`
	Region         string `json:"region,omitempty"`
	PageToken      string `json:"page_token,omitempty"`
	MaxResults     int    `json:"max_results,omitempty"`
	Query          string `json:"query,omitempty"`
	ChannelID      string `json:"channel_id,omitempty"`
	Mine           bool   `json:"mine,omitempty"`
	SearchMine     bool   `json:"search_mine,omitempty"`
	SearchType     string `json:"search_type,omitempty"`
	AccessToken    string `json:"access_token,omitempty"`
	AccessTokenTTL int    `json:"access_token_ttl,omitempty"`
	TokenSource    string `json:"token_source,omitempty"`
}

// Normalize trims inputs, applies defaults, and lowercases action strings.
func (o *Options) Normalize() {
	if o == nil {
		return
	}
	o.Action = strings.ToLower(strings.TrimSpace(o.Action))
	o.Provider = strings.TrimSpace(o.Provider)
	o.ProviderApp = strings.TrimSpace(o.ProviderApp)
	o.AuthMode = strings.TrimSpace(o.AuthMode)
	o.ProviderCode = strings.TrimSpace(o.ProviderCode)
	o.Part = strings.TrimSpace(o.Part)
	o.IDs = strings.TrimSpace(o.IDs)
	o.Chart = strings.TrimSpace(o.Chart)
	o.VideoCategory = strings.TrimSpace(o.VideoCategory)
	o.Region = strings.TrimSpace(o.Region)
	o.PageToken = strings.TrimSpace(o.PageToken)
	o.Query = strings.TrimSpace(o.Query)
	o.ChannelID = strings.TrimSpace(o.ChannelID)
	o.SearchType = strings.TrimSpace(o.SearchType)
	o.ConfigPath = strings.TrimSpace(o.ConfigPath)
	o.AccessToken = strings.TrimSpace(o.AccessToken)
	if o.MaxResults == 0 {
		o.MaxResults = 5
	}
	if o.AccessTokenTTL <= 0 {
		o.AccessTokenTTL = DefaultAccessTokenTTLSeconds
	}
}

// Validate enforces CLI/service level constraints before issuing API calls.
func (o *Options) Validate() error {
	if o == nil {
		return errors.New("options are required")
	}
	if o.Action == "" {
		return errors.New("missing required field: action")
	}
	switch o.Action {
	case "videos.list", "search.list", "playlists.list":
	default:
		return fmt.Errorf("unsupported action %q", o.Action)
	}
	if strings.TrimSpace(o.Part) == "" {
		return errors.New("missing required field: part")
	}
	if o.MaxResults < maxResultsMin || o.MaxResults > maxResultsMax {
		return fmt.Errorf("invalid max_results %d: must be between %d and %d", o.MaxResults, maxResultsMin, maxResultsMax)
	}

	switch o.Action {
	case "videos.list":
		if o.IDs == "" && o.Chart == "" {
			return errors.New("videos.list: provide at least ids or chart")
		}
	case "search.list":
		if o.Query == "" && o.ChannelID == "" && !o.SearchMine {
			return errors.New("search.list: provide query, channel_id, or search_mine")
		}
	case "playlists.list":
		if o.IDs == "" && o.ChannelID == "" && !o.Mine {
			return errors.New("playlists.list: provide at least ids, channel_id, or mine")
		}
		if o.Mine && o.IDs != "" {
			return errors.New("playlists.list: mine and ids are mutually exclusive")
		}
	}
	return nil
}
