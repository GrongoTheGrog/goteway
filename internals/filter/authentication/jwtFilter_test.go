package authentication

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

var context = &filter.Context{
	RequestIp: "111.111.111.111",
	Url:       "http://randomurl",
	Request: &http.Request{
		Header: make(http.Header),
	},
}

var testConfig AuthorizationConfig = AuthorizationConfig{
	Jwt: JwtConfig{
		Algorithm: "HS256",
		Secret:    "E84mnV6OTl7oduTDXxvOPTWS5NYWmrRpQK+xLo3X5Jo=",
		Audience:  "frontend",
		Issuer:    "Backend",
	},
}

func TestIfDefaultConfigWorks(t *testing.T) {
	filter := NewJwtFilter(testConfig)

	context.Request.Header.Add("Authorization", "Bearer "+generateHS256Jwt())

	res := filter.RunFilter(context)

	assert.Equal(t, 500, res.StatusCode)
}

func TestIfInvalidTokenFails(t *testing.T) {
	filter := NewJwtFilter(testConfig)

	context.Request.Header.Set("Authorization", "Bearer super invalid token")

	res := filter.RunFilter(context)

	assert.Equal(t, 401, res.StatusCode)
}

func TestIfWrongAudienceFails(t *testing.T) {
	config := testConfig
	config.Jwt.Audience = "something else"
	filter := NewJwtFilter(config)

	context.Request.Header.Add("Authorization", "Bearer "+generateHS256Jwt())

	res := filter.RunFilter(context)

	assert.Equal(t, 401, res.StatusCode)
}

func TestIfWrongIssuerFails(t *testing.T) {
	config := testConfig
	config.Jwt.Issuer = "something else"
	filter := NewJwtFilter(config)
	context.Request.Header.Add("Authorization", "Bearer "+generateHS256Jwt())

	res := filter.RunFilter(context)

	assert.Equal(t, 401, res.StatusCode)
}

func TestIfRequiredClaimsAreRequired(t *testing.T) {
	config := testConfig
	config.Jwt.RequiredClaims = []string{"missing required claim"}
	filter := NewJwtFilter(config)
	context.Request.Header.Add("Authorization", "Bearer "+generateHS256Jwt())

	res := filter.RunFilter(context)

	assert.Equal(t, res.StatusCode, 401)
}

func TestIfMappedClaimsAreForwarded(t *testing.T) {
	claims := make(map[string]string)
	claims["aud"] = "X-Audience"
	claims["sub"] = "X-User-Id"

	config := testConfig
	config.Jwt.MapHeaderClaims = claims

	jwtFilter := NewJwtFilter(config)

	// setting a filter to run next so i can check
	// if claims were forwarded
	testFilter := filter.NewBasicFilter(func(ctx *filter.Context) *http.Response {
		aud := ctx.Request.Header.Get("X-Audience")
		sub := ctx.Request.Header.Get("X-User-Id")

		assert.Equal(t, aud, globalClaims.Audience)
		assert.Equal(t, sub, globalClaims.Subject)
		return ctx.RunNextFilter()
	})

	context.Request.Header.Set("Authorization", "Bearer "+generateHS256Jwt())

	jwtFilter.SetNext(testFilter)
	res := jwtFilter.RunFilter(context)

	assert.Equal(t, res.StatusCode, 500)
}

func TestIfAllowedRoutesDoNotNeedJwt(t *testing.T) {
	config := testConfig
	config.AllowedRoutes = []string{"/route/*"}

	ctx := context
	ctx.Request.URL = &url.URL{}
	ctx.Request.URL.Path = "/route/allowed"

	jwtFilter := NewJwtFilter(config)

	res := jwtFilter.RunFilter(ctx)

	assert.Equal(t, res.StatusCode, 500)
}

var globalClaims jwt.StandardClaims = jwt.StandardClaims{
	Audience: "frontend",
	Issuer:   "Backend",
	Subject:  "userId123",
}

func generateHS256Jwt() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, globalClaims)
	b, err := base64.StdEncoding.DecodeString(testConfig.Jwt.Secret)

	jwt, err := token.SignedString(b)
	if err != nil {
		panic(err)
	}

	return jwt
}
