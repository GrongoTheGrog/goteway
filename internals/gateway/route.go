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
	pathMatch   *regexp.Regexp
	headerMatch map[string]string
	methodMatch []string
	hostMatch   []string
	filterChain *filter.FilterChain
}

func NewRoute(endpoint string) *Route {
	return &Route{
		Endpoint:    endpoint,
		headerMatch: make(map[string]string),
		filterChain: &filter.FilterChain{ProxyFilter: filter.NewProxyFilter()},
	}
}

func (route *Route) PathPattern(regex string) *Route {
	route.pathMatch = regexp.MustCompile(regex)
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

	if route.pathMatch != nil {
		ok := route.pathMatch.MatchString(request.RequestURI)
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

	if route.pathMatch != nil {
		fmt.Println("\tPath Match: ", route.pathMatch.String())
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

func (route *Route) RemoveLeftPath(pathNum int) *Route {
	route.filterChain.AddFilter(request.NewRemoveLeftPathFilter(pathNum))
	return route
}

func (route *Route) RemoveRightPath(pathNum int) *Route {
	route.filterChain.AddFilter(request.NewRemoveRightPathFilter(pathNum))
	return route
}
