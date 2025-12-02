package gateway

import (
	"fmt"
	"maps"
	"net/http"
	"regexp"
	"slices"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/GrongoTheGrog/goteway/internals/filter/request"
)

type Route struct {
	Endpoint    string
	PathMatch   *regexp.Regexp
	HeaderMatch map[string]string
	MethodMatch []string
	HostMatch   []string
	filterChain *filter.FilterChain
}

func NewRoute(endpoint string) *Route {
	return &Route{
		Endpoint:    endpoint,
		filterChain: &filter.FilterChain{ProxyFilter: filter.NewProxyFilter()},
	}
}

func (route *Route) PathPattern(regex string) *Route {
	route.PathMatch = regexp.MustCompile(regex)
	return route
}

func (route *Route) Headers(headers map[string]string) *Route {
	route.HeaderMatch = headers
	return route
}

func (route *Route) Methods(methods ...string) *Route {
	route.MethodMatch = methods
	return route
}

func (route *Route) Hosts(hosts ...string) *Route {
	route.HostMatch = hosts
	return route
}

func (route *Route) Match(request *http.Request) bool {

	if len(route.HostMatch) > 0 && !slices.Contains(route.HostMatch, request.Host) {
		return false
	}

	if len(route.MethodMatch) > 0 && !slices.Contains(route.MethodMatch, request.Method) {
		return false
	}

	if route.PathMatch != nil {
		ok := route.PathMatch.MatchString(request.RequestURI)
		if !ok {
			return false
		}
	}

	if route.HeaderMatch != nil {
		for key, value := range maps.All(route.HeaderMatch) {
			requestHeader := request.Header.Get(key)
			if requestHeader != value {
				return false
			}
		}
	}
	return true
}

func (route *Route) Print() {
	fmt.Printf("Route added: %s\n", route.Endpoint)

	if route.PathMatch != nil {
		fmt.Println("\tPath Match: ", route.PathMatch.String())
	}

	if len(route.HostMatch) > 0 {
		fmt.Println("\tHosts: ", route.HostMatch)
	}

	if len(route.MethodMatch) > 0 {
		fmt.Println("\tMethods: ", route.MethodMatch)
	}

	if route.HeaderMatch != nil {
		fmt.Println("\tHeaders: ", route.HeaderMatch)
	}
}

func (route *Route) Filter(filter filter.Filter) *Route {
	route.filterChain.AddFilter(filter)
	return route
}

func (route *Route) RemoveLeftPath(pathNum int) *Route {
	route.filterChain.AddFilter(request.NewRemoveLeftPathFilter(pathNum))
	return route
}

func (route *Route) RemoveRightPath(pathNum int) *Route {
	route.filterChain.AddFilter(request.NewRemoveRightPathFilter(pathNum))
	return route
}
