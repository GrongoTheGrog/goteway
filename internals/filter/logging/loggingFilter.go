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

		rId := ctx.GetRequestId()

		beforeRequest := time.Now()
		measureLatency := false
		logStatus := false

		for _, option := range options {
			switch option {
			case Method:
				log.Printf("%s %s", filter.RequestPrefix(rId), ctx.Request.Method)
			case Latency:
				measureLatency = true
			case Status:
				logStatus = true
			case Path:
				log.Printf("%s %s", filter.RequestPrefix(rId), ctx.Request.URL)
			case FullUrl:
				log.Printf("%s %s%s", filter.RequestPrefix(rId), ctx.Url, ctx.Request.URL)
			default:
				log.Panic("Unknown log option.")
			}
		}

		response := ctx.RunNextFilter()

		if measureLatency {
			elapsedTimeMs := (time.Now().Nanosecond() - beforeRequest.Nanosecond()) / 1000
			log.Printf("%s Duration of %v ms", filter.RequestPrefix(rId), elapsedTimeMs)
		}

		if logStatus {
			log.Printf("%s Status %v", filter.RequestPrefix(rId), response.StatusCode)
		}

		return response
	})
}
