package seedlink

import (
	"sync"

	"github.com/bclswl0827/observer/feature"
)

func (s *SeedLink) Run(options *feature.FeatureOptions, waitGroup *sync.WaitGroup) {
	if !options.Config.SeedLink.Enable {
		s.OnStop(options, "service is disabled")
		return
	}
}
