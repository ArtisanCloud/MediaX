package api

import (
	"time"

	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
	"github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/sanitizer"
)

// FlowDTO 表示对外暴露的 Flow 结构。
type FlowDTO struct {
	FlowID          string                  `json:"flow_id"`
	ProviderCode    string                  `json:"provider_code"`
	ProviderAppCode string                  `json:"provider_app_code"`
	TenantUUID      string                  `json:"tenant_uuid"`
	AccountID       string                  `json:"account_id"`
	State           string                  `json:"state"`
	Status          sessiontoken.FlowStatus `json:"status"`
	AuthorizeURL    string                  `json:"authorize_url"`
	ExpiresAt       time.Time               `json:"expires_at"`
	ExpiresIn       int64                   `json:"expires_in"`
	Metadata        map[string]string       `json:"metadata,omitempty"`
	Result          interface{}             `json:"result,omitempty"`
	LastError       string                  `json:"last_error,omitempty"`
	CallbackURL     string                  `json:"callback_url"`
	RetryAttempts   int                     `json:"retry_attempts"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

// GetFlowResponse 表示 Flow 查询响应。
type GetFlowResponse struct {
	Flow *FlowDTO `json:"flow"`
}

// NewGetFlowResponse 根据 Flow 生成响应对象，包含脱敏凭证与剩余 TTL。
func NewGetFlowResponse(flow *sessiontoken.Flow, now time.Time) *GetFlowResponse {
	if flow == nil {
		return &GetFlowResponse{}
	}
	dto := &FlowDTO{
		FlowID:          flow.FlowID,
		ProviderCode:    flow.ProviderCode,
		ProviderAppCode: flow.ProviderAppCode,
		TenantUUID:      flow.TenantUUID,
		AccountID:       flow.AccountID,
		State:           flow.State,
		Status:          flow.Status,
		AuthorizeURL:    flow.AuthorizeURL,
		ExpiresAt:       flow.ExpiresAt,
		Metadata:        cloneMetadata(flow.Metadata),
		Result:          sanitizer.MaskCredentialPayload(flow.Result),
		LastError:       flow.LastError,
		CallbackURL:     flow.CallbackURL,
		RetryAttempts:   flow.RetryAttempts,
		CreatedAt:       flow.CreatedAt,
		UpdatedAt:       flow.UpdatedAt,
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	ttl := flow.ExpiresAt.Sub(now)
	if ttl < 0 {
		ttl = 0
	}
	dto.ExpiresIn = int64(ttl.Seconds())
	return &GetFlowResponse{Flow: dto}
}

func cloneMetadata(meta map[string]string) map[string]string {
	if len(meta) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(meta))
	for k, v := range meta {
		cloned[k] = v
	}
	return cloned
}
