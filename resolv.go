// Package resolvconf provides an interface to read, create and manipulate
// resolv.conf files
package resolvconf

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"io"
	"log"
)

// Add items to the configuration.
//
// Errors are accumulated and can be reinterpreted as
// an multierror type. Logging will occur if logging has
// been setup using the EnableLogging call
func (this *Conf) Add(opts ...interface{}) error {
	var err *multierror.Error
	for _, o := range opts {
		switch opt := o.(type) {
		case nameserver:
			if len(this.Nameservers)+1 > 3 {
				err = multierror.Append(err, fmt.Errorf("Too many nameserver configs, max is 3"))
				break
			}
			// Search if this nameserver is already added
			if this.Find(o) != nil {
				err = multierror.Append(err, fmt.Errorf("Nameserver %s already exists in conf", opt))
				break
			}

			this.logger.Printf("Added nameserver %s", opt.IP)
			this.Nameservers = append(this.Nameservers, opt)
		case domain:
			this.logger.Printf("Added domain %s", opt.Name)
			this.Domain = o.(domain)
		case search:
			this.Search = o.(search)
		case searchDomain:
			// Search if this search domain is already added
			if this.Find(o) != nil {
				err = multierror.Append(err, fmt.Errorf("Search domain %s already exists in conf", opt))
				break
			}
			this.logger.Printf("Added search domain %s", opt.Name)
			this.Search.Domains = append(this.Search.Domains, opt)
		case sortlist:
			this.Sortlist = o.(sortlist)
		case sortlistpair:
			if i := this.Find(o); i != nil {
				err = multierror.Append(err, fmt.Errorf("Searchlist pair %s already exists in conf", opt))
				break
			}
			this.logger.Printf("Added sortlist pair %s", opt)
			this.Sortlist.Pairs = append(this.Sortlist.Pairs, opt)
		case []option:
			this.Options = o.([]option)
		case option:
			if _, e := parseOption(o.(option).String()); e != nil {
				err = multierror.Append(err, fmt.Errorf("Unknown option %s", o.(option)))
				break
			}
			if o := this.Find(opt); o != nil {
				err = multierror.Append(err, fmt.Errorf("Option %s is already present", o.(option).Type))
				break
			}
			this.logger.Printf("Added option %s", opt)
			this.Options = append(this.Options, opt)

		default:
			err = multierror.Append(err, fmt.Errorf("Unknown option type %v", opt))
		}
	}
	return err.ErrorOrNil()
}

// Remove items from the configuration
//
// Errors are accumulated and can be reinterpreted as an multierror type.
// Logging will occur if logging has been setup using the EnableLogging
// call
func (this *Conf) Remove(opts ...interface{}) error {
	var err *multierror.Error
	for _, o := range opts {
		i := this.indexOf(this.Find(o))
		_, isdom := o.(domain)
		if i == -1 && !isdom {
			err = multierror.Append(err, fmt.Errorf("Not found"))
			continue
		}

		switch opt := o.(type) {
		case nameserver:
			this.logger.Printf("Removed nameserver %s", opt.IP)
			this.Nameservers = append(this.Nameservers[:i], this.Nameservers[i+1:]...)
		case domain:
			this.logger.Printf("Removed domain %s", opt.Name)
			this.Domain = Domain("")
		case searchDomain:
			this.logger.Printf("Removed search domain %s", opt.Name)
			this.Search.Domains = append(this.Search.Domains[:i], this.Search.Domains[i+1:]...)
		case sortlistpair:
			this.logger.Printf("Removed sortlist pair %s", opt)
			this.Sortlist.Pairs = append(this.Sortlist.Pairs[:i], this.Sortlist.Pairs[i+1:]...)
		case option:
			this.logger.Printf("Removed option %s", opt)
			this.Options = append(this.Options[:i], this.Options[i+1:]...)
		default:
			err = multierror.Append(err, fmt.Errorf("Unknown option type %v", opt))
		}
	}
	return err.ErrorOrNil()
}

// EnableLogging enables internal logging with given writer as output, currently only one
// writer is supported. This will use LstdFlags for the logging
func (this *Conf) EnableLogging(writer ...io.Writer) error {
	if this.logger == nil {
		return fmt.Errorf("Logging has not been setup properly")
	}
	this.logger.SetFlags(log.LstdFlags)
	this.logger.SetOutput(writer[0])
	return nil
}

// Find an configure item returns nil if item is not found
func (this Conf) Find(o interface{}) interface{} {
	i := this.indexOf(o)
	if i == -1 {
		return nil
	}

	switch o.(type) {
	case nameserver:
		return this.Nameservers[i]
	case searchDomain:
		return this.Search.Domains[i]
	case sortlistpair:
		return this.Sortlist.Pairs[i]
	case option:
		return this.Options[i]
	}
	return nil
}

func (this Conf) indexOf(o interface{}) int {
	switch o.(type) {
	case nameserver:
		for i, item := range this.Nameservers {
			if item.IP.String() == o.(nameserver).IP.String() {
				return i
			}
		}
	case sortlistpair:
		for i, sp := range this.Sortlist.Pairs {
			if sp.Address.String() == o.(sortlistpair).Address.String() {
				return i
			}
		}
	case searchDomain:
		for i, sd := range this.Search.Domains {
			if sd.Name == o.(searchDomain).Name {
				return i
			}
		}
	case option:
		for i, item := range this.Options {
			if item.Type == o.(option).Type {
				return i
			}
		}
	}
	return -1
}
