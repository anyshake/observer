package hardware_test

import (
	"testing"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestNewSelectsProtocol(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	if err := (&config.StationChannelCodesConfigConstraintImpl{}).Init(handler); err != nil {
		t.Fatal(err)
	}
	clock := timesource.New(func() time.Time { return time.Unix(0, 0) })
	options := explorer.ExplorerOptions{Endpoint: "tcp://127.0.0.1:1", ReadTimeout: 1}
	for _, protocol := range []string{"legacy", "v1", "v2", "v3"} {
		options.Protocol = protocol
		device, err := hardware.New(logger.GetLogger("hardware"), clock, handler, options, explorer.NtpOptions{})
		if err != nil || device == nil {
			t.Fatalf("protocol %s: %v", protocol, err)
		}
	}

	options.Protocol = "v9"
	if _, err := hardware.New(logger.GetLogger("hardware"), clock, handler, options, explorer.NtpOptions{}); err == nil {
		t.Fatal("unknown protocol accepted")
	}
	options.Protocol = "v1"
	options.Endpoint = "://bad"
	if _, err := hardware.New(logger.GetLogger("hardware"), clock, handler, options, explorer.NtpOptions{}); err == nil {
		t.Fatal("invalid endpoint accepted")
	}
	_, empty := testsupport.OpenDAO(t)
	options.Endpoint = "tcp://127.0.0.1:1"
	if _, err := hardware.New(logger.GetLogger("hardware"), clock, empty, options, explorer.NtpOptions{}); err == nil {
		t.Fatal("missing channel codes accepted")
	}
}
