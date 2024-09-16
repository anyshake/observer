package jwt

import (
	"errors"
	"net"

	"github.com/google/uuid"
)

func createSecret() ([]byte, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagLoopback == 0 {
			return []byte(uuid.NewSHA1(uuid.NameSpaceOID, []byte(i.HardwareAddr.String())).String()), nil
		}
	}

	return nil, errors.New("no valid network interface found")
}
