package v4

import (
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) handleArticleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	articleID, ok := extractArticleID(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	flow, sessionToken, ok := c.requireFlow(w, r, apiArticleGet)
	if !ok {
		return
	}
	upstream := fmt.Sprintf("/api/v4/articles/%s", url.PathEscape(articleID))
	c.forward(w, r, apiArticleGet, flow, sessionToken, http.MethodGet, upstream, nil, nil)
}
