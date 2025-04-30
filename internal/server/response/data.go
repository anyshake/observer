package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Data(ctx *gin.Context, code int, message string, data any) {
	res := Response{
		Error:   false,
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	}
	ctx.JSON(http.StatusOK, res)
}
