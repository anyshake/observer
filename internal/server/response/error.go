package response

import (
	"github.com/gin-gonic/gin"
)

func Error(ctx *gin.Context, code int, message string) {
	res := Response{
		Error:   true,
		Code:    code,
		Message: message,
		Data:    nil,
	}
	ctx.JSON(code, res)
	ctx.Abort()
}
