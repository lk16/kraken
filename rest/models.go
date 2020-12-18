package rest

type WebSocketTokenResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}
