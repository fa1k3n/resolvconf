# resolvconf

[![Build Status](https://travis-ci.org/Fa1k3n/resolvconf.svg?branch=master)](https://travis-ci.org/Fa1k3n/resolvconf) [![Go Report Card](https://goreportcard.com/badge/github.com/fa1k3n/resolvconf)](https://goreportcard.com/report/github.com/fa1k3n/resolvconf) [![Go Documentation](https://godoc.org/github.com/Fa1k3n/resolvconf?status.svg)](https://godoc.org/github.com/Fa1k3n/resolvconf)

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
	"github.com/fa1k3n/resolvconf"
)

func main() {
	conf := resolvconf.New()
	
	// Add some options
	conf.Add(resolvconf.Option("debug"), resolvconf.Option("ndots", 3))

	// Add a nameservers
	conf.Add(resolvconf.Nameserver(net.ParseIP("8.8.8.8")))

	// Dump to stdout
	conf.Write(os.Stdout)
}
```