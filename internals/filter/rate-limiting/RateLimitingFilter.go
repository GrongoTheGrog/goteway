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

		switch resource {
		case USER:
			allowRequest, err = storeInMap(ctx.RequestIp, maxTokenNumber, tokenCreationTime)
		case ROUTE:
			allowRequest, err = storeInMap(ctx.Url, maxTokenNumber, tokenCreationTime)
		case GATEWAY:
			allowRequest, err = storeInMap("gateway", maxTokenNumber, tokenCreationTime)
		default:
			err = fmt.Errorf("Unknown resource type passed in rate limiting filter.")
		}

		if err != nil {
			ctx.Log("Error in token bucket limiting filter: %s", err.Error())
			return utils.ErrorResponse("Error in Rate Limiting filter. Try again later.", 500)
		}

		if !allowRequest {
			ctx.Log("Too many requests, token bucket was 0. Denying request.")
			return utils.ErrorResponse("Too many requests, try again later.", 429)
		}

		return ctx.RunNextFilter()
	})
}

type Token struct {
	token       int
	lastRequest time.Time
	mu          sync.Mutex
}

func storeInMap(key string, maxTokenNumber int, refillInterval time.Duration) (bool, error) {
	value, _ := tokenMap.LoadOrStore(
		key,
		&Token{
			token:       maxTokenNumber,
			lastRequest: time.Now(),
		},
	)

	bucket := value.(*Token)

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	elapsed := time.Since(bucket.lastRequest)
	newTokens := int(elapsed / refillInterval)

	if newTokens > 0 {
		bucket.token = min(maxTokenNumber, bucket.token+newTokens)
		bucket.lastRequest = time.Now()
	}

	if bucket.token <= 0 {
		return false, nil
	}

	bucket.token--

	return true, nil
}
