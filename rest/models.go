package rest

type Response struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

type WebSocketToken struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}
