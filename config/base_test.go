package config

import (
	"os"
	"path/filepath"
	"testing"

	pkglogger "github.com/anyshake/observer/pkg/logger"
	"github.com/spf13/viper"
)

func TestBaseConfigParseAndMigrate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	contents := []byte(`{
		"location": {"latitude": 25.0, "longitude": 121.5, "elevation": 10},
		"hardware": {"endpoint": "tcp://127.0.0.1:1", "protocol": "v1", "timeout": 1},
		"ntpclient": {"endpoint": "pool.ntp.org", "timeout": 1, "retry": 1},
		"database": {"endpoint": "sqlite3://localhost", "database": "observer", "timeout": 1},
		"server": {"listen": "127.0.0.1:8080"},
		"logger": {"level": "info", "lifecycle": 1, "rotation": 1, "size": 1}
	}`)
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatal(err)
	}

	var cfg BaseConfig
	if err := cfg.Parse(path, "json"); err != nil {
		t.Fatal(err)
	}
	if cfg.Location.Latitude != 25 || cfg.Hardware.Protocol != "v1" || cfg.NtpClient.Endpoint != "pool.ntp.org" {
		t.Fatalf("parsed config = %+v", cfg)
	}
	pkglogger.Init()
	if err := cfg.Migrate(pkglogger.GetLogger("config")); err != nil {
		t.Fatal(err)
	}
	if cfg.NtpClient.Endpoint != "" || len(cfg.NtpClient.Pool) != 1 || cfg.NtpClient.Pool[0] != "pool.ntp.org" {
		t.Fatalf("migrated ntp client = %+v", cfg.NtpClient)
	}
	if err := cfg.Migrate(pkglogger.GetLogger("config")); err != nil {
		t.Fatal(err)
	}
	if len(cfg.NtpClient.Pool) != 1 {
		t.Fatalf("second migrate changed pool: %+v", cfg.NtpClient)
	}

	viper.Reset()
	if err := cfg.Parse(filepath.Join(dir, "missing.json"), "json"); err == nil {
		t.Fatal("missing config accepted")
	}
	invalidPath := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(invalidPath, []byte(`{"server": {}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	viper.Reset()
	if err := (&BaseConfig{}).Parse(invalidPath, "json"); err == nil {
		t.Fatal("invalid config accepted")
	}
}
