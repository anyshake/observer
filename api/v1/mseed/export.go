package mseed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (m *MSeed) getMiniSeedBytes(basePath, fileName string) ([]byte, error) {
	fileName = filepath.Clean(fileName)
	fileName = strings.ReplaceAll(fileName, "/", "")

	filePath := fmt.Sprintf("%s/%s", basePath, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil
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
