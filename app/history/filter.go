package history

import (
	"fmt"

	"github.com/bclswl0827/observer/app"
	"github.com/bclswl0827/observer/driver/dao"
	"github.com/bclswl0827/observer/publisher"
)

func filterHistory(start, end int64, options *app.ServerOptions) ([]publisher.Geophone, error) {
	pdb := options.FeatureOptions.Database
	if pdb == nil {
		return nil, fmt.Errorf("databse is not connected")
	}

	if end-start > DURATION.Milliseconds() {
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
