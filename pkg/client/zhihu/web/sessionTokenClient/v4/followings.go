package v4

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	sessiontoken "github.com/ArtisanCloud/MediaX/pkg/client/sessiontoken"
	"github.com/google/uuid"
)

const defaultFollowingsInclude = "data[*].intro,followers,articles_count,voteup_count,items_count"

func (c *Client) handleMeFollowings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	flow, sessiontoken, ok := c.requireFlow(w, r, apiMeFollowings)
	if !ok {
		return
	}
	rawQuery := r.URL.Query()
	accountID, err := c.resolveAccountID(r.Context(), flow, sessiontoken, rawQuery)
	if err != nil {
		reqID := uuid.NewString()
		status := http.StatusBadRequest
		code := codeBadRequest
		logErr := err
		if typed, ok := err.(*accountIDError); ok {
			if typed.status != 0 {
				status = typed.status
			}
			if strings.TrimSpace(typed.code) != "" {
				code = typed.code
			}
			if typed.logErr != nil {
				logErr = typed.logErr
			}
		}
		c.writeError(w, reqID, status, code, err.Error())
		c.logAPICall(r.Context(), apiMeFollowings, flow, status, code, 0, logErr)
		return
	}
	c.log.InfoF("zhihu.followings: resolved account_id flow_id=%s account_id=%s query_override=%s",
		safeFlowID(flow), accountID, strings.TrimSpace(rawQuery.Get("account_id")))

	query := url.Values{}
	limit := clampInt(parseQueryInt(rawQuery.Get("limit"), 20), 1, 40)
	offset := parseQueryInt(rawQuery.Get("offset"), 0)
	if offset < 0 {
		offset = 0
	}
	query.Set("limit", strconv.Itoa(limit))
	query.Set("offset", strconv.Itoa(offset))
	include := strings.TrimSpace(rawQuery.Get("include"))
	if include == "" {
		include = defaultFollowingsInclude
	}
	query.Set("include", include)
	path := "/api/v4/members/" + url.PathEscape(accountID) + "/following-columns"
	c.log.InfoF("zhihu.followings: forwarding request flow_id=%s account_id=%s path=%s limit=%d offset=%d",
		safeFlowID(flow), accountID, path, limit, offset)
	c.forward(w, r, apiMeFollowings, flow, sessiontoken, http.MethodGet, path, query, nil)
}

func (c *Client) resolveAccountID(ctx context.Context, flow *sessiontoken.Flow, sessiontoken string, query url.Values) (string, error) {
	if candidate := queryAccountIDOverride(query); candidate != "" {
		c.log.InfoF("zhihu.followings: using account_id from query override flow_id=%s account_id=%s",
			safeFlowID(flow), candidate)
		return candidate, nil
	}
	if flow != nil {
		if candidate := accountIDFromMetadata(flow.Metadata); candidate != "" {
			c.log.InfoF("zhihu.followings: using account_id from flow metadata flow_id=%s account_id=%s",
				safeFlowID(flow), candidate)
			return candidate, nil
		}
		if trimmed := sanitizeAccountID(flow.AccountID); trimmed != "" {
			c.log.InfoF("zhihu.followings: using sanitized flow account_id flow_id=%s account_id=%s",
				flow.FlowID, trimmed)
			return trimmed, nil
		}
		c.log.WarnF("zhihu.followings: account_id missing in flow flow_id=%s metadata=%v account_id=%s",
			flow.FlowID, flow.Metadata, flow.AccountID)
	}
	return c.fetchAccountIDFromProfile(ctx, sessiontoken)
}

func queryAccountIDOverride(values url.Values) string {
	if len(values) == 0 {
		return ""
	}
	keys := []string{"account_id", "accountId", "url_token", "urlToken", "uid"}
	for _, key := range keys {
		if value := strings.TrimSpace(values.Get(key)); value != "" {
			return value
		}
	}
	return ""
}

func accountIDFromMetadata(meta map[string]string) string {
	if len(meta) == 0 {
		return ""
	}
	fallbackKeys := []string{
		"account_id",
		"zhihu_account_id",
		"zhihu_account",
		"url_token",
		"urlToken",
		"uid",
		"zhihu_uid",
	}
	for _, key := range fallbackKeys {
		if value := strings.TrimSpace(meta[key]); value != "" {
			return value
		}
	}
	return ""
}

func sanitizeAccountID(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if trimmed == "acct_debug" {
		return ""
	}
	return trimmed
}

func (c *Client) fetchAccountIDFromProfile(ctx context.Context, sessiontoken string) (string, error) {
	c.log.Info("zhihu.followings: fetching account_id via /api/v4/me")
	status, body, upstream, err := c.callZhihu(ctx, http.MethodGet, "/api/v4/me", nil, nil, sessiontoken)
	if err != nil {
		return "", &accountIDError{
			status:  http.StatusBadGateway,
			code:    codeUpstreamError,
			message: "failed to contact Zhihu",
			logErr:  fmt.Errorf("resolve account_id via %s: %w", upstream, err),
		}
	}
	if status >= 200 && status < 300 {
		var profile meProfileResponse
		if err := json.Unmarshal(body, &profile); err != nil {
			return "", &accountIDError{
				status:  http.StatusBadGateway,
				code:    codeUpstreamError,
				message: "failed to decode Zhihu profile response",
				logErr:  fmt.Errorf("decode profile: %w", err),
			}
		}
		accountID := firstNonEmpty(profile.URLToken, profile.URLTokenCamel, profile.ID)
		if accountID == "" {
			return "", &accountIDError{
				status:  http.StatusBadRequest,
				code:    codeBadRequest,
				message: "account_id not found, please pass ?account_id=<url_token>",
				logErr:  fmt.Errorf("profile missing url_token/id body=%s", truncateBody(body)),
			}
		}
		c.log.InfoF("zhihu.followings: resolved account_id from profile url_token=%s id=%s",
			profile.URLToken, profile.ID)
		return accountID, nil
	}
	mappedStatus, code, _ := classifyStatus(status)
	message := c.extractErrorMessage(body, defaultErrorMessage(code))
	return "", &accountIDError{
		status:  mappedStatus,
		code:    code,
		message: message,
		logErr:  fmt.Errorf("resolve account_id via %s status=%d message=%s body=%s", upstream, status, message, truncateBody(body)),
	}
}

type meProfileResponse struct {
	URLToken      string `json:"url_token"`
	URLTokenCamel string `json:"urlToken"`
	ID            string `json:"id"`
}

type accountIDError struct {
	status  int
	code    string
	message string
	logErr  error
}

func (e *accountIDError) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

func safeFlowID(flow *sessiontoken.Flow) string {
	if flow == nil {
		return ""
	}
	return flow.FlowID
}

func truncateBody(body []byte) string {
	max := 512
	if len(body) <= max {
		return string(body)
	}
	return string(body[:max]) + "...(truncated)"
}
