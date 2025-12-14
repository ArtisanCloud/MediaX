package v4

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/config"
	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
	sessionmiddleware "github.com/ArtisanCloud/MediaX/server/middleware/session_token"
	zhmiddleware "github.com/ArtisanCloud/MediaX/server/zhihu/sessionToken/middleware"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/google/uuid"
)

const (
	defaultZhihuAPIBase = "https://www.zhihu.com"
	defaultUserAgent    = "Mozilla/5.0 (compatible; MediaX-SessionToken/1.0)"
)

const (
	codeBadRequest    = "ZH_BAD_REQUEST"
	codeCookieExpired = "ZH_COOKIE_EXPIRED"
	codeRiskBlock     = "ZH_RISK_BLOCK"
	codeUpstreamError = "ZH_UPSTREAM_ERROR"
)

const (
	apiMeFollowings     = "me.followings"
	apiChannelsArticles = "channels.articles"
	apiArticleGet       = "articles.get"
	apiArticlePost      = "articles.post"
	apiSanityCheck      = "sanity.check"
)

var errSessionTokenMissing = errors.New("sessiontoken: session token missing")

// Client 负责代理 zhihu/v1/* API。
type Client struct {
	manager     sessiontoken.SessionTokenClient
	log         *logger.Logger
	httpClient  *http.Client
	baseURL     string
	userAgent   string
	retryDelays []time.Duration
	proxyDesc   string
}

// Option allows overriding client defaults (used in tests).
type Option func(*Client)

// WithHTTPClient overrides the HTTP client used for upstream requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
			c.proxyDesc = "custom"
		}
	}
}

// WithBaseURL overrides the Zhihu upstream base url.
func WithBaseURL(base string) Option {
	return func(c *Client) {
		if trimmed := strings.TrimSpace(base); trimmed != "" {
			c.baseURL = strings.TrimRight(trimmed, "/")
		}
	}
}

// WithRetrySchedule overrides API retry backoff strategy (mainly for tests).
func WithRetrySchedule(delays []time.Duration) Option {
	return func(c *Client) {
		if len(delays) == 0 {
			c.retryDelays = nil
			return
		}
		c.retryDelays = delays
	}
}

// NewClient 构建 Zhihu API handler。
func NewClient(cfg *config.ZhihuSessionTokenConfig, manager sessiontoken.SessionTokenClient, log *logger.Logger, opts ...Option) (*Client, error) {
	if cfg == nil {
		return nil, errors.New("zhihu.sessiontoken: config is nil")
	}
	if manager == nil {
		return nil, errors.New("zhihu.sessiontoken: session token client is nil")
	}
	httpClient, proxyDesc := buildHTTPClientWithProxy(cfg)
	client := &Client{
		manager:     manager,
		log:         log,
		baseURL:     resolveZhihuBaseURL(cfg),
		userAgent:   resolveUserAgent(cfg),
		httpClient:  httpClient,
		retryDelays: resolveRetrySchedule(cfg),
		proxyDesc:   proxyDesc,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}
	if client.httpClient == nil {
		client.httpClient, client.proxyDesc = buildHTTPClientWithProxy(cfg)
	}
	if client.baseURL == "" {
		client.baseURL = defaultZhihuAPIBase
	}
	if client.userAgent == "" {
		client.userAgent = defaultUserAgent
	}
	return client, nil
}

// RegisterRoutes 在 mux 上挂载全部 zhihu/v1 routes。
func (c *Client) RegisterRoutes(mux *http.ServeMux, apiToken string) error {
	if mux == nil {
		return errors.New("zhihu.sessiontoken: mux is nil")
	}
	if strings.TrimSpace(apiToken) == "" {
		return errors.New("zhihu.sessiontoken: api token is empty")
	}
	wrap := c.composeMiddleware(apiToken)
	mux.Handle("/zhihu/v1/me/followings", wrap(http.HandlerFunc(c.handleMeFollowings)))
	mux.Handle("/zhihu/v1/channels/", wrap(http.HandlerFunc(c.handleChannelArticles)))
	mux.Handle("/zhihu/v1/articles/", wrap(http.HandlerFunc(c.handleArticleGet)))
	mux.Handle("/zhihu/v1/articles", wrap(http.HandlerFunc(c.handleArticlePost)))
	mux.Handle("/zhihu/v1/sanity/check", wrap(http.HandlerFunc(c.handleSanityCheck)))
	return nil
}

// Version 返回当前版本标识。
func (c *Client) Version() string {
	return "v4"
}

func (c *Client) composeMiddleware(apiToken string) func(http.Handler) http.Handler {
	auth := sessionmiddleware.SessionTokenAuthMiddleware(strings.TrimSpace(apiToken), c.log)
	token := zhmiddleware.SessionTokenHeaderMiddleware(c.log)
	return func(next http.Handler) http.Handler {
		return auth(token(next))
	}
}

func resolveZhihuBaseURL(cfg *config.ZhihuSessionTokenConfig) string {
	if env := strings.TrimSpace(os.Getenv("SESSIONTOKEN_ZHIHU_API_BASE_URL")); env != "" {
		return strings.TrimRight(env, "/")
	}
	if cfg != nil {
		candidate := strings.TrimSpace(cfg.Service.BaseURL)
		if strings.Contains(strings.ToLower(candidate), "zhihu") {
			return strings.TrimRight(candidate, "/")
		}
	}
	return defaultZhihuAPIBase
}

func resolveUserAgent(cfg *config.ZhihuSessionTokenConfig) string {
	if cfg == nil {
		return defaultUserAgent
	}
	if ua := strings.TrimSpace(cfg.Authenticator.DefaultUserAgent); ua != "" {
		return ua
	}
	return defaultUserAgent
}

func buildHTTPClient(cfg *config.ZhihuSessionTokenConfig) *http.Client {
	client, _ := buildHTTPClientWithProxy(cfg)
	return client
}

func buildHTTPClientWithProxy(cfg *config.ZhihuSessionTokenConfig) (*http.Client, string) {
	timeout := 30
	if cfg != nil && cfg.Network.RequestTimeout > 0 {
		timeout = cfg.Network.RequestTimeout
	}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	pool := ""
	if cfg != nil {
		pool = strings.TrimSpace(cfg.Network.ProxyPool)
	}
	proxyCandidate := strings.TrimSpace(os.Getenv("SESSIONTOKEN_ZHIHU_PROXY"))
	if proxyCandidate == "" && cfg != nil {
		proxyCandidate = strings.TrimSpace(cfg.Network.Proxy)
	}
	proxyDesc := describeDirectProxy(pool)
	switch {
	case proxyCandidate == "", strings.EqualFold(proxyCandidate, "none"), strings.EqualFold(proxyCandidate, "direct"):
		transport.Proxy = nil
	case func() bool {
		parsed, err := url.Parse(proxyCandidate)
		if err != nil || parsed == nil {
			return false
		}
		if parsed.Scheme == "" || parsed.Host == "" {
			return false
		}
		transport.Proxy = http.ProxyURL(parsed)
		proxyDesc = parsed.String()
		return true
	}():
		// proxy configured successfully, nothing to do
	default:
		transport.Proxy = nil
	}
	return &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: transport,
	}, proxyDesc
}

func describeDirectProxy(pool string) string {
	if strings.TrimSpace(pool) == "" {
		return "direct"
	}
	return fmt.Sprintf("direct(pool:%s)", pool)
}

func resolveRetrySchedule(cfg *config.ZhihuSessionTokenConfig) []time.Duration {
	if cfg == nil {
		return nil
	}
	var delays []time.Duration
	for _, second := range cfg.API.RetryBackoff {
		if second <= 0 {
			continue
		}
		delays = append(delays, time.Duration(second)*time.Second)
	}
	if len(delays) == 0 {
		return nil
	}
	return delays
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func (c *Client) requireFlow(w http.ResponseWriter, r *http.Request, apiName string) (*sessiontoken.Flow, string, bool) {
	token := strings.TrimSpace(zhmiddleware.SessionTokenFromContext(r.Context()))
	reqID := uuid.NewString()
	if token == "" {
		c.writeError(w, reqID, http.StatusUnauthorized, codeCookieExpired, "session token missing")
		c.logAPICall(r.Context(), apiName, nil, http.StatusUnauthorized, codeCookieExpired, 0, errSessionTokenMissing)
		return nil, "", false
	}
	flow, err := c.manager.FindFlowBySessionToken(r.Context(), token)
	if err != nil {
		if errors.Is(err, sessiontoken.ErrFlowNotFound) {
			if c.log != nil {
				c.log.WithContext(r.Context()).InfoF(
					"zhihu_sessiontoken: session token has no flow binding, fallback to raw cookie api=%s",
					apiName,
				)
			}
			return nil, token, true
		} else {
			c.writeError(w, reqID, http.StatusInternalServerError, codeUpstreamError, "failed to resolve session token")
			c.logAPICall(r.Context(), apiName, nil, http.StatusInternalServerError, codeUpstreamError, 0, err)
		}
		return nil, "", false
	}
	return flow, token, true
}

func (c *Client) forward(
	w http.ResponseWriter,
	r *http.Request,
	apiName string,
	flow *sessiontoken.Flow,
	sessionToken string,
	method string,
	upstreamPath string,
	query url.Values,
	payload []byte,
) {
	reqID := uuid.NewString()
	start := time.Now()
	status, body, upstreamURL, err := c.callZhihu(r.Context(), method, upstreamPath, query, payload, sessionToken)
	if err != nil {
		c.writeError(w, reqID, http.StatusBadGateway, codeUpstreamError, "upstream request failed")
		c.logAPICall(r.Context(), apiName, flow, http.StatusBadGateway, codeUpstreamError, time.Since(start), err)
		return
	}
	if status >= 200 && status < 300 {
		c.writeSuccess(w, reqID, status, flow, fmt.Sprintf("%s %s", method, upstreamURL), body)
		c.logAPICall(r.Context(), apiName, flow, status, "", time.Since(start), nil)
		return
	}
	mappedStatus, code, markInvalid := classifyStatus(status)
	message := c.extractErrorMessage(body, defaultErrorMessage(code))
	c.writeError(w, reqID, mappedStatus, code, message)
	if markInvalid && flow != nil {
		if _, err := c.manager.MarkTokenInvalid(r.Context(), flow.FlowID, code, message, fmt.Sprintf("%s %s", method, upstreamURL)); err != nil && c.log != nil {
			c.log.WithContext(r.Context()).WarnF(
				"zhihu_sessiontoken: mark token invalid failed flow_id=%s error=%v",
				sanitizeLogValue(flow.FlowID), err,
			)
		}
	}
	c.logAPICall(
		r.Context(),
		apiName,
		flow,
		mappedStatus,
		code,
		time.Since(start),
		fmt.Errorf("upstream status %d", status),
	)
}

func (c *Client) callZhihu(ctx context.Context, method, path string, query url.Values, payload []byte, sessionToken string) (int, []byte, string, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	endpoint := c.baseURL + path
	if len(query) > 0 {
		endpoint = endpoint + "?" + query.Encode()
	}
	status, body, err := c.invokeZhihu(ctx, method, endpoint, payload, sessionToken)
	return status, body, endpoint, err
}

func (c *Client) invokeZhihu(ctx context.Context, method, endpoint string, payload []byte, sessionToken string) (int, []byte, error) {
	attempts := len(c.retryDelays) + 1
	var lastErr error
	var status int
	var body []byte
	for attempt := 1; attempt <= attempts; attempt++ {
		status, body, lastErr = c.performRequest(ctx, method, endpoint, payload, sessionToken)
		if lastErr != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return 0, nil, ctxErr
			}
			if attempt == attempts {
				return 0, nil, lastErr
			}
			if err := c.pauseForRetry(ctx, attempt); err != nil {
				return 0, nil, err
			}
			continue
		}
		if status >= 500 && attempt < attempts {
			if err := c.pauseForRetry(ctx, attempt); err != nil {
				return status, body, err
			}
			continue
		}
		return status, body, nil
	}
	if lastErr != nil {
		return 0, nil, lastErr
	}
	return status, body, errors.New("zhihu.sessiontoken: exhausted retries")
}

func (c *Client) pauseForRetry(ctx context.Context, attempt int) error {
	if attempt <= 0 || attempt > len(c.retryDelays) {
		return nil
	}
	delay := c.retryDelays[attempt-1]
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (c *Client) performRequest(ctx context.Context, method, endpoint string, payload []byte, sessionToken string) (int, []byte, error) {
	var body io.Reader
	if len(payload) > 0 {
		body = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return 0, nil, err
	}
	if len(payload) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Cookie", sessionToken)
	if req.Header.Get("Referer") == "" {
		req.Header.Set("Referer", "https://www.zhihu.com/")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, data, nil
}

type apiSuccessResponse struct {
	RequestID string          `json:"request_id"`
	Data      json.RawMessage `json:"data"`
	Meta      responseMeta    `json:"meta"`
}

type responseMeta struct {
	Source      string `json:"source"`
	FlowID      string `json:"flow_id,omitempty"`
	TenantUUID  string `json:"tenant_uuid,omitempty"`
	ProviderApp string `json:"provider_app_code,omitempty"`
	Upstream    string `json:"upstream,omitempty"`
	CapturedAt  string `json:"captured_at"`
}

type apiError struct {
	RequestID string `json:"request_id,omitempty"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

func (c *Client) writeSuccess(w http.ResponseWriter, requestID string, status int, flow *sessiontoken.Flow, upstream string, body []byte) {
	response := apiSuccessResponse{
		RequestID: requestID,
		Data:      normalizeJSONBody(body),
		Meta: responseMeta{
			Source:      "zhihu",
			FlowID:      flowIDFromFlow(flow),
			TenantUUID:  sanitizeLogValue(flowTenant(flow)),
			ProviderApp: sanitizeLogValue(flowProviderApp(flow)),
			Upstream:    upstream,
			CapturedAt:  time.Now().UTC().Format(time.RFC3339),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func (c *Client) writeError(w http.ResponseWriter, requestID string, status int, code, message string) {
	if code == "" {
		code = codeUpstreamError
	}
	if message == "" {
		message = defaultErrorMessage(code)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiError{
		RequestID: requestID,
		Code:      code,
		Message:   message,
	})
}

func normalizeJSONBody(body []byte) json.RawMessage {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return json.RawMessage("null")
	}
	if json.Valid(trimmed) {
		return json.RawMessage(trimmed)
	}
	quoted, _ := json.Marshal(string(trimmed))
	return json.RawMessage(quoted)
}

func flowIDFromFlow(flow *sessiontoken.Flow) string {
	if flow == nil {
		return ""
	}
	return sanitizeLogValue(flow.FlowID)
}

func flowTenant(flow *sessiontoken.Flow) string {
	if flow == nil {
		return ""
	}
	return flow.TenantUUID
}

func flowProviderApp(flow *sessiontoken.Flow) string {
	if flow == nil {
		return ""
	}
	return flow.ProviderAppCode
}

func classifyStatus(status int) (int, string, bool) {
	switch status {
	case http.StatusUnauthorized:
		return http.StatusUnauthorized, codeCookieExpired, true
	case http.StatusForbidden:
		return http.StatusForbidden, codeRiskBlock, true
	case http.StatusBadRequest, http.StatusNotFound, http.StatusUnprocessableEntity:
		return http.StatusBadRequest, codeBadRequest, false
	default:
		if status >= 500 {
			return http.StatusBadGateway, codeUpstreamError, false
		}
		return http.StatusBadGateway, codeUpstreamError, false
	}
}

func defaultErrorMessage(code string) string {
	switch code {
	case codeCookieExpired:
		return "session token expired"
	case codeRiskBlock:
		return "request blocked by Zhihu risk control"
	case codeBadRequest:
		return "invalid request payload"
	default:
		return "upstream service error"
	}
}

func (c *Client) extractErrorMessage(body []byte, fallback string) string {
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return fallback
	}
	var parsed map[string]interface{}
	if json.Unmarshal(body, &parsed) == nil {
		if msg, ok := parsed["message"].(string); ok && strings.TrimSpace(msg) != "" {
			return msg
		}
		if errVal, ok := parsed["error"].(string); ok && strings.TrimSpace(errVal) != "" {
			return errVal
		}
		if desc, ok := parsed["error_description"].(string); ok && strings.TrimSpace(desc) != "" {
			return desc
		}
	}
	if len(trimmed) > 256 {
		trimmed = trimmed[:256]
	}
	return trimmed
}

func (c *Client) logAPICall(ctx context.Context, api string, flow *sessiontoken.Flow, status int, code string, latency time.Duration, cause error) {
	if c.log == nil {
		return
	}
	flowID := "-"
	tenant := "-"
	providerApp := "-"
	if flow != nil {
		flowID = sanitizeLogValue(flow.FlowID)
		tenant = sanitizeLogValue(flow.TenantUUID)
		providerApp = sanitizeLogValue(flow.ProviderAppCode)
	}
	base := sanitizeLogValue(c.baseURL)
	proxy := sanitizeLogValue(c.proxyDesc)
	message := "sessiontoken_api: provider=zhihu api=%s version=%s base_url=%s proxy=%s flow_id=%s tenant_uuid=%s provider_app=%s status=%d latency_ms=%d code=%s"
	logger := c.log.WithContext(ctx)
	if cause != nil {
		logger.WarnF(message+" error=%v", api, c.Version(), base, proxy, flowID, tenant, providerApp, status, latency.Milliseconds(), sanitizeLogValue(code), cause)
		return
	}
	logger.InfoF(message, api, c.Version(), base, proxy, flowID, tenant, providerApp, status, latency.Milliseconds(), sanitizeLogValue(code))
}

func sanitizeLogValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}

func methodNotAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func parseQueryInt(raw string, def int) int {
	if strings.TrimSpace(raw) == "" {
		return def
	}
	if v, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil {
		return v
	}
	return def
}

func extractChannelID(path string) (string, bool) {
	trimmed := strings.TrimPrefix(path, "/zhihu/v1/channels/")
	if trimmed == path {
		return "", false
	}
	if !strings.HasSuffix(trimmed, "/articles") {
		return "", false
	}
	trimmed = strings.TrimSuffix(trimmed, "/articles")
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "" {
		return "", false
	}
	return trimmed, true
}

func extractArticleID(path string) (string, bool) {
	trimmed := strings.TrimPrefix(path, "/zhihu/v1/articles/")
	if trimmed == path {
		return "", false
	}
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "" || strings.Contains(trimmed, "/") {
		return "", false
	}
	return trimmed, true
}
