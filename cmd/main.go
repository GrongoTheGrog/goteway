package main

import (
	"time"

	"github.com/GrongoTheGrog/goteway/internals/filter/logging"
	rate_limiting "github.com/GrongoTheGrog/goteway/internals/filter/rate-limiting"
	"github.com/GrongoTheGrog/goteway/internals/gateway"
)

func main() {

	gateway := gateway.NewGateway()

	gateway.
		TokenBucketFilter(100, 1*time.Second, rate_limiting.USER).
		LogFilter(
			logging.Path,
			logging.Status,
			logging.Latency,
		)

	gateway.NewRoute("/user-service/*", "http://localhost:8082").
		RemoveLeftPath(1)

	gateway.Start(":9000")
}
