package filter

import (
	"log"
	"net/http"
)

type FilterChain struct {
	Last        Filter
	First       Filter
	ProxyFilter *ProxyFilter
	EntryFilter *EntryFilter
}

func (fc *FilterChain) Execute(writer http.ResponseWriter, request *http.Request, endpoint string) {
	fc.EntryFilter.StartChain(writer, request, endpoint)
}

func (fc *FilterChain) AddFilter(filter Filter) {
	filter.SetNext(fc.ProxyFilter)

	if fc.First == nil {
		if fc.EntryFilter != nil {
			fc.EntryFilter.SetNext(filter)
		}
		fc.First = filter
		fc.Last = filter
		return
	}

	fc.Last.SetNext(filter)
	fc.Last = filter
}

func (fc *FilterChain) AddFilterAfter(newFilter, afterFilter Filter) {
	if _, ok := afterFilter.(*ProxyFilter); ok {
		log.Fatal("Can't put a filter after the proxy filter.")
	}

	if fc.First == nil {
		log.Fatal("The filter provided as after reference does not exist.")
	}

	if fc.Last == afterFilter {
		fc.Last.SetNext(newFilter)
		fc.Last = newFilter
		newFilter.SetNext(fc.ProxyFilter)
	}

	cur := fc.First

	for cur != nil {
		if cur == afterFilter {
			newFilter.SetNext(cur.Next())
			cur.SetNext(newFilter)
			return
		}
	}

	log.Fatal("The filter provided as after reference does not exist.")

}

func (fc *FilterChain) AddFilterBefore(newFilter, beforeFilter Filter) {
	if _, ok := beforeFilter.(*EntryFilter); ok {
		log.Fatal("Can't put a filter before the entry filter.")
	}

	if fc.First == nil {
		log.Fatal("The filter provided as before reference does not exist.")
	}

	if fc.First == beforeFilter {
		if fc.EntryFilter != nil {
			fc.EntryFilter.SetNext(newFilter)
		}
		newFilter.SetNext(fc.First)
		fc.First = newFilter
	}

	cur := fc.First

	for cur.Next() != nil {
		if cur.Next() == beforeFilter {
			newFilter.SetNext(cur.Next())
			cur.SetNext(newFilter)
			return
		}
	}

	log.Fatal("The filter provided as before reference does not exist.")
}

// CombineFilterChains appends all filters from chain to the end of fc.
//
// It sets the last filter of fc to point to the first filter of chain.
// CombineFilterChains terminates the program using log.Fatal if fc contains
// a ProxyFilter or if chain contains an EntryFilter.
func (fc *FilterChain) CombineFilterChains(chain *FilterChain) {
	if fc.ProxyFilter != nil {
		log.Fatal("Calling filter chain can't have a proxy filter.")
	}

	if chain.EntryFilter != nil {
		log.Fatal("Referenced filter chain can't have an entry filter.")
	}

	if chain.First == nil {
		fc.Last.SetNext(chain.ProxyFilter)
		return
	}

	fc.Last.SetNext(chain.First)
}
