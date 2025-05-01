package response

import (
	"github.com/gin-gonic/gin"
)

func Error(ctx *gin.Context, code int, message string) {
	ctx.AbortWithStatusJSON(code, Response{
		Error:   true,
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
