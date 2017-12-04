package resolvconf_test

import (
	"."
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestNewConf(t *testing.T) {
	conf := resolvconf.New()
	assert.NotNil(t, conf)
}

func TestAddNameserver(t *testing.T) {
	conf := resolvconf.New()
	ns, _ := resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))
	err := conf.Add(ns)
	assert.Nil(t, err)
	assert.Equal(t, "8.8.8.8", conf.Nameservers[0].IP.String())
	assert.NotNil(t, conf.Find(ns))
}

func TestRemoveNameserver(t *testing.T) {
	conf := resolvconf.New()
	ns, _ := resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))
	conf.Add(ns)
	err := conf.Remove(ns)
	assert.Nil(t, err)
	assert.Nil(t, conf.Find(ns))
}

func TestRemoveNonExistingNameserver(t *testing.T) {
	conf := resolvconf.New()
	ip := net.ParseIP("8.8.8.8")
	err := conf.Remove(resolvconf.Nameserver{ip})
	assert.NotNil(t, err)
}

func TestAddSecondDomainReplacesFirst(t *testing.T) {
	conf := resolvconf.New()
	foo, _ := resolvconf.NewDomain("foo.com")
	bar, _ := resolvconf.NewDomain("bar.com")
	conf.Add(foo)
	conf.Add(bar)
	assert.Equal(t, "bar.com", conf.Domain.Name)
}

func TestRemoveDomain(t *testing.T) {
	conf := resolvconf.New()
	foo, _ := resolvconf.NewDomain("foo.com")
	conf.Add(foo)
	assert.Equal(t, "foo.com", conf.Domain.Name)
	conf.Remove(foo)
	assert.Equal(t, "", conf.Domain.Name)
}

func TestBasicSearchDomain(t *testing.T) {
	conf := resolvconf.New()
	dom, _ := resolvconf.NewSearchDomain("foo.com")
	// Add a search domain
	err := conf.Add(dom)
	assert.Nil(t, err)
	assert.Equal(t, dom.Name, conf.Search.Domains[0].Name)

	// Test that search domain exists
	assert.NotNil(t, conf.Find(dom))

	// Add search domain again should yield error
	err = conf.Add(dom)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.Search.Domains))

	// Remove search domain
	err = conf.Remove(dom)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.Search.Domains))

	// Test that search domain does not exists
	assert.Nil(t, conf.Find(dom))

	// Remove non existing yields error
	err = conf.Remove(dom)
	assert.NotNil(t, err)
}

func TestBasicSortlist(t *testing.T) {
	conf := resolvconf.New()
	sp, _ := resolvconf.NewSortlistPair(net.ParseIP("8.8.8.8"), net.ParseIP("255.255.255.0"))

	// Add a pair
	err := conf.Add(sp)
	assert.Nil(t, err)
	assert.Equal(t, sp.Address.String(), conf.Sortlist.Pairs[0].Address.String())
	assert.Equal(t, sp.Netmask.String(), conf.Sortlist.Pairs[0].Netmask.String())

	// Check if pair exists
	assert.NotNil(t, conf.Find(sp))

	// Add pair again should yield error
	err = conf.Add(sp)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.Sortlist.Pairs))

	// Remove sortlist pair
	err = conf.Remove(sp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.Sortlist.Pairs))

	// Test that sortlistpair  does not exists
	assert.Nil(t, conf.Find(sp))

	// Remove non existing yields error
	err = conf.Remove(sp)
	assert.NotNil(t, err)
}

func TestNewOption(t *testing.T) {
	// New boolean option
	opt, err := resolvconf.NewOption("debug")
	assert.Nil(t, err)
	assert.Equal(t, "debug", opt.Type)
	assert.Equal(t, -1, opt.Value)

	// New integer option
	opt, err = resolvconf.NewOption("ndots", 3)
	assert.Nil(t, err)
	assert.Equal(t, "ndots", opt.Type)
	assert.Equal(t, 3, opt.Value)

	// Too many values
	opt, err = resolvconf.NewOption("ndots", 3, 4)
	assert.NotNil(t, err)
	assert.Equal(t, "", opt.Type)
	assert.Equal(t, -1, opt.Value)

	// Bad value
	opt, err = resolvconf.NewOption("ndots", -3)
	assert.NotNil(t, err)

	// Unknown option
	opt, err = resolvconf.NewOption("foo")
	assert.NotNil(t, err)
	assert.Equal(t, "", opt.Type)
	assert.Equal(t, -1, opt.Value)
}

func TestBasicOption(t *testing.T) {
	conf := resolvconf.New()

	// Test to set option
	opt, _ := resolvconf.NewOption("debug")
	err := conf.Add(opt)
	assert.Nil(t, err)
	assert.Equal(t, "debug", conf.Options[0].Type)
	assert.Equal(t, 1, len(conf.Options))

	// Test if option is set
	o := conf.Find(opt)
	assert.NotNil(t, o)

	// Test to set option again should yiled error
	err = conf.Add(opt)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.Options))

	// Test to remove option
	err = conf.Remove(opt)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.Options))

	// Test that option is not set
	o = conf.Find(opt)
	assert.Nil(t, o)

	// Remove non existing option
	err = conf.Remove(opt)
	assert.NotNil(t, err)
}

func TestOptionWithValue(t *testing.T) {
	conf := resolvconf.New()

	// Test to set option
	opt, _ := resolvconf.NewOption("ndots", 4)
	err := conf.Add(opt)
	assert.Nil(t, err)
	assert.Equal(t, "ndots", conf.Options[0].Type)
	assert.Equal(t, 4, conf.Options[0].Value)
	assert.Equal(t, 1, len(conf.Options))
}

func TestAddMultipleItems(t *testing.T) {
	conf := resolvconf.New()
	opt, _ := resolvconf.NewOption("ndots", 4)
	ns, _ := resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))
	err := conf.Add(opt, ns)
	assert.Nil(t, err)
	assert.Equal(t, "ndots", conf.Options[0].Type)
	assert.Equal(t, 4, conf.Options[0].Value)
	assert.Equal(t, "8.8.8.8", conf.Nameservers[0].IP.String())
}
