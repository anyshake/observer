package fdsnws

import (
	"context"
	"path"
	"strings"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
)

type provider struct {
	fdsnwsgo.UnimplementedProvider
	handler *action.Handler
	device  hardware.IHardware
	now     func() time.Time
}

func (p *provider) Now() time.Time { return p.now().UTC() }

func (p *provider) stationCodes(ctx context.Context) (network, station, location string, err error) {
	constraints := []config.IConstraint{
		&config.StationNetworkCodeConfigConstraintImpl{},
		&config.StationStationCodeConfigConstraintImpl{},
		&config.StationLocationCodeConfigConstraintImpl{},
	}
	values := make([]string, len(constraints))
	for i, constraint := range constraints {
		value, getErr := constraint.Get(ctx, p.handler)
		if getErr != nil {
			return "", "", "", getErr
		}
		values[i] = strings.TrimSpace(value.(string))
	}
	return values[0], values[1], values[2], nil
}

func matches(selector, value string) bool {
	if selector == "" {
		return true
	}
	for _, pattern := range strings.Split(selector, ",") {
		if pattern == "--" || pattern == "  " {
			pattern = ""
		}
		if matched, _ := path.Match(pattern, strings.TrimSpace(value)); matched {
			return true
		}
	}
	return false
}
