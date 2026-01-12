package config

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/GrongoTheGrog/goteway/internals/filter"
	"github.com/GrongoTheGrog/goteway/internals/filter/authentication"
	"github.com/GrongoTheGrog/goteway/internals/filter/logging"
	"github.com/GrongoTheGrog/goteway/internals/filter/rateLimiting"
	"github.com/GrongoTheGrog/goteway/internals/gateway"
	"gopkg.in/yaml.v3"
)

func LoadConfig(gw *gateway.Gateway) {
	b, err := os.ReadFile("goteway.yml")
	if err != nil {
		log.Println("Unable to open config file")
	}

	config := &GeneralConfig{}
	err = yaml.Unmarshal(b, config)
	if err != nil {
		log.Panicf("Error parsing yaml file: %s", err.Error())
	}

	loadGatewayConfig(gw, config.Gateway)
	loadRoute(gw, config.Routes)
}

func loadGatewayConfig(gw *gateway.Gateway, config GatewayConfig) {
	if config.LogFilter != nil {
		var options []logging.LogOption

		for _, option := range config.LogFilter {
			switch option {
			case "method":
				options = append(options, logging.Method)
			case "path":
				options = append(options, logging.Path)
			case "latency":
				options = append(options, logging.Latency)
			case "full_url":
				options = append(options, logging.FullUrl)
			case "status":
				options = append(options, logging.Status)
			default:
				panic("Unknown option selected.")
			}
		}
		gw.LogFilter(options...)
	}

	loadAuthConfig(gw, config.Auth)
}

func loadAuthConfig(gw *gateway.Gateway, config authentication.AuthorizationConfig) {
	if config.Jwt.Enabled == false {
		return
	}

	gw.AddFilter(authentication.NewJwtFilter(config))
}

func loadRoute(gw *gateway.Gateway, configs []RouteConfig) {
	for _, routeConfig := range configs {
		if !routeConfig.Enabled {
			continue
		}

		route := gateway.NewRoute(routeConfig.Endpoint)
		route.PathPattern(routeConfig.Paths...)

		var resource rateLimiting.ResourceLimiting
		switch routeConfig.RateLimiting.Resource {
		case "user":
			resource = rateLimiting.USER
		case "route":
			resource = rateLimiting.ROUTE
		case "gateway":
			resource = rateLimiting.GATEWAY
		default:
			panic("Wrong type of limiting resource provided")
		}

		if len(routeConfig.Headers) != 0 {
			route.Filter(filter.NewBasicFilter(
				func(ctx *filter.Context) *http.Response {
					response := ctx.RunNextFilter()
					for _, header := range routeConfig.Headers {
						ctx.Log("Added header %s=%s", header.Name, header.Value)
						response.Header.Add(header.Name, header.Value)
					}
					return response
				},
			))
		}

		if routeConfig.RateLimiting.Enabled {
			route.RateLimit(
				routeConfig.RateLimiting.MaxRequests,
				time.Duration(routeConfig.RateLimiting.IntervalSeconds)*time.Second,
				resource,
			)
		}

		if routeConfig.Status != 0 {
			route.Filter(filter.NewBasicFilter(func(ctx *filter.Context) *http.Response {

				response := ctx.RunNextFilter()

				response.StatusCode = routeConfig.Status
				response.Status = strconv.Itoa(routeConfig.Status)

				return response
			}))
		}

		gw.AddRoute(route)
	}
}
