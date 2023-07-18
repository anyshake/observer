package response

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func File(c *gin.Context, fileName string, dataBytes []byte) error {
	dataReader := bytes.NewReader(dataBytes)
	c.Header("Content-Disposition", "attachment; filename="+fileName)

	_, err := io.Copy(c.Writer, dataReader)
	if err != nil {
		Error(c, http.StatusInternalServerError)
		return err
	}

	return nil
}
