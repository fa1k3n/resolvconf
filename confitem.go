package resolvconf

import (
	"fmt"
)

// ConfItem is the generic interface all items in a resolv.conf file
// must implement.
type ConfItem interface {
	fmt.Stringer
	applyLimits(conf *Conf) (bool, error)
	Equal(b ConfItem) bool
}
