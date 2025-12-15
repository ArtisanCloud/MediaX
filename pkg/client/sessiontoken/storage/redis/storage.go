package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken"
)

// FlowStore 使用 redis 存储 Flow。
type FlowStore struct {
	client redis.Cmdable
	now    func() time.Time
}

// NewFlowStore 创建新的 redis 存储驱动。
func NewFlowStore(client redis.Cmdable) *FlowStore {
	return &FlowStore{
		client: client,
		now:    time.Now,
	}
}

// WithClock 允许在测试中覆写时间函数。
func (s *FlowStore) WithClock(now func() time.Time) {
	if now != nil {
		s.now = now
	}
}

func (s *FlowStore) Save(ctx context.Context, flow *sessiontoken.Flow, ttl time.Duration) error {
	return s.persist(ctx, flow, ttl, true)
}

func (s *FlowStore) Update(ctx context.Context, flow *sessiontoken.Flow, ttl time.Duration) error {
	return s.persist(ctx, flow, ttl, false)
}

func (s *FlowStore) persist(ctx context.Context, flow *sessiontoken.Flow, ttl time.Duration, isCreate bool) error {
	if flow == nil {
		return errors.New("sessiontoken/storage: flow is nil")
	}
	if flow.FlowID == "" {
		return errors.New("sessiontoken/storage: flow id is empty")
	}

	expiration := sessiontoken.NormalizeTTL(ttl) + sessiontoken.DefaultAuditTTL
	now := s.now().UTC()
	if flow.CreatedAt.IsZero() {
		flow.CreatedAt = now
	}
	flow.UpdatedAt = now

	data, err := json.Marshal(flow)
	if err != nil {
		return err
	}

	flowKey := sessiontoken.FlowKey(flow.FlowID)
	stateKey := sessiontoken.FlowStateIndexKey(flow.TenantUUID, flow.State)

	pipe := s.client.TxPipeline()
	var flowSetCmd *redis.BoolCmd
	var stateSetCmd *redis.BoolCmd
	if isCreate {
		flowSetCmd = pipe.SetNX(ctx, flowKey, data, expiration)
		stateSetCmd = pipe.SetNX(ctx, stateKey, flow.FlowID, expiration)
	} else {
		pipe.Set(ctx, flowKey, data, expiration)
		pipe.Set(ctx, stateKey, flow.FlowID, expiration)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}
	if isCreate {
		if ok, err := flowSetCmd.Result(); err != nil {
			return err
		} else if !ok {
			return sessiontoken.ErrFlowAlreadyExists
		}
		if ok, err := stateSetCmd.Result(); err != nil {
			return err
		} else if !ok {
			return errors.New("sessiontoken/storage: flow state already exists")
		}
	}
	return nil
}

func (s *FlowStore) Get(ctx context.Context, flowID string) (*sessiontoken.Flow, error) {
	key := sessiontoken.FlowKey(flowID)
	data, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, sessiontoken.ErrFlowNotFound
	}
	if err != nil {
		return nil, err
	}

	var flow sessiontoken.Flow
	if err := json.Unmarshal(data, &flow); err != nil {
		return nil, err
	}
	return &flow, nil
}

func (s *FlowStore) GetByState(ctx context.Context, tenantUUID, state string) (*sessiontoken.Flow, error) {
	key := sessiontoken.FlowStateIndexKey(tenantUUID, state)
	flowID, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, sessiontoken.ErrFlowNotFound
	}
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, flowID)
}

func (s *FlowStore) Delete(ctx context.Context, flowID string) error {
	flow, err := s.Get(ctx, flowID)
	if err != nil {
		if errors.Is(err, sessiontoken.ErrFlowNotFound) {
			return nil
		}
		return err
	}

	flowKey := sessiontoken.FlowKey(flowID)
	stateKey := sessiontoken.FlowStateIndexKey(flow.TenantUUID, flow.State)
	pipe := s.client.TxPipeline()
	pipe.Del(ctx, flowKey)
	pipe.Del(ctx, stateKey)
	_, err = pipe.Exec(ctx)
	return err
}
