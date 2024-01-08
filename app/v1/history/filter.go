package history

import (
	"fmt"
	"time"

	"github.com/anyshake/observer/app"
	"github.com/anyshake/observer/driver/dao"
	"github.com/anyshake/observer/publisher"
)

func filterHistory(start, end int64, limit time.Duration, options *app.ServerOptions) ([]publisher.Geophone, error) {
	pdb := options.FeatureOptions.Database
	if pdb == nil {
		return nil, fmt.Errorf("databse is not connected")
	}

	if end-start > limit.Milliseconds() {
		return nil, fmt.Errorf("duration is too large")
	}

	data, err := dao.Query(pdb, start, end)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return data, nil
}
