package fdsnws

import (
	"context"

	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
	"github.com/bclswl0827/fdsnws-go/service/dataselect"
	"github.com/bclswl0827/fdsnws-go/service/station"
)

// stationService binds request contexts to the library's context-free provider API.
type stationService struct {
	*station.Component
	provider *provider
}

func newStationService(provider *provider) *stationService {
	return &stationService{
		Component: station.New(provider),
		provider:  provider,
	}
}

func (s *stationService) Handle(ctx context.Context, endpoint string, request fdsnwsgo.Request) (fdsnwsgo.Outcome, error) {
	requestProvider := &stationProvider{provider: s.provider, ctx: ctx}
	return station.New(requestProvider).Handle(ctx, endpoint, request)
}

// dataSelectService binds request contexts to the library's context-free provider API.
type dataSelectService struct {
	*dataselect.Component
	provider *provider
}

func newDataSelectService(provider *provider) *dataSelectService {
	return &dataSelectService{
		Component: dataselect.New(provider),
		provider:  provider,
	}
}

func (s *dataSelectService) Handle(ctx context.Context, endpoint string, request fdsnwsgo.Request) (fdsnwsgo.Outcome, error) {
	requestProvider := &dataSelectProvider{provider: s.provider, ctx: ctx}
	return dataselect.New(requestProvider).Handle(ctx, endpoint, request)
}
