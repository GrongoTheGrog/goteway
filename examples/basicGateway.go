package examples

import (
	"net/http"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/GrongoTheGrog/goteway/internals/gateway"
)

func main() {
	gw := gateway.NewGateway()

	service1 := gateway.NewRoute("http://service"). // Creates a new route with the final endpoint
							PathPattern("/service1/*", "/v1/service1/*"). // Regex containing the path to be routed to the endpoint
							Header("X-Service", "service1").              // Requests being routed must contain that header
							Hosts("frontend.com", "www.frontend.com").    // Hosts allowed to access the route
							Methods("GET", "OPTIONS")                     // Methods allowed

	service1.RemoveLeftPath(1) // Adds a basic filter that removes one path segment in the left

	customFilter := filter.NewBasicFilter(func(ctx *filter.Context) *http.Response {
		ctx.Log("Logging method from custom filter: %s", ctx.Request.Method)
		return ctx.RunNextFilter()
	})

	service1.Filter(customFilter) // Binds the filter to the route
	gw.AddFilter(customFilter)    // Binds the filter to the entire gateway

	gw.AddRoute(service1) // Assigns route to the gateway

	gw.Start(":9000")
}
