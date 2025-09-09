package metadata

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed library
var library embed.FS

func getMetadataFromLibrary(model, template string) (string, error) {
	data, err := library.ReadFile(fmt.Sprintf("library/%s/%s", model, template))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getMetadataFromLocalDisk(model, template string) (string, error) {
	data, err := os.ReadFile(filepath.Join(model, template))
	if err != nil {
		return "", err
	}

	return string(data), nil
}
