package main

import (
	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/pkg/logger"
)

func migrateConfig(cfg *config.BaseConfig) error {
	// 2025-08-22
	// Starting from v4.2.0, NTP Client configuration has deprecated the `endpoint` field.
	if len(cfg.NtpClient.Endpoint) > 0 && len(cfg.NtpClient.Pool) == 0 {
		cfg.NtpClient.Pool = []string{cfg.NtpClient.Endpoint}
		cfg.NtpClient.Endpoint = ""
		logger.GetLogger(main).Warnln("configuration field `ntpclient.endpoint` has been deprecated, please use the `ntpclient.pool` field instead")
	}

	return nil
}
