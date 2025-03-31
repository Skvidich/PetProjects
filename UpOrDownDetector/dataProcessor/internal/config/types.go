package config

type AppConfig struct {
	ErrLog string `ini:"ErrLog"`
	KafkaReporterConfig
	KafkaConsumerConfig
	ProcessEngineConfig
	PostgresStorageConfig
}

type KafkaReporterConfig struct {
	KafkaBrokers []string `ini:"KafkaBrokers"`
	StartTopic   string   `ini:"StartTopic"`
	EndTopic     string   `ini:"EndTopic"`
}
type KafkaConsumerConfig struct {
	KafkaBrokers []string `ini:"KafkaBrokers"`
	Topic        string   `ini:"Topic"`
}

type PostgresStorageConfig struct {
	Host     string `ini:"Host"`
	User     string `ini:"User"`
	Password string `ini:"Password"`
}

type ProcessEngineConfig struct {
	AggregationInterval string `ini:"AggregationInterval"`
	ReadTimeout         string `ini:"ReadTimeout"`
}
