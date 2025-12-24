package gateway

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/GrongoTheGrog/goteway/internals/filter/logging"
	"github.com/GrongoTheGrog/goteway/internals/filter/rateLimiting"
)

type Gateway struct {
	routes      []*Route
	port        int
	FilterChain *filter.FilterChain
}

func NewGateway() *Gateway {
	return &Gateway{
		port:        0,
		routes:      make([]*Route, 0),
		FilterChain: &filter.FilterChain{EntryFilter: filter.NewEntryFilter()},
	}
}

func (gateway *Gateway) Start(port string) {

	mux := http.DefaultServeMux

	for _, route := range gateway.routes {
		route.Print()
	}

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		for _, route := range gateway.routes {
			match := route.Match(request)

			if match {
				gateway.FilterChain.CombineFilterChains(route.filterChain)
				gateway.FilterChain.Execute(writer, request, route.Endpoint)
				break
			}
		}
	})

	truePort := port
	if gateway.port == 0 {
		truePort = ":" + strconv.Itoa(gateway.port)
	}

	log.Printf("Starting gateway at port %s", truePort)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error while starting the gateway: %s", err.Error())
	}
}

func (gateway *Gateway) AddRoute(route *Route) *Gateway {
	gateway.routes = append(gateway.routes, route)
	return gateway
}

// AddFilter adds the filter at the end of the filter chain
// and before the proxy filter, if existing
func (gateway *Gateway) AddFilter(filter filter.Filter) *Gateway {
	gateway.FilterChain.AddFilter(filter)
	return gateway
}

// AddFilterBefore adds the filter before the beforeFilter.
// Panics if beforeFilter does not exist, or it is a filter.EntryFilter.
func (gateway *Gateway) AddFilterBefore(filter, beforeFilter filter.Filter) *Gateway {
	gateway.FilterChain.AddFilterBefore(filter, beforeFilter)
	return gateway
}

// AddFilterAfter adds the filter after the afterFilter.
// Panics if afterFilter does not exist, or it is a filter.ProxyFilter
func (gateway *Gateway) AddFilterAfter(filter, afterFilter filter.Filter) *Gateway {
	gateway.FilterChain.AddFilterAfter(filter, afterFilter)
	return gateway
}

func (gateway *Gateway) LogFilter(options ...logging.LogOption) *Gateway {
	gateway.FilterChain.AddFilter(logging.NewLogFilter(options...))
	return gateway
}

// TokenBucketFilter rate limits the request using the Token Bucket algorithm.
// See more at https://en.wikipedia.org/wiki/Token_bucket
func (gateway *Gateway) TokenBucketFilter(
	maxTokenNumber int,
	tokenCreationTime time.Duration,
	resource rateLimiting.ResourceLimiting,
) *Gateway {
	gateway.AddFilter(rateLimiting.NewTokenBucketFilter(maxTokenNumber, tokenCreationTime, resource))
	return gateway
}

func (gateway *Gateway) SlidingWindowCounterFilter(
	maxRequests int,
	windowTime time.Duration,
	limiting rateLimiting.ResourceLimiting,
) *Gateway {

	gateway.AddFilterAfter(
		rateLimiting.NewSlidingWindowCounterFilter(maxRequests, windowTime, limiting),
		gateway.FilterChain.EntryFilter)

	return gateway
}
