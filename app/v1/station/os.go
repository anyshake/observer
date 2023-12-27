package station

import (
	"os"
	"runtime"

	"github.com/wille/osutil"
)

func getOS() osModel {
	hostname, _ := os.Hostname()
	osInfo := osModel{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Distro:   osutil.Name,
		Hostname: hostname,
	}

	return osInfo
}
