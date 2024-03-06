package seedlink

import (
	"time"

	"github.com/anyshake/observer/driver/seedlink"
	"github.com/anyshake/observer/publisher"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

func (s *SeedLink) handleBuffer(gp *publisher.Geophone, buffer *seedlink.SeedLinkBuffer) error {
	currentTime := time.UnixMilli(gp.TS).UTC()
	if currentTime.Minute()%10 == 0 && currentTime.Second() == 0 {
		expireThreshold := currentTime.Add(-buffer.Duration).UnixMilli()
		buffer.Database.Delete(query.NewQuery(buffer.Collection).Where(query.Field("ts").Lt(expireThreshold)))
	}

	var (
		ehz, _ = gp.EHZ.Value()
		ehe, _ = gp.EHE.Value()
		ehn, _ = gp.EHN.Value()
	)
	doc := document.NewDocument()
	doc.Set("ehz", ehz.(string))
	doc.Set("ehe", ehe.(string))
	doc.Set("ehn", ehn.(string))
	doc.Set("ts", gp.TS)

	_, err := buffer.Database.InsertOne(buffer.Collection, doc)
	if err != nil {
		return err
	}

	s.OnReady(nil, "1 record added to buffer")
	return nil
}
