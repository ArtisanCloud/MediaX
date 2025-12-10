package sessiontoken

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// Manager 负责协调 Flow 状态机、存储与 adapter。
type Manager struct {
	baseClient    *kernel.BaseClient
	logger        *logger.Logger
	cache         cache.ICache
	store         FlowStore
	stateMachine  *StateMachine
	flowTTL       time.Duration
	clock         func() time.Time
	authenticator Authenticator
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
func (m *Manager) CreateFlow(ctx context.Context, opts *CreateFlowOptions) (*Flow, error) {
	if opts == nil {
		return nil, errors.New("sessiontoken: create options is nil")
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
		return nil, errors.New("sessiontoken: authenticator is not configured")
	}
	flow, ttl := m.newFlowFromOptions(opts, now)
	authorizeURL, err := m.authenticator.BuildAuthorizeURL(ctx, flow)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(authorizeURL) == "" {
		return nil, errors.New("sessiontoken: authenticator returned empty authorize url")
	}
	flow.AuthorizeURL = authorizeURL
	flow.ExpiresAt = now.Add(ttl)
	if err := m.store.Save(ctx, flow, ttl); err != nil {
		return nil, err
	}
	return flow, nil
}

// GetFlow 通过 flowID 查询 Flow。
func (m *Manager) GetFlow(ctx context.Context, flowID string) (*Flow, error) {
	if strings.TrimSpace(flowID) == "" {
		return nil, errors.New("sessiontoken: flow_id is empty")
	}
	return m.store.Get(ctx, flowID)
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
