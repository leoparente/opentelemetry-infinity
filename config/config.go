package config

import "time"

type Status struct {
	StartTime time.Time     `json:"start_time"`
	UpTime    time.Duration `json:"up_time"`
	Version   string        `json:"version"`
}

type Policy struct {
	FeatureGates []string               `yaml:"feature_gates"`
	Set          map[string]string      `yaml:"set"`
	Config       map[string]interface{} `yaml:"config"`
}

type Config struct {
	Debug         bool   `mapstructure:"debug"`
	SelfTelemetry bool   `mapstructure:"self_telemetry"`
	ServerHost    string `mapstructure:"server_host"`
	ServerPort    uint64 `mapstructure:"server_port"`
}
