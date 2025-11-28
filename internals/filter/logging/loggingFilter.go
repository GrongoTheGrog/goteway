package logging

import (
	"log"
	"net/http"
	"time"

	"github.com/GrongoTheGrog/goteway/internals/filter"
)

type LogOption int

const (
	Method LogOption = iota
	Path
	Status
	Latency
	FullUrl
)

func NewLogFilter(options ...LogOption) filter.Filter {
	if len(options) == 0 {
		log.Panic("Log filter must have at least one log option enabled")
	}

	return filter.NewBasicFilter(func(ctx *filter.Context) *http.Response {

		beforeRequest := time.Now()
		measureLatency := false
		logStatus := false

		for _, option := range options {
			switch option {
			case Method:
				ctx.Log("%s", ctx.Request.Method)
			case Latency:
				measureLatency = true
			case Status:
				logStatus = true
			case Path:
				ctx.Log("%s", ctx.Request.URL)
			case FullUrl:
				ctx.Log("%s%s", ctx.Request.URL)
			default:
				log.Panic("Unknown log option.")
			}
		}

		response := ctx.RunNextFilter()

		if measureLatency {
			elapsedTimeMs := time.Since(beforeRequest)
			ctx.Log("Duration of %v ms", elapsedTimeMs)
		}

		if logStatus {
			ctx.Log("Status %v", response.StatusCode)
		}

		return response
	})
}
