package station

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func GetUptime() int64 {
	if runtime.GOOS == "windows" {
		return -1
	}

	file, _ := os.Open("/proc/uptime")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return -1
	}

	fields := strings.Fields(scanner.Text())
	uptime, _ := strconv.ParseFloat(fields[0], 64)
	return int64(uptime)
}
