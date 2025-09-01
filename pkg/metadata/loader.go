package metadata

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed library
var library embed.FS

func getMetadataFromLibrary(model, template string) (string, error) {
	data, err := library.ReadFile(filepath.Join("library", model, template))
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
