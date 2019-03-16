package main

import "fmt"

type rectangle struct {
	width  float64
	height float64
}

func (r rectangle) area() float64 {
	return r.width * r.height
}

func (r rectangle) perimeter() float64 {
	return 2 * (r.width + r.height)
}

func main() {
	r := rectangle{
		width:  3,
		height: 4,
	}

	fmt.Println("Area: ", r.area())
	fmt.Println("Perimeter: ", r.perimeter())
}
