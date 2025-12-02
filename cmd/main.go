package main

import (
	"time"

	"github.com/GrongoTheGrog/goteway/internals/filter/logging"
	"github.com/GrongoTheGrog/goteway/internals/filter/rateLimiting"
	"github.com/GrongoTheGrog/goteway/internals/gateway"
)

func main() {

	appGateway := gateway.NewGateway()

	userServiceRoute := gateway.NewRoute("http://localhost:8082").
		PathPattern("/user-service/*").
		Methods("POST", "GET").
		RemoveLeftPath(1)

	videoServiceRoute := gateway.NewRoute("https://amazon.com").
		PathPattern("/video-service/*").
		Methods("GET").
		RemoveLeftPath(1)

	appGateway.
		SlidingWindowCounterFilter(100, 1*time.Minute, rateLimiting.USER).
		LogFilter(logging.Path, logging.Status, logging.Latency).
		AddRoute(userServiceRoute).
		AddRoute(videoServiceRoute)

	appGateway.Start(":9000")
}
