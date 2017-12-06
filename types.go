package resolvconf

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

// Represents a configuration object
type Conf struct {
	Nameservers []nameserver // List of added nameservers
	Domain      domain       // 'domain' item
	Search      search       // List of the search domains
	Sortlist    sortlist     // 'sortlist' items
	Options     []option     // 'options' items
	logger      *log.Logger
}

type nameserver struct {
	IP net.IP
}

type domain struct {
	Name string
}

type search struct {
	Domains []searchDomain
}

type searchDomain struct {
	Name string
}

type sortlist struct {
	Pairs []sortlistpair
}

type sortlistpair struct {
	Address net.IP
	Netmask net.IP
}

func (s sortlistpair) String() string {
	if len(s.Netmask) > 0 {
		return fmt.Sprintf("%s/%s", s.Address, s.Netmask)
	}
	return fmt.Sprintf("%s", s.Address)
}

type option struct {
	Type  string
	Value int
}

func (o option) String() string {
	if o.Value == -1 {
		return fmt.Sprintf("%s", o.Type)
	}
	return fmt.Sprintf("%s:%d", o.Type, o.Value)
}

// New creates a new configuration
func New() *Conf {
	c := new(Conf)
	c.logger = log.New(ioutil.Discard, "[resolvconf] ", 0)
	return c
}

// Nameserver creates a new nameserver item
func Nameserver(IP net.IP) nameserver {
	return nameserver{IP}
}

// Domain creates a new domain that will be used
// as value for the 'domain' option in the generated file
func Domain(dom string) domain {
	return domain{dom}
}

// SearchDomain creates a new search domain that will be added
// to the 'search' list in the generated file
func SearchDomain(dom string) searchDomain {
	return searchDomain{dom}
}

// SortlistPair creates a new sortlist that will be added to the
// 'sort' item in the resolv.conf file.
// If mask is given the output will be IP/mask e.g.
// 8.8.8.8/255.255.255.0 otherwise output will be
// IP only, e.g. 8.8.8.8
func SortlistPair(addr net.IP, mask ...net.IP) sortlistpair {
	if len(mask) > 1 {
		return sortlistpair{nil, nil}
	} else if len(mask) == 0 {
		return sortlistpair{addr, nil}
	}
	return sortlistpair{addr, mask[0]}
}

// Option creates a new option, val must be a positive number if used.
// Witout val the option will be interpreted as a bolean e.g.
// debug , with a val the option will be interpreted as an
// setvalue, e.g. ndots:5
func Option(t string, val ...int) option {
	// Check va
	opt := option{t, -1}
	if len(val) > 1 {
		return option{"", -1}
	} else if len(val) == 1 {
		if val[0] < 0 {
			return option{"", -1}
		}
		opt.Value = val[0]
	}
	if _, err := parseOption(opt.String()); err != nil {
		return option{"", -1}
	}

	return opt
}
