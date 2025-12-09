package config

type GeneralConfig struct {
	Gateway GatewayConfig
	Routes  []RouteConfig
}

type GatewayConfig struct {
	LogFilter []string `yaml:"log_filter"`
}

type RouteConfig struct {
	Name         string
	Enabled      bool
	Endpoint     string
	Paths        []string
	RateLimiting RateLimitingConfig `yaml:"rate_limiting"`
	StripSuffix  int                `yaml:"strip_suffix"`
	Headers      []HeaderConfig     `yaml:"header_config"`
	Status       int
	ApiKey       string `yaml:"api_key"`
}

type RateLimitingConfig struct {
	Enabled         bool
	IntervalSeconds int `yaml:"interval_seconds"`
	MaxRequests     int `yaml:"max_requests"`
	Resource        string
}

type HeaderConfig struct {
	Name  string
	Value string
}
