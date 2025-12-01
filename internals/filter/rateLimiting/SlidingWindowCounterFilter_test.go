package rateLimiting

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIfSlidingWindowCounterFilterAllowsRequests(t *testing.T) {
	filter := NewSlidingWindowCounterFilter(
		3,
		time.Second,
		USER,
	)

	response := filter.RunFilter(context)
	assert.Equal(t, response.StatusCode, 500)
}

func TestIfSlidingWindowCounterFilterBlocksRequests(t *testing.T) {
	filter := NewSlidingWindowCounterFilter(
		2,
		time.Minute,
		USER,
	)

	filter.RunFilter(context)
	filter.RunFilter(context)
	response := filter.RunFilter(context)

	assert.Equal(t, response.StatusCode, 429)
}

func TestIfSlidingWindowCounterFilterBlocksConcurrentRequests(t *testing.T) {
	for i := 0; i < 20; i++ {
		t.Run(fmt.Sprintf("Test Case: %v", i), func(t *testing.T) {
			filter := NewSlidingWindowCounterFilter(
				2,
				time.Minute,
				USER,
			)

			wg := sync.WaitGroup{}
			wg.Add(2)

			go func() {
				filter.RunFilter(context)
				wg.Done()
			}()

			go func() {
				filter.RunFilter(context)
				wg.Done()
			}()

			wg.Wait()

			response := filter.RunFilter(context)
			assert.Equal(t, 429, response.StatusCode)
		})
	}

}
