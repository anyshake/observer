package fdsnws

import (
	"context"

	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
	"github.com/bclswl0827/fdsnws-go/service/dataselect"
)

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
