package sessiontoken

import (
	"context"
	"errors"
	"fmt"
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

	if err := m.dispatchAndRecordCallback(ctx, flow, credentials); err != nil {
		return flow, err
	}
	return flow, nil
}

// CompleteFlowFailed 在凭证采集失败时更新 Flow 并触发失败回调。
func (m *Manager) CompleteFlowFailed(ctx context.Context, flowID string, reason string) (*Flow, error) {
	start := time.Now()
	var flow *Flow
	var err error
	defer func() {
		m.logFlowMetric(ctx, "complete_flow_failed", flow, start, err)
	}()
	if strings.TrimSpace(reason) == "" {
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
	flow.LastError = reason
	flow.Result = nil

	if err := m.dispatchAndRecordCallback(ctx, flow, nil); err != nil {
		return flow, err
	}
	return flow, nil
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
	for k, v := range meta {
		cloned[k] = v
	}
	return cloned
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
