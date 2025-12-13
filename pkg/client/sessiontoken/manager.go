package sessiontoken

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/callback"
	"github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/sanitizer"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// Manager 负责协调 Flow 状态机、存储与 adapter。
type Manager struct {
	baseClient         *kernel.BaseClient
	logger             *logger.Logger
	cache              cache.ICache
	store              FlowStore
	stateMachine       *StateMachine
	flowTTL            time.Duration
	clock              func() time.Time
	authenticator      Authenticator
	callbackDispatcher CallbackDispatcher
}

const (
	metadataSessionTokenKey        = "session_token"
	metadataCredentialsNote        = "credentials_note"
	metadataCredentialsExpiresHint = "credentials_expires_hint"
	metadataLastFailedAPI          = "last_failed_api"
	metadataReuseSessionKey        = "reuse_session"
)

const (
	reuseSessionCachePrefix = "sessionToken:reuse:"
)

// ManagerOption 用于配置 Manager。
type ManagerOption func(*Manager)

// WithFlowTTL 自定义 Flow TTL。
func WithFlowTTL(ttl time.Duration) ManagerOption {
	return func(m *Manager) {
		if ttl > 0 {
			m.flowTTL = ttl
		}
	}
}

// WithStateMachine 指定自定义状态机。
func WithStateMachine(sm *StateMachine) ManagerOption {
	return func(m *Manager) {
		if sm != nil {
			m.stateMachine = sm
		}
	}
}

// WithClock 自定义时间函数，方便测试。
func WithClock(clock func() time.Time) ManagerOption {
	return func(m *Manager) {
		if clock != nil {
			m.clock = clock
		}
	}
}

// WithAuthenticator 指定 Authenticator。
func WithAuthenticator(auth Authenticator) ManagerOption {
	return func(m *Manager) {
		m.authenticator = auth
	}
}

// WithCallbackDispatcher 指定回调派发器。
func WithCallbackDispatcher(dispatcher CallbackDispatcher) ManagerOption {
	return func(m *Manager) {
		m.callbackDispatcher = dispatcher
	}
}

// NewManager 创建一个新的 Manager。
func NewManager(baseClient *kernel.BaseClient, log *logger.Logger, cache cache.ICache, store FlowStore, opts ...ManagerOption) *Manager {
	m := &Manager{
		baseClient:   baseClient,
		logger:       log,
		cache:        cache,
		store:        store,
		stateMachine: NewStateMachine(),
		flowTTL:      DefaultFlowTTL,
		clock:        time.Now,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// FlowTTL 返回当前使用的 Flow TTL。
func (m *Manager) FlowTTL() time.Duration {
	if m.flowTTL <= 0 {
		return DefaultFlowTTL
	}
	return m.flowTTL
}

// Store 提供外部访问底层存储的能力。
func (m *Manager) Store() FlowStore {
	return m.store
}

// Logger 返回注入的 logger。
func (m *Manager) Logger() *logger.Logger {
	return m.logger
}

// CreateFlow 根据请求创建新的 Flow。
func (m *Manager) CreateFlow(ctx context.Context, opts *CreateFlowOptions) (flow *Flow, err error) {
	start := time.Now()
	defer func() {
		m.logFlowMetric(ctx, "create_flow", flow, start, err)
	}()
	if opts == nil {
		err = errors.New("sessiontoken: create options is nil")
		return nil, err
	}
	if err := m.validateCreateOptions(opts); err != nil {
		return nil, err
	}
	now := m.clock().UTC()
	if existing, err := m.store.GetByState(ctx, opts.TenantUUID, opts.State); err == nil && existing != nil {
		if existing.ExpiresAt.After(now) {
			return existing, nil
		}
	} else if err != nil && !errors.Is(err, ErrFlowNotFound) {
		return nil, err
	}
	if m.authenticator == nil {
		err = errors.New("sessiontoken: authenticator is not configured")
		return nil, err
	}
	var ttl time.Duration
	flow, ttl = m.newFlowFromOptions(opts, now)
	var authorizeURL string
	authorizeURL, err = m.authenticator.BuildAuthorizeURL(ctx, flow)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(authorizeURL) == "" {
		err = errors.New("sessiontoken: authenticator returned empty authorize url")
		return nil, err
	}
	flow.AuthorizeURL = authorizeURL
	flow.ExpiresAt = now.Add(ttl)
	if err := m.store.Save(ctx, flow, ttl); err != nil {
		return nil, err
	}
	return flow, nil
}

// GetFlow 通过 flowID 查询 Flow。
func (m *Manager) GetFlow(ctx context.Context, flowID string) (flow *Flow, err error) {
	start := time.Now()
	defer func() {
		m.logFlowMetric(ctx, "get_flow", flow, start, err)
	}()
	if strings.TrimSpace(flowID) == "" {
		err = errors.New("sessiontoken: flow_id is empty")
		return nil, err
	}
	flow, err = m.store.Get(ctx, flowID)
	if err != nil {
		return nil, err
	}
	if flow == nil {
		return nil, ErrFlowNotFound
	}
	now := m.clock().UTC()
	if flow.ExpiresAt.Before(now) {
		return nil, ErrFlowNotFound
	}
	return flow, nil
}

// CompleteFlowSuccess 在凭证采集成功后更新 Flow 状态并触发回调。
func (m *Manager) CompleteFlowSuccess(ctx context.Context, flowID string, credentials *callback.CredentialPayload) (*Flow, error) {
	start := time.Now()
	var flow *Flow
	var err error
	defer func() {
		m.logFlowMetric(ctx, "complete_flow_success", flow, start, err)
	}()
	if credentials == nil {
		err = errors.New("sessiontoken: credential payload is nil")
		return nil, err
	}
	flow, err = m.GetFlow(ctx, flowID)
	if err != nil {
		return nil, err
	}
	if err := m.stateMachine.Transition(flow, FlowStatusSucceeded); err != nil {
		return nil, err
	}
	flow.Result = sanitizer.MaskCredentialPayload(credentials)
	flow.LastError = ""
	flow.Code = ""
	flow.Message = ""
	flow.LastFailedAPI = ""
	flow.CredentialsNote = firstNonEmpty(credentials.Note, metadataValue(flow, metadataCredentialsNote))
	flow.CredentialsExpiresHint = firstNonEmpty(credentials.ExpiresAt, metadataValue(flow, metadataCredentialsExpiresHint))
	setMetadataValue(flow, metadataSessionTokenKey, credentials.SessionToken)
	setMetadataValue(flow, metadataCredentialsNote, flow.CredentialsNote)
	setMetadataValue(flow, metadataCredentialsExpiresHint, flow.CredentialsExpiresHint)
	setMetadataValue(flow, metadataLastFailedAPI, "")
	m.bindSessionTokenIndex(ctx, flow, credentials.SessionToken)
	m.storeReusableSession(ctx, flow, credentials)

	if err := m.dispatchAndRecordCallback(ctx, flow, credentials); err != nil {
		return flow, err
	}
	return flow, nil
}

// CompleteFlowFailed 在凭证采集失败时更新 Flow 并触发失败回调。
func (m *Manager) CompleteFlowFailed(ctx context.Context, flowID string, failure *FlowFailure) (*Flow, error) {
	start := time.Now()
	var flow *Flow
	var err error
	defer func() {
		m.logFlowMetric(ctx, "complete_flow_failed", flow, start, err)
	}()
	if failure == nil {
		err = errors.New("sessiontoken: failure payload is nil")
		return nil, err
	}
	if strings.TrimSpace(failure.Reason) == "" {
		err = errors.New("sessiontoken: failure reason is empty")
		return nil, err
	}
	flow, err = m.GetFlow(ctx, flowID)
	if err != nil {
		return nil, err
	}
	if err := m.stateMachine.Transition(flow, FlowStatusFailed); err != nil {
		return nil, err
	}
	flow.LastError = strings.TrimSpace(failure.Reason)
	flow.Result = nil
	flow.Code = strings.TrimSpace(failure.Code)
	flow.Message = strings.TrimSpace(failure.Message)
	flow.LastFailedAPI = strings.TrimSpace(failure.LastFailedAPI)
	setMetadataValue(flow, metadataLastFailedAPI, flow.LastFailedAPI)

	if err := m.dispatchAndRecordCallback(ctx, flow, nil); err != nil {
		return flow, err
	}
	return flow, nil
}

// MarkTokenInvalid 简化 API 调用触发的凭证失效流程。
func (m *Manager) MarkTokenInvalid(ctx context.Context, flowID string, code string, message string, api string) (*Flow, error) {
	reason := strings.TrimSpace(message)
	if reason == "" && strings.TrimSpace(code) != "" {
		reason = fmt.Sprintf("session token invalid: %s", strings.TrimSpace(code))
	}
	if reason == "" {
		reason = "session token invalid"
	}
	return m.CompleteFlowFailed(ctx, flowID, &FlowFailure{
		Reason:        reason,
		Code:          strings.TrimSpace(code),
		Message:       strings.TrimSpace(message),
		LastFailedAPI: strings.TrimSpace(api),
	})
}

// FindFlowBySessionToken 根据 session_token 查询对应的 Flow。
func (m *Manager) FindFlowBySessionToken(ctx context.Context, token string) (*Flow, error) {
	key := sessionTokenIndexKey(token)
	if key == "" {
		return nil, ErrFlowNotFound
	}
	if m.cache == nil {
		return nil, errors.New("sessiontoken: cache is not configured")
	}
	data, err := m.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, ErrFlowNotFound
	}
	return m.GetFlow(ctx, string(data))
}

func (m *Manager) validateCreateOptions(opts *CreateFlowOptions) error {
	if strings.TrimSpace(opts.ProviderCode) == "" {
		return errors.New("sessiontoken: provider_code is required")
	}
	if strings.TrimSpace(opts.ProviderAppCode) == "" {
		return errors.New("sessiontoken: provider_app_code is required")
	}
	if strings.TrimSpace(opts.TenantUUID) == "" {
		return errors.New("sessiontoken: tenant_uuid is required")
	}
	if strings.TrimSpace(opts.State) == "" {
		return errors.New("sessiontoken: state is required")
	}
	if strings.TrimSpace(opts.CallbackURL) == "" {
		return errors.New("sessiontoken: callback_url is required")
	}
	return nil
}

func (m *Manager) newFlowFromOptions(opts *CreateFlowOptions, now time.Time) (*Flow, time.Duration) {
	flow := &Flow{
		FlowID:          m.generateFlowID(),
		ProviderCode:    strings.TrimSpace(opts.ProviderCode),
		ProviderAppCode: strings.TrimSpace(opts.ProviderAppCode),
		TenantUUID:      strings.TrimSpace(opts.TenantUUID),
		AccountID:       strings.TrimSpace(opts.AccountID),
		State:           strings.TrimSpace(opts.State),
		Status:          FlowStatusPending,
		CallbackURL:     strings.TrimSpace(opts.CallbackURL),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if len(opts.Metadata) > 0 {
		flow.Metadata = make(map[string]string, len(opts.Metadata))
		for k, v := range opts.Metadata {
			flow.Metadata[k] = v
		}
	}
	flowTTL := opts.TTL
	if flowTTL <= 0 {
		flowTTL = m.FlowTTL()
	}
	return flow, flowTTL
}

func (m *Manager) generateFlowID() string {
	return fmt.Sprintf("stf_%s", strings.ReplaceAll(uuid.New().String(), "-", ""))
}

func (m *Manager) dispatchAndRecordCallback(ctx context.Context, flow *Flow, credentials *callback.CredentialPayload) error {
	if flow == nil {
		return errors.New("sessiontoken: flow is nil")
	}
	if err := m.persistFlow(ctx, flow); err != nil {
		return err
	}
	attempts, dispatchErr := m.dispatchCallback(ctx, flow, credentials)
	flow.RetryAttempts = attempts
	if dispatchErr == nil {
		if flow.Status == FlowStatusSucceeded {
			flow.LastError = ""
		}
	} else {
		if flow.Status == FlowStatusFailed && strings.TrimSpace(flow.LastError) != "" {
			flow.LastError = fmt.Sprintf("%s; callback: %v", flow.LastError, dispatchErr)
		} else {
			flow.LastError = dispatchErr.Error()
		}
	}
	if err := m.persistFlow(ctx, flow); err != nil {
		return err
	}
	return dispatchErr
}

func (m *Manager) persistFlow(ctx context.Context, flow *Flow) error {
	if flow == nil {
		return errors.New("sessiontoken: flow is nil")
	}
	if m.store == nil {
		return errors.New("sessiontoken: flow store is not configured")
	}
	now := m.clock().UTC()
	ttl := flow.ExpiresAt.Sub(now)
	if ttl <= 0 {
		ttl = time.Second
	}
	return m.store.Update(ctx, flow, ttl)
}

func (m *Manager) dispatchCallback(ctx context.Context, flow *Flow, credentials *callback.CredentialPayload) (int, error) {
	if m.callbackDispatcher == nil {
		return 0, errors.New("sessiontoken: callback dispatcher is not configured")
	}
	if flow == nil {
		return 0, errors.New("sessiontoken: flow is nil")
	}
	if strings.TrimSpace(flow.CallbackURL) == "" {
		return 0, errors.New("sessiontoken: callback url is empty")
	}
	payload := &callback.Payload{
		FlowID:          flow.FlowID,
		State:           flow.State,
		Status:          string(flow.Status),
		ProviderCode:    flow.ProviderCode,
		ProviderAppCode: flow.ProviderAppCode,
		TenantUUID:      flow.TenantUUID,
		Code:            flow.Code,
		Message:         flow.Message,
		LastFailedAPI:   flow.LastFailedAPI,
		Credentials:     credentials,
		Metadata:        cloneMetadata(flow.Metadata),
	}
	req := &callback.Request{
		URL:     flow.CallbackURL,
		Payload: payload,
	}
	return m.callbackDispatcher.Dispatch(ctx, req)
}

func cloneMetadata(meta map[string]string) map[string]string {
	if len(meta) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(meta))
	keys := make([]string, 0, len(meta))
	for k := range meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		cloned[k] = meta[k]
	}
	return cloned
}

func metadataValue(flow *Flow, key string) string {
	if flow == nil || len(flow.Metadata) == 0 {
		return ""
	}
	return strings.TrimSpace(flow.Metadata[key])
}

func setMetadataValue(flow *Flow, key, value string) {
	if flow == nil || key == "" {
		return
	}
	if flow.Metadata == nil {
		flow.Metadata = make(map[string]string)
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		delete(flow.Metadata, key)
		return
	}
	flow.Metadata[key] = trimmed
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func (m *Manager) logFlowMetric(ctx context.Context, action string, flow *Flow, start time.Time, err error) {
	if m.logger == nil {
		return
	}
	latency := time.Since(start).Milliseconds()
	provider := "-"
	providerApp := "-"
	tenant := "-"
	account := "-"
	state := "-"
	flowID := "-"
	status := ""
	retry := 0
	if flow != nil {
		provider = sanitizeLogValue(flow.ProviderCode)
		providerApp = sanitizeLogValue(flow.ProviderAppCode)
		tenant = sanitizeLogValue(flow.TenantUUID)
		account = sanitizeLogValue(flow.AccountID)
		state = sanitizeLogValue(flow.State)
		flowID = sanitizeLogValue(flow.FlowID)
		status = string(flow.Status)
		retry = flow.RetryAttempts
	}
	logger := m.logger.WithContext(ctx)
	message := "sessiontoken_metric: action=%s provider=%s provider_app=%s tenant_uuid=%s account_id=%s state=%s flow_id=%s status=%s latency_ms=%d retry=%d"
	if err != nil {
		logger.ErrorF(message+" error=%v", action, provider, providerApp, tenant, account, state, flowID, status, latency, retry, err)
		return
	}
	logger.InfoF(message, action, provider, providerApp, tenant, account, state, flowID, status, latency, retry)
}

func sanitizeLogValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}

func (m *Manager) bindSessionTokenIndex(ctx context.Context, flow *Flow, token string) {
	if m.cache == nil || flow == nil {
		return
	}
	key := sessionTokenIndexKey(token)
	if key == "" || strings.TrimSpace(flow.FlowID) == "" {
		return
	}
	ttl := flow.ExpiresAt.Sub(m.clock().UTC()) + DefaultAuditTTL
	if ttl <= DefaultAuditTTL {
		ttl = DefaultAuditTTL
	}
	if err := m.cache.Set(ctx, key, []byte(flow.FlowID), ttl); err != nil && m.logger != nil {
		m.logger.WithContext(ctx).WarnF("sessiontoken: cache session token index failed flow_id=%s error=%v", sanitizeLogValue(flow.FlowID), err)
	}
}

// FetchReusableSession 从缓存中读取可复用的 session 凭证。
func (m *Manager) FetchReusableSession(ctx context.Context, flow *Flow) (*callback.CredentialPayload, error) {
	if m.cache == nil || flow == nil {
		return nil, ErrFlowNotFound
	}
	key := reuseSessionCacheKey(flow)
	if key == "" {
		return nil, ErrFlowNotFound
	}
	data, err := m.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, ErrFlowNotFound
	}
	var payload callback.CredentialPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	if strings.TrimSpace(payload.SessionToken) == "" {
		return nil, errors.New("sessiontoken: cached session token is empty")
	}
	return &payload, nil
}

func (m *Manager) storeReusableSession(ctx context.Context, flow *Flow, credentials *callback.CredentialPayload) {
	if m.cache == nil || flow == nil || credentials == nil {
		return
	}
	if strings.TrimSpace(credentials.SessionToken) == "" {
		return
	}
	key := reuseSessionCacheKey(flow)
	if key == "" {
		return
	}
	clone := cloneCredentialPayload(credentials)
	data, err := json.Marshal(clone)
	if err != nil {
		if m.logger != nil {
			m.logger.WithContext(ctx).WarnF("sessiontoken: marshal reusable session failed flow_id=%s error=%v", sanitizeLogValue(flow.FlowID), err)
		}
		return
	}
	ttl := flow.ExpiresAt.Sub(m.clock().UTC()) + DefaultAuditTTL
	if ttl <= DefaultAuditTTL {
		ttl = DefaultAuditTTL
	}
	if err := m.cache.Set(ctx, key, data, ttl); err != nil && m.logger != nil {
		m.logger.WithContext(ctx).WarnF("sessiontoken: cache reusable session failed flow_id=%s error=%v", sanitizeLogValue(flow.FlowID), err)
	}
}

func reuseSessionCacheKey(flow *Flow) string {
	if flow == nil {
		return ""
	}
	parts := []string{
		strings.ToLower(strings.TrimSpace(flow.ProviderCode)),
		strings.ToLower(strings.TrimSpace(flow.ProviderAppCode)),
		strings.ToLower(strings.TrimSpace(flow.TenantUUID)),
		strings.ToLower(strings.TrimSpace(flow.AccountID)),
	}
	base := strings.Join(parts, "|")
	sum := sha256.Sum256([]byte(base))
	return reuseSessionCachePrefix + hex.EncodeToString(sum[:])
}

func cloneCredentialPayload(payload *callback.CredentialPayload) *callback.CredentialPayload {
	if payload == nil {
		return nil
	}
	cloned := *payload
	if len(payload.Cookies) > 0 {
		cloned.Cookies = make([]callback.Cookie, len(payload.Cookies))
		copy(cloned.Cookies, payload.Cookies)
	}
	if len(payload.Headers) > 0 {
		cloned.Headers = make(map[string]string, len(payload.Headers))
		for k, v := range payload.Headers {
			cloned.Headers[k] = v
		}
	}
	return &cloned
}

func parseBoolFlag(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "t", "yes", "y", "on":
		return true
	default:
		return false
	}
}
