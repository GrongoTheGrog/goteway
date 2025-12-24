package authentication

import (
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/GrongoTheGrog/goteway/internals/config"
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

var testConfig config.JwtConfig = config.JwtConfig{
	Algorithm: "HS256",
	Secret:    "E84mnV6OTl7oduTDXxvOPTWS5NYWmrRpQK+xLo3X5Jo=",
	Audience:  "frontend",
	Issuer:    "Backend",
}

func TestIfDefaultConfigWorks(t *testing.T) {
	filter := NewJwtFilter(testConfig)

	context.Request.Header.Add("Authorization", "Bearer "+generateHS256Jwt())

	res := filter.RunFilter(context)

	assert.Equal(t, res.StatusCode, 500)
}

func TestIfInvalidTokenFails(t *testing.T) {
	filter := NewJwtFilter(testConfig)

	context.Request.Header.Add("Authorization", "Bearer "+"super invalid token")

	res := filter.RunFilter(context)

	assert.Equal(t, res.StatusCode, 401)
}

func TestIfWrongAudienceFails(t *testing.T) {
	testConfig.Audience = "something else"
	filter := NewJwtFilter(testConfig)

	context.Request.Header.Add("Authorization", "Bearer "+generateHS256Jwt())

	res := filter.RunFilter(context)

	assert.Equal(t, res.StatusCode, 401)
}

func TestIfWrongIssuerFails(t *testing.T) {
	testConfig.Issuer = "something else"
	filter := NewJwtFilter(testConfig)
	context.Request.Header.Add("Authorization", "Bearer "+generateHS256Jwt())

	res := filter.RunFilter(context)

	assert.Equal(t, res.StatusCode, 401)
}

func TestIfRequiredClaimsAreRequired(t *testing.T) {
	testConfig.RequiredClaims = []string{"missing required claim"}
	filter := NewJwtFilter(testConfig)
	context.Request.Header.Add("Authorization", "Bearer "+generateHS256Jwt())

	res := filter.RunFilter(context)

	assert.Equal(t, res.StatusCode, 401)
}

func TestIfMappedClaimsAreForwarded(t *testing.T) {
	claims := make(map[string]string)
	claims["aud"] = "X-Audience"
	claims["sub"] = "X-User-Id"
	testConfig.MapHeaderClaims = claims

	jwtFilter := NewJwtFilter(testConfig)

	// setting a filter to run next so i can check
	// if claims were forwarded
	testFilter := filter.NewBasicFilter(func(ctx *filter.Context) *http.Response {
		aud := ctx.Request.Header.Get("X-Audience")
		sub := ctx.Request.Header.Get("X-User-Id")

		assert.Equal(t, aud, globalClaims.Audience)
		assert.Equal(t, sub, globalClaims.Subject)
		return ctx.RunNextFilter()
	})

	context.Request.Header.Add("Authorization", "Bearer "+generateHS256Jwt())

	jwtFilter.SetNext(testFilter)
	res := jwtFilter.RunFilter(context)

	assert.Equal(t, res.StatusCode, 500)
}

var globalClaims jwt.StandardClaims = jwt.StandardClaims{
	Audience: "frontend",
	Issuer:   "Backend",
	Subject:  "userId123",
}

func generateHS256Jwt() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, globalClaims)
	b, err := base64.StdEncoding.DecodeString(testConfig.Secret)

	jwt, err := token.SignedString(b)
	if err != nil {
		panic(err)
	}

	return jwt
}
