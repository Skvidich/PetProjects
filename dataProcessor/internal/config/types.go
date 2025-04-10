package config

import "time"

type Config struct {
	AppConfig
	KafkaReporterConfig
	KafkaConsumerConfig
	ProcessEngineConfig
	PostgresStorageConfig
}

type KafkaReporterConfig struct {
	Brokers    []string `ini:"Brokers"`
	StartTopic string   `ini:"StartTopic"`
	EndTopic   string   `ini:"EndTopic"`
}
type KafkaConsumerConfig struct {
	Brokers []string `ini:"Brokers"`
	Topic   string   `ini:"Topic"`
}

type PostgresStorageConfig struct {
	DSN            string `ini:"DSN"`
	ReportTable    string `ini:"ReportTable"`
	ComponentTable string `ini:"ComponentTable"`
	IncidentTable  string `ini:"IncidentTable"`
}

type ProcessEngineConfig struct {
	AggregationInterval time.Duration `ini:"AggregationInterval"`
	ReadTimeout         time.Duration `ini:"ReadTimeout"`
}

type AppConfig struct {
	ErrLog   string `ini:"ErrLog"`
	RetryMax int    `ini:"RetryMax"`
}
