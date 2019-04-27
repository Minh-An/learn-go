package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	var host string
	flag.StringVar(&host, "host", ":4040", "Host Server")
	flag.Parse()
	text := flag.Arg(0)

	raddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	w, err := conn.Write([]byte(text))
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 1024)
	r, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	if r != w {
		log.Fatalln("Bytes read not equal to bytes written")
	}
	fmt.Println(string(buf[:r]))
}
