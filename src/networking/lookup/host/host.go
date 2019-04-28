package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var host string

func main() {
	flag.StringVar(&host, "host", "localhost", "hostname to resolve")
	flag.Parse()

	address, err := net.LookupHost(host)
	if err != nil {
		log.Fatalln("err")
	}

	fmt.Println(len(address), address)
}
