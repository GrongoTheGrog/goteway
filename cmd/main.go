package main

import (
	"time"

	"github.com/GrongoTheGrog/goteway/internals/filter/logging"
	"github.com/GrongoTheGrog/goteway/internals/filter/rateLimiting"
	"github.com/GrongoTheGrog/goteway/internals/gateway"
)

func main() {

	appGateway := gateway.NewGateway()

	appGateway.
		SlidingWindowCounterFilter(100, 1*time.Second, rateLimiting.USER).
		LogFilter(logging.Path, logging.Status, logging.Latency)

	appGateway.NewRoute("/user-service/*", "http://localhost:8082").
		RemoveLeftPath(1)

	appGateway.NewRoute("/video-service/*", "http://localhost:8081").
		RemoveLeftPath(1)

	appGateway.Start(":9000")
}
