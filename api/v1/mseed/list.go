package mseed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (m *MSeed) getMiniSeedList(basePath, stationCode, networkCode string, lifeCycle int) ([]miniSeedFileInfo, error) {
	var files []miniSeedFileInfo
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() &&
			filepath.Ext(info.Name()) == ".mseed" &&
			strings.Contains(info.Name(), stationCode) &&
			strings.Contains(info.Name(), networkCode) {
			modTime := info.ModTime().UTC()

			// Calculate file TTL
			var fileTTL int
			if lifeCycle > 0 {
				fileTTL = lifeCycle - int(time.Since(modTime).Hours()/24)
			} else {
				fileTTL = -1
			}

			files = append(files, miniSeedFileInfo{
				TTL:  fileTTL,
				Name: info.Name(),
				Time: modTime.UTC().UnixMilli(),
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
