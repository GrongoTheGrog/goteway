package rate_limiting

import (
	"net/http"
	"testing"
	"time"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/stretchr/testify/assert"
)

var context = &filter.Context{
	RequestIp: "111.111.111.111",
	Url:       "http://randomurl",
}

func TestIfRateLimitingFilterBlocksRequests(t *testing.T) {
	rateLimiterFilter := NewTokenBucketFilter(
		10,
		1*time.Second,
		USER,
	)

	var response *http.Response

	for i := 0; i < 11; i++ {
		response = rateLimiterFilter.RunFilter(context)
	}

	assert.Equal(t, response.StatusCode, 429)
}

func TestIfRateLimitingFilterWorksWithRoutes(t *testing.T) {
	rateLimiterFilter := NewTokenBucketFilter(
		10,
		1*time.Second,
		ROUTE,
	)

	var response *http.Response

	for i := 0; i < 11; i++ {
		response = rateLimiterFilter.RunFilter(context)
	}

	assert.Equal(t, response.StatusCode, 429)
}

func TestIfRateLimitingFilterWorksWithWholeGateway(t *testing.T) {
	rateLimiterFilter := NewTokenBucketFilter(
		10,
		1*time.Second,
		GATEWAY,
	)

	var response *http.Response

	for i := 0; i < 11; i++ {
		response = rateLimiterFilter.RunFilter(context)
	}

	assert.Equal(t, response.StatusCode, 429)
}

func TestIfRateLimitingFilterCanRegenerateTokens(t *testing.T) {
	rateLimiterFilter := NewTokenBucketFilter(
		10,
		100*time.Millisecond,
		GATEWAY,
	)

	var response *http.Response

	for i := 0; i < 11; i++ {
		response = rateLimiterFilter.RunFilter(context)
	}

	assert.Equal(t, response.StatusCode, 429)

	time.Sleep(300 * time.Millisecond)

	response = rateLimiterFilter.RunFilter(context)

	assert.Equal(t, response.StatusCode, 500)
}
