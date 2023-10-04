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
			Message: "Unable to process this request",
		})
	case 401:
		c.JSON(http.StatusUnauthorized, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized access",
		})
	case 403:
		c.JSON(http.StatusForbidden, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusForbidden,
			Message: "Request is forbidden",
		})
	case 404:
		c.JSON(http.StatusNotFound, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusNotFound,
			Message: "Could not find this resource",
		})
	case 405:
		c.JSON(http.StatusMethodNotAllowed, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusMethodNotAllowed,
			Message: "Request method is not allowed",
		})
	case 410:
		c.JSON(http.StatusGone, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusGone,
			Message: "The requested resource is not available",
		})
	case 413:
		c.JSON(http.StatusRequestEntityTooLarge, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusMethodNotAllowed,
			Message: "Request exceeds data limit",
		})
	case 429:
		c.JSON(http.StatusTooManyRequests, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusTooManyRequests,
			Message: "Too many requests",
		})
	case 500:
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusInternalServerError,
			Message: "Server internal error",
		})
	case 502:
		c.JSON(http.StatusBadGateway, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusBadGateway,
			Message: "Server gateway error",
		})
	case 503:
		c.JSON(http.StatusServiceUnavailable, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusServiceUnavailable,
			Message: "Service is currently unavailable",
		})
	case 504:
		c.JSON(http.StatusGatewayTimeout, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusGatewayTimeout,
			Message: "Server gateway timeout",
		})
	default:
		c.JSON(http.StatusSeeOther, HttpResponse{
			Error:   true,
			Path:    currentPath,
			Time:    currentTime,
			Status:  http.StatusSeeOther,
			Message: "Unknown error occurred",
		})
	}

	c.Abort()
}
