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
func (conf *Conf) Add(opts ...interface{}) error {
	var err *multierror.Error
	for _, o := range opts {
		switch opt := o.(type) {
		case nameserver:
			if len(conf.Nameservers)+1 > 3 {
				err = multierror.Append(err, fmt.Errorf("Too many nameserver configs, max is 3"))
				break
			}
			// Search if conf nameserver is already added
			if conf.Find(o) != nil {
				err = multierror.Append(err, fmt.Errorf("Nameserver %s already exists in conf", opt))
				break
			}

			conf.logger.Printf("Added nameserver %s", opt.IP)
			conf.Nameservers = append(conf.Nameservers, opt)
		case domain:
			conf.logger.Printf("Added domain %s", opt.Name)
			conf.Domain = o.(domain)
		case search:
			conf.Search = o.(search)
		case searchDomain:
			// Search if conf search domain is already added
			if conf.Find(o) != nil {
				err = multierror.Append(err, fmt.Errorf("Search domain %s already exists in conf", opt))
				break
			}
			conf.logger.Printf("Added search domain %s", opt.Name)
			conf.Search.Domains = append(conf.Search.Domains, opt)
		case sortlist:
			conf.Sortlist = o.(sortlist)
		case sortlistpair:
			if i := conf.Find(o); i != nil {
				err = multierror.Append(err, fmt.Errorf("Searchlist pair %s already exists in conf", opt))
				break
			}
			conf.logger.Printf("Added sortlist pair %s", opt)
			conf.Sortlist.Pairs = append(conf.Sortlist.Pairs, opt)
		case []option:
			conf.Options = o.([]option)
		case option:
			if _, e := parseOption(o.(option).String()); e != nil {
				err = multierror.Append(err, fmt.Errorf("Unknown option %s", o.(option)))
				break
			}
			if o := conf.Find(opt); o != nil {
				err = multierror.Append(err, fmt.Errorf("Option %s is already present", o.(option).Type))
				break
			}
			conf.logger.Printf("Added option %s", opt)
			conf.Options = append(conf.Options, opt)

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
func (conf *Conf) Remove(opts ...interface{}) error {
	var err *multierror.Error
	for _, o := range opts {
		i := conf.indexOf(conf.Find(o))
		_, isdom := o.(domain)
		if i == -1 && !isdom {
			err = multierror.Append(err, fmt.Errorf("Not found"))
			continue
		}

		switch opt := o.(type) {
		case nameserver:
			conf.logger.Printf("Removed nameserver %s", opt.IP)
			conf.Nameservers = append(conf.Nameservers[:i], conf.Nameservers[i+1:]...)
		case domain:
			conf.logger.Printf("Removed domain %s", opt.Name)
			conf.Domain = Domain("")
		case searchDomain:
			conf.logger.Printf("Removed search domain %s", opt.Name)
			conf.Search.Domains = append(conf.Search.Domains[:i], conf.Search.Domains[i+1:]...)
		case sortlistpair:
			conf.logger.Printf("Removed sortlist pair %s", opt)
			conf.Sortlist.Pairs = append(conf.Sortlist.Pairs[:i], conf.Sortlist.Pairs[i+1:]...)
		case option:
			conf.logger.Printf("Removed option %s", opt)
			conf.Options = append(conf.Options[:i], conf.Options[i+1:]...)
		default:
			err = multierror.Append(err, fmt.Errorf("Unknown option type %v", opt))
		}
	}
	return err.ErrorOrNil()
}

// EnableLogging enables internal logging with given writer as output, currently only one
// writer is supported. conf will use LstdFlags for the logging
func (conf *Conf) EnableLogging(writer ...io.Writer) error {
	if conf.logger == nil {
		return fmt.Errorf("Logging has not been setup properly")
	}
	conf.logger.SetFlags(log.LstdFlags)
	conf.logger.SetOutput(writer[0])
	return nil
}

// Find an configure item returns nil if item is not found
func (conf Conf) Find(o interface{}) interface{} {
	i := conf.indexOf(o)
	if i == -1 {
		return nil
	}

	switch o.(type) {
	case nameserver:
		return conf.Nameservers[i]
	case searchDomain:
		return conf.Search.Domains[i]
	case sortlistpair:
		return conf.Sortlist.Pairs[i]
	case option:
		return conf.Options[i]
	}
	return nil
}

func (conf Conf) indexOf(o interface{}) int {
	switch o.(type) {
	case nameserver:
		for i, item := range conf.Nameservers {
			if item.IP.String() == o.(nameserver).IP.String() {
				return i
			}
		}
	case sortlistpair:
		for i, sp := range conf.Sortlist.Pairs {
			if sp.Address.String() == o.(sortlistpair).Address.String() {
				return i
			}
		}
	case searchDomain:
		for i, sd := range conf.Search.Domains {
			if sd.Name == o.(searchDomain).Name {
				return i
			}
		}
	case option:
		for i, item := range conf.Options {
			if item.Type == o.(option).Type {
				return i
			}
		}
	}
	return -1
}
