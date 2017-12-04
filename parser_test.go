package resolvconf_test

import (
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"testing"

	"." // import the main package
)

func TestReadNameserver(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("nameserver 8.8.8.8"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(conf.Nameservers))
	assert.Equal(t, net.ParseIP("8.8.8.8"), conf.Nameservers[0].IP)

	conf, err = resolvconf.ReadConf(strings.NewReader("nameserver 8.8.8.9"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(conf.Nameservers))
	assert.Equal(t, net.ParseIP("8.8.8.9"), conf.Nameservers[0].IP)
}

func TestReadFaultyNameserver(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("nameserver 8.8.8"))
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(conf.Nameservers))

	conf, err = resolvconf.ReadConf(strings.NewReader("nameserver 8.8.8.8.8"))
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(conf.Nameservers))

	conf, err = resolvconf.ReadConf(strings.NewReader("nameserver www.golang.org"))
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(conf.Nameservers))
}

func TestReadUnknownConfOption(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("nameserv 8.8.8.9"))
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(conf.Nameservers))
}

func TestReadSeveralNameservers(t *testing.T) {
	conf_str := "nameserver 8.8.8.8\n" +
		"nameserver 8.8.8.9\n"
	conf, _ := resolvconf.ReadConf(strings.NewReader(conf_str))
	assert.Equal(t, 2, len(conf.Nameservers))
}

func TestMaxThreeNameservers(t *testing.T) {
	conf_str := "nameserver 8.8.8.8\n" +
		"nameserver 8.8.8.9\n" +
		"nameserver 8.8.8.10\n" +
		"nameserver 8.8.8.11\n"
	conf, err := resolvconf.ReadConf(strings.NewReader(conf_str))
	assert.NotNil(t, err)
	assert.Equal(t, 3, len(conf.Nameservers))
}

func TestAddSameNameserverGivesError(t *testing.T) {
	conf_str := "nameserver 8.8.8.8\n" +
		"nameserver 8.8.8.8\n"
	conf, err := resolvconf.ReadConf(strings.NewReader(conf_str))
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(conf.Nameservers))
}

func TestCommentsAndBlankLinesAreSkipped(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("# This is a comment"))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.Nameservers))

	conf, err = resolvconf.ReadConf(strings.NewReader("#This is another comment"))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.Nameservers))

	conf, err = resolvconf.ReadConf(strings.NewReader("  #This is a third comment"))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.Nameservers))

	conf, err = resolvconf.ReadConf(strings.NewReader("; This is a forth comment"))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.Nameservers))

	// Empty line
	conf, err = resolvconf.ReadConf(strings.NewReader("\n"))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(conf.Nameservers))
}

func TestReadDomain(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("domain foo.com"))
	assert.Nil(t, err)
	assert.Equal(t, "foo.com", conf.Domain.Name)

	conf, err = resolvconf.ReadConf(strings.NewReader("domain     foo.com"))
	assert.Nil(t, err)
	assert.Equal(t, "foo.com", conf.Domain.Name)

	conf, err = resolvconf.ReadConf(strings.NewReader("    domain     foo.com"))
	assert.Nil(t, err)
	assert.Equal(t, "foo.com", conf.Domain.Name)
}

func TestReadSearch(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("search foo.com"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(conf.Search.Domains))
	assert.Equal(t, "foo.com", conf.Search.Domains[0].Name)
}

func TestReadMultiSearch(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("search foo.com bar.com     baz.com"))
	assert.Nil(t, err)
	assert.Equal(t, 3, len(conf.Search.Domains))
	for i, dom := range []string{"foo.com", "bar.com", "baz.com"} {
		assert.Equal(t, dom, conf.Search.Domains[i].Name)
	}
}

func TestReadSortlist(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("sortlist 130.155.160.0"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(conf.Sortlist.Pairs))
	assert.Equal(t, "130.155.160.0", conf.Sortlist.Pairs[0].Address.String())
	assert.Nil(t, conf.Sortlist.Pairs[0].Netmask)
}

func TestReadSortlistFaultyAddress(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("sortlist 130.155.160"))
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(conf.Sortlist.Pairs))
}

func TestReadMultiSortlist(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("sortlist 130.155.160.0 130.155.0.0"))
	assert.Nil(t, err)
	assert.Equal(t, 2, len(conf.Sortlist.Pairs))
	assert.Equal(t, "130.155.160.0", conf.Sortlist.Pairs[0].Address.String())
	assert.Equal(t, "130.155.0.0", conf.Sortlist.Pairs[1].Address.String())
	assert.Nil(t, conf.Sortlist.Pairs[0].Netmask)
}

func TestReadSortlistWithNetmask(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("sortlist 130.155.160.0/255.255.240.0"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(conf.Sortlist.Pairs))
	assert.Equal(t, "130.155.160.0", conf.Sortlist.Pairs[0].Address.String())
	assert.Equal(t, "255.255.240.0", conf.Sortlist.Pairs[0].Netmask.String())
}

func TestReadSortlistWithBadNetmask(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("sortlist 130.155.160.0/255.255.240"))
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(conf.Sortlist.Pairs))
}

func TestMaxTenSortlistPairsMayBeDefined(t *testing.T) {
	conf_str := "sortlist 1.1.1.0 1.1.1.1 " +
		"1.1.1.2 1.1.1.3 1.1.1.4 1.1.1.5 1.1.1.6 " +
		"1.1.1.7 1.1.1.8 1.1.1.9 1.1.1.10"
	conf, err := resolvconf.ReadConf(strings.NewReader(conf_str))
	assert.NotNil(t, err)
	assert.Equal(t, 10, len(conf.Sortlist.Pairs))
}

func TestOptions(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("options debug"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(conf.Options))
	assert.Equal(t, "debug", conf.Options[0].Type)

	conf, err = resolvconf.ReadConf(strings.NewReader("options debug rotate"))
	assert.Nil(t, err)
	assert.Equal(t, 2, len(conf.Options))
	assert.Equal(t, "rotate", conf.Options[1].Type)

	conf, err = resolvconf.ReadConf(strings.NewReader("options debug rotate ndots:12"))
	assert.Nil(t, err)
	assert.Equal(t, 3, len(conf.Options))
	assert.Equal(t, "ndots", conf.Options[2].Type)
	assert.Equal(t, 12, conf.Options[2].Value)
}

func TestUnknownOption(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("options foo"))
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(conf.Options))
}

func TestBadOption(t *testing.T) {
	conf, err := resolvconf.ReadConf(strings.NewReader("options ndots:"))
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(conf.Options))

	conf, err = resolvconf.ReadConf(strings.NewReader("options ndots:foos"))
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(conf.Options))
}

func TestAllOptions(t *testing.T) {
	conf_str := "options debug   ndots:3 timeout:5 attempts:4 " +
		"rotate no-check-names inet6 ip6-bytestring ip6-dotint " +
		"no-ip6-dotint edns0 single-request single-request-reopen " +
		"no-tld-query use-vc"
	conf, err := resolvconf.ReadConf(strings.NewReader(conf_str))
	assert.Nil(t, err)
	assert.Equal(t, 15, len(conf.Options))
}
