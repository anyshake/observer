package archiver

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func WriteMessage(rdb *redis.Client, options *ArchiverOptions) {
	if !options.Enable {
		return
	}

	data, err := json.Marshal(options.Message)
	if err != nil {
		options.OnErrorCallback(err)
		return
	}

	err = rdb.Set(rdb.Context(), fmt.Sprintf("%d", options.Message.Acceleration[0].Timestamp), data, 0).Err()
	if err != nil {
		options.OnErrorCallback(err)
		return
	}

	options.OnCompleteCallback()
}
