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
