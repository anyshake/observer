package quakesense

import (
	"context"
	"sync"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/eewgo"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const ID = "service_quakesense"

const (
	FILTER_NUM_TAPS  = 101
	NO_FILTER        = "nofilter"
	LOW_PASS_FILTER  = "lowpass"
	BAND_PASS_FILTER = "bandpass"
	HIGH_PASS_FILTER = "highpass"
)

const (
	CLASSIC_STA_LTA   = "classic_sta_lta"
	RECURSIVE_STA_LTA = "recursive_sta_lta"
	DELAYED_STA_LTA   = "delayed_sta_lta"
	Z_DETECT          = "z_detect"
)

type QuakeSenseServiceImpl struct {
	mu     sync.Mutex
	status service.Status

	wg       sync.WaitGroup
	ctx      context.Context
	cancelFn context.CancelFunc

	hardwareDev   hardware.IHardware
	timeSource    *timesource.Source
	actionHandler *action.Handler

	channelBuffer  *ringbuf.Buffer[float64]
	prevSamplerate int

	mqttBroker   string
	mqttTopic    string
	mqttUsername string
	mqttPassword string
	mqttClient   mqtt.Client

	stationName        string
	stationDescription string
	stationPlace       string
	stationCountry     string
	stationAffiliation string

	stationCode  string
	networkCode  string
	locationCode string

	monitorChannel  string
	throttleSeconds int

	triggerMethod string
	staWindow     float64
	ltaWindow     float64
	trigOn        float64
	trigOff       float64

	filterType   string
	maxFreq      float64
	minFreq      float64
	filterKernel *eewgo.FIRFilter
}
