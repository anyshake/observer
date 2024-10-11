package helicorder

import "time"

const (
	HELICORDER_TIME_SPAN         = 15 * time.Minute
	HELICORDER_IMAGE_SIZE        = 1000
	HELICORDER_DOWNSAMPLE_FACTOR = 5000
	HELICORDER_SCALE_FACTOR      = 2.2
	HELICORDER_LINE_WIDTH        = 0.8
)

type HelicorderService struct {
	basePath     string
	lifeCycle    int
	stationCode  string
	networkCode  string
	locationCode string
}
