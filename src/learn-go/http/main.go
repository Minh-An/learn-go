package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type logWriter struct{}

func (logWriter) Write(bs []byte) (int, error) {
	fmt.Println("LogWriter:", string(bs))
	return len(bs), nil
}

func main() {
	resp, err := http.Get("http://google.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	lw := logWriter{}

	n, err := io.Copy(lw, resp.Body)
	fmt.Println("Bytes written:", n, "Error:", err)
}
