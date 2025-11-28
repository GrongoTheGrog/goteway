package rate_limiting

import (
	"net/http"
	"testing"
	"time"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/stretchr/testify/assert"
)

func TestIfRateLimitingFilterBlocksRequests(t *testing.T) {
	rateLimiterFilter := NewTokenBucketFilter(
		10,
		1*time.Second,
		USER,
	)

	var response *http.Response

	for i := 0; i < 11; i++ {
		context := &filter.Context{
			RequestIp: "111.111.111.111",
			Url:       "http://randomurl",
		}
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
		context := &filter.Context{
			RequestIp: "111.111.111.111",
			Url:       "http://randomurl",
		}
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
		context := &filter.Context{
			RequestIp: "111.111.111.111",
			Url:       "http://randomurl",
		}
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

	context := &filter.Context{
		RequestIp: "111.111.111.111",
		Url:       "http://randomurl",
	}

	for i := 0; i < 11; i++ {
		response = rateLimiterFilter.RunFilter(context)
	}

	assert.Equal(t, response.StatusCode, 429)

	time.Sleep(300 * time.Millisecond)

	response = rateLimiterFilter.RunFilter(context)

	// I'm testing against 500 because it's supposed to return that response,
	// since there is no filter set after the run filter.
	assert.Equal(t, response.StatusCode, 500)
}
