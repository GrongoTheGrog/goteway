package filter

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GrongoTheGrog/goteway/internals/utils"
)

type RunFilter func(context *Context) *http.Response
type Filter interface {
	RunFilter(context *Context) *http.Response
	Next() Filter
	SetNext(filter Filter)
}

type Context struct {
	Request    *http.Request
	Url        string
	RequestIp  string
	requestId  string
	attributes map[string]interface{}
	next       Filter
}

func (context *Context) GetRequestId() string {
	return context.requestId
}

func (context *Context) SetAttribute(name string, attribute interface{}) {
	context.attributes[name] = attribute
}

func (context *Context) GetAttribute(name string) (interface{}, bool) {
	value, exists := context.attributes[name]
	return value, exists
}

func (context *Context) RunNextFilter() *http.Response {
	if context.next == nil || context.next.Next() == nil {
		context.Log("The filter chain is broken, there is no filter set to run next.")
		return utils.ErrorResponse("Broken filter chain.", 500)
	}

	context.next = context.next.Next()
	return context.next.RunFilter(context)
}

func (context *Context) Log(message string, more ...any) {
	preString := fmt.Sprintf("[%s] [%s] [%s]: ", context.GetRequestId(), context.Url, context.RequestIp)
	log.Printf(preString+message, more...)
}
