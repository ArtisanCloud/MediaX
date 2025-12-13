package v4

import (
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) handleMeFollowings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	flow, sessionToken, ok := c.requireFlow(w, r, apiMeFollowings)
	if !ok {
		return
	}
	query := url.Values{}
	limit := clampInt(parseQueryInt(r.URL.Query().Get("limit"), 20), 1, 40)
	offset := parseQueryInt(r.URL.Query().Get("offset"), 0)
	if offset < 0 {
		offset = 0
	}
	query.Set("limit", strconv.Itoa(limit))
	query.Set("offset", strconv.Itoa(offset))
	c.forward(w, r, apiMeFollowings, flow, sessionToken, http.MethodGet, "/api/v4/me/following-columns", query, nil)
}
