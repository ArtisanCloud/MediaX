package sessiontoken

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken/callback"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

// FlowOrchestrator 将 SessionToken manager 与 CredentialHarvester 织入在一起，
// 确保 Flow 创建后能够自动启动抓取流程并在成功/失败时触发回调。
// 它实现 SessionTokenClient 接口，可直接作为 HTTP handler 的依赖。
type FlowOrchestrator struct {
	manager   SessionTokenClient
	harvester CredentialHarvester
	logger    *logger.Logger
	ctx       context.Context

	mu      sync.Mutex
	running map[string]struct{}

	metadataPollInterval time.Duration
	metadataWaitTimeout  time.Duration
}

const (
	defaultMetadataPollInterval = 2 * time.Second
	defaultMetadataWaitTimeout  = 2 * time.Minute
)

// NewFlowOrchestrator 创建一个带 Harvester 能力的 SessionTokenClient。
func NewFlowOrchestrator(ctx context.Context, manager SessionTokenClient, harvester CredentialHarvester, log *logger.Logger) *FlowOrchestrator {
	if ctx == nil {
		ctx = context.Background()
	}
	return &FlowOrchestrator{
		manager:              manager,
		harvester:            harvester,
		logger:               log,
		ctx:                  ctx,
		running:              make(map[string]struct{}),
		metadataPollInterval: defaultMetadataPollInterval,
		metadataWaitTimeout:  defaultMetadataWaitTimeout,
	}
}

// CreateFlow 代理 manager 创建 Flow，并在成功后启动 harvester。
func (o *FlowOrchestrator) CreateFlow(ctx context.Context, opts *CreateFlowOptions) (*Flow, error) {
	flow, err := o.manager.CreateFlow(ctx, opts)
	if err != nil {
		return nil, err
	}
	o.startHarvest(flow)
	return flow, nil
}

// GetFlow 直接代理 manager 能力。
func (o *FlowOrchestrator) GetFlow(ctx context.Context, flowID string) (*Flow, error) {
	return o.manager.GetFlow(ctx, flowID)
}

// CompleteFlowSuccess 直接代理 manager 能力。
func (o *FlowOrchestrator) CompleteFlowSuccess(ctx context.Context, flowID string, credentials *callback.CredentialPayload) (*Flow, error) {
	return o.manager.CompleteFlowSuccess(ctx, flowID, credentials)
}

// CompleteFlowFailed 直接代理 manager 能力。
func (o *FlowOrchestrator) CompleteFlowFailed(ctx context.Context, flowID string, failure *FlowFailure) (*Flow, error) {
	return o.manager.CompleteFlowFailed(ctx, flowID, failure)
}

// MarkTokenInvalid 直接代理 manager 能力。
func (o *FlowOrchestrator) MarkTokenInvalid(ctx context.Context, flowID string, code string, message string, api string) (*Flow, error) {
	return o.manager.MarkTokenInvalid(ctx, flowID, code, message, api)
}

// FindFlowBySessionToken 根据 session_token 查询 Flow。
func (o *FlowOrchestrator) FindFlowBySessionToken(ctx context.Context, token string) (*Flow, error) {
	return o.manager.FindFlowBySessionToken(ctx, token)
}

func (o *FlowOrchestrator) startHarvest(flow *Flow) {
	if flow == nil || flow.FlowID == "" {
		return
	}
	if flow.Status == FlowStatusSucceeded || flow.Status == FlowStatusFailed {
		return
	}
	if shouldReuseSession(flow) && o.tryReuseSession(flow) {
		return
	}
	if o.harvester == nil {
		return
	}
	o.mu.Lock()
	if _, exists := o.running[flow.FlowID]; exists {
		o.mu.Unlock()
		return
	}
	o.running[flow.FlowID] = struct{}{}
	o.mu.Unlock()

	go o.runHarvest(flow.FlowID)
}

func (o *FlowOrchestrator) runHarvest(flowID string) {
	defer func() {
		o.mu.Lock()
		delete(o.running, flowID)
		o.mu.Unlock()
	}()

	ctx := o.ctx
	flow, err := o.waitForMetadata(ctx, flowID)
	if err != nil {
		o.logError(ctx, "sessiontoken_orchestrator: metadata wait failed flow_id=%s error=%v", flowID, err)
		o.failFlow(ctx, flowID, fmt.Errorf("metadata wait: %w", err))
		return
	}
	if flow == nil {
		return
	}
	if flow.Status == FlowStatusSucceeded || flow.Status == FlowStatusFailed {
		// Flow 已经在等待期间被更新为成功/失败，直接退出。
		return
	}

	payload, err := o.harvester.Watch(ctx, flow)
	if err != nil {
		failure := &FlowFailure{Reason: fmt.Sprintf("harvester: %v", err)}
		if _, failErr := o.manager.CompleteFlowFailed(ctx, flowID, failure); failErr != nil {
			o.logError(ctx, "sessiontoken_orchestrator: mark flow failed flow_id=%s error=%v", flowID, failErr)
		}
		return
	}
	if _, err := o.manager.CompleteFlowSuccess(ctx, flowID, payload); err != nil {
		o.logError(ctx, "sessiontoken_orchestrator: mark flow success flow_id=%s error=%v", flowID, err)
	}
}

func (o *FlowOrchestrator) logError(ctx context.Context, format string, args ...interface{}) {
	if o.logger == nil {
		return
	}
	o.logger.WithContext(ctx).ErrorF(format, args...)
}

func shouldReuseSession(flow *Flow) bool {
	if flow == nil {
		return false
	}
	return parseBoolFlag(metadataValue(flow, metadataReuseSessionKey))
}

func (o *FlowOrchestrator) tryReuseSession(flow *Flow) bool {
	fetcher, ok := o.manager.(ReusableSessionFetcher)
	if !ok {
		return false
	}
	payload, err := fetcher.FetchReusableSession(o.ctx, flow)
	if err != nil || payload == nil {
		return false
	}
	if _, err := o.manager.CompleteFlowSuccess(o.ctx, flow.FlowID, payload); err != nil {
		o.logError(o.ctx, "sessiontoken_orchestrator: reuse session failed flow_id=%s error=%v", flow.FlowID, err)
		return false
	}
	return true
}

func (o *FlowOrchestrator) waitForMetadata(ctx context.Context, flowID string) (*Flow, error) {
	if flowID == "" {
		return nil, fmt.Errorf("sessiontoken: flow id is empty")
	}
	timeout := o.metadataWaitTimeout
	if timeout <= 0 {
		return o.manager.GetFlow(ctx, flowID)
	}
	interval := o.metadataPollInterval
	if interval <= 0 {
		interval = defaultMetadataPollInterval
	}
	deadline := time.Now().Add(timeout)
	for {
		flow, err := o.manager.GetFlow(ctx, flowID)
		if err != nil {
			return nil, err
		}
		if flow == nil {
			return nil, ErrFlowNotFound
		}
		if flow.Status == FlowStatusSucceeded || flow.Status == FlowStatusFailed {
			return flow, nil
		}
		if isMetadataReady(flow) {
			return flow, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("sessiontoken: metadata session_token not ready before timeout")
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
		}
	}
}

func isMetadataReady(flow *Flow) bool {
	if flow == nil {
		return false
	}
	return metadataValue(flow, metadataSessionTokenKey) != ""
}

func (o *FlowOrchestrator) failFlow(ctx context.Context, flowID string, reason error) {
	if reason == nil {
		return
	}
	failure := &FlowFailure{Reason: reason.Error()}
	if _, err := o.manager.CompleteFlowFailed(ctx, flowID, failure); err != nil {
		o.logError(ctx, "sessiontoken_orchestrator: mark flow failed flow_id=%s error=%v", flowID, err)
	}
}
