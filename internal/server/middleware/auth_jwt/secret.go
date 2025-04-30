package auth_jwt

import (
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"

	"github.com/anyshake/observer/pkg/system"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func createJwtSecret() ([]byte, error) {
	if interfaces, err := net.Interfaces(); len(interfaces) > 0 && err == nil {
		macAddrArr := lo.FilterMap(interfaces, func(iface net.Interface, _ int) (string, bool) {
			if iface.Flags&net.FlagLoopback != 0 || len(iface.HardwareAddr) == 0 {
				return "", false
			}
			return iface.HardwareAddr.String(), true
		})
		if len(macAddrArr) > 0 {
			secret := uuid.NewSHA1(uuid.NameSpaceOID, []byte(strings.Join(macAddrArr, ";"))).String()
			return []byte(secret), nil
		}
	}

	if hostname, err := os.Hostname(); len(hostname) > 0 && err == nil {
		cpuModel, err := system.GetCpuModel()
		if err != nil {
			return nil, fmt.Errorf("failed to get CPU model when creating JWT secret: %w", err)
		}
		memUsage, err := mem.VirtualMemory()
		if err != nil {
			return nil, fmt.Errorf("failed to get memory usage when creating JWT secret: %w", err)
		}
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current executable path when creating JWT secret: %w", err)
		}
		diskObj, err := disk.Usage(cwd)
		if err != nil {
			return nil, fmt.Errorf("failed to get disk usage when creating JWT secret: %w", err)
		}
		secretSeed := fmt.Sprintf(
			"%s/%s on %s, CPU model %s, RAM size %d bytes, disk size %d bytes, current run path %s",
			runtime.GOOS,
			runtime.GOARCH,
			hostname,
			cpuModel,
			memUsage.Total,
			diskObj.Total,
			cwd,
		)
		secret := uuid.NewSHA1(uuid.NameSpaceOID, []byte(secretSeed)).String()
		return []byte(secret), nil
	}

	return nil, errors.New("no suitable secret seed found")
}
