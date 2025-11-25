package filter

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/GrongoTheGrog/goteway/internals/utils"
)

type ProxyFilter struct {
}

func NewProxyFilter() *ProxyFilter {
	return &ProxyFilter{}
}

func (proxyFilter *ProxyFilter) RunFilter(context *Context) *http.Response {

	body, err := io.ReadAll(context.Request.Body)
	newBuffer := io.NopCloser(bytes.NewBuffer(body))

	newRequest, err := http.NewRequest(context.Request.Method, context.Url+context.Request.URL.Path, newBuffer)
	if err != nil {
		log.Printf("Could not form request: %s", err.Error())
		return utils.ErrorResponse("Could not form request", 500)
	}

	response, err := http.DefaultClient.Do(newRequest)
	if err != nil || response == nil {
		log.Print("Gateway could not finish the request.")
		return utils.ErrorResponse("Gateway could not perform request.", 500)
	}

	return response
}

func (proxyFilter *ProxyFilter) Next() Filter {
	return nil
}

func (proxyFilter *ProxyFilter) SetNext(filter Filter) {
}
