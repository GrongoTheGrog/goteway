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
		log.Printf("Logging context from other filter: %s", token)

		return ctx.RunNextFilter()
	})

	gateway := gateway.NewGateway()

	gateway.
		AddFilter(filter1).
		AddFilter(filter2).
		LogFilter(
			filter.Path,
		)

	gateway.NewRoute("/user-service/*", "http://localhost:8082").
		RemoveLeftPath(1)

	gateway.Start(":9000")
}
