package main

import (
	"github.com/GrongoTheGrog/goteway/internals/config"
	"github.com/GrongoTheGrog/goteway/internals/gateway"
)

func main() {

	appGateway := gateway.NewGateway()

	config.LoadConfig(appGateway)

	appGateway.Start(":9000")
}

