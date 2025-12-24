package config

type GeneralConfig struct {
	Gateway GatewayConfig `yaml:"gateway"`
	Routes  []RouteConfig `yaml:"routes"`
}

type GatewayConfig struct {
	LogFilter []string `yaml:"log_filter"`
	Port      int      `yaml:"port"`
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

type JwtConfig struct {
	Audience string `yaml:"audience"`
	Issuer   string `yaml:"issuer"`

	Cookie string `yaml:"cookie"`
	Header string `yaml:"header"`
	Prefix string `yaml:"prefix"`

	MapHeaderClaims map[string]string `yaml:"claims"`
	RequiredClaims  []string          `yaml:"required_claims"`

	// either RS256 or HS256
	Algorithm string `yaml:"algorithm"`

	// HS256
	Secret string `yaml:"secret"`

	// RS256
	PublicKey string `yaml:"public_key"`
	JwksUrl   string `yaml:"jwks_url"`
}
