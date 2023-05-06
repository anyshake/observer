package history

import (
	"encoding/json"
	"fmt"

	"com.geophone.observer/app"
	"com.geophone.observer/common/postgres"
	"com.geophone.observer/features/geophone"
)

func FilterHistory(timestamp int64, options *app.ServerOptions) ([]geophone.Acceleration, error) {
	if options.ConnPostgres == nil {
		return nil, fmt.Errorf("postgres connection is nil")
	}

	data, err := postgres.SelectData(options.ConnPostgres, timestamp)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	var acceleration []geophone.Acceleration
	for _, v := range data {
		var acc geophone.Acceleration
		err := json.Unmarshal([]byte(v["data"].(string)), &acc)
		if err != nil {
			return nil, err
		}

		acceleration = append(acceleration, acc)
	}

	return acceleration, nil
}
