package gateway

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutePathMatch(t *testing.T) {
	route := NewRoute("random").PathPattern("/test")

	t.Run("valid path", func(t *testing.T) {
		request := &http.Request{RequestURI: "/test"}
		match := route.Match(request)

		assert.True(t, match)
	})

	t.Run("invalid path", func(t *testing.T) {
		request := &http.Request{RequestURI: "/nonvalid"}
		match := route.Match(request)

		assert.False(t, match)
	})
}

func TestRouteMethodMatch(t *testing.T) {
	route := NewRoute("random").Methods("GET")

	t.Run("valid method", func(t *testing.T) {
		request := &http.Request{Method: "GET"}
		match := route.Match(request)

		assert.True(t, match)
	})

	t.Run("invalid method", func(t *testing.T) {
		request := &http.Request{Method: "POST"}
		match := route.Match(request)

		assert.False(t, match)
	})
}

func TestRouteHostMatch(t *testing.T) {
	route := NewRoute("random").Hosts("host")

	t.Run("valid path", func(t *testing.T) {
		request := &http.Request{Host: "host"}
		match := route.Match(request)

		assert.True(t, match)
	})

	t.Run("invalid path", func(t *testing.T) {
		request := &http.Request{Host: "invalid"}
		match := route.Match(request)

		assert.False(t, match)
	})
}

func TestRouteHeaderMatch(t *testing.T) {
	route := NewRoute("random").Header("X-Header-Test", "test")

	t.Run("valid header", func(t *testing.T) {
		request := &http.Request{Header: http.Header{"X-Header-Test": {"test"}}}
		match := route.Match(request)

		assert.True(t, match)
	})

	t.Run("invalid path", func(t *testing.T) {
		request := &http.Request{Header: http.Header{"X-Header-Test": {"nonvalid"}}}
		match := route.Match(request)

		assert.False(t, match)
	})
}
