package dao

import (
	"fmt"
	"runtime"

	"gorm.io/gorm"
)

func Open(host string, port int, engineName, username, password, database string) (*gorm.DB, error) {
	engines := createEngines()
	databaseDriver, ok := engines[engineName]
	if !ok {
		var availableEngines []string
		for engine := range engines {
			availableEngines = append(availableEngines, engine)
		}
		return nil, fmt.Errorf("database engine %s is not supported on %s/%s, available engines: %v", engineName, runtime.GOOS, runtime.GOARCH, availableEngines)
	}

	return databaseDriver.Open(host, port, username, password, database, TIMEOUT_THRESHOLD)
}
