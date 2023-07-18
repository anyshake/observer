package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func Message(c *gin.Context, message string, data any) error {
	currentTime := time.Now().UTC().Format(time.RFC3339)
	currentPath := c.Request.URL.Path

	response := HttpResponse{
		Error:   false,
		Path:    currentPath,
		Time:    currentTime,
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}

	jsonBytes, err := jsoniter.Marshal(response)
	if err != nil {
		Error(c, http.StatusInternalServerError)
		return err
	}

	c.Data(http.StatusOK, "application/json", jsonBytes)
	return nil
}
