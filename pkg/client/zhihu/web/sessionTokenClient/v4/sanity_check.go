package v4

import "net/http"

func (c *Client) handleSanityCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	flow, sessiontoken, ok := c.requireFlow(w, r, apiSanityCheck)
	if !ok {
		return
	}
	c.forward(w, r, apiSanityCheck, flow, sessiontoken, http.MethodGet, "/api/v4/me", nil, nil)
}
