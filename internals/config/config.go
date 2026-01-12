package config

import "github.com/GrongoTheGrog/goteway/internals/filter/authentication"

type GeneralConfig struct {
	Gateway GatewayConfig `yaml:"gateway"`
	Routes  []RouteConfig `yaml:"routes"`
}

type GatewayConfig struct {
	Auth      authentication.AuthorizationConfig `yaml:"auth"`
	LogFilter []string                           `yaml:"log_filter"`
	Port      int                                `yaml:"port"`
}

type RouteConfig struct {
	Name         string             `yaml:"name"`
	Enabled      bool               `yaml:"enabled"`
	Endpoint     string             `yaml:"endpoint"`
	Paths        []string           `yaml:"paths"`
	RateLimiting RateLimitingConfig `yaml:"rate_limiting"`
	StripSuffix  int                `yaml:"strip_suffix"`
	Headers      []HeaderConfig     `yaml:"header_config"`
	Status       int                `yaml:"status"`
	ApiKey       string             `yaml:"api_key"`
}

type RateLimitingConfig struct {
	Enabled         bool
	IntervalSeconds int    `yaml:"interval_seconds"`
	MaxRequests     int    `yaml:"max_requests"`
	Resource        string `yaml:"resource"`
}

type HeaderConfig struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
