package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	curr "github.com/vladimirvivien/go-networking/currency/lib"
)

var currencies = curr.Load("../data.csv")

// Implement simple lookup service over TCP
// Loads ISO currency info usuing package lib
// Uses simple text-based protocol to interact
// with client and send the data

// Clients send currency search requests as JSON objects
// Command: {"Get":"<currency, country, or code>"}
// Data is then unmarshalled to curr.CurrencyRequest

//request used to search currency list

// Focus:
// use encoding packages to serialize data to/from GO data types
// to JSON representation.
// Uses encoding/json package Encoder/Decoder types that accept
// io.Writer/Reader so they can be used directly with io.Conn

// Usage: server [options]
// options:
// 	-host host endpoint, default ":4040"
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

	if _, err := conn.Write([]byte("Connected...\nUsage: {\"Get\":\"<currency, country, or code>\"}\n")); err != nil {
		log.Println("Error writing:", err)
		return
	}

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	// Command Loop
	for {

		// Decode incoming data -> curr.CurrencyRequest
		var req curr.CurrencyRequest
		if err := decoder.Decode(&req); err != nil {
			log.Println("Failed to decode request", err)
			return
		}

		result := curr.Find(currencies, req.Get)

		if err := encoder.Encode(&result); err != nil {
			log.Println("Failed to encode data", err)
			return
		}
	}
}
