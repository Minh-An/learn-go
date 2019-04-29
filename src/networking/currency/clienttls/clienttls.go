package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	curr "github.com/vladimirvivien/go-networking/currency/lib"
)

const prompt = "currency"

// Client for currency service project
// JSON requests : {"Get":"USD"}
// receives JSON encoded currency info over TCP

// Focus:
// IO streaming, data serialization, and client-side error handling.
// Configure the dialer to setup settings such as timeout and KeepAlive values.
// Also implements a simple connection-retry strategy when connecting.

// Usage: client [options]
// options:
//  - e service endpoint or socket path, default localhost:4040
func main() {
	// setup flags
	var addr string
	var cert, key, ca string
	flag.StringVar(&addr, "e", "localhost:4040", "service endpoint [ip addr or socket path]")
	flag.StringVar(&cert, "cert", "../ssl/client-cert.pem", "Public certificate")
	flag.StringVar(&key, "key", "../ssl/client-key.pem", "Private key")
	flag.StringVar(&ca, "ca", "../ssl/ca-cert.pem", "CA certificate")
	flag.Parse()

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

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            certificatePool,
		Certificates:       []tls.Certificate{certificate},
	}

	network := "tcp"

	// create a dialer to configure settings
	dialer := &net.Dialer{
		Timeout:   time.Second * 300,
		KeepAlive: time.Minute * 5,
	}

	// simple dialing strategy with retry with a simple backoff
	var (
		conn           net.Conn
		connTries      = 0
		connMaxRetries = 3
		connSleepRetry = time.Second * 1
	)
	for connTries < connMaxRetries {
		fmt.Println("creating connection socket to", addr)
		conn, err = tls.DialWithDialer(dialer, network, addr, tlsConfig)
		if err != nil {
			fmt.Println("failed to create socket:", err)
			switch nerr := err.(type) {
			case net.Error:
				// attempt to retry
				if nerr.Temporary() {
					connTries++
					fmt.Println("trying again in:", connSleepRetry)
					time.Sleep(connSleepRetry)
					continue
				}
				// non-temporary error
				fmt.Println("unable to recover")
				os.Exit(1)

			default: // non networking error
				os.Exit(1)
			}
		}
		// no error break
		break
	}

	if conn == nil {
		fmt.Println("failed to create a connection successfully")
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("connected to currency service: ", addr)

	var param string

	// repl
	for {
		fmt.Print(prompt, "> ")
		_, err = fmt.Scanf("%s", &param)
		if err != nil {
			fmt.Println("Usage: <search string or *>")
			continue
		}

		req := curr.CurrencyRequest{Get: param}

		// Send request:
		// use json encoder to encode value of type curr.CurrencyRequest
		// and stream it to the server via net.Conn.
		if err := json.NewEncoder(conn).Encode(&req); err != nil {
			switch err := err.(type) {
			case net.Error:
				fmt.Println("failed to send request:", err)
				os.Exit(1)
			default:
				fmt.Println("failed to encode request:", err)
				continue
			}
		}

		// Display response
		var currencies []curr.Currency
		err = json.NewDecoder(conn).Decode(&currencies)
		if err != nil {
			switch err := err.(type) {
			case net.Error:
				fmt.Println("failed to receive response:", err)
				os.Exit(1)
			default:
				fmt.Println("failed to decode response:", err)
				continue
			}
		}

		for i, c := range currencies {
			fmt.Printf("%2d. %s[%s]\t%s, %s\n", i, c.Code, c.Number, c.Name, c.Country)
		}
	}

}
