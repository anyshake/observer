package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context, res int) {
	switch res {
	case 400:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不合法的请求",
			"status":  "error",
		})
	case 401:
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未经授权的访问",
			"status":  "error",
		})
	case 403:
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "本次请求被禁止",
			"status":  "error",
		})
	case 404:
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "找不到这个页面",
			"status":  "error",
		})
	case 429:
		c.JSON(http.StatusTooManyRequests, gin.H{
			"code":    http.StatusTooManyRequests,
			"message": "请求过多，现在限流中",
			"status":  "error",
		})
	case 500:
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "服务器内部错误",
			"status":  "error",
		})
	case 502:
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "服务器网关错误",
			"status":  "error",
		})
	case 503:
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    http.StatusServiceUnavailable,
			"message": "服务目前不可用",
			"status":  "error",
		})
	case 504:
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"code":    http.StatusGatewayTimeout,
			"message": "服务器网关超时",
			"status":  "error",
		})
	default:
		c.JSON(http.StatusSeeOther, gin.H{
			"code":    http.StatusSeeOther,
			"message": "其它类型错误",
			"status":  "error",
		})
	}
}
