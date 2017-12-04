package resolvconf

import (
	"fmt"
	"net"
)

type Conf struct {
	Nameservers []Nameserver
	Domain      Domain
	Search      Search
	Sortlist    Sortlist
	Options     []Option
}

type Nameserver struct {
	IP net.IP
}

type Domain struct {
	Name string
}

type Search struct {
	Domains []SearchDomain
}

type SearchDomain struct {
	Name string
}

type Sortlist struct {
	Pairs []Sortlistpair
}

type Sortlistpair struct {
	Address net.IP
	Netmask net.IP
}

func (s Sortlistpair) String() string {
	if len(s.Netmask) > 0 {
		return fmt.Sprintf("%s/%s", s.Address, s.Netmask)
	}
	return fmt.Sprintf("%s", s.Address)
}

type Option struct {
	Type  string
	Value int
}

func (o Option) String() string {
	if o.Value == -1 {
		return fmt.Sprintf("%s", o.Type)
	}
	return fmt.Sprintf("%s:%d", o.Type, o.Value)
}

func New() *Conf {
	return new(Conf)
}

func NewNameserver(IP net.IP) (Nameserver, error) {
	return Nameserver{IP}, nil
}

func NewDomain(dom string) (Domain, error) {
	return Domain{dom}, nil
}

func NewSearchDomain(dom string) (SearchDomain, error) {
	return SearchDomain{dom}, nil
}

// @TODO: Make variadic and netmask optional
func NewSortlistPair(addr, mask net.IP) (Sortlistpair, error) {
	return Sortlistpair{addr, mask}, nil
}

func NewOption(t string, val ...int) (Option, error) {
	// Check va
	opt := Option{t, -1}
	if len(val) > 1 {
		return Option{"", -1}, fmt.Errorf("Only a single value is allowed, found %v", val)
	} else if len(val) == 1 {
		if val[0] < 0 {
			return Option{"", -1}, fmt.Errorf("Only positive values are allowed, found %d", val[0])
		}
		opt.Value = val[0]
	}
	if _, err := parseOption(opt.String()); err != nil {
		return Option{"", -1}, err
	}

	return opt, nil
}
