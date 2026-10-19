package winston

import (
	"context"
	"fmt"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/pkg/ringbuf"
)

func (s *WinstonServiceImpl) Init(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return fmt.Errorf("service %s cannot be initialized while running", ID)
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(ctx, s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	stationCode, err := (&config.StationStationCodeConfigConstraintImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.stationCode = stationCode.(string)

	networkCode, err := (&config.StationNetworkCodeConfigConstraintImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.networkCode = networkCode.(string)

	locationCode, err := (&config.StationLocationCodeConfigConstraintImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.locationCode = locationCode.(string)

	listenHost, err := (&winstonConfigListenHostImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.listenHost = listenHost.(string)

	listenPort, err := (&winstonConfigListenPortImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.listenPort = listenPort.(int)

	bufferSize, err := (&winstonConfigBufferSizeImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.bufferSize = bufferSize.(int)
	s.ringBuffer = ringbuf.New[winstonRingBuffer](s.bufferSize)

	return nil
}
