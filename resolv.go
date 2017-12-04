package resolvconf

import (
	"fmt"
)

func (this *Conf) Add(opts ...interface{}) error {
	for _, o := range opts {
		switch opt := o.(type) {
		case Nameserver:
			if len(this.Nameservers)+1 > 3 {
				return fmt.Errorf("Too many nameserver configs, max is 3")
			}
			// Search if this nameserver is already added
			if this.Find(o.(Nameserver)) != nil {
				return fmt.Errorf("Nameserver %s already exists in conf", opt)
			}

			this.Nameservers = append(this.Nameservers, opt)
		case Domain:
			this.Domain = o.(Domain)
		case Search:
			this.Search = o.(Search)
		case SearchDomain:
			// Search if this search domain is already added
			if this.Find(o) != nil {
				return fmt.Errorf("Search domain %s already exists in conf", opt)
			}
			this.Search.Domains = append(this.Search.Domains, opt)
		case Sortlist:
			this.Sortlist = o.(Sortlist)
		case Sortlistpair:
			if i := this.Find(o); i != nil {
				return fmt.Errorf("Searchlist pair %s already exists in conf", opt)
			}
			this.Sortlist.Pairs = append(this.Sortlist.Pairs, opt)
		case []Option:
			this.Options = o.([]Option)
		case Option:
			if o := this.Find(opt); o != nil {
				return fmt.Errorf("Option %s is already present", o.(Option).Type)
			}
			this.Options = append(this.Options, opt)

		default:
			return fmt.Errorf("Unknown option type %v", opt)
		}
	}
	return nil
}

func (this *Conf) Remove(o interface{}) error {
	i := this.indexOf(this.Find(o))
	_, isdom := o.(Domain)
	if i == -1 && !isdom {
		return fmt.Errorf("Not found")
	}

	switch opt := o.(type) {
	case Nameserver:
		this.Nameservers = append(this.Nameservers[:i], this.Nameservers[i+1:]...)
	case Domain:
		this.Domain = Domain{""}
	case SearchDomain:
		this.Search.Domains = append(this.Search.Domains[:i], this.Search.Domains[i+1:]...)
	case Sortlistpair:
		this.Sortlist.Pairs = append(this.Sortlist.Pairs[:i], this.Sortlist.Pairs[i+1:]...)
	case Option:
		this.Options = append(this.Options[:i], this.Options[i+1:]...)
	default:
		return fmt.Errorf("Unknown option type %v", opt)
	}
	return nil
}

func (this Conf) Find(o interface{}) interface{} {
	i := this.indexOf(o)
	if i == -1 {
		return nil
	}

	switch o.(type) {
	case Nameserver:
		return this.Nameservers[i]
	case SearchDomain:
		return this.Search.Domains[i]
	case Sortlistpair:
		return this.Sortlist.Pairs[i]
	case Option:
		return this.Options[i]
	}
	return nil
}

func (this Conf) indexOf(o interface{}) int {
	switch o.(type) {
	case Nameserver:
		for i, item := range this.Nameservers {
			if item.IP.String() == o.(Nameserver).IP.String() {
				return i
			}
		}
	case Sortlistpair:
		for i, sp := range this.Sortlist.Pairs {
			if sp.Address.String() == o.(Sortlistpair).Address.String() {
				return i
			}
		}
	case SearchDomain:
		for i, sd := range this.Search.Domains {
			if sd.Name == o.(SearchDomain).Name {
				return i
			}
		}
	case Option:
		for i, item := range this.Options {
			if item.Type == o.(Option).Type {
				return i
			}
		}
	}
	return -1
}
