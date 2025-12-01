package rateLimiting

import (
	"net/http"
	"sync"
	"time"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/GrongoTheGrog/goteway/internals/utils"
)

type WindowCounter struct {
	time     time.Time
	requests int
}

func newWindowCounter(time time.Time, requests int) *WindowCounter {
	return &WindowCounter{
		time:     time,
		requests: requests,
	}
}

type WindowFrame struct {
	prev *WindowCounter
	cur  *WindowCounter
}

var mu sync.Mutex

func (windowFrame *WindowFrame) addRequest(interval time.Duration) {

	elapsedTime := time.Since(windowFrame.cur.time)

	if elapsedTime <= interval {
		windowFrame.cur.requests++
	} else if elapsedTime > interval && elapsedTime <= elapsedTime+interval {
		windowFrame.prev = windowFrame.cur
		windowFrame.cur = newWindowCounter(windowFrame.prev.time.Add(interval), 1)
	} else {
		windowFrame.prev.requests = 0
		windowFrame.cur = newWindowCounter(time.Now().Truncate(interval), 0)
	}
}

func (windowFrame *WindowFrame) count(interval time.Duration) float64 {
	elapsedTime := time.Since(windowFrame.cur.time)
	ratio := float64(elapsedTime) / float64(interval)

	if ratio < 0 {
		ratio = 0
	}

	if ratio > 1 {
		ratio = 1
	}

	return (1-ratio)*float64(windowFrame.prev.requests) + float64(windowFrame.cur.requests)
}

var windowMap = make(map[any]*WindowFrame)

func NewSlidingWindowCounterFilter(
	maxRequests int,
	interval time.Duration,
	resource ResourceLimiting,
) filter.Filter {
	return filter.NewBasicFilter(func(ctx *filter.Context) *http.Response {

		mu.Lock()

		key := getKeyForResource(resource, ctx)

		windowFrame, ok := windowMap[key]

		if !ok {
			windowFrame = &WindowFrame{
				prev: newWindowCounter(time.Now().Truncate(interval).Add(-interval), 0),
				cur:  newWindowCounter(time.Now().Truncate(interval), 0),
			}
			windowMap[key] = windowFrame
		}

		windowFrame.addRequest(interval)
		requests := windowFrame.count(interval)

		mu.Unlock()

		if requests > float64(maxRequests) {
			ctx.Log("Too many requests. Returning 429 in response.")
			response := utils.ErrorResponse("Too many requests", 429)
			writeRateLimitingHeaders(response, maxRequests, 0)
			return response
		}

		ctx.Log("%v/%v requests remaining.", requests, maxRequests)

		response := ctx.RunNextFilter()
		writeRateLimitingHeaders(response, maxRequests, 0)
		return response
	})
}
