// Package resolvconf provides an interface to read, create and manipulate
// resolv.conf files
package resolvconf

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"io"
	"log"
	"reflect"
	"strings"
)

// Add items to the configuration.
//
// Errors are accumulated and can be reinterpreted as
// an multierror type. Logging will occur if logging has
// been setup using the EnableLogging call
func (conf *Conf) Add(opts ...ConfItem) error {
	var err *multierror.Error
	for _, o := range opts {
		if o == nil {
			err = multierror.Append(err, fmt.Errorf("Trying to add nil element"))
			continue
		}
		if ok, e := o.applyLimits(conf); e != nil {
			err = multierror.Append(err, e)
		} else if ok {
			typeName := reflect.TypeOf(o).Elem().Name()
			conf.logger.Printf("Added %s %s", strings.ToLower(typeName), o)
			conf.items = append(conf.items, o)
		}
	}
	return err.ErrorOrNil()
}

// Remove items from the configuration
//
// Errors are accumulated and can be reinterpreted as an multierror type.
// Logging will occur if logging has been setup using the EnableLogging
// call
func (conf *Conf) Remove(opts ...ConfItem) error {
	var err *multierror.Error
	for _, o := range opts {
		if o == nil {
			err = multierror.Append(err, fmt.Errorf("Trying to remove nil element"))
			continue
		}
		i := conf.indexOf(o)
		//_, isdom := o.(Domain)
		if i == -1 {
			err = multierror.Append(err, fmt.Errorf("Not found"))
			continue
		}
		typeName := reflect.TypeOf(conf.items[i]).Elem().Name()
		conf.logger.Printf("Removed %s %s", strings.ToLower(typeName), conf.items[i])
		conf.items = append(conf.items[:i], conf.items[i+1:]...)
	}
	return err.ErrorOrNil()
}

// EnableLogging enables internal logging with given writer as output, currently only one
// writer is supported. conf will use LstdFlags for the logging
func (conf *Conf) EnableLogging(writer io.Writer) error {
	if conf.logger == nil {
		return fmt.Errorf("Logging has not been setup properly")
	}
	conf.logger.SetFlags(log.LstdFlags)
	conf.logger.SetOutput(writer)
	return nil
}

// Find an configure item returns nil if item is not found. Returned will be
// a pointer to the actual item that can be converted into expected type
func (conf Conf) Find(o ConfItem) ConfItem {
	i := conf.indexOf(o)
	if i == -1 {
		return nil
	}
	return conf.items[i]
}

func (conf Conf) indexOf(o ConfItem) int {
	for i, item := range conf.items {
		if o.Equal(item) {
			return i
		}
	}
	return -1
}
