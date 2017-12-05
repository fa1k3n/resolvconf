package main

import (
	"os"
	"log"
	"net"
	"fmt"
	"net/http"
	".."
)

func main() {
	fmt.Println("=== EX1: reading from local file ===\n")
	file, err := os.Open("resolv.conf.ex1")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r, err := resolvconf.ReadConf(file)	
	if err != nil {
		log.Fatal(err)
	}
	r.Add(resolvconf.NewOption("debug"), resolvconf.NewOption("ndots", 3))
	r.Remove(resolvconf.NewNameserver(net.ParseIP("202.54.1.10")))

	r.Write(os.Stdout)

	outf, err := os.OpenFile("/tmp/resolv.conf.ex1.out", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer outf.Close()
	r.Write(outf)
	
	fmt.Println("=== EX2: reading from http request ===\n")
	res, err := http.Get("https://gist.githubusercontent.com/turadg/7876784/raw/c7f2500fa4762cfe443e30c64c6ed8a888f6ac74/resolv.conf")
	if err != nil {
		log.Fatal(err)
	}
	conf, err := resolvconf.ReadConf(res.Body)
	res.Body.Close()
	conf.Write(os.Stdout)
}