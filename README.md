# resolvconf

[![Build Status](https://travis-ci.org/fa1k3n/resolvconf.svg?branch=master)](https://travis-ci.org/fa1k3n/resolvconf) [![Go Report Card](https://goreportcard.com/badge/github.com/fa1k3n/resolvconf)](https://goreportcard.com/report/github.com/fa1k3n/resolvconf) [![Go Documentation](https://godoc.org/github.com/fa1k3n/resolvconf?status.svg)](https://godoc.org/github.com/fa1k3n/resolvconf)

Go package that simplifies manipulating resolv.conf files

The package provides a way to read and parse existing resolv.conf files from an io.Reader or to create a new file. The read objects can then be manipulated and written to a io.Writer object of your choice. 

Examples:

```go
package main

import (
	"net"
	"fmt"
	"os"
	"bytes"
	"github.com/Fa1k3n/resolvconf"
)

func main() {
	conf := resolvconf.New()
	
	// Add some options
	conf.Add(resolvconf.NewOption("debug"), resolvconf.NewOption("ndots").Set(3))

	// Add a nameservers
	conf.Add(resolvconf.NewNameserver(net.ParseIP("8.8.8.8")))

	// Add a sortlist
	conf.Add(resolvconf.NewSortItem(net.ParseIP("130.155.160.0")).SetNetmask("255.255.240.0"))

	// Dump to stdout
	conf.Write(os.Stdout)
}
```