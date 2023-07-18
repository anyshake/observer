package station

import (
	"os"
	"runtime"

	"github.com/wille/osutil"
)

func getOS() OS {
	hostname, _ := os.Hostname()
	osInfo := OS{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Distro:   osutil.Name,
		Hostname: hostname,
	}

	return osInfo
}
