package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessHandler(c *gin.Context, res int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": msg,
		"status":  "success",
	})
}
