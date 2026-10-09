package dao

import "testing"

func TestOpenWithoutDriver(t *testing.T) {
	t.Parallel()

	if err := (&DAO{}).Open("observer.db"); err == nil {
		t.Fatal("open without a driver succeeded")
	}
}
