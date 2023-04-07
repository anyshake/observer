package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MessageHandler(c *gin.Context, message string, data any) HttpResponse {
	currentTime := time.Now().UTC().Format(time.RFC3339)
	currentPath := c.Request.URL.Path

	return HttpResponse{
		Error:   false,
		Path:    currentPath,
		Time:    currentTime,
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
}
