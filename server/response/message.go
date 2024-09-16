package response

import (
	"time"

	"github.com/anyshake/observer/utils/timesource"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func Message(ctx *gin.Context, timeSource *timesource.Source, message string, code int, data any) error {
	var currentTime string
	if timeSource != nil {
		currentTime = timeSource.Get().Format(time.RFC3339)
	} else {
		currentTime = time.Now().UTC().Format(time.RFC3339)
	}

	response := HttpResponse{
		Error:   code >= 400 && code < 600,
		Path:    ctx.Request.URL.Path,
		Time:    currentTime,
		Status:  code,
		Message: message,
		Data:    data,
	}

	jsonBytes, _ := jsoniter.Marshal(response)
	ctx.Data(code, "application/json", jsonBytes)
	ctx.Abort()
	return nil
}
