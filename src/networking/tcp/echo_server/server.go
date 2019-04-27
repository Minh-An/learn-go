package main

import (
	"flag"
	"log"
	"net"
)

func main() {
	var host string
	flag.StringVar(&host, "host", ":4040", "Echo Server Host:Port")
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Fatalln(err)
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Fatalln(err)
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	r, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	w, err := conn.Write(buf[:r])
	if err != nil {
		log.Fatalln(err)
	}

	if r != w {
		log.Fatalf("Bytes read not equal (%d) to bytes written (%d)", r, w)
	}

}
