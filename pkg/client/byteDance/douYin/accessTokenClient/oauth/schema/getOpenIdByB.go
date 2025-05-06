package schema

// DouYinOAuthGetOpenIdByBReq 通过B端OpenID获取C端OpenID请求结构体
type DouYinOAuthGetOpenIdByBReq struct {
	BClientKey string `json:"b_client_key"` // B端移动网站应用对应的应用ID
	CClientKey string `json:"c_client_key"` // 绑定的C端应用ID
	OpenID     string `json:"open_id"`      // 用户在绑定的C端应用下的OpenID
}

// DouYinOAuthGetOpenIdByBRes 通过B端OpenID获取C端OpenID响应结构体
type DouYinOAuthGetOpenIdByBRes struct {
	BClientKey string `json:"b_client_key"` // B端移动网站应用对应的应用ID
	CClientKey string `json:"c_client_key"` // 绑定的C端应用ID
	OpenId     string `json:"open_id"`      // 用户在绑定的C端应用下的OpenID
}
