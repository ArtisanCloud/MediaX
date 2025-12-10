package session_token

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken"
	sessionapi "github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/api"
	"github.com/ArtisanCloud/MediaX/pkg/client/sessionToken/sanitizer"
	sessionmiddleware "github.com/ArtisanCloud/MediaX/pkg/server/middleware/session_token"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
)

const (
	// SessionTokenFlowCreatePath 提供 Flow 创建 API 的 HTTP path。
	SessionTokenFlowCreatePath = "/session-token/flows"
	logAPISessionTokenCreate   = "session_token.create_flow"
)

var (
	errNilSessionTokenManager = errors.New("sessiontoken: session token manager is nil")
	errNilServeMux            = errors.New("sessiontoken: http mux is nil")
	errMissingAPIToken        = errors.New("sessiontoken: api token is required")
)

// SessionTokenFlowCreateHandler 处理 Flow 创建请求。
type SessionTokenFlowCreateHandler struct {
	Manager sessiontoken.SessionTokenClient
	Logger  *logger.Logger
}

// RegisterSessionTokenFlowCreateRoute 将 Flow 创建 handler 注册到 mux 上，并强制校验 Bearer Token。
func RegisterSessionTokenFlowCreateRoute(mux *http.ServeMux, manager sessiontoken.SessionTokenClient, apiToken string, log *logger.Logger) error {
	if mux == nil {
		return errNilServeMux
	}
	if manager == nil {
		return errNilSessionTokenManager
	}
	if strings.TrimSpace(apiToken) == "" {
		return errMissingAPIToken
	}
	handler := NewSessionTokenFlowCreateHandler(manager, log)
	mux.Handle(SessionTokenFlowCreatePath, sessionmiddleware.SessionTokenAuthMiddleware(apiToken, log)(handler))
	return nil
}

// NewSessionTokenFlowCreateHandler 创建 handler。
func NewSessionTokenFlowCreateHandler(manager sessiontoken.SessionTokenClient, log *logger.Logger) *SessionTokenFlowCreateHandler {
	return &SessionTokenFlowCreateHandler{Manager: manager, Logger: log}
}

// ServeHTTP 实现 http.Handler 接口。
func (h *SessionTokenFlowCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if r.Body == nil {
		h.logError(r, nil, "", errors.New("sessiontoken: empty request body"))
		h.writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	defer r.Body.Close()

	var req sessionapi.CreateFlowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		wrapped := fmt.Errorf("sessiontoken: decode request body: %w", err)
		h.logError(r, nil, "", wrapped)
		h.writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	opts, err := req.ToOptions()
	if err != nil {
		h.logError(r, &req, "", fmt.Errorf("sessiontoken: invalid request payload: %w", err))
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	flow, err := h.Manager.CreateFlow(r.Context(), opts)
	if err != nil {
		h.logError(r, &req, "", fmt.Errorf("sessiontoken: manager create flow failed: %w", err))
		h.writeError(w, http.StatusInternalServerError, "failed to create flow")
		return
	}
	h.logSuccess(r, &req, flow.FlowID)
	h.writeSuccess(w, flow)
}

func (h *SessionTokenFlowCreateHandler) writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (h *SessionTokenFlowCreateHandler) writeSuccess(w http.ResponseWriter, flow *sessiontoken.Flow) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resFlow := *flow
	resFlow.Result = sanitizer.MaskCredentialPayload(flow.Result)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"flow": &resFlow})
}

func (h *SessionTokenFlowCreateHandler) logSuccess(r *http.Request, req *sessionapi.CreateFlowRequest, flowID string) {
	if h.Logger == nil {
		return
	}
	provider, providerApp, tenant, account, state := extractRequestFields(req)
	h.Logger.WithContext(r.Context()).InfoF(
		"api=%s action=create_flow status=success provider=%s provider_app=%s tenant_uuid=%s account_id=%s state=%s flow_id=%s",
		logAPISessionTokenCreate,
		provider, providerApp, tenant, account, state, sanitizeLogValue(flowID),
	)
}

func (h *SessionTokenFlowCreateHandler) logError(r *http.Request, req *sessionapi.CreateFlowRequest, flowID string, err error) {
	if h.Logger == nil || err == nil {
		return
	}
	provider, providerApp, tenant, account, state := extractRequestFields(req)
	h.Logger.WithContext(r.Context()).ErrorF(
		"api=%s action=create_flow status=failed provider=%s provider_app=%s tenant_uuid=%s account_id=%s state=%s flow_id=%s error=%v",
		logAPISessionTokenCreate,
		provider, providerApp, tenant, account, state, sanitizeLogValue(flowID), err,
	)
}

func extractRequestFields(req *sessionapi.CreateFlowRequest) (provider, providerApp, tenant, account, state string) {
	if req == nil {
		return "-", "-", "-", "-", "-"
	}
	return sanitizeLogValue(req.ProviderCode),
		sanitizeLogValue(req.ProviderAppCode),
		sanitizeLogValue(req.TenantUUID),
		sanitizeLogValue(req.AccountID),
		sanitizeLogValue(req.State)
}

func sanitizeLogValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}
