package resolvconf_test

import (
	"." // import the main package
	"bytes"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func GetConf(conf *resolvconf.Conf) (string, error) {
	buf := new(bytes.Buffer)
	err := conf.Write(buf)
	return buf.String(), err
}

func TestNameserverGeneration(t *testing.T) {
	conf := resolvconf.New()
	ns := resolvconf.NewNameserver(net.ParseIP("8.8.8.8"))

	// Test write a Nameserver
	conf.Add(ns)
	str, err := GetConf(conf)
	assert.Nil(t, err)
	assert.Contains(t, str, "nameserver 8.8.8.8")

	// Remove it and expect empty string
	conf.Remove(ns)
	str, err = GetConf(conf)
	assert.Nil(t, err)
	assert.Equal(t, "", str)

	// Add two Nameservers
	ns2 := resolvconf.NewNameserver(net.ParseIP("8.8.8.9"))
	conf.Add(ns, ns2)
	str, err = GetConf(conf)
	assert.Nil(t, err)
	assert.Contains(t, str, "nameserver 8.8.8.8")
	assert.Contains(t, str, "nameserver 8.8.8.9")
}

func TestDomainGeneration(t *testing.T) {
	conf := resolvconf.New()
	dom := resolvconf.NewDomain("foo.com")
	conf.Add(dom)
	str, err := GetConf(conf)
	assert.Nil(t, err)
	assert.Contains(t, str, "domain foo.com")

	// When no domain is given then "domain" should not be printed
	conf.Remove(dom)
	str, err = GetConf(conf)
	assert.Nil(t, err)
	assert.NotContains(t, str, "domain")
}

func TestOptionsGeneration(t *testing.T) {
	conf := resolvconf.New()
	dbg := resolvconf.NewOption("debug")
	err := conf.Add(dbg)
	assert.Nil(t, err)
	str, err := GetConf(conf)
	assert.Nil(t, err)
	assert.Contains(t, str, "options debug")

	conf.Remove(dbg)
	str, err = GetConf(conf)
	assert.Nil(t, err)
	assert.Equal(t, "", str)

	rotate := resolvconf.NewOption("rotate")
	conf.Add(dbg, rotate)
	str, err = GetConf(conf)
	assert.Nil(t, err)
	assert.Contains(t, str, "options debug rotate")

	ndots := resolvconf.NewOption("ndots").Set(3)
	conf.Add(ndots)
	conf.Remove(dbg)
	str, err = GetConf(conf)
	assert.Nil(t, err)
	assert.Contains(t, str, "options rotate ndots:3")
}

func TestSortlistGeneration(t *testing.T) {
	conf := resolvconf.New()
	addr1 := resolvconf.NewSortItem(net.ParseIP("8.8.8.8")).SetNetmask(net.ParseIP("255.255.255.0"))
	conf.Add(addr1)
	str, err := GetConf(conf)
	assert.Nil(t, err)
	assert.Contains(t, str, "sortlist 8.8.8.8/255.255.255.0")

	addr2 := resolvconf.NewSortItem(net.ParseIP("8.8.8.7"))
	conf.Remove(*addr1)
	conf.Add(addr2)
	str, err = GetConf(conf)
	assert.Nil(t, err)
	assert.Contains(t, str, "sortlist 8.8.8.7")

	conf.Add(addr1)
	str, err = GetConf(conf)
	assert.Nil(t, err)
	assert.Contains(t, str, "sortlist 8.8.8.7 8.8.8.8/255.255.255.0")
}

func TestSearchGeneration(t *testing.T) {
	conf := resolvconf.New()
	dom1 := resolvconf.NewSearchDomain("foo.bar")
	conf.Add(dom1)
	str, err := GetConf(conf)
	assert.Nil(t, err)
	assert.Contains(t, str, "search foo.bar")
}
