package static

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/utils/timesource"
	"github.com/gin-gonic/gin"
)

func NewEmbed(timeSource *timesource.Source, fs *LocalFileSystem) gin.HandlerFunc {
	fileserver := http.FileServer(fs.FileSystem)
	if len(fs.Prefix) > 0 {
		fileserver = http.StripPrefix(fs.Prefix, fileserver)
	}

	return func(c *gin.Context) {
		_, err := fs.FileSystem.Open(c.Request.URL.Path)
		if err != nil {
			response.Message(c, timeSource, "requested resource is not found", http.StatusNotFound, nil)
			return
		}

		fileserver.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func NewFilesystem(src embed.FS, dir string) http.FileSystem {
	fs := func(path string, f fs.FS) fs.FS {
		p, _ := fs.Sub(f, path)
		return p
	}

	return http.FS(fs(dir, src))
}
