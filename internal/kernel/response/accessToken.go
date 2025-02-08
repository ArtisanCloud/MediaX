package response

type AccessTokenRes struct {
	AccessToken string  `json:"access_token,omitempty"`
	ExpiresIn   float64 `json:"expires_in"`
}
