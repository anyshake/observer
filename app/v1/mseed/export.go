package mseed

import (
	"fmt"
	"os"
	"strings"

	"github.com/anyshake/observer/config"
)

func getMiniSEEDBytes(conf *config.Conf, fileName string) ([]byte, error) {
	// Remove slash in file name to avoid path traversal
	fileName = strings.ReplaceAll(fileName, "\\", "")
	fileName = strings.ReplaceAll(fileName, "/", "")
	basePath := conf.MiniSEED.Path

	// Check if file exists, return nil if not exists to avoid 500 error
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
