package history

import (
	"fmt"

	"com.geophone.observer/app"
	"com.geophone.observer/common/postgres"
	"com.geophone.observer/handler"
)

func FilterHistory(start, end int64, options *app.ServerOptions) ([]handler.Geophone, error) {
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

	var count []handler.Geophone
	for _, v := range data {
		count = append(count, handler.Geophone{
			TS:  v["ts"].(int64),
			EHZ: v["ehz"].([]int32),
			EHE: v["ehe"].([]int32),
			EHN: v["ehn"].([]int32),
		})
	}

	return count, nil
}
