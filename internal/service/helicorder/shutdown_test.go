package helicorder

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/timesource"
)

type panickingConfigHardware struct{ *testsupport.Hardware }

func (panickingConfigHardware) GetConfig() *explorer.DeviceConfig {
	panic("hardware configuration unavailable")
}

func TestPanicAfterTimerFiresCanStop(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	synctest.Test(t, func(t *testing.T) {
		hw := panickingConfigHardware{testsupport.NewHardware()}
		svc := New(t.Context(), hw, handler, timesource.New(nil))
		if err := svc.Start(t.Context()); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Minute)
		synctest.Wait()
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		if err := svc.Stop(ctx); err != nil {
			t.Fatalf("Stop after timer panic: %v", err)
		}
		if svc.GetStatus().GetIsRunning() {
			t.Fatal("panic left the service running")
		}
	})
}
