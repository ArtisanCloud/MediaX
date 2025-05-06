package schema

// DouYinOAuthGetOpenIdByCReq 通过C端OpenID获取B端OpenID请求结构体
type DouYinOAuthGetOpenIdByCReq struct {
	BClientKey string `json:"b_client_key"` // 绑定的B端应用ID
	CClientKey string `json:"c_client_key"` // C端移动网站应用对应的应用ID
	OpenID     string `json:"open_id"`      // 用户在当前C端移动网站应用下的OpenID
}

// OAuthGetOpenIdByC 通过C端OpenID获取B端OpenID响应结构体
type OAuthGetOpenIdByC struct {
	OpenId string `json:"open_id"` // 用户在B端应用下的OpenID
}

// DouYinOAuthGetOpenIdByCRes 通过C端OpenID获取B端OpenID响应结构体
type DouYinOAuthGetOpenIdByCRes struct {
	ErrNo  int               `json:"err_no"`
	ErrMsg string            `json:"err_msg"`
	LogId  string            `json:"log_id"`
	Data   OAuthGetOpenIdByC `json:"data"`
}
