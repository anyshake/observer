package mseed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anyshake/observer/config"
)

func getMiniSEEDList(conf *config.Conf) ([]MiniSEEDFile, error) {
	var (
		basePath  = conf.MiniSEED.Path
		station   = conf.Station.Station
		network   = conf.Station.Network
		LifeCycle = conf.MiniSEED.LifeCycle
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
			modTime := info.ModTime().UTC()

			var fileTTL int
			if LifeCycle > 0 {
				fileTTL = LifeCycle - int(time.Since(modTime).Hours()/24)
			} else {
				fileTTL = -1
			}

			files = append(files, MiniSEEDFile{
				TTL:  fileTTL,
				Name: info.Name(),
				Time: modTime.UnixMilli(),
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
