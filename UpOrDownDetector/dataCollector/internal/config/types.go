package config

import "time"

type Config struct {
	ReqDelay     time.Duration `ini:"ReqDelay"`
	Getters      []string      `ini:"Getters,omitempty"`
	LogRelay     bool          `ini:"LogRelay"`
	ResendRelay  bool          `ini:"ResendRelay"`
	ErrLog       string        `ini:"ErrLog"`
	StatLog      string        `ini:"StatLog"`
	KafkaBrokers []string      `ini:"KafkaBrokers,omitempty"`
	KafkaTopic   string        `ini:"KafkaTopic"`
}
