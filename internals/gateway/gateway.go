package gateway

import (
	"log"
	"net/http"
	"regexp"

	"github.com/GrongoTheGrog/goteway/internals/filter"
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

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Error while starting the gateway: %s", err.Error())
	}
}

func (gateway *Gateway) NewRoute(pattern, endpoint string) *Route {
	route := NewRoute(pattern, endpoint)
	gateway.routes = append(gateway.routes, route)
	return route
}
