package resolvconf_test

import (
	"."
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
)

func TestNewConf(t *testing.T) {
	conf := resolvconf.New()
	assert.NotNil(t, conf)
}

func TestAddNewNameserver(t *testing.T) {
	conf := resolvconf.New()
	ns := resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))
	err := conf.Add(ns)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Find(ns))
}

func TestRemoveNewNameserver(t *testing.T) {
	conf := resolvconf.New()
	ns := resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))
	conf.Add(ns)
	err := conf.Remove(ns)
	assert.Nil(t, err)
	assert.Nil(t, conf.Find(ns))
}

func TestRemoveNonExistingNewNameserver(t *testing.T) {
	conf := resolvconf.New()
	ip := net.ParseIP("8.8.8.8")
	err := conf.Remove(resolvconf.NewNameserver(ip))
	assert.NotNil(t, err)
}

func TestAddSecondDomainReplacesFirst(t *testing.T) {
	conf := resolvconf.New()
	conf.Add(resolvconf.NewDomain("foo.com"), resolvconf.NewDomain("bar.com"))
	assert.Equal(t, "bar.com", conf.Domain().Name)
}

func TestRemoveDomain(t *testing.T) {
	conf := resolvconf.New()
	foo := resolvconf.NewDomain("foo.com")
	conf.Add(foo)
	assert.Equal(t, "foo.com", conf.Domain().Name)
	conf.Remove(foo)
	assert.Equal(t, "", conf.Domain().Name)
}

func TestBasicSearchDomain(t *testing.T) {
	conf := resolvconf.New()
	dom := resolvconf.NewSearchDomain("foo.com")
	// Add a search domain
	err := conf.Add(dom)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Find(dom))

	// Test that search domain exists
	assert.NotNil(t, conf.Find(dom))

	// Add search domain again should yield error
	err = conf.Add(dom)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.Search()))

	// Remove search domain
	err = conf.Remove(dom)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.Search()))

	// Test that search domain does not exists
	assert.Nil(t, conf.Find(dom))

	// Remove non existing yields error
	err = conf.Remove(dom)
	assert.NotNil(t, err)
}

func TestBasicSortlist(t *testing.T) {
	conf := resolvconf.New()
	sp := resolvconf.NewSortlistPair(net.ParseIP("8.8.8.8")).SetNetmask(net.ParseIP("255.255.255.0"))

	// Add a pair
	err := conf.Add(sp)
	assert.Nil(t, err)
	assert.Equal(t, sp.Address.String(), conf.Sortlist()[0].Address.String())
	assert.Equal(t, sp.Netmask.String(), conf.Sortlist()[0].Netmask.String())

	// Check if pair exists
	assert.NotNil(t, conf.Find(*sp))

	// Add pair again should yield error
	err = conf.Add(sp)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.Sortlist()))

	// Remove sortlist pair
	err = conf.Remove(*sp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.Sortlist()))

	// Test that sortlistpair  does not exists
	assert.Nil(t, conf.Find(sp))

	// Remove non existing yields error
	err = conf.Remove(sp)
	assert.NotNil(t, err)
}

func TestNewOption(t *testing.T) {
	// New boolean option
	opt := resolvconf.NewOption("debug")
	assert.Equal(t, "debug", opt.Type)
	assert.Equal(t, -1, opt.Value)

	// New integer option
	opt = resolvconf.NewOption("ndots").Set(3)
	assert.Equal(t, "ndots", opt.Type)
	assert.Equal(t, 3, opt.Value)

	// Bad value
	opt = resolvconf.NewOption("ndots").Set(-3)
	assert.NotEqual(t, -3, opt.Value)

	// Unknown option
	//opt = resolvconf.NewOption("foo")
	//assert.Equal(t, "foo", opt.Type)
	//assert.Equal(t, -1, opt.Value)
}

func TestBasicNewOption(t *testing.T) {
	conf := resolvconf.New()

	// Test to set option
	opt := resolvconf.NewOption("debug")
	err := conf.Add(opt)
	assert.Nil(t, err)
	assert.Equal(t, "debug", conf.Options()[0].Type)
	assert.Equal(t, 1, len(conf.Options()))

	// Test if option is set
	o := conf.Find(*opt)
	assert.NotNil(t, o)

	// Test to set option again should yiled error
	err = conf.Add(opt)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.Options()))

	// Test to remove option
	err = conf.Remove(*opt)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.Options()))

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
	opt := resolvconf.NewOption("ndots").Set(4)
	err := conf.Add(opt)
	assert.Nil(t, err)
	assert.Equal(t, "ndots", conf.Options()[0].Type)
	assert.Equal(t, 4, conf.Options()[0].Value)
	assert.Equal(t, 1, len(conf.Options()))
}

func TestAddMultipleItems(t *testing.T) {
	conf := resolvconf.New()
	opt := resolvconf.NewOption("ndots").Set(4)
	ns := resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))
	err := conf.Add(opt, ns)
	assert.Nil(t, err)
	assert.Equal(t, "ndots", conf.Options()[0].Type)
	assert.Equal(t, 4, conf.Options()[0].Value)
	assert.NotNil(t, conf.Find(ns))
}

func TestAddItemsWithoutVariable(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")),
		resolvconf.NewOption("debug"))
	assert.Nil(t, err)
	assert.Equal(t, "debug", conf.Options()[0].Type)
	assert.NotNil(t, conf.Find(resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))))
}

func TestAddBadOptionInList(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")),
		resolvconf.NewOption("ndots").Set(-3),
		resolvconf.NewOption("debug"))

	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.Options()))
	assert.Equal(t, "debug", conf.Options()[0].Type)
	assert.NotNil(t, conf.Find(resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))))
}

func TestRemoveMultipleItems(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.NewOption("ndots").Set(4), resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.NotNil(t, conf.Find(resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))))
	assert.Equal(t, 1, len(conf.Options()))

	err = conf.Remove(*resolvconf.NewOption("ndots").Set(4), resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.Nil(t, conf.Find(resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))))
	assert.Equal(t, 0, len(conf.Options()))
}

func TestVariadicStorlistPair(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.NewSortlistPair(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.Equal(t, net.ParseIP("8.8.8.8"), conf.Sortlist()[0].Address)
	assert.Equal(t, net.ParseIP(""), conf.Sortlist()[0].Netmask)

	conf = resolvconf.New()
	err = conf.Add(resolvconf.NewSortlistPair(net.ParseIP("8.8.8.8")).SetNetmask(net.ParseIP("255.255.255.0")))
	assert.Nil(t, err)
	assert.Equal(t, net.ParseIP("8.8.8.8"), conf.Sortlist()[0].Address)
	assert.Equal(t, net.ParseIP("255.255.255.0"), conf.Sortlist()[0].Netmask)
}

func TestLogging(t *testing.T) {

	// Nothing is logged if not enabeled
	conf := resolvconf.New()
	buf := new(bytes.Buffer)
	err := conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.NotContains(t, buf.String(), fmt.Sprintf("Added Nameserver %s", net.ParseIP("8.8.8.8")))

	// Enable logging, test add Nameserver
	conf = resolvconf.New()
	buf.Reset()
	conf.EnableLogging(buf)
	assert.Nil(t, err)
	conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))
	assert.Contains(t, buf.String(), fmt.Sprintf("Added Nameserver %s", net.ParseIP("8.8.8.8")))

	// Enable logging, test remove Nameserver
	buf.Reset()
	conf.Remove(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))
	assert.Contains(t, buf.String(), fmt.Sprintf("Removed nameserver %s", net.ParseIP("8.8.8.8")))

	// Add & remove domain
	buf.Reset()
	conf.Add(resolvconf.NewDomain("foo.bar"))
	assert.Contains(t, buf.String(), "Added domain foo.bar")
	conf.Remove(resolvconf.NewDomain("foo.bar"))
	assert.Contains(t, buf.String(), "Removed domain foo.bar")

	// Add & remove search domain
	buf.Reset()
	conf.Add(resolvconf.NewSearchDomain("foo.bar"))
	assert.Contains(t, buf.String(), "Added search domain foo.bar")
	conf.Remove(resolvconf.NewSearchDomain("foo.bar"))
	assert.Contains(t, buf.String(), "Removed searchdomain foo.bar")

	// Add & remove sort list pair
	buf.Reset()
	conf.Add(resolvconf.NewSortlistPair(net.ParseIP("8.8.8.8")))
	assert.Contains(t, buf.String(), "Added sortlist pair 8.8.8.8")
	conf.Remove(*resolvconf.NewSortlistPair(net.ParseIP("8.8.8.8")))
	assert.Contains(t, buf.String(), "Removed sortlistpair 8.8.8.8")

	// Add & remove option
	buf.Reset()
	conf.Add(resolvconf.NewOption("debug"))
	assert.Contains(t, buf.String(), "Added option debug")
	conf.Remove(*resolvconf.NewOption("debug"))
	assert.Contains(t, buf.String(), "Removed option debug")
}

func ExampleConf_Add() {
	conf := resolvconf.New()
	conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))
	conf.Write(os.Stdout)
	// Output: nameserver 8.8.8.8
}

func ExampleConf_Add_second() {
	conf := resolvconf.New()
	conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")), resolvconf.NewOption("debug"))
	conf.Write(os.Stdout)
	// Output: nameserver 8.8.8.8
	//
	// options debug
}

func ExampleConf_Remove() {
	conf := resolvconf.New()
	ns := resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))
	conf.Add(ns, resolvconf.NewNameserver(net.ParseIP("8.8.8.9")))
	conf.Remove(ns)
	conf.Write(os.Stdout)
	// Output: nameserver 8.8.8.9
}

func Example() {
	res, err := http.Get("https://gist.githubusercontent.com/turadg/7876784/raw/c7f2500fa4762cfe443e30c64c6ed8a888f6ac74/resolv.conf")
	if err != nil {
		log.Fatal(err)
	}
	conf, err := resolvconf.ReadConf(res.Body)
	res.Body.Close()
	conf.Remove(resolvconf.NewNameserver(net.ParseIP("8.8.4.4")))
	conf.Add(resolvconf.NewDomain("foo.bar"), resolvconf.NewSortlistPair(net.ParseIP("130.155.160.0")).SetNetmask(net.ParseIP("255.255.240.0")))
	conf.Write(os.Stdout)
	// Output: domain foo.bar
	// nameserver 8.8.8.8
	//
	// sortlist 130.155.160.0/255.255.240.0
}
