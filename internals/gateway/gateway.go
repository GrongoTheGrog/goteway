package gateway

import (
	"log"
	"net/http"
	"regexp"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/GrongoTheGrog/goteway/internals/filter/logging"
)

type Gateway struct {
	routes      []*Route
	FilterChain *filter.FilterChain
}

func NewGateway() *Gateway {
	return &Gateway{
		routes:      make([]*Route, 0),
		FilterChain: &filter.FilterChain{EntryFilter: filter.NewEntryFilter()},
	}
}

func (gateway *Gateway) Start(port string) {

	log.Printf("Starting gateway at port %s", port)

	mux := http.DefaultServeMux

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		path := request.RequestURI

		for _, route := range gateway.routes {
			match, err := regexp.MatchString(route.Match, path)
			if err != nil {
				log.Printf("Error comparing regex: " + err.Error())
			}

			if match {
				gateway.FilterChain.CombineFilterChains(route.filterChain)
				gateway.FilterChain.Execute(writer, request, route.Endpoint)
				break
			}
		}
	})

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error while starting the gateway: %s", err.Error())
	}
}

func (gateway *Gateway) NewRoute(pattern, endpoint string) *Route {
	route := NewRoute(pattern, endpoint)
	gateway.routes = append(gateway.routes, route)
	return route
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
