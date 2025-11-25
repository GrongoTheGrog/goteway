package filter

import (
	"io"
	"log"
	"net/http"
)

type EntryFilter struct {
	next Filter
}

func NewEntryFilter() *EntryFilter {
	return &EntryFilter{}
}

func (entryFilter *EntryFilter) StartChain(writer http.ResponseWriter, request *http.Request, endpoint string) {
	context := &Context{
		Request:    request,
		Url:        endpoint,
		Attributes: make(map[string]interface{}),
		next:       entryFilter,
	}

	response := entryFilter.RunFilter(context)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading the response body from request: %s", request.RequestURI)
		writer.WriteHeader(500)
		writer.Write([]byte("Error reading the response body from request"))
		return
	}

	writer.WriteHeader(response.StatusCode)
	_, err = writer.Write(body)
	if err != nil {
		log.Printf("Failed to write response body stream to client response.")
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to write response body stream to client response."))
		return
	}

	for name, values := range response.Header {
		for _, value := range values {
			writer.Header().Add(name, value)
		}
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
