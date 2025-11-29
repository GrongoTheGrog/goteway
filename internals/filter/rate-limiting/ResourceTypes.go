package rate_limiting

import (
	"log"

	"github.com/GrongoTheGrog/goteway/internals/filter"
)

type ResourceLimiting int

const (
	USER ResourceLimiting = iota
	ROUTE
	GATEWAY
)

func getKeyForResource(resource ResourceLimiting, ctx *filter.Context) string {
	switch resource {
	case USER:
		return ctx.RequestIp
	case ROUTE:
		return ctx.Url
	case GATEWAY:
		return "GATEWAY"
	default:
		log.Panic("Unknown resource type.")
		return ""
	}
}
