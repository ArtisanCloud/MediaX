package v4

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (c *Client) handleArticlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	flow, sessionToken, ok := c.requireFlow(w, r, apiArticlePost)
	if !ok {
		return
	}
	defer r.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		reqID := uuid.NewString()
		c.writeError(w, reqID, http.StatusBadRequest, codeBadRequest, "failed to read request body")
		c.logAPICall(r.Context(), apiArticlePost, flow, http.StatusBadRequest, codeBadRequest, 0, err)
		return
	}
	if len(bytes.TrimSpace(raw)) == 0 {
		reqID := uuid.NewString()
		c.writeError(w, reqID, http.StatusBadRequest, codeBadRequest, "request body is empty")
		c.logAPICall(r.Context(), apiArticlePost, flow, http.StatusBadRequest, codeBadRequest, 0, errors.New("empty body"))
		return
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		reqID := uuid.NewString()
		c.writeError(w, reqID, http.StatusBadRequest, codeBadRequest, "invalid JSON payload")
		c.logAPICall(r.Context(), apiArticlePost, flow, http.StatusBadRequest, codeBadRequest, 0, err)
		return
	}
	title := extractStringField(payload["title"])
	content := extractStringField(payload["content"])
	if title == "" || content == "" {
		reqID := uuid.NewString()
		c.writeError(w, reqID, http.StatusBadRequest, codeBadRequest, "title and content are required")
		c.logAPICall(r.Context(), apiArticlePost, flow, http.StatusBadRequest, codeBadRequest, 0, errors.New("missing title/content"))
		return
	}
	c.forward(w, r, apiArticlePost, flow, sessionToken, http.MethodPost, "/api/v4/articles", nil, raw)
}

func extractStringField(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	default:
		return ""
	}
}
