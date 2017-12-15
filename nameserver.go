package resolvconf

import (
	"fmt"
	"net"
)

// Nameserver is the nameserver type
type Nameserver struct {
	IP net.IP // IP address
}

// NewNameserver creates a new Nameserver item
func NewNameserver(IP net.IP) *Nameserver {
	return &Nameserver{IP}
}

func (ns Nameserver) applyLimits(conf *Conf) (bool, error) {

	if len(conf.GetNameservers()) == nameserversMaxCount {
		return false, fmt.Errorf("Too many Nameserver configs, max is %d", nameserversMaxCount)
	}
	// Search if conf Nameserver is already added
	if conf.Find(ns) != nil {
		return false, fmt.Errorf("Nameserver %s already exists in conf", ns.IP)
	}

	return true, nil
}

// Equal compares to nameservers with eachother, returns true if equal
func (ns Nameserver) Equal(b ConfItem) bool {
	if item, ok := b.(*Nameserver); ok {
		return ns.IP.Equal(item.IP)
	}
	return false
}

func (ns Nameserver) String() string {
	return ns.IP.String()
}
