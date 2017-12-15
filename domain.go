package resolvconf

// Domain is the single domain in a resolv.conf file
type Domain struct {
	Name string
}

// NewDomain creates a new domain that will be used
// as value for the 'domain' option in the generated file
func NewDomain(dom string) *Domain {
	return &Domain{dom}
}

func (dom Domain) applyLimits(conf *Conf) (bool, error) {
	i := conf.indexOf(conf.GetDomain())
	if i != -1 {
		// Found it, update and return not ok to add
		conf.items[i] = &Domain{dom.Name}
		return false, nil
	}

	// Ok to add
	return true, nil
}

func (dom Domain) String() string {
	return dom.Name
}

// Equal compares two domains with each other, returns true if equal
func (dom Domain) Equal(b ConfItem) bool {
	if item, ok := b.(*Domain); ok {
		return dom.Name == item.Name
	}
	return false
}
