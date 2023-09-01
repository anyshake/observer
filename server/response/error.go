package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, code int) {
	currentTime := time.Now().UTC().Format(time.RFC3339)
	currentPath := c.Request.URL.Path

	switch code {
	case 400:
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusBadRequest,
			Message: "无法处理此请求",
		})
	case 401:
		c.JSON(http.StatusUnauthorized, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusUnauthorized,
			Message: "未经授权的访问",
		})
	case 403:
		c.JSON(http.StatusForbidden, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusForbidden,
			Message: "请求被禁止",
		})
	case 404:
		c.JSON(http.StatusNotFound, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusNotFound,
			Message: "无法找到此资源",
		})
	case 405:
		c.JSON(http.StatusMethodNotAllowed, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusMethodNotAllowed,
			Message: "请求方法不被允许",
		})
	case 410:
		c.JSON(http.StatusMethodNotAllowed, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusMethodNotAllowed,
			Message: "所请求资源不可用",
		})
	case 413:
		c.JSON(http.StatusMethodNotAllowed, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusMethodNotAllowed,
			Message: "请求超出数据限制",
		})
	case 429:
		c.JSON(http.StatusTooManyRequests, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusTooManyRequests,
			Message: "请求过于频繁",
		})
	case 500:
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusInternalServerError,
			Message: "服务器内部错误",
		})
	case 502:
		c.JSON(http.StatusBadGateway, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusBadGateway,
			Message: "服务器网关错误",
		})
	case 503:
		c.JSON(http.StatusServiceUnavailable, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusServiceUnavailable,
			Message: "服务目前不可用",
		})
	case 504:
		c.JSON(http.StatusGatewayTimeout, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusGatewayTimeout,
			Message: "服务器网关超时",
		})
	default:
		c.JSON(http.StatusSeeOther, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusSeeOther,
			Message: "未知错误",
		})
	}

	c.Abort()
}
