package history

import (
	"fmt"

	"github.com/bclswl0827/observer/app"
	"github.com/bclswl0827/observer/driver/postgres"
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

	data, err := postgres.Query(pdb, start, end)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	var count []publisher.Geophone
	for _, v := range data {
		count = append(count, publisher.Geophone{
			TS:  v["ts"].(int64),
			EHZ: v["ehz"].([]int32),
			EHE: v["ehe"].([]int32),
			EHN: v["ehn"].([]int32),
		})
	}

	return count, nil
}
