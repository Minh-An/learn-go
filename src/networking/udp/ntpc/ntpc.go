package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	var host string
	flag.StringVar(&host, "e", "us.pool.ntp.org:123", "NTP Host")
	flag.Parse()

	req := make([]byte, 48)

	req[0] = 0x1B

	rsp := make([]byte, 48)

	raddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalln("Failed while closing connection", err)
		}
	}()

	fmt.Printf("Time from (udp) %s\n", conn.RemoteAddr())

	if _, err = conn.Write(req); err != nil {
		log.Fatalf("Failed to send request %v\n", err)
	}

	read, err := conn.Read(rsp)
	if err != nil {
		log.Fatalf("Failed to recieve response %v\n", err)
	}

	if read != 48 {
		log.Fatalf("Didn't get all bytes from server\n")
	}

	secs := binary.BigEndian.Uint32(rsp[40:])
	frac := binary.BigEndian.Uint32(rsp[44:])

	offset := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Sub(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)).Seconds()
	now := float64(secs) - offset
	fmt.Println(time.Unix(int64(now), int64(frac)))
}
