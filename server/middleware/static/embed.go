package static

import (
	"net/http"

	"github.com/bclswl0827/observer/server/response"
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
			response.Error(c, http.StatusNotFound)
			return
		}

		fileserver.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
