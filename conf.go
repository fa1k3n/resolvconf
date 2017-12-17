package resolvconf

import (
	"io/ioutil"
	"log"
)

// Limits
const (
	searchDomainMaxCount     = 6   // Maximum count of search domains
	searchDomainMaxCharCount = 256 // Maximum total number of chars in search domains
	nameserversMaxCount      = 3   // Maximum number of nameservers
	sortListMaxCount         = 10  // Maximum number of items in sortlist
	optionNdotsMax           = 15  // Maximum ndots value, silently capped
	optionTimeoutMax         = 30  // Maximum timeout value, silently capped
	optionAttemptsMax        = 5   // Maximum attempts value, silently capped
)

// Conf represents a configuration object
type Conf struct {
	items  []ConfItem
	logger *log.Logger
}

// New creates a new configuration
func New() *Conf {
	c := new(Conf)
	c.logger = log.New(ioutil.Discard, "[resolvconf] ", 0)
	return c
}

// GetNameservers returns a list of all added nameservers
func (conf *Conf) GetNameservers() []Nameserver {
	var ret []Nameserver
	for _, item := range conf.items {
		if _, ok := item.(*Nameserver); ok {
			ret = append(ret, *item.(*Nameserver))
		}
	}
	return ret
}

// GetSortItems returns list of all added sortitems
func (conf *Conf) GetSortItems() []SortItem {
	var ret []SortItem
	for _, item := range conf.items {
		if _, ok := item.(*SortItem); ok {
			ret = append(ret, *item.(*SortItem))
		}
	}
	return ret
}

// GetDomain returns current domain
func (conf *Conf) GetDomain() Domain {
	for _, item := range conf.items {
		if d, ok := item.(*Domain); ok {
			return *d
		}
	}
	return Domain{}
}

// GetSearchDomains returns a list of all added SearchDomains
func (conf *Conf) GetSearchDomains() []SearchDomain {
	var ret []SearchDomain
	for _, item := range conf.items {
		if _, ok := item.(*SearchDomain); ok {
			ret = append(ret, *item.(*SearchDomain))
		}
	}
	return ret
}

// GetOptions returns a list of all added options
func (conf *Conf) GetOptions() []Option {
	var ret []Option
	for _, item := range conf.items {
		if _, ok := item.(*Option); ok {
			ret = append(ret, *item.(*Option))
		}
	}
	return ret
}
