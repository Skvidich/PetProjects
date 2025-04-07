package config

import (
	"fmt"
	"gopkg.in/ini.v1"
)

func LoadConfig(path string) (*Config, error) {
	iniFile, err := ini.Load(path)
	if err != nil {

		return nil, fmt.Errorf("can't load ini file %v", err)
	}

	cfg := &Config{}
	section := iniFile.Section("Logger")
	err = section.MapTo(&cfg.LoggerConfig)
	if err != nil {
		return nil, fmt.Errorf("config logger mapping failed: %w", err)
	}

	section = iniFile.Section("Producer")
	err = section.MapTo(&cfg.KafkaProducerConfig)
	if err != nil {
		return nil, fmt.Errorf("config producer mapping failed: %w", err)
	}

	section = iniFile.Section("Coordinator")
	err = section.MapTo(&cfg.CoordinatorConfig)
	if err != nil {
		return nil, fmt.Errorf("config coordinator mapping failed: %w", err)
	}

	section = iniFile.Section("Storage")
	err = section.MapTo(&cfg.StorageConfig)
	if err != nil {
		return nil, fmt.Errorf("config storage mapping failed: %w", err)
	}

	section = iniFile.Section("Relay")
	err = section.MapTo(&cfg.RelayConfig)
	if err != nil {
		return nil, fmt.Errorf("config relay mapping failed: %w", err)
	}

	section = iniFile.Section("Server")
	err = section.MapTo(&cfg.ServerConfig)
	if err != nil {
		return nil, fmt.Errorf("config server mapping failed: %w", err)
	}

	return cfg, nil
}
