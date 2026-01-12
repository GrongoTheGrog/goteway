package authentication

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"maps"
	"net/http"
	"regexp"

	"github.com/GrongoTheGrog/goteway/internals/filter"
)

var byteDecodedSecret []byte

func NewJwtFilter(c AuthorizationConfig) *filter.BasicFilter {

	applyDefaults(&c.Jwt)

	if c.Jwt.Secret != "" {
		key, _ := base64.StdEncoding.DecodeString(c.Jwt.Secret)
		byteDecodedSecret = key
	}

	allowedRoutesRegex := make([]*regexp.Regexp, 0)
	for _, route := range c.AllowedRoutes {
		allowedRoutesRegex = append(allowedRoutesRegex, regexp.MustCompile(route))
	}

	return filter.NewBasicFilter(func(context *filter.Context) *http.Response {

		for _, regex := range allowedRoutesRegex {
			if regex.Match([]byte(context.Request.URL.Path)) {
				return context.RunNextFilter()
			}
		}

		tokenString, err := getToken(context, c.Jwt)
		if err != nil {
			context.Log("Invalid token provided: %s", err.Error())
			return invalidToken(err.Error())
		}
		context.Log("Jwt token retrieved.")

		claims, err := decodeToken(tokenString, c.Jwt)
		if err != nil {
			context.Log("Invalid token provided: %s", err.Error())
			return invalidToken(err.Error())
		}
		context.Log("Jwt token decoded.")

		if c.Jwt.RequiredClaims != nil {
			for _, claim := range c.Jwt.RequiredClaims {
				_, ok := claims[claim]
				if !ok {
					context.Log("Invalid token provided ")
					return invalidToken(fmt.Sprintf("Required claim '%v' not present in jwt.", claim))
				}
			}
		}

		for claim, header := range maps.All(c.Jwt.MapHeaderClaims) {
			value, ok := claims[claim]
			if ok {
				castedValue, _ := value.(string)
				context.Request.Header.Add(header, castedValue)
			}
		}

		return context.RunNextFilter()
	})
}

func applyDefaults(c *JwtConfig) {
	if c.Header == "" {
		c.Header = "Authorization"
	}

	if c.Prefix == "" {
		c.Prefix = "Bearer "
	}
}

func invalidToken(message string) *http.Response {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewBuffer([]byte(message))),
		Status:     "401 UNAUTHORIZED",
		StatusCode: 401,
	}
}
