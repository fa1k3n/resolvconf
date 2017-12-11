package resolvconf

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"reflect"
)

// Conf represents a configuration object
type Conf struct {
	items  []ConfItem
	logger *log.Logger
}

func (conf *Conf) Nameservers() []Nameserver {
	var ret []Nameserver
	for _, item := range conf.items {
		if reflect.TypeOf(item).Name() == "Nameserver" {
			ret = append(ret, item.(Nameserver))
		}
	}
	return ret
}

func (conf *Conf) Sortlist() []SortlistPair {
	var ret []SortlistPair
	for _, item := range conf.items {
		if reflect.TypeOf(item).Name() == "SortlistPair" {
			ret = append(ret, item.(SortlistPair))
		}
	}
	return ret
}

func (conf *Conf) Domain() Domain {
	for _, item := range conf.items {
		if reflect.TypeOf(item).Name() == "Domain" {
			return item.(Domain)
		}
	}
	return Domain{}
}

func (conf *Conf) Search() []SearchDomain {
	var ret []SearchDomain
	for _, item := range conf.items {
		if reflect.TypeOf(item).Name() == "SearchDomain" {
			ret = append(ret, item.(SearchDomain))
		}
	}
	return ret
}

func (conf *Conf) Options() []Option {
	var ret []Option
	for _, item := range conf.items {
		if reflect.TypeOf(item).Name() == "Option" {
			ret = append(ret, item.(Option))
		}
	}
	return ret
}

type ConfItem interface {
	fmt.Stringer
	AddToConf(conf *Conf) error
	Equal(b ConfItem) bool
}

type Nameserver struct {
	IP net.IP
}

func (ns Nameserver) AddToConf(conf *Conf) error {

	if len(conf.Nameservers())+1 > 3 {
		return fmt.Errorf("Too many Nameserver configs, max is 3")
	}
	// Search if conf Nameserver is already added
	if conf.Find(ConfItem(ns)) != nil {
		return fmt.Errorf("Nameserver %s already exists in conf", ns.IP)
	}
	conf.logger.Printf("Added Nameserver %s", ns.IP)
	conf.items = append(conf.items, ns)
	return nil
}

func (ns Nameserver) Equal(b ConfItem) bool {
	if item, ok := b.(Nameserver); ok {
		return ns.IP.String() == item.IP.String()
	}
	return false
}

func (ns Nameserver) String() string {
	return ns.IP.String()
}

type Domain struct {
	Name string
}

func (dom Domain) AddToConf(conf *Conf) error {
	conf.logger.Printf("Added domain %s", dom.Name)
	i := conf.indexOf(conf.Domain())
	if i != -1 {
		conf.items[i] = dom
	} else {
		conf.items = append(conf.items, dom)
	}
	return nil
}

func (dom Domain) String() string {
	return dom.Name
}

func (dom Domain) Equal(b ConfItem) bool {
	if item, ok := b.(Domain); ok {
		return dom.Name == item.Name
	}
	return false
}

type SearchDomain struct {
	Name string
}

func (sd SearchDomain) AddToConf(conf *Conf) error {
	// Search if conf search domain is already added
	if conf.Find(sd) != nil {
		return fmt.Errorf("Search domain %s already exists in conf", sd.Name)
	}
	conf.logger.Printf("Added search domain %s", sd.Name)
	conf.items = append(conf.items, sd)
	return nil
}

func (sd SearchDomain) String() string {
	return sd.Name
}

func (sd SearchDomain) Equal(b ConfItem) bool {
	if item, ok := b.(SearchDomain); ok {
		return sd.Name == item.Name
	}
	return false
}

type SortlistPair struct {
	Address net.IP
	Netmask net.IP
}

func (pair SortlistPair) AddToConf(conf *Conf) error {
	if i := conf.Find(pair); i != nil {
		return fmt.Errorf("Sortlist pair %s already exists in conf", pair)
	}
	if len(conf.Sortlist()) == 10 {
		return fmt.Errorf("Too long sortlist, 10 is maximum")
	}
	conf.logger.Printf("Added sortlist pair %s", pair)
	conf.items = append(conf.items, pair)
	return nil
}

func (slp SortlistPair) Equal(b ConfItem) bool {
	if item, ok := b.(SortlistPair); ok {
		return slp.Address.String() == item.Address.String()
	}

	return false
}

func (slp *SortlistPair) SetNetmask(nm net.IP) *SortlistPair {
	slp.Netmask = nm
	return slp
}

func (slp SortlistPair) GetNetmask() net.IP {
	return slp.Netmask
}

func (s SortlistPair) String() string {
	if len(s.Netmask) > 0 {
		return fmt.Sprintf("%s/%s", s.Address, s.Netmask)
	}
	return fmt.Sprintf("%s", s.Address)
}

type Option struct {
	Type  string
	Value int
}

func (opt Option) AddToConf(conf *Conf) error {
	if opt.Type == "ndots" && opt.Value < 0 {
		return fmt.Errorf("Bad value %d", opt.Value)
	}
	if _, e := parseOption(opt.String()); e != nil {
		return fmt.Errorf("Unknown option %s", opt)
	}
	if o := conf.Find(opt); o != nil {
		return fmt.Errorf("Option %s is already present")
	}
	conf.logger.Printf("Added option %s", opt)
	conf.items = append(conf.items, opt)
	return nil
}

func (opt Option) Equal(b ConfItem) bool {
	if o, ok := b.(Option); ok {
		return opt.Type == o.Type
	}
	return false
}

func (opt *Option) Set(value int) *Option {
	if value < 0 {
		return opt
	}
	opt.Value = value
	return opt
}

func (opt Option) Get() int {
	return opt.Value
}

func (o Option) String() string {
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

// Nameserver creates a new Nameserver item
func NewNameserver(IP net.IP) *Nameserver {
	return &Nameserver{IP}
}

// Domain creates a new domain that will be used
// as value for the 'domain' option in the generated file
func NewDomain(dom string) *Domain {
	return &Domain{dom}
}

// SearchDomain creates a new search domain that will be added
// to the 'search' list in the generated file
func NewSearchDomain(dom string) *SearchDomain {
	return &SearchDomain{dom}
}

// SortlistPair creates a new sortlist that will be added to the
// 'sort' item in the resolv.conf file.
// If mask is given the output will be IP/mask e.g.
// 8.8.8.8/255.255.255.0 otherwise output will be
// IP only, e.g. 8.8.8.8
func NewSortlistPair(addr net.IP) *SortlistPair {
	/*if len(mask) > 1 {
		return SortlistPair{nil, nil}
	} else if len(mask) == 0 {
		return SortlistPair{addr, nil}
	}
	return SortlistPair{addr, mask[0]}*/
	slp := new(SortlistPair)
	slp.Address = addr
	return slp
}

// Option creates a new option, val must be a positive number if used.
// Witout val the option will be interpreted as a bolean e.g.
// debug , with a val the option will be interpreted as an
// setvalue, e.g. ndots:5
func NewOption(t string) *Option {
	// Check va
	opt := new(Option)
	opt.Type = t
	opt.Value = -1

	//if _, err := parseOption(opt.String()); err != nil {
	//	return Option{"", -1}
	//}

	return opt
}
