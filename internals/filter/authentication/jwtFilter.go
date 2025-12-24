package authentication

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"maps"
	"net/http"

	"github.com/GrongoTheGrog/goteway/internals/config"
	"github.com/GrongoTheGrog/goteway/internals/filter"
)

var byteDecodedSecret []byte

func NewJwtFilter(c config.JwtConfig) *filter.BasicFilter {

	applyDefaults(&c)
	if c.Secret != "" {
		key, _ := base64.StdEncoding.DecodeString(c.Secret)
		byteDecodedSecret = key
	}

	return filter.NewBasicFilter(func(context *filter.Context) *http.Response {

		tokenString, err := getToken(context, c)
		if err != nil {
			context.Log("Invalid token provided: %s", err.Error())
			return invalidToken(err.Error())
		}
		context.Log("Jwt token retrieved.")

		claims, err := decodeToken(tokenString, c)
		if err != nil {
			context.Log("Invalid token provided: %s", err.Error())
			return invalidToken(err.Error())
		}
		context.Log("Jwt token decoded.")

		if c.RequiredClaims != nil {
			for _, claim := range c.RequiredClaims {
				_, ok := claims[claim]
				if !ok {
					context.Log("Invalid token provided ")
					return invalidToken(fmt.Sprintf("Required claim '%v' not present in jwt.", claim))
				}
			}
		}

		for claim, header := range maps.All(c.MapHeaderClaims) {
			value, ok := claims[claim]
			if ok {
				castedValue, _ := value.(string)
				context.Request.Header.Add(header, castedValue)
			}
		}

		return context.RunNextFilter()
	})
}

func applyDefaults(c *config.JwtConfig) {
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
