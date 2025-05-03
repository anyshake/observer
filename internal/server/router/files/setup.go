package files

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anyshake/observer/internal/server/response"
	"github.com/anyshake/observer/internal/service"
	"github.com/gin-gonic/gin"
)

func Setup(routerGroup *gin.RouterGroup, serviceMap map[string]service.IService, jwtMiddleware gin.HandlerFunc) {
	secretKey, _ := getSecretKey(32)

	routerGroup.POST("/files", jwtMiddleware, func(ctx *gin.Context) {
		var requestModel struct {
			FilePath string `form:"file_path" json:"file_path" xml:"file_path" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&requestModel); err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid request")
			return
		}
		expireAt := time.Now().Add(TOKEN_LIFETIME).UnixMilli()
		token := generateToken(requestModel.FilePath, secretKey, expireAt)
		response.Data(ctx, http.StatusOK, "successfully generated token for this asset", token)
	})
	routerGroup.GET("/files", func(ctx *gin.Context) {
		var requestModel struct {
			Namespace string `form:"namespace" json:"namespace" xml:"namespace" binding:"required"`
			FilePath  string `form:"file_path" json:"file_path" xml:"file_path" binding:"required"`
			Token     string `form:"token" json:"token" xml:"token" binding:"required"`
		}
		if err := ctx.ShouldBind(&requestModel); err != nil {
			response.Error(ctx, http.StatusBadRequest, "request body is not valid")
			return
		}
		if !validateToken(requestModel.FilePath, requestModel.Token, secretKey) {
			response.Error(ctx, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		service, ok := serviceMap[requestModel.Namespace]
		if !ok {
			response.Error(ctx, http.StatusBadRequest, fmt.Sprintf("service ID %s was not found", requestModel.Namespace))
			return
		}
		asset, err := service.GetAssetData(requestModel.FilePath)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}
		response.Blob(ctx, asset.FileName, asset.ContentType, asset.Data)
	})
}
