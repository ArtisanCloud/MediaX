package sessiontoken

import (
	"time"

	"github.com/ArtisanCloud/MediaX/internal/kernel"
	"github.com/ArtisanCloud/MediaXCore/pkg/cache"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// Manager 负责协调 Flow 状态机、存储与 adapter。
type Manager struct {
	baseClient   *kernel.BaseClient
	logger       *logger.Logger
	cache        cache.ICache
	store        FlowStore
	stateMachine *StateMachine
	flowTTL      time.Duration
	clock        func() time.Time
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
