package rate_limiting

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/GrongoTheGrog/goteway/internals/utils"
)

type ResourceLimiting int

const (
	USER = iota
	ROUTE
	GATEWAY
)

var tokenMap sync.Map

func NewTokenBucketFilter(
	maxTokenNumber int,
	tokenCreationTime time.Duration,
	resource ResourceLimiting,
) filter.Filter {
	return filter.NewBasicFilter(func(ctx *filter.Context) *http.Response {

		var allowRequest bool
		var err error
		var remaining int

		switch resource {
		case USER:
			allowRequest, remaining, err = storeInMap(ctx.RequestIp, maxTokenNumber, tokenCreationTime)
		case ROUTE:
			allowRequest, remaining, err = storeInMap(ctx.Url, maxTokenNumber, tokenCreationTime)
		case GATEWAY:
			allowRequest, remaining, err = storeInMap("gateway", maxTokenNumber, tokenCreationTime)
		default:
			err = fmt.Errorf("Unknown resource type passed in rate limiting filter.")
		}

		if err != nil {
			ctx.Log("Error in token bucket limiting filter: %s", err.Error())
			return utils.ErrorResponse("Error in Rate Limiting filter. Try again later.", 500)
		}

		if !allowRequest {
			ctx.Log("Too many requests, token bucket was 0. Denying request.")
			response := utils.ErrorResponse("Too many requests, try again later.", 429)

			writeRateLimitingHeaders(response, maxTokenNumber, remaining)
			return response
		}

		ctx.Log("%v/%v requests remaining.", remaining, maxTokenNumber)

		response := ctx.RunNextFilter()

		writeRateLimitingHeaders(response, maxTokenNumber, remaining)
		return response
	})
}

type Bucket struct {
	token       int
	lastRequest time.Time
	mu          sync.Mutex
}

func storeInMap(key string, maxTokenNumber int, refillInterval time.Duration) (bool, int, error) {
	value, _ := tokenMap.LoadOrStore(
		key,
		&Bucket{
			token:       maxTokenNumber,
			lastRequest: time.Now(),
		},
	)

	bucket := value.(*Bucket)

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	elapsed := time.Since(bucket.lastRequest)
	newTokens := int(elapsed / refillInterval)

	if newTokens > 0 {
		bucket.token = min(maxTokenNumber, bucket.token+newTokens)
		bucket.lastRequest = time.Now()
	}

	if bucket.token <= 0 {
		return false, 0, nil
	}

	bucket.token--

	return true, bucket.token, nil
}
