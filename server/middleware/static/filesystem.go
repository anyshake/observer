package static

import (
	"embed"
	"io/fs"
	"net/http"
)

func CreateFilesystem(src embed.FS, dir string) http.FileSystem {
	fs := func(path string, f fs.FS) fs.FS {
		p, _ := fs.Sub(f, path)
		return p
	}

	return http.FS(fs(dir, src))
}
