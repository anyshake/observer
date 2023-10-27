package seedlink

import (
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/publisher"
)

func (s *SeedLink) handleMessage(gp *publisher.Geophone, options *feature.FeatureOptions) error {
	// Append channel to geophone
	s.OnReady(options)
	return nil
}
