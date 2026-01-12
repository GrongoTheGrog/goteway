package authentication

type AuthorizationConfig struct {
	Jwt           JwtConfig `yaml:"jwt"`
	AllowedRoutes []string  `yaml:"allowed_routes"`
}

type JwtConfig struct {
	Enabled  bool `yaml:"enabled"`
	Required bool `yaml:"required"`

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
