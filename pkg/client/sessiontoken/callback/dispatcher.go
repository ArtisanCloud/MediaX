package callback

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// Dispatcher 抽象回调投递行为。
type Dispatcher interface {
	Dispatch(ctx context.Context, req *Request) (int, error)
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

// HTTPDispatcher 使用 HTTP POST 投递回调请求。
type HTTPDispatcher struct {
	client      *http.Client
	secret      string
	logger      *logger.Logger
	retryDelays []time.Duration
	httpDo      func(req *http.Request) (*http.Response, error)
}

// HTTPDispatcherOption 配置 dispatcher。
type HTTPDispatcherOption func(*HTTPDispatcher)

// WithRetrySchedule 自定义重试间隔。
func WithRetrySchedule(delays []time.Duration) HTTPDispatcherOption {
	return func(d *HTTPDispatcher) {
		if len(delays) > 0 {
			d.retryDelays = delays
		}
	}
}

// WithHTTPClient 设置自定义 HTTP 客户端。
func WithHTTPClient(client *http.Client) HTTPDispatcherOption {
	return func(d *HTTPDispatcher) {
		if client != nil {
			d.client = client
		}
	}
}

var defaultRetrySchedule = []time.Duration{2 * time.Second, 4 * time.Second, 8 * time.Second}

// NewHTTPDispatcher 创建默认的 HTTP dispatcher。
func NewHTTPDispatcher(secret string, log *logger.Logger, opts ...HTTPDispatcherOption) *HTTPDispatcher {
	dispatcher := &HTTPDispatcher{
		client:      http.DefaultClient,
		secret:      secret,
		logger:      log,
		retryDelays: defaultRetrySchedule,
	}
	for _, opt := range opts {
		opt(dispatcher)
	}
	dispatcher.httpDo = dispatcher.client.Do
	return dispatcher
}

// Dispatch 将回调请求发送到业务侧，失败时按照 2s→4s→8s 重试。
func (d *HTTPDispatcher) Dispatch(ctx context.Context, req *Request) (int, error) {
	if req == nil || req.Payload == nil {
		return 0, errors.New("callback: request or payload is nil")
	}
	if strings.TrimSpace(req.URL) == "" {
		return 0, errors.New("callback: callback url is empty")
	}
	if req.Payload.Timestamp == 0 {
		req.Payload.Timestamp = time.Now().UTC().Unix()
	}
	if req.Payload.Nonce == "" {
		req.Payload.Nonce = uuid.NewString()
	}
	body, err := req.Payload.Body()
	if err != nil {
		return 0, err
	}
	sig, err := BuildSignature(d.secret, req.Payload.Timestamp, req.Payload.Nonce, body)
	if err != nil {
		return 0, err
	}
	headers := cloneHeaders(req.Headers)
	headers["Content-Type"] = "application/json"
	headers["X-MediaX-Timestamp"] = fmt.Sprintf("%d", req.Payload.Timestamp)
	headers["X-MediaX-Nonce"] = req.Payload.Nonce
	headers["X-MediaX-Signature"] = sig

	start := time.Now()
	targetURL := strings.TrimSpace(req.URL)
	attempts := 0
	maxAttempts := len(d.retryDelays) + 1
	var lastErr error
	for attempts = 1; attempts <= maxAttempts; attempts++ {
		if attempts > 1 {
			delay := d.retryDelays[attempts-2]
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return attempts - 1, ctx.Err()
			case <-timer.C:
			}
		}

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, req.URL, bytes.NewReader(body))
		if err != nil {
			lastErr = err
			break
		}
		for k, v := range headers {
			httpReq.Header.Set(k, v)
		}
		resp, err := d.httpDo(httpReq)
		if err != nil {
			lastErr = err
			d.logRetry(ctx, req.Payload, attempts-1, err, targetURL)
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			d.logSuccess(ctx, req.Payload, attempts-1, resp.StatusCode, time.Since(start), targetURL)
			return attempts - 1, nil
		}
		lastErr = fmt.Errorf("callback: unexpected status code %d", resp.StatusCode)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			d.logFailure(ctx, req.Payload, attempts-1, lastErr, time.Since(start), targetURL)
			return attempts - 1, lastErr
		}
		d.logRetry(ctx, req.Payload, attempts-1, lastErr, targetURL)
	}
	if lastErr == nil {
		lastErr = errors.New("callback: dispatcher exhausted retries")
	}
	d.logFailure(ctx, req.Payload, attempts-1, lastErr, time.Since(start), targetURL)
	return attempts - 1, lastErr
}

func (d *HTTPDispatcher) logSuccess(ctx context.Context, payload *Payload, retry int, status int, latency time.Duration, callbackURL string) {
	if d.logger == nil {
		return
	}
	d.logger.WithContext(ctx).InfoF(
		"sessiontoken_callback: success provider=%s provider_app=%s tenant_uuid=%s callback_url=%s flow_id=%s state=%s flow_status=%s code=%s retry_count=%d http_status=%d latency_ms=%d",
		valueOrDash(payload.ProviderCode), valueOrDash(payload.ProviderAppCode), valueOrDash(payload.TenantUUID),
		valueOrDash(callbackURL), payload.FlowID, payload.State, valueOrDash(payload.Status), valueOrDash(payload.Code),
		retry, status, latency.Milliseconds(),
	)
}

func (d *HTTPDispatcher) logRetry(ctx context.Context, payload *Payload, retry int, err error, callbackURL string) {
	if d.logger == nil {
		return
	}
	d.logger.WithContext(ctx).WarnF(
		"sessiontoken_callback: retry provider=%s provider_app=%s tenant_uuid=%s callback_url=%s flow_id=%s state=%s flow_status=%s code=%s retry_count=%d error=%v",
		valueOrDash(payload.ProviderCode), valueOrDash(payload.ProviderAppCode), valueOrDash(payload.TenantUUID),
		valueOrDash(callbackURL), payload.FlowID, payload.State, valueOrDash(payload.Status), valueOrDash(payload.Code),
		retry, err,
	)
}

func (d *HTTPDispatcher) logFailure(ctx context.Context, payload *Payload, retry int, err error, latency time.Duration, callbackURL string) {
	if d.logger == nil {
		return
	}
	d.logger.WithContext(ctx).ErrorF(
		"sessiontoken_callback: failed provider=%s provider_app=%s tenant_uuid=%s callback_url=%s flow_id=%s state=%s flow_status=%s code=%s retry_count=%d latency_ms=%d error=%v",
		valueOrDash(payload.ProviderCode), valueOrDash(payload.ProviderAppCode), valueOrDash(payload.TenantUUID),
		valueOrDash(callbackURL), payload.FlowID, payload.State, valueOrDash(payload.Status), valueOrDash(payload.Code),
		retry, latency.Milliseconds(), err,
	)
}
