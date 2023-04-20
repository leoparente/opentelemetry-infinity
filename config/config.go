package config

import "time"

type Status struct {
	StartTime      time.Time     `json:"start_time"`
	UpTime         time.Duration `json:"up_time"`
	InfVersion     string        `json:"otlpinf_version"`
	ContribVersion string        `json:"otel_contrib_version"`
}

type Policy struct {
	FeatureGates []string               `yaml:"feature_gates"`
	Set          map[string]string      `yaml:"set"`
	Config       map[string]interface{} `yaml:"config"`
}

type OtlpInf struct {
	Debug         bool   `mapstructure:"debug"`
	SelfTelemetry bool   `mapstructure:"self_telemetry"`
	ServerHost    string `mapstructure:"server_host"`
	ServerPort    uint64 `mapstructure:"server_port"`
}

type Config struct {
	Version string  `mapstructure:"version"`
	OtlpInf OtlpInf `mapstructure:"otlp_inf"`
}
