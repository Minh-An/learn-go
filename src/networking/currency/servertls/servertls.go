package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

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
// improves robustness of server code by introducing configuration
// for read/write timeout values. Ensures that a client can't hold
// a connection hostage by taking a long time to send/receive data

// use encoding packages to serialize data to/from GO data types
// to JSON representation.
// Uses encoding/json package Encoder/Decoder types that accept
// io.Writer/Reader so they can be used directly with io.Conn

// Usage: server [options]
// options:
// 	-host host endpoint, default ":4040"
func main() {
	var addr string
	var cert, key, ca string
	flag.StringVar(&addr, "addr", ":4040", "Service endpoint [IP Addr or Socket Path]")
	flag.StringVar(&cert, "cert", "../ssl/t480-cert.pem", "Public certificate")
	flag.StringVar(&key, "key", "../ssl/t480-key.pem", "Private key")
	flag.StringVar(&ca, "ca", "../ssl/ca-cert.pem", "CA certificate")
	flag.Parse()

	network := "tcp"

	// Load server certificate by providing private key
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA certificate
	caCertificate, err := ioutil.ReadFile(ca)
	if err != nil {
		log.Fatal("failed to read CA cert", err)
	}

	certificatePool := x509.NewCertPool()
	certificatePool.AppendCertsFromPEM(caCertificate)

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certificatePool,
		Certificates: []tls.Certificate{certificate},
	}

	//Create Listener for Network + Host Address
	ln, err := tls.Listen(network, addr, tlsConfig)
	if err != nil {
		log.Fatal("Failed to create listener:", err)
	}
	defer ln.Close()

	log.Println("**** Global Currency Service ***")
	log.Printf("Service started: (%s) %s; Server certificate %s\n", network, addr, cert)

	// delay to sleep when accept fails w/ temporary error
	acceptDelay := time.Millisecond * 10
	acceptCount := 0

	// Connection loop - Handle incoming requests
	for {
		conn, err := ln.Accept()
		if err != nil {
			switch e := err.(type) {
			case net.Error:
				//if temp error, attempt to connect again
				if e.Temporary() {
					if acceptCount > 5 {
						log.Printf("Unable to connect after %d retries: %v", acceptCount, err)
						return
					}
					acceptCount++
					acceptDelay *= 2
					time.Sleep(acceptDelay)
					continue
				}
			default:
				fmt.Println(err)
				conn.Close()
				continue
			}
			acceptDelay = time.Millisecond * 10
			acceptCount = 0
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

	// 45 second deadline
	if err := conn.SetDeadline(time.Now().Add(time.Second * 45)); err != nil {
		log.Println("Failed to set deadline:", err)
		return
	}

	// Command Loop
	for {

		// Decode incoming data -> curr.CurrencyRequest
		var req curr.CurrencyRequest
		decoder := json.NewDecoder(conn)
		if err := decoder.Decode(&req); err != nil {
			switch err := err.(type) {
			case net.Error:
				if err.Timeout() {
					log.Println("Deadline reached, disconnecting...")
				}
				log.Println("Network error: ", err)
				return
			default:
				if err == io.EOF {
					log.Println("Closing connection", err)
					return
				}
				encoder := json.NewEncoder(conn)
				if encerr := encoder.Encode(&curr.CurrencyError{Error: err.Error()}); encerr != nil {
					fmt.Println("Failed error encoding:", encerr)
					return
				}
				continue
			}

		}

		result := curr.Find(currencies, req.Get)

		encoder := json.NewEncoder(conn)
		if err := encoder.Encode(&result); err != nil {
			switch err := err.(type) {
			case net.Error:
				log.Println("Failed to send response: ", err)
				return
			default:
				if encerr := encoder.Encode(&curr.CurrencyError{Error: err.Error()}); encerr != nil {
					fmt.Println("Failed error encoding:", encerr)
					return
				}
				continue
			}
		}

		//renew dealine for 45 sec later
		if err := conn.SetDeadline(time.Now().Add(time.Second * 45)); err != nil {
			log.Println("Failed to set deadline:", err)
			return
		}
	}
}
