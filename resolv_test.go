package resolvconf_test

import (
	"."
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
)

func TestNewConf(t *testing.T) {
	conf := resolvconf.New()
	assert.NotNil(t, conf)
}

func TestAddNameserver(t *testing.T) {
	conf := resolvconf.New()
	ns := resolvconf.Nameserver(net.ParseIP("8.8.8.8"))
	err := conf.Add(ns)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Find(ns))
}

func TestRemoveNameserver(t *testing.T) {
	conf := resolvconf.New()
	ns := resolvconf.Nameserver(net.ParseIP("8.8.8.8"))
	conf.Add(ns)
	err := conf.Remove(ns)
	assert.Nil(t, err)
	assert.Nil(t, conf.Find(ns))
}

func TestRemoveNonExistingNameserver(t *testing.T) {
	conf := resolvconf.New()
	ip := net.ParseIP("8.8.8.8")
	err := conf.Remove(resolvconf.Nameserver(ip))
	assert.NotNil(t, err)
}

func TestAddSecondDomainReplacesFirst(t *testing.T) {
	conf := resolvconf.New()
	foo := resolvconf.Domain("foo.com")
	bar := resolvconf.Domain("bar.com")
	conf.Add(foo)
	conf.Add(bar)
	assert.Equal(t, "bar.com", conf.Domain.Name)
}

func TestRemoveDomain(t *testing.T) {
	conf := resolvconf.New()
	foo := resolvconf.Domain("foo.com")
	conf.Add(foo)
	assert.Equal(t, "foo.com", conf.Domain.Name)
	conf.Remove(foo)
	assert.Equal(t, "", conf.Domain.Name)
}

func TestBasicSearchDomain(t *testing.T) {
	conf := resolvconf.New()
	dom := resolvconf.SearchDomain("foo.com")
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
	sp := resolvconf.SortlistPair(net.ParseIP("8.8.8.8"), net.ParseIP("255.255.255.0"))

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

func TestOption(t *testing.T) {
	// New boolean option
	opt := resolvconf.Option("debug")
	assert.Equal(t, "debug", opt.Type)
	assert.Equal(t, -1, opt.Value)

	// New integer option
	opt = resolvconf.Option("ndots", 3)
	assert.Equal(t, "ndots", opt.Type)
	assert.Equal(t, 3, opt.Value)

	// Too many values
	opt = resolvconf.Option("ndots", 3, 4)
	assert.Equal(t, "", opt.Type)
	assert.Equal(t, -1, opt.Value)

	// Bad value
	opt = resolvconf.Option("ndots", -3)
	assert.Equal(t, "", opt.Type)
	assert.Equal(t, -1, opt.Value)

	// Unknown option
	opt = resolvconf.Option("foo")
	assert.Equal(t, "", opt.Type)
	assert.Equal(t, -1, opt.Value)
}

func TestBasicOption(t *testing.T) {
	conf := resolvconf.New()

	// Test to set option
	opt := resolvconf.Option("debug")
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
	opt := resolvconf.Option("ndots", 4)
	err := conf.Add(opt)
	assert.Nil(t, err)
	assert.Equal(t, "ndots", conf.Options[0].Type)
	assert.Equal(t, 4, conf.Options[0].Value)
	assert.Equal(t, 1, len(conf.Options))
}

func TestAddMultipleItems(t *testing.T) {
	conf := resolvconf.New()
	opt := resolvconf.Option("ndots", 4)
	ns := resolvconf.Nameserver(net.ParseIP("8.8.8.8"))
	err := conf.Add(opt, ns)
	assert.Nil(t, err)
	assert.Equal(t, "ndots", conf.Options[0].Type)
	assert.Equal(t, 4, conf.Options[0].Value)
	assert.NotNil(t, conf.Find(ns))
}

func TestAddItemsWithoutVariable(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.Nameserver(net.ParseIP("8.8.8.8")),
		resolvconf.Option("debug"))
	assert.Nil(t, err)
	assert.Equal(t, "debug", conf.Options[0].Type)
	assert.NotNil(t, conf.Find(resolvconf.Nameserver(net.ParseIP("8.8.8.8"))))
}

func TestAddBadOptionInList(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.Nameserver(net.ParseIP("8.8.8.8")),
		resolvconf.Option("ndots", -3),
		resolvconf.Option("debug"))

	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.Options))
	assert.Equal(t, "debug", conf.Options[0].Type)
	assert.NotNil(t, conf.Find(resolvconf.Nameserver(net.ParseIP("8.8.8.8"))))
}

func TestRemoveMultipleItems(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.Option("ndots", 4), resolvconf.Nameserver(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.NotNil(t, conf.Find(resolvconf.Nameserver(net.ParseIP("8.8.8.8"))))
	assert.Equal(t, 1, len(conf.Options))

	err = conf.Remove(resolvconf.Option("ndots", 4), resolvconf.Nameserver(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.Nil(t, conf.Find(resolvconf.Nameserver(net.ParseIP("8.8.8.8"))))
	assert.Equal(t, 0, len(conf.Options))
}

func TestVariadicStorlistPair(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.SortlistPair(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.Equal(t, net.ParseIP("8.8.8.8"), conf.Sortlist.Pairs[0].Address)
	assert.Equal(t, net.ParseIP(""), conf.Sortlist.Pairs[0].Netmask)

	conf = resolvconf.New()
	err = conf.Add(resolvconf.SortlistPair(net.ParseIP("8.8.8.8"), net.ParseIP("255.255.255.0")))
	assert.Nil(t, err)
	assert.Equal(t, net.ParseIP("8.8.8.8"), conf.Sortlist.Pairs[0].Address)
	assert.Equal(t, net.ParseIP("255.255.255.0"), conf.Sortlist.Pairs[0].Netmask)
}

func TestLogging(t *testing.T) {

	// Nothing is logged if not enabeled
	conf := resolvconf.New()
	buf := new(bytes.Buffer)
	err := conf.Add(resolvconf.Nameserver(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.NotContains(t, buf.String(), fmt.Sprintf("Added nameserver %s", net.ParseIP("8.8.8.8")))

	// Enable logging, test add nameserver
	conf = resolvconf.New()
	buf.Reset()
	conf.EnableLogging(buf)
	assert.Nil(t, err)
	conf.Add(resolvconf.Nameserver(net.ParseIP("8.8.8.8")))
	assert.Contains(t, buf.String(), fmt.Sprintf("Added nameserver %s", net.ParseIP("8.8.8.8")))

	// Enable logging, test remove nameserver
	buf.Reset()
	conf.Remove(resolvconf.Nameserver(net.ParseIP("8.8.8.8")))
	assert.Contains(t, buf.String(), fmt.Sprintf("Removed nameserver %s", net.ParseIP("8.8.8.8")))

	// Add & remove domain
	buf.Reset()
	conf.Add(resolvconf.Domain("foo.bar"))
	assert.Contains(t, buf.String(), "Added domain foo.bar")
	conf.Remove(resolvconf.Domain("foo.bar"))
	assert.Contains(t, buf.String(), "Removed domain foo.bar")

	// Add & remove search domain
	buf.Reset()
	conf.Add(resolvconf.SearchDomain("foo.bar"))
	assert.Contains(t, buf.String(), "Added search domain foo.bar")
	conf.Remove(resolvconf.SearchDomain("foo.bar"))
	assert.Contains(t, buf.String(), "Removed search domain foo.bar")

	// Add & remove sort list pair
	buf.Reset()
	conf.Add(resolvconf.SortlistPair(net.ParseIP("8.8.8.8")))
	assert.Contains(t, buf.String(), "Added sortlist pair 8.8.8.8")
	conf.Remove(resolvconf.SortlistPair(net.ParseIP("8.8.8.8")))
	assert.Contains(t, buf.String(), "Removed sortlist pair 8.8.8.8")

	// Add & remove option
	buf.Reset()
	conf.Add(resolvconf.Option("debug"))
	assert.Contains(t, buf.String(), "Added option debug")
	conf.Remove(resolvconf.Option("debug"))
	assert.Contains(t, buf.String(), "Removed option debug")
}

func ExampleConf_Add() {
	conf := resolvconf.New()
	conf.Add(resolvconf.Nameserver(net.ParseIP("8.8.8.8")))
	conf.Write(os.Stdout)
	// Output: nameserver 8.8.8.8
}
