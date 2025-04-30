package dao

import (
	"fmt"
	"net/url"
	"runtime"
	"time"

	"github.com/anyshake/observer/pkg/dbengine"
)

func New(endpoint, username, password, prefix string, timeout time.Duration) (*DAO, error) {
	urlObj, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database endpoint: %w", err)
	}

	engines := dbengine.New()
	drvier, ok := engines[urlObj.Scheme]
	if !ok {
		var availableEngines []string
		for engine := range engines {
			availableEngines = append(availableEngines, engine)
		}
		return nil, fmt.Errorf("database engine %s is not supported on %s/%s, available engines: %v", urlObj.Scheme, runtime.GOOS, runtime.GOARCH, availableEngines)
	}

	return &DAO{
		address:  urlObj.Host,
		timeout:  timeout,
		username: username,
		password: password,
		prefix:   prefix,
		driver:   drvier,
	}, nil
}
