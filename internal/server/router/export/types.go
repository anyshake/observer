package export

import (
	"time"

	"github.com/anyshake/observer/internal/dao/model"
)

const LOG_PREFIX = "restful_api_export"

type seismicDataEncoder interface {
	GetName() string
	Encode(records []model.SeisRecord, channel string) ([]byte, error)
	GetFileName(startTime time.Time, channelCode string) (string, error)
}
