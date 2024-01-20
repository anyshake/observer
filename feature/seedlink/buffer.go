package seedlink

import (
	"github.com/anyshake/observer/driver/seedlink"
	"github.com/anyshake/observer/publisher"
)

func (s *SeedLink) handleBuffer(gp *publisher.Geophone, buffer *seedlink.SeedLinkBuffer) error {
	if len(buffer.Data) < buffer.Size {
		buffer.Data = append(buffer.Data, *gp)
	} else {
		buffer.Data = append(buffer.Data[1:], *gp)
	}

	s.OnReady(nil, "1 record added to buffer")
	return nil
}
