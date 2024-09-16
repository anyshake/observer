package static

import (
	"net/http"
)

type LocalFileSystem struct {
	Root       string
	Prefix     string
	FileSystem http.FileSystem
}
