package metadata

import (
	"embed"
	"fmt"
	"maps"
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

func mergeMap(a, b map[string]string) map[string]string {
	out := make(map[string]string, len(a)+len(b))
	maps.Copy(out, a)
	maps.Copy(out, b)
	return out
}
