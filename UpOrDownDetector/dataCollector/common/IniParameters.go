package common

import (
	"fmt"
	"gopkg.in/ini.v1"

	"time"
)

type Config struct {
	StatusRequestDelay time.Duration `ini:"StatusRequestDelay"`
	GetterNames        []string      `ini:"GetterNames,omitempty"`
	RelayIsLog         bool          `ini:"RelayIsLog"`
	RelayIsResend      bool          `ini:"RelayIsResend"`
	ErrorLogPath       string        `ini:"ErrorLogPath"`
	StatusLogPath      string        `ini:"StatusLogPath"`
	KafkaAddrs         []string      `ini:"KafkaAddrs,omitempty"`
	KafkaTopic         string        `ini:"KafkaTopic"`
}

func GetIni(path string) (*Config, error) {
	cfg, err := ini.Load(path)
	if err != nil {

		return nil, fmt.Errorf("can't load ini file %v", err)
	}

	config := &Config{
		ErrorLogPath:  "logs/error.log",
		StatusLogPath: "logs/status.log",
	}
	section := cfg.Section("")
	err = section.MapTo(config)
	if err != nil {
		return nil, fmt.Errorf("config mapping failed: %w", err)
	}

	return config, nil
}
