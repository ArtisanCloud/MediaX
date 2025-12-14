package api

import (
	"errors"
	"net/url"
	"strings"
	"time"

	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
)

// CreateFlowRequest 描述插件创建 Flow 时提交的 payload。
type CreateFlowRequest struct {
	ProviderCode    string            `json:"provider_code"`
	ProviderAppCode string            `json:"provider_app_code"`
	AccountID       string            `json:"account_id"`
	TenantUUID      string            `json:"tenant_uuid"`
	State           string            `json:"state"`
	CallbackURL     string            `json:"callback_url"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	TTLSeconds      int64             `json:"ttl_seconds,omitempty"`
}

// Validate 对请求字段做基础校验。
func (r *CreateFlowRequest) Validate() error {
	if r == nil {
		return errors.New("sessiontoken: request is nil")
	}
	if strings.TrimSpace(r.ProviderCode) == "" {
		return errors.New("sessiontoken: provider_code is required")
	}
	if strings.TrimSpace(r.ProviderAppCode) == "" {
		return errors.New("sessiontoken: provider_app_code is required")
	}
	if strings.TrimSpace(r.TenantUUID) == "" {
		return errors.New("sessiontoken: tenant_uuid is required")
	}
	if strings.TrimSpace(r.State) == "" {
		return errors.New("sessiontoken: state is required")
	}
	if strings.TrimSpace(r.CallbackURL) == "" {
		return errors.New("sessiontoken: callback_url is required")
	}
	if _, err := url.ParseRequestURI(r.CallbackURL); err != nil {
		return errors.New("sessiontoken: callback_url is invalid")
	}
	return nil
}

// ToOptions 转换为 Manager 可以消费的 CreateFlowOptions。
func (r *CreateFlowRequest) ToOptions() (*sessiontoken.CreateFlowOptions, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	metadata := make(map[string]string, len(r.Metadata))
	for k, v := range r.Metadata {
		metadata[k] = v
	}
	var ttl time.Duration
	if r.TTLSeconds > 0 {
		ttl = time.Duration(r.TTLSeconds) * time.Second
	}
	return &sessiontoken.CreateFlowOptions{
		ProviderCode:    strings.TrimSpace(r.ProviderCode),
		ProviderAppCode: strings.TrimSpace(r.ProviderAppCode),
		TenantUUID:      strings.TrimSpace(r.TenantUUID),
		AccountID:       strings.TrimSpace(r.AccountID),
		State:           strings.TrimSpace(r.State),
		CallbackURL:     strings.TrimSpace(r.CallbackURL),
		Metadata:        metadata,
		TTL:             ttl,
	}, nil
}
