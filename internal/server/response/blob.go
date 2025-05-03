package response

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Blob(ctx *gin.Context, fileName, contentType string, dataBytes []byte) {
	ctx.DataFromReader(
		http.StatusOK,
		int64(len(dataBytes)),
		contentType,
		bytes.NewReader(dataBytes),
		map[string]string{"Content-Disposition": fmt.Sprintf("attachment; filename=%s", fileName)},
	)
}
