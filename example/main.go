package main

import (
	"os"
	"os/exec"
	"log"
	"net"
	"fmt"
	"bytes"
	".."
)

func main() {
	file, err := os.Open("resolv.conf.ex1")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r, err := resolvconf.ReadConf(file)	
	if err != nil {
		log.Fatal(err)
	}
	dbg, _ := resolvconf.NewOption("debug")
	ndots, _ := resolvconf.NewOption("ndots", 3)
	r.Add(dbg, ndots)

	ns, _ := resolvconf.NewNameserver(net.ParseIP("202.54.1.10"))
	r.Remove(ns)

	fmt.Println("=== New resolv.conf content ===\n")
	r.Write(os.Stdout)

	outf, err := os.OpenFile("/tmp/resolv.conf.ex1.out", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer outf.Close()
	r.Write(outf)

	
}