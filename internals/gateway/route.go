package gateway

import (
	"github.com/GrongoTheGrog/goteway/internals/filter"
)

type Route struct {
	Endpoint    string
	Match       string
	filterChain *filter.FilterChain
}

func NewRoute(match, endpoint string) *Route {
	return &Route{
		Match:       match,
		Endpoint:    endpoint,
		filterChain: &filter.FilterChain{ProxyFilter: filter.NewProxyFilter()},
	}
}

func (route *Route) Filter(filter filter.Filter) *Route {
	route.filterChain.AddFilter(filter)
	return route
}

func (route *Route) RemoveLeftPath(pathNum int) *Route {
	route.filterChain.AddFilter(filter.NewRemoveLeftPathFilter(pathNum))
	return route
}

func (route *Route) RemoveRightPath(pathNum int) *Route {
	route.filterChain.AddFilter(filter.NewRemoveRightPathFilter(pathNum))
	return route
}

//TODO Rewrite path
//TODO Logging filter
//TODO Jwt Filter
//TODO Api Key filter
//TODO
