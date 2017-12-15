package resolvconf

import (
	"fmt"
	"unicode/utf8"
)

// SearchDomain is one of the items in the search list
type SearchDomain struct {
	Name string
}

// NewSearchDomain creates a new search domain that will be added
// to the 'search' list in the generated file
func NewSearchDomain(dom string) *SearchDomain {
	return &SearchDomain{dom}
}

func (sd SearchDomain) applyLimits(conf *Conf) (bool, error) {
	// Search if conf search domain is already added
	if conf.Find(sd) != nil {
		return false, fmt.Errorf("Search domain %s already exists in conf", sd.Name)
	}
	// Check max limit
	doms := conf.GetSearchDomains()
	if len(doms) == searchDomainMaxCount {
		return false, fmt.Errorf("Too many search domains, %d is maximum", searchDomainMaxCount)
	}
	// Check max char count limit
	var charcount int
	for _, str := range doms {
		charcount += utf8.RuneCountInString(str.Name)
	}
	if charcount+utf8.RuneCountInString(sd.Name) > searchDomainMaxCharCount {
		return false, fmt.Errorf("Too many charactes is search domain list, %d is maximum", searchDomainMaxCharCount)
	}

	return true, nil
}

func (sd SearchDomain) String() string {
	return sd.Name
}

// Equal compares two search domains with each other, returns true if equal
func (sd SearchDomain) Equal(b ConfItem) bool {
	if item, ok := b.(*SearchDomain); ok {
		return sd.Name == item.Name
	}
	return false
}
