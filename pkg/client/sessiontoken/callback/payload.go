package callback

import (
	"encoding/json"
	"strings"
)

// CredentialPayload 表示需要回传给业务的凭证。
type CredentialPayload struct {
	SessionToken string            `json:"session_token"`
	Cookies      []Cookie          `json:"cookies_json,omitempty"`
	Headers      map[string]string `json:"headers_json,omitempty"`
	ExpiresAt    string            `json:"expires_at,omitempty"`
	CapturedAt   string            `json:"captured_at,omitempty"`
	Note         string            `json:"note,omitempty"`
}

// Cookie 是回传时的 cookie 结构。
type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain,omitempty"`
}

// Payload 是发送到业务侧的回调信息。
type Payload struct {
	FlowID          string             `json:"flow_id"`
	State           string             `json:"state"`
	Status          string             `json:"status"`
	ProviderCode    string             `json:"provider_code,omitempty"`
	ProviderAppCode string             `json:"provider_app_code,omitempty"`
	TenantUUID      string             `json:"tenant_uuid,omitempty"`
	Code            string             `json:"code,omitempty"`
	Message         string             `json:"message,omitempty"`
	LastFailedAPI   string             `json:"last_failed_api,omitempty"`
	Credentials     *CredentialPayload `json:"credentials,omitempty"`
	Metadata        map[string]string  `json:"metadata,omitempty"`
	Timestamp       int64              `json:"timestamp"`
	Nonce           string             `json:"nonce"`
}

// Body 返回 JSON 编码后的 payload。
func (p *Payload) Body() ([]byte, error) {
	return json.Marshal(p)
}

// Request 定义了回调请求。
type Request struct {
	URL     string            `json:"-"`
	Payload *Payload          `json:"payload"`
	Headers map[string]string `json:"headers"`
}

func cloneHeaders(in map[string]string) map[string]string {
	if len(in) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func valueOrDash(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}
