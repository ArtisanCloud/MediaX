package sessiontoken

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/callback"
)

// FlowStatus 表示 Flow 生命周期中的状态。
type FlowStatus string

const (
	FlowStatusPending     FlowStatus = "pending"
	FlowStatusAuthorizing FlowStatus = "authorizing"
	FlowStatusSucceeded   FlowStatus = "succeeded"
	FlowStatusFailed      FlowStatus = "failed"
)

const (
	FlowKeyPrefix      = "sessionToken:flow:"
	FlowStateKeyPrefix = "sessionToken:flow:state:"
)

var (
	// DefaultFlowTTL 默认 Flow 的生存时间，后续可通过配置覆写。
	DefaultFlowTTL = 30 * time.Minute
	// DefaultAuditTTL Flow 完成后仍在 Redis 中保留的审计时间。
	DefaultAuditTTL = 6 * time.Hour

	// ErrFlowNotFound 表示 Flow 缺失或已过期。
	ErrFlowNotFound = errors.New("sessiontoken: flow not found")
	// ErrFlowAlreadyExists 表示 Flow ID 冲突。
	ErrFlowAlreadyExists = errors.New("sessiontoken: flow already exists")
)

// Flow 代表一次 SessionToken 登录流程。
type Flow struct {
	FlowID          string                      `json:"flow_id"`
	ProviderCode    string                      `json:"provider_code"`
	ProviderAppCode string                      `json:"provider_app_code"`
	TenantUUID      string                      `json:"tenant_uuid"`
	AccountID       string                      `json:"account_id"`
	State           string                      `json:"state"`
	Status          FlowStatus                  `json:"status"`
	AuthorizeURL    string                      `json:"authorize_url"`
	ExpiresAt       time.Time                   `json:"expires_at"`
	Metadata        map[string]string           `json:"metadata,omitempty"`
	Result          *callback.CredentialPayload `json:"result,omitempty"`
	LastError       string                      `json:"last_error,omitempty"`
	CallbackURL     string                      `json:"callback_url"`
	RetryAttempts   int                         `json:"retry_attempts"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at"`
}

// FlowStore 定义了 Flow 持久化层需要实现的接口。
type FlowStore interface {
	Save(ctx context.Context, flow *Flow, ttl time.Duration) error
	Update(ctx context.Context, flow *Flow, ttl time.Duration) error
	Get(ctx context.Context, flowID string) (*Flow, error)
	GetByState(ctx context.Context, tenantUUID, state string) (*Flow, error)
	Delete(ctx context.Context, flowID string) error
}

// FlowKey 返回 Flow 在 Redis 中的存储 key。
func FlowKey(flowID string) string {
	return FlowKeyPrefix + flowID
}

// FlowStateIndexKey 返回 tenant + state 的索引 key。
func FlowStateIndexKey(tenantUUID, state string) string {
	hash := sha256.Sum256([]byte(strings.ToLower(fmt.Sprintf("%s|%s", tenantUUID, state))))
	return FlowStateKeyPrefix + hex.EncodeToString(hash[:])
}

// NormalizeTTL 若传入 ttl<=0 则统一回退到默认值。
func NormalizeTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return DefaultFlowTTL
	}
	return ttl
}
