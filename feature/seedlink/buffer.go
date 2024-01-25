package seedlink

import (
	"os"
	"time"

	"github.com/anyshake/observer/driver/seedlink"
	"github.com/anyshake/observer/publisher"
	jsonutil "github.com/multiprocessio/go-json"
)

func (s *SeedLink) handleBuffer(gp *publisher.Geophone, buffer *seedlink.SeedLinkBuffer) error {
	if len(buffer.Data) < buffer.Size {
		buffer.Data = append(buffer.Data, *gp)
	} else {
		buffer.Data = append(buffer.Data[1:], *gp)
	}

	// Write buffer to file every 10 minutes
	if time.Now().UTC().Minute()%10 == 0 {
		file, err := os.OpenFile(buffer.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		err = jsonutil.Encode(file, buffer.Data)
		if err != nil {
			return err
		}

		s.OnReady(nil, "buffer file updated successfully")
	}

	s.OnReady(nil, "1 record added to buffer")
	return nil
}
