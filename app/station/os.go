package station

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

func GetOS() OS {
	hostname, _ := os.Hostname()
	osInfo := OS{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Hostname: hostname,
	}

	if osInfo.OS == "linux" {
		file, _ := os.Open("/etc/os-release")
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "ID=") {
				osInfo.Distro = strings.TrimPrefix(line, "ID=")
			}
		}
	}

	return osInfo
}
