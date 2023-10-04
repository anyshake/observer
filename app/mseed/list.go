package mseed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bclswl0827/observer/config"
)

func getMiniSEEDList(conf *config.Conf) ([]MiniSEEDFile, error) {
	var (
		basePath = conf.MiniSEED.Path
		station  = conf.MiniSEED.Station
		network  = conf.MiniSEED.Network
	)

	var files []MiniSEEDFile
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() &&
			filepath.Ext(info.Name()) == ".mseed" &&
			strings.Contains(info.Name(), station) &&
			strings.Contains(info.Name(), network) {
			files = append(files, MiniSEEDFile{
				Name: info.Name(),
				Time: info.ModTime().UTC().Format(time.RFC3339),
				Size: fmt.Sprintf("%d MB", info.Size()/1024/1024),
			})
		}

		return nil
	}

	err := filepath.Walk(basePath, walkFn)
	if err != nil {
		return nil, err
	}

	return files, nil
}
