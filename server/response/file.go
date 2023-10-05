package response

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func File(c *gin.Context, fileName string, dataBytes []byte) {
	c.DataFromReader(
		http.StatusOK,
		int64(len(dataBytes)),
		"application/octet-stream",
		bytes.NewReader(dataBytes),
		map[string]string{
			"Content-Disposition": "attachment; filename=" + fileName,
		},
	)
}
