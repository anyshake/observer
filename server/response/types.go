package response

type HttpResponse struct {
	Time    string `json:"time"`
	Status  int    `json:"status"`
	Error   bool   `json:"error"`
	Path    string `json:"path"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
