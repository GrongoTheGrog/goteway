package filter

import "net/http"

type RunFilter func(context *Context) *http.Response
type Filter interface {
	RunFilter(context *Context) *http.Response
	Next() Filter
	SetNext(filter Filter)
}

type Context struct {
	Request    *http.Request
	Url        string
	Attributes map[string]interface{}
	next       Filter
}

func (context *Context) RunNextFilter() *http.Response {
	context.next = context.next.Next()
	return context.next.RunFilter(context)
}
