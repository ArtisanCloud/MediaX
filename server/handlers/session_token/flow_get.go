package session_token

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken"
	sessionapi "github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken/api"
	sessionmiddleware "github.com/ArtisanCloud/MediaX/server/middleware/session_token"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

const (
	// SessionTokenFlowGetPrefix 是 Flow 查询接口的路径前缀。
	SessionTokenFlowGetPrefix = "/session-token/flows/"
	logAPISessionTokenGet     = "session_token.get_flow"
)

// SessionTokenFlowGetHandler 处理 Flow 查询请求。
type SessionTokenFlowGetHandler struct {
	Manager sessiontoken.SessionTokenClient
	Logger  *logger.Logger
	clock   func() time.Time
}

// RegisterSessionTokenFlowGetRoute 将 GET handler 注册到 mux 上。
func RegisterSessionTokenFlowGetRoute(mux *http.ServeMux, manager sessiontoken.SessionTokenClient, apiToken string, log *logger.Logger) error {
	if mux == nil {
		return errNilServeMux
	}
	if manager == nil {
		return errNilSessionTokenManager
	}
	if strings.TrimSpace(apiToken) == "" {
		return errMissingAPIToken
	}
	handler := NewSessionTokenFlowGetHandler(manager, log)
	mux.Handle(SessionTokenFlowGetPrefix, sessionmiddleware.SessionTokenAuthMiddleware(apiToken, log)(handler))
	return nil
}

// NewSessionTokenFlowGetHandler 创建 GET handler。
func NewSessionTokenFlowGetHandler(manager sessiontoken.SessionTokenClient, log *logger.Logger) *SessionTokenFlowGetHandler {
	return &SessionTokenFlowGetHandler{
		Manager: manager,
		Logger:  log,
		clock:   time.Now,
	}
}

// ServeHTTP 实现 http.Handler。
func (h *SessionTokenFlowGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	flowID, ok := h.extractFlowID(r.URL.Path)
	if !ok {
		h.writeError(w, http.StatusNotFound, "flow not found")
		return
	}
	flow, err := h.Manager.GetFlow(r.Context(), flowID)
	if err != nil {
		if errors.Is(err, sessiontoken.ErrFlowNotFound) {
			h.logError(r, nil, flowID, fmt.Errorf("sessiontoken: flow not found: %w", err))
			h.writeError(w, http.StatusNotFound, "flow not found")
			return
		}
		h.logError(r, nil, flowID, fmt.Errorf("sessiontoken: manager get flow failed: %w", err))
		h.writeError(w, http.StatusInternalServerError, "failed to fetch flow")
		return
	}
	resp := sessionapi.NewGetFlowResponse(flow, h.clock().UTC())
	h.logSuccess(r, flow)
	h.writeSuccess(w, resp)
}

func (h *SessionTokenFlowGetHandler) writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (h *SessionTokenFlowGetHandler) writeSuccess(w http.ResponseWriter, resp *sessionapi.GetFlowResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *SessionTokenFlowGetHandler) extractFlowID(path string) (string, bool) {
	if !strings.HasPrefix(path, SessionTokenFlowGetPrefix) {
		return "", false
	}
	flowID := strings.TrimSpace(strings.TrimPrefix(path, SessionTokenFlowGetPrefix))
	if flowID == "" {
		return "", false
	}
	return flowID, true
}

func (h *SessionTokenFlowGetHandler) logSuccess(r *http.Request, flow *sessiontoken.Flow) {
	if h.Logger == nil {
		return
	}
	provider, providerApp, tenant, account, state := extractFlowFields(flow)
	h.Logger.WithContext(r.Context()).InfoF(
		"api=%s action=get_flow status=success provider=%s provider_app=%s tenant_uuid=%s account_id=%s state=%s flow_id=%s",
		logAPISessionTokenGet,
		provider, providerApp, tenant, account, state, sanitizeLogValue(flow.FlowID),
	)
}

func (h *SessionTokenFlowGetHandler) logError(r *http.Request, flow *sessiontoken.Flow, flowID string, err error) {
	if h.Logger == nil || err == nil {
		return
	}
	provider, providerApp, tenant, account, state := extractFlowFields(flow)
	h.Logger.WithContext(r.Context()).ErrorF(
		"api=%s action=get_flow status=failed provider=%s provider_app=%s tenant_uuid=%s account_id=%s state=%s flow_id=%s error=%v",
		logAPISessionTokenGet,
		provider, providerApp, tenant, account, state, sanitizeLogValue(flowID), err,
	)
}

func extractFlowFields(flow *sessiontoken.Flow) (provider, providerApp, tenant, account, state string) {
	if flow == nil {
		return "-", "-", "-", "-", "-"
	}
	return sanitizeLogValue(flow.ProviderCode),
		sanitizeLogValue(flow.ProviderAppCode),
		sanitizeLogValue(flow.TenantUUID),
		sanitizeLogValue(flow.AccountID),
		sanitizeLogValue(flow.State)
}
