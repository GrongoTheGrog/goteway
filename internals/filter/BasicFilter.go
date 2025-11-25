package filter

import "net/http"

type BasicFilter struct {
	next       Filter
	filterFunc RunFilter
}

func NewBasicFilter(filterFunc RunFilter) *BasicFilter {
	return &BasicFilter{
		next:       nil,
		filterFunc: filterFunc,
	}
}

func (b *BasicFilter) RunFilter(context *Context) *http.Response {
	return b.filterFunc(context)
}

func (b *BasicFilter) Next() Filter {
	return b.next
}

func (b *BasicFilter) SetNext(filter Filter) {
	b.next = filter
}
