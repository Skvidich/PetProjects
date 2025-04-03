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
	section := iniFile.Section("")
	err = section.MapTo(&cfg.AppConfig)
	if err != nil {
		return nil, fmt.Errorf("config mapping base config: %w", err)
	}

	section = iniFile.Section("KafkaConsumer")
	err = section.MapTo(&cfg.KafkaConsumerConfig)
	if err != nil {
		return nil, fmt.Errorf("config mapping kafka consumer: %w", err)
	}

	section = iniFile.Section("KafkaReporter")
	err = section.MapTo(&cfg.KafkaReporterConfig)
	if err != nil {
		return nil, fmt.Errorf("config mapping kafka reporter: %w", err)
	}

	section = iniFile.Section("ProcessEngine")
	err = section.MapTo(&cfg.ProcessEngineConfig)
	if err != nil {
		return nil, fmt.Errorf("config mapping process engine: %w", err)
	}

	section = iniFile.Section("PostgresStorage")
	err = section.MapTo(&cfg.PostgresStorageConfig)
	if err != nil {
		return nil, fmt.Errorf("config mapping postgress storage: %w", err)
	}
	return cfg, nil
}
