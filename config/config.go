package config

type Policy struct {
	FeatureGates []string               `mapstructure:"feature_gates"`
	Set          map[string]string      `mapstructure:"set"`
	Config       map[string]interface{} `mapstructure:"config"`
}

type OtlpInf struct {
	Debug      bool   `mapstructure:"debug"`
	ServerHost string `mapstructure:"server_host"`
	ServerPort uint64 `mapstructure:"server_port"`
}

type Config struct {
	Version float64 `mapstructure:"version"`
	OtlpInf OtlpInf `mapstructure:"otlp_inf"`
}
