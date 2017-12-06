package resolvconf

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

func parseOption(o string) (option, error) {
	keyval := strings.Split(o, ":")

	switch opt := keyval[0]; opt {
	case "debug", "rotate", "no-check-names", "inet6",
		"ip6-bytestring", "ip6-dotint", "no-ip6-dotint",
		"edns0", "single-request", "single-request-reopen",
		"no-tld-query", "use-vc":
		return option{o, -1}, nil
	case "ndots", "timeout", "attempts":
		val, err := strconv.Atoi(keyval[1])
		if err != nil {
			return option{"", -1}, fmt.Errorf("%s unable to parse option value %s", opt, keyval[1])
		}
		return option{opt, val}, nil
	default:
		return option{"", -1}, fmt.Errorf("Unknown option %s", opt)
	}
}

func parseLine(line string) (interface{}, error) {
	toks := strings.Fields(line)
	switch keyword := toks[0]; keyword {
	case "nameserver":
		var ns nameserver
		if ns.IP = net.ParseIP(toks[1]); ns.IP == nil {
			return nil, fmt.Errorf("Malformed IP address: %s", toks[1])
		}
		return ns, nil
	case "domain":
		return Domain(toks[1]), nil
	case "search":
		var doms []searchDomain
		for _, dom := range toks[1:] {
			doms = append(doms, SearchDomain(dom))
		}
		return search{doms}, nil
	case "sortlist":
		var pairs []sortlistpair
		for i, pair := range toks[1:] {
			var addr, nm net.IP
			if i == 10 {
				return sortlist{pairs}, fmt.Errorf("Too long sortlist, 10 is maximum")
			}
			addr_nm_str := strings.Split(pair, "/")
			if addr = net.ParseIP(addr_nm_str[0]); addr == nil {
				return nil, fmt.Errorf("Malformed IP address %s in searchlist", pair)
			}
			if len(addr_nm_str) > 1 {
				if nm = net.ParseIP(addr_nm_str[1]); nm == nil {
					return nil, fmt.Errorf("Malformed netmask %s in searchlist", pair)
				}
			}
			pairs = append(pairs, SortlistPair(addr, nm))

		}
		return sortlist{pairs}, nil
	case "options":
		var opts []option
		for _, opt_str := range toks[1:] {
			opt, err := parseOption(opt_str)
			if err != nil {
				return nil, err
			}
			opts = append(opts, opt)
		}
		return opts, nil
	default:
		return nil, fmt.Errorf("Unknown keyword %s", keyword)
	}
}

// ReadConf will read a configuration from given io.Reader
//
// Returns a new Conf object when succesful otherwise 
// nil and an error
func ReadConf(r io.Reader) (*Conf, error) {
	var stored_err error
	stored_err = nil
	conf := New()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return conf, err
	}
	conf_file := strings.TrimSpace(string(b[:]))
	lines := strings.Split(conf_file, "\n")
	for _, line := range lines {
		// Check if this line is a comment or empty
		if len(line) == 0 || line[0] == byte('#') || line[0] == byte(';') {
			continue
		}
		// Otherwise decode line
		opt, err := parseLine(line)
		if err != nil {
			if opt == nil {
				// Only if there is an error and no option
				return conf, err
			}
			// Otherwise add this error to stored errors
			// and continue
			stored_err = fmt.Errorf("%s\n%s", err, stored_err)
		}

		if err := conf.Add(opt); err != nil {
			return conf, err
		}
	}
	return conf, stored_err
}
