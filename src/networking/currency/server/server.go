package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	curr "networking/currency/lib"
	"strings"
)

var currencies = curr.Load("../data.csv")

// Implement simple lookup service over TCP
// Loads ISO currency info usuing package lib
// Uses simple text-based protocol to interact
// with client and send the data

// Command: GET <currency, country, or code>
// search result is printed line by line back to client

// Focus: streaming strategy when receiving data from client
// to avoid dropping data when the request is larger than interal buffer

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":4040", "Serpice endpoid [IP Addr or Socket Path]")
	flag.Parse()

	network := "tcp"

	//Create Listener for Network + Host Address
	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal("Failed to create listener:", err)
	}
	defer ln.Close()

	log.Println("**** Global Currency Service ***")
	log.Printf("Service started: (%s) %s\n", network, addr)

	// Connection loop - Handle incoming requests
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println("Failed to close connection:", err)
			}
			continue
		}
		log.Println("Connected to", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
	}()

	if _, err := conn.Write([]byte("Connected...\nUsage: GET <currency, country, or code>\n")); err != nil {
		log.Println("Error writing:", err)
		return
	}

	// appendBytes is a func stimulating EOF marker error using '\n'
	appendBytes := func(dest, src []byte) ([]byte, error) {
		for _, b := range src {
			if b == '\n' {
				return dest, io.EOF
			}
			dest = append(dest, b)
		}
		return dest, nil
	}

	// Loop for client
	for {
		var cmdLine []byte

		// stream with 4-bytes until EOF
		for {
			block := make([]byte, 4)
			n, err := conn.Read(block)
			if err != nil {
				if err == io.EOF {
					cmdLine, _ = appendBytes(cmdLine, block[:n])
				}
				log.Println("Connection read error:", err)
				return
			}
			if cmdLine, err = appendBytes(cmdLine, block[:n]); err == io.EOF {
				break
			}
		}

		cmd, param := parseCommand(string(cmdLine))

		// execute the command
		switch strings.ToUpper(cmd) {
		case "GET":
			result := curr.Find(currencies, param)
			if len(result) == 0 {
				if _, err := conn.Write([]byte("Nothing found\n")); err != nil {
					log.Println("Error writing:", err)
				}
				continue
			}

			for _, cur := range result {
				_, err := conn.Write([]byte(
					fmt.Sprintf("%s %s %s %s\n", cur.Name, cur.Code, cur.Number, cur.Country)))
				if err != nil {
					log.Println("Error writing:", err)
					return
				}
			}
		default:
			if _, err := conn.Write([]byte("Invalid command\n")); err != nil {
				log.Println("Error writing:", err)
				return
			}
		}
	}
}

func parseCommand(cmdLine string) (cmd, param string) {
	parts := strings.Split(cmdLine, " ")
	if len(parts) != 2 {
		return "", ""
	}
	cmd = strings.TrimSpace(parts[0])
	param = strings.TrimSpace(parts[1])
	return
}
