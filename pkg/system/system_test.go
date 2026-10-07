package system

import "testing"

func TestHostMetrics(t *testing.T) {
	if _, err := GetCpuModel(); err != nil {
		t.Fatal(err)
	}
	if _, err := GetCpuPercent(); err != nil {
		t.Fatal(err)
	}
	if _, err := GetMemoryPercent(); err != nil {
		t.Fatal(err)
	}
	if percent, err := GetDiskPercent(); err != nil || percent < 0 {
		t.Fatalf("disk percent = %v, %v", percent, err)
	}
	if uptime, err := GetOsUptime(); err != nil || uptime <= 0 {
		t.Fatalf("uptime = %d, %v", uptime, err)
	}
}
