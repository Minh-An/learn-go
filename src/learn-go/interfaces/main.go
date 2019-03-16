package main

import "fmt"

func main() {
	eb := englishBot{}
	printGreeting(eb)
	sb := spanishBot{}
	printGreeting(sb)
}

type bot interface {
	getGreeting() string
}

type englishBot struct{}
type spanishBot struct{}

func (englishBot) getGreeting() string {
	return "Why hello there!"
}

func (spanishBot) getGreeting() string {
	return "¡Hola! ¿Comó estas?"
}

func printGreeting(b bot) {
	fmt.Println(b.getGreeting())
}
