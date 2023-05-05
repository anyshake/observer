package archiver

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"com.geophone.observer/common/postgres"
)

func WriteMessage(pdb *sql.DB, options *ArchiverOptions) {
	if !options.Enable {
		return
	}

	if pdb == nil {
		options.OnErrorCallback(fmt.Errorf("database not connected"))
		return
	}

	acceleration, err := json.Marshal(options.Message.Acceleration)
	if err != nil {
		options.OnErrorCallback(err)
		return
	}

	err = postgres.InsertData(
		pdb, options.Message.Acceleration.Timestamp,
		options.Message.Station, options.Message.UUID, acceleration,
	)
	if err != nil {
		options.OnErrorCallback(err)
		return
	}

	options.OnCompleteCallback()
}
