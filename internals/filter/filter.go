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
	attributes map[string]interface{}
	next       Filter
}

func (context *Context) SetAttribute(name string, attribute interface{}) {
	context.attributes[name] = attribute
}

func (context *Context) GetAttribute(name string) (interface{}, bool) {
	value, exists := context.attributes[name]
	return value, exists
}

func (context *Context) RunNextFilter() *http.Response {
	context.next = context.next.Next()
	return context.next.RunFilter(context)
}
