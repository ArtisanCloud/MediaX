package callback

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
	FlowID      string             `json:"flow_id"`
	State       string             `json:"state"`
	Status      string             `json:"status"`
	Credentials *CredentialPayload `json:"credentials,omitempty"`
	Metadata    map[string]string  `json:"metadata,omitempty"`
	Timestamp   int64              `json:"timestamp"`
	Nonce       string             `json:"nonce"`
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

// Dispatcher 抽象回调投递行为。
type Dispatcher interface {
	Dispatch(ctx context.Context, req *Request) error
}

// Signer 负责根据 secret 生成 HMAC-SHA256 签名。
type Signer struct {
	secret []byte
}

// NewSigner 创建新的 HMAC 签名器。
func NewSigner(secret string) *Signer {
	return &Signer{secret: []byte(secret)}
}

// Sign 根据约定的格式生成签名：HMAC(secret, "timestamp:nonce:" + body)。
func (s *Signer) Sign(timestamp int64, nonce string, body []byte) (string, error) {
	if len(s.secret) == 0 {
		return "", errors.New("callback: secret is empty")
	}
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(fmt.Sprintf("%d:%s:", timestamp, nonce)))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// BuildSignature 是 Signer.Sign 的便捷封装。
func BuildSignature(secret string, timestamp int64, nonce string, body []byte) (string, error) {
	return NewSigner(secret).Sign(timestamp, nonce, body)
}
