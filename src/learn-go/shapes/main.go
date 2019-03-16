package main

import "fmt"

type shape interface {
	getArea() float64
}

type triangle struct {
	base   float64
	height float64
}

type square struct {
	length float64
}

func (t triangle) getArea() float64 {
	return t.base * t.height / 2
}

func (s square) getArea() float64 {
	return s.length * s.length
}

func printArea(s shape) {
	fmt.Println("Area:", s.getArea())
}

func main() {
	t := triangle{base: 5, height: 6}
	s := square{5}
	printArea(t)
	printArea(s)
}
