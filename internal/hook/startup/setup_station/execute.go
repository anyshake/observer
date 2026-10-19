package setup_station

import (
	"context"
	"fmt"
)

func (t *SetupStationStartupImpl) Execute(ctx context.Context) error {
	for _, constraint := range t.StationConfigConstraints {
		if err := constraint.Init(ctx, t.ActionHandler); err != nil {
			return fmt.Errorf("failed to initialize station settings, namespace %s, key %s: %w", constraint.GetNamespace(), constraint.GetKey(), err)
		}
	}

	return nil
}
