package seedlink

import (
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/publisher"
)

func (s *SeedLink) handleMessage(gp *publisher.Geophone, options *feature.FeatureOptions) error {
	// Append channel to geophone
	s.OnReady(options)
	return nil
}
