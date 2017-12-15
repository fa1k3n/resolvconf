package resolvconf

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

func parseOption(o string) (*Option, error) {
	keyval := strings.Split(o, ":")

	switch opt := keyval[0]; opt {
	case "debug", "rotate", "no-check-names", "inet6",
		"ip6-bytestring", "ip6-dotint", "no-ip6-dotint",
		"edns0", "single-request", "single-request-reopen",
		"no-tld-query", "use-vc":
		return &Option{o, -1}, nil
	case "ndots", "timeout", "attempts":
		val, err := strconv.Atoi(keyval[1])
		if err != nil {
			return nil, fmt.Errorf("%s unable to parse option value %s", opt, keyval[1])
		}
		return &Option{opt, val}, nil
	default:
		return nil, fmt.Errorf("Unknown option %s", opt)
	}
}

func parseLine(line string) ([]ConfItem, error) {
	toks := strings.Fields(line)
	var items []ConfItem
	var err error
	switch keyword := toks[0]; keyword {
	case "nameserver":
		ns := new(Nameserver)
		if ns.IP = net.ParseIP(toks[1]); ns.IP == nil {
			err = fmt.Errorf("Malformed IP address: %s", toks[1])
			break
		}
		items = append(items, ns)
	case "domain":
		items = append(items, NewDomain(toks[1]))
	case "search":
		for _, dom := range toks[1:] {
			items = append(items, NewSearchDomain(dom))
		}
	case "sortlist":
		for _, pair := range toks[1:] {
			var addr, nm net.IP
			addrNmStr := strings.Split(pair, "/")
			if addr = net.ParseIP(addrNmStr[0]); addr == nil {
				err = fmt.Errorf("Malformed IP address %s in searchlist", pair)
				break
			}
			if len(addrNmStr) > 1 {
				if nm = net.ParseIP(addrNmStr[1]); nm == nil {
					err = fmt.Errorf("Malformed netmask %s in searchlist", pair)
					break
				}
			}
			items = append(items, NewSortItem(addr).SetNetmask(nm))
		}
	case "options":
		for _, optStr := range toks[1:] {
			opt, e := parseOption(optStr)
			if e != nil {
				err = e
				break
			}
			items = append(items, opt)
		}
	default:
		err = fmt.Errorf("Unknown keyword %s", keyword)
	}

	return items, err
}

// ReadConf will read a configuration from given io.Reader
//
// Returns a new Conf object when successful otherwise
// nil and an error
func ReadConf(r io.Reader) (*Conf, error) {
	var res *multierror.Error
	conf := New()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		res = multierror.Append(res, err)
		return nil, res
	}
	confFile := strings.TrimSpace(string(b[:]))
	lines := strings.Split(confFile, "\n")
	for _, line := range lines {
		// Check if this line is a comment or empty
		if len(line) == 0 || line[0] == byte('#') || line[0] == byte(';') {
			continue
		}
		// Otherwise decode line
		opt, err := parseLine(line)
		if err != nil {
			res = multierror.Append(res, err)
			continue
		}

		for _, o := range opt {
			if err := conf.Add(o); err != nil {
				res = multierror.Append(res, err)
			}
		}
	}
	return conf, res.ErrorOrNil()
}
