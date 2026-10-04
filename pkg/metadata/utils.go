package metadata

import (
	"bytes"
	"embed"
	"encoding/xml"
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

func escapeXMLMap(data map[string]string) map[string]string {
	for key, value := range data {
		var escaped bytes.Buffer
		_ = xml.EscapeText(&escaped, []byte(value))
		data[key] = escaped.String()
	}
	return data
}
