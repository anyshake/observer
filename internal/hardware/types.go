package hardware

import (
	"context"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/metadata"
)

type IHardware interface {
	Open(context.Context) (context.Context, context.CancelFunc, error)
	Close() error
	Flush() error

	Subscribe(clientId string, handler explorer.EventHandler) error
	Unsubscribe(clientId string) error

	GetConfig() explorer.DeviceConfig
	GetStatus() explorer.DeviceStatus

	GetCoordinates(fuzzy bool) (float64, float64, float64, error)
	GetTemperature() (float64, error)
	GetDeviceId() string

	GetMetadata(
		stationAffiliation,
		stationDescription,
		stationCountry,
		stationPlace,
		networkCode,
		stationCode,
		locationCode string,
		fuzzyCoordinates bool,
	) (*metadata.Render, error)
}
