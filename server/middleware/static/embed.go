package static

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServeEmbed(fs *LocalFileSystem) gin.HandlerFunc {
	fileserver := http.FileServer(fs.FileSystem)
	if len(fs.Prefix) > 0 {
		fileserver = http.StripPrefix(fs.Prefix, fileserver)
	}

	return func(c *gin.Context) {
		_, err := fs.FileSystem.Open(c.Request.URL.Path)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "404 page not found",
			})
			c.Abort()

			return
		}

		fileserver.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
