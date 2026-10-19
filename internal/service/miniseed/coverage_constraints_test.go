package miniseed

import (
	"testing"

	"github.com/anyshake/observer/internal/testsupport"
)

func TestConfigConstraintRoundTrip(t *testing.T) {
	_, handler := testsupport.OpenDAO(t)
	for _, constraint := range New(t.Context(), nil, nil, nil).GetConfigConstraint() {
		testsupport.ExerciseConstraint(t, handler, constraint)
	}
}
