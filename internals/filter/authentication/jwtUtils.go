package authentication

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/GrongoTheGrog/goteway/internals/config"
	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/golang-jwt/jwt"
)

func getToken(ctx *filter.Context, c config.JwtConfig) (string, error) {
	if c.Cookie == "" {
		header := ctx.Request.Header.Get(c.Header)
		if header == "" {
			return "", fmt.Errorf("Header '%s' required for jwt token not found.", c.Header)
		}
		return header[len(c.Prefix):], nil
	} else {
		cookieJar := ctx.Request.CookiesNamed(c.Cookie)
		if len(cookieJar) == 0 {
			return "", fmt.Errorf("Cookie '%s' required for jwt token not found.", c.Cookie)
		}
		return cookieJar[0].Value, nil
	}
}

func decodeToken(rawToken string, c config.JwtConfig) (jwt.MapClaims, error) {

	claims := jwt.MapClaims{}

	switch c.Algorithm {
	case "HS256":
		_, err := jwt.ParseWithClaims(rawToken, claims, func(t *jwt.Token) (interface{}, error) {
			return byteDecodedSecret, nil
		})

		if err != nil {
			return nil, err
		}
	case "RS256":
		if c.PublicKey == "" && c.JwksUrl == "" {
			return nil, fmt.Errorf("Neither public key nor jwks url have been set.")
		}

		_, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
			if c.PublicKey != "" {
				return c.PublicKey, nil
			}

			rawKid, ok := t.Header["kid"]
			if !ok {
				return nil, fmt.Errorf("Missing Kid in JWT token.")
			}
			kid, _ := rawKid.(string)
			pub, err := getJwksPublicKey(kid, c.JwksUrl)

			if err != nil {
				return nil, err
			}

			return pub, nil
		})

		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Algorithm '%s' not supported.", c.Algorithm)
	}

	if !claims.VerifyAudience(c.Audience, true) ||
		!claims.VerifyIssuer(c.Issuer, true) {
		return nil, fmt.Errorf("Invalid token provided.")
	}

	return claims, nil
}

var keyCache map[string]*rsa.PublicKey = make(map[string]*rsa.PublicKey)

func getJwksPublicKey(kid, url string) (*rsa.PublicKey, error) {
	if pub, ok := keyCache[kid]; ok {
		return pub, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the keys from endpoint; %s", err.Error())
	}

	b := []byte{}
	_, err = res.Body.Read(b)
	keyRes := keyResponse{}
	json.Unmarshal(b, &keyRes)

	for _, key := range keyRes.Keys {
		if key.Kid == kid {
			return buildPub(key.N, key.E)
		}
	}

	return nil, fmt.Errorf("No public keys matched the Kid.")
}

func buildPub(n, e string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(n)
	eBytes, err := base64.RawURLEncoding.DecodeString(e)

	if err != nil {
		return nil, err
	}

	modulus := new(big.Int).SetBytes(nBytes)

	exponent := 0
	for _, b := range eBytes {
		exponent = exponent<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: modulus,
		E: exponent,
	}, nil
}

type keyResponse struct {
	Keys []keyInstance `json:"keys"`
}

type keyInstance struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}
