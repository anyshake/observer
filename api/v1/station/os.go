package station

import (
	"os"
	"runtime"

	"github.com/anyshake/observer/utils/timesource"
	"github.com/mackerelio/go-osstat/uptime"
	"github.com/wille/osutil"
)

func (o *osInfo) get(timeSource *timesource.Source) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	up, err := uptime.Get()
	if err != nil {
		return err
	}

	timestamp, err := timeSource.Get()
	if err != nil {
		return err
	}

	o.Uptime = int64(up.Seconds())
	o.OS = runtime.GOOS
	o.Arch = runtime.GOARCH
	o.Distro = osutil.Name
	o.Hostname = hostname
	o.Timestamp = timestamp.UnixMilli()
	return nil
}
