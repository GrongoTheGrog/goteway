package main

import (
	"log"
	"net/http"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/GrongoTheGrog/goteway/internals/gateway"
)

func main() {

	filter1 := filter.NewBasicFilter(func(ctx *filter.Context) *http.Response {

		ctx.SetAttribute("token", "token")
		return ctx.RunNextFilter()
	})

	filter2 := filter.NewBasicFilter(func(ctx *filter.Context) *http.Response {
		token, _ := ctx.GetAttribute("token")
		log.Printf("Logging context frmo other filter: %s", token)

		response := ctx.RunNextFilter()

		response.StatusCode = 200
		return response
	})

	gateway := gateway.NewGateway()

	//TODO Hide filter chain with fluent
	gateway.FilterChain.AddFilter(filter1)
	gateway.FilterChain.AddFilterAfter(filter2, filter1)

	//TODO Make new Route return route
	gateway.NewRoute("/user-service/*", "http://localhost:8082")

	gateway.Start("9000")
}
