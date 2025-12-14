package v4

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) handleChannelArticles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	channelID, ok := extractChannelID(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	flow, sessionToken, ok := c.requireFlow(w, r, apiChannelsArticles)
	if !ok {
		return
	}
	query := url.Values{}
	limit := clampInt(parseQueryInt(r.URL.Query().Get("limit"), 20), 1, 50)
	offset := parseQueryInt(r.URL.Query().Get("offset"), 0)
	if offset < 0 {
		offset = 0
	}
	query.Set("limit", strconv.Itoa(limit))
	query.Set("offset", strconv.Itoa(offset))
	upstream := fmt.Sprintf("/api/v4/columns/%s/items", url.PathEscape(channelID))
	c.forward(w, r, apiChannelsArticles, flow, sessionToken, http.MethodGet, upstream, query, nil)
}
