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

	cfg := &Config{
		ErrLog:  "logs/error.log",
		StatLog: "logs/status.log",
	}
	section := iniFile.Section("")
	err = section.MapTo(cfg)
	if err != nil {
		return nil, fmt.Errorf("config mapping failed: %w", err)
	}

	return cfg, nil
}
