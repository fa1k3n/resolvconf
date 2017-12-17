package resolvconf

import (
	"fmt"
)

// Option represents an option item which must have a Type
// and some options must have a value
type Option struct {
	Type  string
	Value int
}

// NewOption creates a new option, val must be a positive number if used.
// Witout val the option will be interpreted as a bolean e.g.
// debug , with a val the option will be interpreted as an
// setvalue, e.g. ndots:5
func NewOption(t string) *Option {
	// Check va
	opt := new(Option)
	opt.Type = t
	opt.Value = -1

	if _, err := parseOption(opt.String()); err != nil {
		return nil
	}

	return opt
}

func (opt *Option) applyLimits(conf *Conf) (bool, error) {
	if opt.Type == "ndots" && opt.Value < 0 {
		return false, fmt.Errorf("Bad value %d", opt.Value)
	}
	if _, e := parseOption(opt.String()); e != nil {
		return false, fmt.Errorf("Unknown option %s", opt)
	}
	// Check limits
	newVal := -1
	switch opt.Type {
	case "ndots":
		if opt.Value > optionNdotsMax {
			newVal = optionNdotsMax
		}
	case "timeout":
		if opt.Value > optionTimeoutMax {
			newVal = optionTimeoutMax
		}
	case "attempts":
		if opt.Value > optionAttemptsMax {
			newVal = optionAttemptsMax
		}
	}
	if newVal > -1 {
		conf.logger.Printf("[WARN] Option %s is capped to %d, set value is %d", opt.Type, newVal, opt.Value)
		opt.Value = newVal
	}
	if o := conf.Find(opt); o != nil {
		// If option has a value then update otherwise error
		if o.(*Option).Value != -1 {
			i := conf.indexOf(o)
			conf.items[i].(*Option).Value = opt.Value
			return false, nil // Dont add
		}
		return false, fmt.Errorf("Option %s is already present", opt)
	}
	return true, nil
}

// Equal compares two Option, return true if equal
func (opt Option) Equal(b ConfItem) bool {
	if o, ok := b.(*Option); ok {
		return opt.Type == o.Type
	}
	return false
}

// Set sets the value of an option
func (opt *Option) Set(value int) *Option {
	if value < 0 {
		return opt
	}
	opt.Value = value
	return opt
}

// Get returns the option value
func (opt Option) Get() int {
	return opt.Value
}

func (opt Option) String() string {
	switch opt.Type {
	case "debug", "rotate", "no-check-names", "inet6",
		"ip6-bytestring", "ip6-dotint", "no-ip6-dotint",
		"edns0", "single-request", "single-request-reopen",
		"no-tld-query", "use-vc":
		return fmt.Sprintf("%s", opt.Type)
	case "ndots", "timeout", "attempts":
		return fmt.Sprintf("%s:%d", opt.Type, opt.Value)
	}
	return ""
}
