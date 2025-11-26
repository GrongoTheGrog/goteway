package request

import (
	"log"
	"net/http"
	"strings"

	"github.com/GrongoTheGrog/goteway/internals/filter"
)

func NewRemoveLeftPathFilter(pathCount int) filter.Filter {
	return filter.NewBasicFilter(func(ctx *filter.Context) *http.Response {

		fullPath := ctx.Request.URL.Path
		paths := strings.Split(fullPath, "/")
		paths = paths[pathCount+1 : len(paths)]

		fullPath = "/" + strings.Join(paths, "/")

		log.Printf("%s Changed path from %s to %s",
			filter.RequestPrefix(ctx.GetRequestId()),
			ctx.Request.URL.Path,
			fullPath,
		)

		ctx.Request.URL.Path = fullPath

		return ctx.RunNextFilter()
	})
}

func NewRemoveRightPathFilter(pathCount int) filter.Filter {
	return filter.NewBasicFilter(func(ctx *filter.Context) *http.Response {

		fullPath := ctx.Request.URL.Path
		paths := strings.Split(fullPath, "/")
		paths = paths[0 : len(paths)-pathCount]

		fullPath = "/" + strings.Join(paths, "/")

		log.Printf("%s Changed path from %s to %s",
			filter.RequestPrefix(ctx.GetRequestId()),
			ctx.Request.URL.Path,
			fullPath,
		)

		ctx.Request.URL.Path = fullPath

		return ctx.RunNextFilter()
	})
}
