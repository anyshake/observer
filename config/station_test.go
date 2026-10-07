package config_test

import (
	"testing"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/testsupport"
)

func TestStationConstraintRoundTrip(t *testing.T) {
	_, handler := testsupport.OpenDAO(t)
	for _, constraint := range config.NewStationConstraints() {
		testsupport.ExerciseConstraint(t, handler, constraint)
	}
}
