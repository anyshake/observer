package helicorder

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (m *HeliCorder) handleList(basePath, stationCode, networkCode string, lifeCycle int) ([]heliCorderFileInfo, error) {
	var files []heliCorderFileInfo

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".svg" &&
			strings.Contains(entry.Name(), stationCode) &&
			strings.Contains(entry.Name(), networkCode) {

			info, err := entry.Info()
			if err != nil {
				return nil, err
			}

			modTime := info.ModTime().UTC()

			var fileTTL int
			if lifeCycle > 0 {
				fileTTL = lifeCycle - int(time.Since(modTime).Hours()/24)
			} else {
				fileTTL = -1
			}

			files = append(files, heliCorderFileInfo{
				TTL:  fileTTL,
				Name: info.Name(),
				Time: modTime.UTC().UnixMilli(),
				Size: info.Size(),
			})
		}
	}

	return files, nil
}
