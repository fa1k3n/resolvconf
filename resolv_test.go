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
	"strconv"
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

func TestIPv6Nameserver(t *testing.T) {
	conf := resolvconf.New()
	ns := resolvconf.NewNameserver(net.ParseIP("2001:0db8:0000:0000:0000:0000:1428:07ab"))
	err := conf.Add(ns)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Find(ns))
	assert.Equal(t, net.ParseIP("2001:0db8:0000:0000:0000:0000:1428:07ab"), conf.Find(ns).(*resolvconf.Nameserver).IP)
}

func TestAddSecondDomainReplacesFirst(t *testing.T) {
	conf := resolvconf.New()
	conf.Add(resolvconf.NewDomain("foo.com"), resolvconf.NewDomain("bar.com"))
	assert.Equal(t, "bar.com", conf.GetDomain().Name)
}

func TestRemoveDomain(t *testing.T) {
	conf := resolvconf.New()
	foo := resolvconf.NewDomain("foo.com")
	conf.Add(foo)
	assert.Equal(t, "foo.com", conf.GetDomain().Name)
	conf.Remove(foo)
	assert.Equal(t, "", conf.GetDomain().Name)
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
	assert.Equal(t, 1, len(conf.GetSearchDomains()))

	// Remove search domain
	err = conf.Remove(dom)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.GetSearchDomains()))

	// Test that search domain does not exists
	assert.Nil(t, conf.Find(dom))

	// Remove non existing yields error
	err = conf.Remove(dom)
	assert.NotNil(t, err)
}

func TestBasicSortlist(t *testing.T) {
	conf := resolvconf.New()
	sp := resolvconf.NewSortItem(net.ParseIP("8.8.8.8")).SetNetmask(net.ParseIP("255.255.255.0"))

	// Add a pair
	err := conf.Add(sp)
	assert.Nil(t, err)
	assert.Equal(t, sp.Address.String(), conf.GetSortItems()[0].Address.String())
	assert.Equal(t, sp.Netmask.String(), conf.GetSortItems()[0].Netmask.String())

	// Check if pair exists
	assert.NotNil(t, conf.Find(sp))

	// Add pair again should yield error
	err = conf.Add(sp)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.GetSortItems()))

	// Remove sortlist pair
	err = conf.Remove(sp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.GetSortItems()))

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
	opt = resolvconf.NewOption("foo")
	assert.Nil(t, opt)
}

func TestBasicNewOption(t *testing.T) {
	conf := resolvconf.New()

	// Test to set option
	opt := resolvconf.NewOption("debug")
	err := conf.Add(opt)
	assert.Nil(t, err)
	assert.Equal(t, "debug", conf.GetOptions()[0].Type)
	assert.Equal(t, 1, len(conf.GetOptions()))

	// Test if option is set
	o := conf.Find(opt)
	assert.NotNil(t, o)

	// Test to set option again should yiled error
	err = conf.Add(opt)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.GetOptions()))

	// Test to remove option
	err = conf.Remove(opt)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.GetOptions()))

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
	assert.Equal(t, "ndots", conf.GetOptions()[0].Type)
	assert.Equal(t, 4, conf.GetOptions()[0].Value)
	assert.Equal(t, 1, len(conf.GetOptions()))
}

func TestAddMultipleItems(t *testing.T) {
	conf := resolvconf.New()
	opt := resolvconf.NewOption("ndots").Set(4)
	ns := resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))
	err := conf.Add(opt, ns)
	assert.Nil(t, err)
	assert.Equal(t, "ndots", conf.GetOptions()[0].Type)
	assert.Equal(t, 4, conf.GetOptions()[0].Value)
	assert.NotNil(t, conf.Find(ns))
}

func TestAddItemsWithoutVariable(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")),
		resolvconf.NewOption("debug"))
	assert.Nil(t, err)
	assert.Equal(t, "debug", conf.GetOptions()[0].Type)
	assert.NotNil(t, conf.Find(resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))))
}

func TestAddBadOptionInList(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")),
		resolvconf.NewOption("ndots").Set(-3),
		resolvconf.NewOption("debug"))

	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.GetOptions()))
	assert.Equal(t, "debug", conf.GetOptions()[0].Type)
	assert.NotNil(t, conf.Find(resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))))
}

func TestRemoveMultipleItems(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.NewOption("ndots").Set(4), resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.NotNil(t, conf.Find(resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))))
	assert.Equal(t, 1, len(conf.GetOptions()))

	err = conf.Remove(resolvconf.NewOption("ndots").Set(4), resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.Nil(t, conf.Find(resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))))
	assert.Equal(t, 0, len(conf.GetOptions()))
}

func TestVariadicStorlistPair(t *testing.T) {
	conf := resolvconf.New()
	err := conf.Add(resolvconf.NewSortItem(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.Equal(t, net.ParseIP("8.8.8.8"), conf.GetSortItems()[0].Address)
	assert.Equal(t, net.ParseIP(""), conf.GetSortItems()[0].Netmask)

	conf = resolvconf.New()
	err = conf.Add(resolvconf.NewSortItem(net.ParseIP("8.8.8.8")).SetNetmask(net.ParseIP("255.255.255.0")))
	assert.Nil(t, err)
	assert.Equal(t, net.ParseIP("8.8.8.8"), conf.GetSortItems()[0].Address)
	assert.Equal(t, net.ParseIP("255.255.255.0"), conf.GetSortItems()[0].Netmask)
}

func TestThatOptionsWithValueUpdatesExistingItems(t *testing.T) {
	conf := resolvconf.New()

	// Start with ndots, timeout and attempts, e.g. the once that should work
	for i, opt := range []string{"ndots", "timeout", "attempts"} {
		conf.Add(resolvconf.NewOption(opt).Set(3))
		assert.Equal(t, 3, conf.GetOptions()[i].Get())
		err := conf.Add(resolvconf.NewOption(opt).Set(5))
		assert.Nil(t, err)
		assert.Equal(t, 5, conf.GetOptions()[i].Get())
	}

	// Now test with one that should not work
	err := conf.Add(resolvconf.NewOption("debug"))
	assert.Nil(t, err)
	assert.NotNil(t, conf.Find(resolvconf.NewOption("debug")))
	err = conf.Add(resolvconf.NewOption("debug"))
	assert.NotNil(t, err)
}

func TestThatSortItemWithDifferentNetmaskToSortItemUpdatesItem(t *testing.T) {
	conf := resolvconf.New()
	conf.Add(resolvconf.NewSortItem(net.ParseIP("130.155.160.0")).SetNetmask(net.ParseIP("255.255.240.0")))
	si := conf.Find(resolvconf.NewSortItem(net.ParseIP("130.155.160.0")))
	assert.NotNil(t, si)
	assert.Equal(t, net.ParseIP("255.255.240.0"), si.(*resolvconf.SortItem).GetNetmask())

	err := conf.Add(resolvconf.NewSortItem(net.ParseIP("130.155.160.0")).SetNetmask(net.ParseIP("255.255.240.100")))
	assert.Nil(t, err)
	si = conf.Find(resolvconf.NewSortItem(net.ParseIP("130.155.160.0")))
	assert.NotNil(t, si)
	assert.Equal(t, net.ParseIP("255.255.240.100"), si.(*resolvconf.SortItem).GetNetmask())
}

func TestSearchDomainLimit(t *testing.T) {
	conf := resolvconf.New()
	for i := 0; i < 6; i++ {
		err := conf.Add(resolvconf.NewSearchDomain("foo.bar" + strconv.Itoa(i)))
		assert.Nil(t, err)
	}
	// Too many SearchDomain
	err := conf.Add(resolvconf.NewSearchDomain("foo.bar7"))
	assert.NotNil(t, err)
}

func TestSearchDomainCharLimit(t *testing.T) {
	conf := resolvconf.New()
	var dom string
	for i := 0; i < 256; i++ {
		dom = dom + "1"
	}
	err := conf.Add(resolvconf.NewSearchDomain(dom))
	assert.Nil(t, err)
	// Adding one more should break maximum number of chars limit
	err = conf.Add(resolvconf.NewSearchDomain("2"))
	assert.NotNil(t, err)
}

func TestLogging(t *testing.T) {

	// Nothing is logged if not enabeled
	conf := resolvconf.New()
	buf := new(bytes.Buffer)
	err := conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))
	assert.Nil(t, err)
	assert.NotContains(t, buf.String(), fmt.Sprintf("Added nameserver %s", net.ParseIP("8.8.8.8")))

	// Enable logging, test add Nameserver
	conf = resolvconf.New()
	buf.Reset()
	conf.EnableLogging(buf)
	assert.Nil(t, err)
	conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))
	assert.Contains(t, buf.String(), fmt.Sprintf("Added nameserver %s", net.ParseIP("8.8.8.8")))

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
	assert.Contains(t, buf.String(), "Added searchdomain foo.bar")
	conf.Remove(resolvconf.NewSearchDomain("foo.bar"))
	assert.Contains(t, buf.String(), "Removed searchdomain foo.bar")

	// Add & remove sort list pair
	buf.Reset()
	conf.Add(resolvconf.NewSortItem(net.ParseIP("8.8.8.8")))
	assert.Contains(t, buf.String(), "Added sortitem 8.8.8.8")
	conf.Remove(*resolvconf.NewSortItem(net.ParseIP("8.8.8.8")))
	assert.Contains(t, buf.String(), "Removed sortitem 8.8.8.8")

	// Add & remove option
	buf.Reset()
	conf.Add(resolvconf.NewOption("debug"))
	assert.Contains(t, buf.String(), "Added option debug")
	conf.Remove(resolvconf.NewOption("debug"))
	assert.Contains(t, buf.String(), "Removed option debug")
}

func TestIfFindReturnsPointer(t *testing.T) {
	conf := resolvconf.New()
	conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))
	ns := conf.Find(resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))).(*resolvconf.Nameserver)
	ns.IP = net.ParseIP("8.8.8.9")
	assert.NotNil(t, conf.Find(resolvconf.NewNameserver(net.ParseIP("8.8.8.9"))))

	// Try to remove after find
	err := conf.Remove(conf.Find(resolvconf.NewNameserver(net.ParseIP("8.8.8.9"))))
	assert.Nil(t, err)
}

func TestAddNilElements(t *testing.T) {
	conf := resolvconf.New()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Add nil element causes panic")
		}
	}()
	conf.Add(nil)
}

func TestRemoveNilElements(t *testing.T) {
	conf := resolvconf.New()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Remove nil element causes panic")
		}
	}()
	conf.Remove(nil)
}

func TestFindNilElements(t *testing.T) {
	conf := resolvconf.New()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Find nil element causes panic")
		}
	}()
	conf.Find(nil)
}

func TestOptionsCanBeCapped(t *testing.T) {
	conf := resolvconf.New()
	buf := new(bytes.Buffer)
	conf.EnableLogging(buf)

	conf.Add(resolvconf.NewOption("ndots").Set(16),
		resolvconf.NewOption("timeout").Set(31),
		resolvconf.NewOption("attempts").Set(6))
	ndots := conf.Find(resolvconf.NewOption("ndots"))
	timeout := conf.Find(resolvconf.NewOption("timeout"))
	attempts := conf.Find(resolvconf.NewOption("attempts"))
	assert.Equal(t, 15, ndots.(*resolvconf.Option).Get())
	assert.Equal(t, 30, timeout.(*resolvconf.Option).Get())
	assert.Equal(t, 5, attempts.(*resolvconf.Option).Get())

	assert.Contains(t, buf.String(), fmt.Sprintf("[WARN] Option ndots is capped to 15, set value is 16"))
	assert.Contains(t, buf.String(), fmt.Sprintf("[WARN] Option timeout is capped to 30, set value is 31"))
	assert.Contains(t, buf.String(), fmt.Sprintf("[WARN] Option attempts is capped to 5, set value is 6"))
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
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	conf.Remove(resolvconf.NewNameserver(net.ParseIP("8.8.4.4")))
	conf.Add(resolvconf.NewDomain("foo.bar"), resolvconf.NewSortItem(net.ParseIP("130.155.160.0")).SetNetmask(net.ParseIP("255.255.240.0")))
	conf.Write(os.Stdout)
	// Output: domain foo.bar
	// nameserver 8.8.8.8
	//
	// sortlist 130.155.160.0/255.255.240.0
}
