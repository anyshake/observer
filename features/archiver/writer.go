package archiver

import (
	"database/sql"
	"encoding/json"

	"com.geophone.observer/common/postgres"
)

func WriteMessage(pdb *sql.DB, options *ArchiverOptions) {
	if !options.Enable {
		return
	}

	acceleration, err := json.Marshal(options.Message.Acceleration)
	if err != nil {
		options.OnErrorCallback(err)
		return
	}

	postgres.InsertData(
		pdb, options.Message.Acceleration[0].Timestamp,
		options.Message.Station, options.Message.UUID, acceleration,
	)

	options.OnCompleteCallback()
}
