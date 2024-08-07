package transport

import "fmt"

func New(dsn *TransportDependency) (TransportDriver, error) {
	engines := map[string]TransportDriver{
		"serial": &TransportDriverSerialImpl{},
		"tcp":    &TransportDriverTcpImpl{},
	}
	engine, ok := engines[dsn.Engine]
	if !ok {
		return nil, fmt.Errorf("engine %s is not supported", dsn.Engine)
	}

	return engine, nil
}
