package main

import "fmt"

func main() {
	// var colors map[string]string
	colors := map[string]string{
		"red":   "ff0000",
		"green": "00ff00",
	}

	// colors := make(map[string]string)
	colors["white"] = "ffffff"
	// delete(colors, "white")

	printMap(colors)
}

func printMap(c map[string]string) {
	for color, hex := range c {
		fmt.Printf("Hex code for %v is %v\n", color, hex)
	}
}
