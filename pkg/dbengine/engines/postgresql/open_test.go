package postgresql

import (
	"testing"
	"time"
)

func TestOpenRejectsInvalidAddress(t *testing.T) {
	if _, err := (&PostgreSQL{}).Open("not-a-host", "user", "secret", "observer", "t_", time.Second); err == nil {
		t.Fatal("address without a port was accepted")
	}
}
