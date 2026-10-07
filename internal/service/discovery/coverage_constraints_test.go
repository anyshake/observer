package mdns_discovery

import (
	"testing"

	"github.com/anyshake/observer/internal/testsupport"
)

func TestConfigConstraintRoundTrip(t *testing.T) {
	_, handler := testsupport.OpenDAO(t)
	for _, constraint := range New("127.0.0.1:8080", nil, nil, nil).GetConfigConstraint() {
		testsupport.ExerciseConstraint(t, handler, constraint)
	}
}
