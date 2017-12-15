package resolvconf

import (
	"fmt"
	"net"
)

// SortItem is one of the items in the sort list, it must have an address and
// may have an netmask
type SortItem struct {
	Address net.IP
	Netmask net.IP
}

// NewSortItem creates a new sortlist that will be added to the
// 'sort' item in the resolv.conf file.
// If mask is given the output will be IP/mask e.g.
// 8.8.8.8/255.255.255.0 otherwise output will be
// IP only, e.g. 8.8.8.8
func NewSortItem(addr net.IP) *SortItem {
	slp := new(SortItem)
	slp.Address = addr
	return slp
}

func (si SortItem) applyLimits(conf *Conf) (bool, error) {
	if i := conf.Find(si); i != nil {
		// Check if netmask is different otherwise error
		if si.Netmask.Equal((*i).(*SortItem).Netmask) {
			return false, fmt.Errorf("Sortlist pair %s already exists in conf", si)
		}
		(*i).(*SortItem).Netmask = si.Netmask
	}
	if len(conf.GetSortItems()) == sortListMaxCount {
		return false, fmt.Errorf("Too long sortlist, %d is maximum", sortListMaxCount)
	}
	return true, nil
}

// Equal compares two SortItems, return true if equal
func (si SortItem) Equal(b ConfItem) bool {
	if item, ok := b.(*SortItem); ok {
		return si.Address.String() == item.Address.String()
	}

	return false
}

// SetNetmask sets the netmask for an SortItem
func (si *SortItem) SetNetmask(nm net.IP) *SortItem {
	si.Netmask = nm
	return si
}

// GetNetmask returns netmask from an SortItems
func (si SortItem) GetNetmask() net.IP {
	return si.Netmask
}

func (si SortItem) String() string {
	if len(si.Netmask) > 0 {
		return fmt.Sprintf("%s/%s", si.Address, si.Netmask)
	}
	return fmt.Sprintf("%s", si.Address)
}
