package typ

type AppHeader struct {
	Token string `json:"token"` // 令牌
	Nonce string `json:"nonce"` // 随机字符串
	Sign  string `json:"sign"`  // 签名
}

type WsHeader struct {
	Token string `json:"token"` // 令牌
}
