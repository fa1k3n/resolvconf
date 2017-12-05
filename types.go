package resolvconf

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

type Conf struct {
	Nameservers []nameserver
	Domain      domain
	Search      search
	Sortlist    sortlist
	Options     []option
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

func New() *Conf {
	c := new(Conf)
	c.logger = log.New(ioutil.Discard, "[resolvconf] ", 0)
	return c
}

func Nameserver(IP net.IP) nameserver {
	return nameserver{IP}
}

func Domain(dom string) domain {
	return domain{dom}
}

func SearchDomain(dom string) searchDomain {
	return searchDomain{dom}
}

func SortlistPair(addr net.IP, mask ...net.IP) sortlistpair {
	if len(mask) > 1 {
		return sortlistpair{nil, nil}
	} else if len(mask) == 0 {
		return sortlistpair{addr, nil}
	}
	return sortlistpair{addr, mask[0]}
}

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
