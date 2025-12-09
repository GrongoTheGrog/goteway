package gateway

import (
	"fmt"
	"maps"
	"net/http"
	"regexp"
	"slices"
	"time"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/GrongoTheGrog/goteway/internals/filter/rateLimiting"
	"github.com/GrongoTheGrog/goteway/internals/filter/request"
)

type Route struct {
	Endpoint    string
	pathMatches []*regexp.Regexp
	headerMatch map[string]string
	methodMatch []string
	hostMatch   []string
	filterChain *filter.FilterChain
}

func NewRoute(endpoint string) *Route {
	return &Route{
		pathMatches: make([]*regexp.Regexp, 0),
		Endpoint:    endpoint,
		headerMatch: make(map[string]string),
		filterChain: &filter.FilterChain{ProxyFilter: filter.NewProxyFilter()},
	}
}

func (route *Route) PathPattern(regexes ...string) *Route {
	for _, regex := range regexes {
		route.pathMatches = append(route.pathMatches, regexp.MustCompile(regex))
	}
	return route
}

func (route *Route) Header(key, value string) *Route {
	route.headerMatch[key] = value
	return route
}

func (route *Route) Methods(methods ...string) *Route {
	route.methodMatch = methods
	return route
}

func (route *Route) Hosts(hosts ...string) *Route {
	route.hostMatch = hosts
	return route
}

func (route *Route) Match(request *http.Request) bool {

	if len(route.hostMatch) > 0 && !slices.Contains(route.hostMatch, request.Host) {
		return false
	}

	if len(route.methodMatch) > 0 && !slices.Contains(route.methodMatch, request.Method) {
		return false
	}

	for _, regex := range route.pathMatches {
		ok := regex.MatchString(request.RequestURI)
		if !ok {
			return false
		}
	}

	if route.headerMatch != nil {
		for key, value := range maps.All(route.headerMatch) {
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

	if route.pathMatches != nil {
		fmt.Println("\tPath Match: ", route.pathMatches)
	}

	if len(route.hostMatch) > 0 {
		fmt.Println("\tHosts: ", route.hostMatch)
	}

	if len(route.methodMatch) > 0 {
		fmt.Println("\tMethods: ", route.methodMatch)
	}

	if route.headerMatch != nil {
		fmt.Println("\tHeaders: ", route.headerMatch)
	}
}

func (route *Route) Filter(filter filter.Filter) *Route {
	route.filterChain.AddFilter(filter)
	return route
}

func (route *Route) RateLimit(maxRequests int, window time.Duration, resource rateLimiting.ResourceLimiting) {
	route.filterChain.AddFilter(rateLimiting.NewSlidingWindowCounterFilter(maxRequests, window, resource))
}

func (route *Route) RemoveLeftPath(pathNum int) *Route {
	route.filterChain.AddFilter(request.NewRemoveLeftPathFilter(pathNum))
	return route
}

func (route *Route) RemoveRightPath(pathNum int) *Route {
	route.filterChain.AddFilter(request.NewRemoveRightPathFilter(pathNum))
	return route
}
