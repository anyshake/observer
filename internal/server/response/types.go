package response

type Response struct {
	Code    int    `json:"code"`
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
