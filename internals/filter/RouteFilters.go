package filter

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func NewRemoveLeftPathFilter(pathCount int) Filter {
	return NewBasicFilter(func(ctx *Context) *http.Response {

		fullPath := ctx.Request.URL.Path
		paths := strings.Split(fullPath, "/")
		paths = paths[pathCount+1 : len(paths)]

		fullPath = "/" + strings.Join(paths, "/")
		ctx.Request.URL.Path = fullPath

		return ctx.RunNextFilter()
	})
}

func NewRemoveRightPathFilter(pathCount int) Filter {
	return NewBasicFilter(func(ctx *Context) *http.Response {

		fullPath := ctx.Request.URL.Path
		paths := strings.Split(fullPath, "/")
		paths = paths[0 : len(paths)-pathCount]

		fullPath = "/" + strings.Join(paths, "/")
		ctx.Request.URL.Path = fullPath

		return ctx.RunNextFilter()
	})
}

type LogOption int

const (
	Method LogOption = iota
	Path
	Status
	Latency
	FullUrl
)

func NewLogFilter(options ...LogOption) Filter {
	if len(options) == 0 {
		log.Panic("Log filter must have at least one log option enabled")
	}

	return NewBasicFilter(func(ctx *Context) *http.Response {

		rId := ctx.GetRequestId()

		beforeRequest := time.Now()
		measureLatency := false
		logStatus := false

		for _, option := range options {
			switch option {
			case Method:
				log.Printf("%s %s", requestPrefix(rId), ctx.Request.Method)
			case Latency:
				measureLatency = true
			case Status:
				logStatus = true
			case Path:
				log.Printf("%s %s", requestPrefix(rId), ctx.Request.URL)
			case FullUrl:
				log.Printf("%s %s%s", requestPrefix(rId), ctx.Url, ctx.Request.URL)
			default:
				log.Panic("Unknown log option.")
			}
		}

		response := ctx.RunNextFilter()

		if measureLatency {
			elapsedTimeMs := (time.Now().Nanosecond() - beforeRequest.Nanosecond()) / 1000
			log.Printf("%s Duration of %v ms", requestPrefix(rId), elapsedTimeMs)
		}

		if logStatus {
			log.Printf("%s Status %v", requestPrefix(rId), response.StatusCode)
		}

		return response
	})
}
