package main

import (
	"fmt"
	"net/http"
)

func main() {
	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
	}

	c := make(chan string)

	for _, link := range links {
		go checkLink(link, c)
	}

	var s string
	for {
		s = <-c
		if s == "" {
			return
		}
		fmt.Println(s)
	}

}

func checkLink(link string, c chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println("Error getting", link)
		c <- "Error getting " + link
		return
	}
	fmt.Println(link, "is up!")
	c <- link + " is up!"
}
