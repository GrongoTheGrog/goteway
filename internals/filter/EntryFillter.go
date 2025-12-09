package filter

import (
	"io"
	"net"
	"net/http"

	"github.com/google/uuid"
)

type EntryFilter struct {
	next Filter
}

func NewEntryFilter() *EntryFilter {
	return &EntryFilter{}
}

func (entryFilter *EntryFilter) StartChain(writer http.ResponseWriter, request *http.Request, endpoint string) {
	ip, _, _ := net.SplitHostPort(request.RemoteAddr)

	context := &Context{
		Request:    request,
		RequestIp:  ip,
		Url:        endpoint,
		requestId:  uuid.NewString(),
		attributes: make(map[string]interface{}),
		next:       entryFilter,
	}
	context.Log("Routing %s to %s", request.RequestURI, endpoint)

	request.Header.Set("X-Request-ID", context.requestId)

	response := entryFilter.RunFilter(context)

	for name, values := range response.Header {
		for _, value := range values {
			writer.Header().Add(name, value)
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		context.Log("Error reading the response body from request: %s", request.RequestURI)
		writer.WriteHeader(500)
		writer.Write([]byte("Error reading the response body from request"))
		return
	}

	writer.WriteHeader(response.StatusCode)
	_, err = writer.Write(body)
	if err != nil {
		context.Log("Failed to write response body stream to client response.")
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to write response body stream to client response."))
		return
	}

}

func (entryFilter *EntryFilter) RunFilter(context *Context) *http.Response {
	return context.RunNextFilter()
}

func (entryFilter *EntryFilter) Next() Filter {
	return entryFilter.next
}

func (entryFilter *EntryFilter) SetNext(filter Filter) {
	entryFilter.next = filter
}
