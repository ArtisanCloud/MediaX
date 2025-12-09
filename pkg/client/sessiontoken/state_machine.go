package sessiontoken

import (
	"errors"
	"time"
)

var (
	// ErrInvalidTransition 在 Flow 状态变化非法时返回。
	ErrInvalidTransition = errors.New("sessiontoken: invalid flow status transition")
)

// StateMachine 管理 Flow 的状态迁移关系。
type StateMachine struct {
	transitions map[FlowStatus][]FlowStatus
	now         func() time.Time
}

// NewStateMachine 创建默认的状态机。
func NewStateMachine() *StateMachine {
	return &StateMachine{
		now: time.Now,
		transitions: map[FlowStatus][]FlowStatus{
			FlowStatusPending:     {FlowStatusPending, FlowStatusAuthorizing, FlowStatusFailed},
			FlowStatusAuthorizing: {FlowStatusAuthorizing, FlowStatusSucceeded, FlowStatusFailed},
			FlowStatusSucceeded:   {FlowStatusSucceeded},
			FlowStatusFailed:      {FlowStatusFailed},
		},
	}
}

// CanTransit 判断是否允许从 current 切换到 target。
func (sm *StateMachine) CanTransit(current, target FlowStatus) bool {
	options, ok := sm.transitions[current]
	if !ok {
		return false
	}
	for _, candidate := range options {
		if candidate == target {
			return true
		}
	}
	return false
}

// Transition 会尝试更新 Flow 状态，不合法时返回 ErrInvalidTransition。
func (sm *StateMachine) Transition(flow *Flow, target FlowStatus) error {
	if flow == nil {
		return errors.New("sessiontoken: flow is nil")
	}
	if flow.Status == target {
		flow.UpdatedAt = sm.now().UTC()
		return nil
	}
	if !sm.CanTransit(flow.Status, target) {
		return ErrInvalidTransition
	}
	flow.Status = target
	flow.UpdatedAt = sm.now().UTC()
	return nil
}

// WithClock 允许在测试中注入自定义时间函数。
func (sm *StateMachine) WithClock(now func() time.Time) {
	if now != nil {
		sm.now = now
	}
}
