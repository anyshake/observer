package setup_station

import "fmt"

func (t *SetupStationStartupImpl) Execute() error {
	for _, constraint := range t.StationConfigConstraints {
		if err := constraint.Init(t.ActionHandler); err != nil {
			return fmt.Errorf("failed to initialize station settings, namespace %s, key %s: %w", constraint.GetNamespace(), constraint.GetKey(), err)
		}
	}

	return nil
}
