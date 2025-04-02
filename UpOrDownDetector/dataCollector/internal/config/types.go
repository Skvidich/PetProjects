package config

import "time"

type Config struct {
	CoordinatorConfig
	RelayConfig
	LoggerConfig
	KafkaProducerConfig
	StorageConfig
}
type CoordinatorConfig struct {
	ReqDelay time.Duration `ini:"ReqDelay"`
	Getters  []string      `ini:"Getters"`
}

type RelayConfig struct {
	Save   bool `ini:"Save"`
	Resend bool `ini:"Resend"`
}

type LoggerConfig struct {
	ErrLog string `ini:"ErrLog"`
}

type KafkaProducerConfig struct {
	Brokers []string `ini:"Brokers"`
	Topic   string   `ini:"Topic"`
}

type StorageConfig struct {
	DSN            string `ini:"DSN"`
	ReportTable    string `ini:"ReportTable"`
	ComponentTable string `ini:"ComponentTable"`
}
