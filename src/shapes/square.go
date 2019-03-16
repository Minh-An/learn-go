package main

type square float64

func (s square) perimeter() float64 {
	return 4 * float64(s)
}

func (s square) area() float64 {
	return float64(s * s)
}

/*
func main() {
	var s square = 4.5
	fmt.Println(s.perimeter())
	fmt.Println(s.area())
}
*/
