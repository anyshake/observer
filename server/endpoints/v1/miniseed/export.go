package miniseed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (m *MiniSEED) handleExport(basePath, fileName string) ([]byte, error) {
	filePath := fmt.Sprintf("%s/%s", basePath, strings.ReplaceAll(filepath.Clean(fileName), "/", ""))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, fileInfo.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
