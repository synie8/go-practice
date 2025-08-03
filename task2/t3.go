package task2

import "fmt"

type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
}
type Circle struct {
}

func (r *Rectangle) Area() {
	fmt.Println("Rectangle Area")
}
func (r *Circle) Area() {
	fmt.Println("Circle Area")
}
func (r *Rectangle) Perimeter() {
	fmt.Println("Rectangle Perimeter")
}
func (r *Circle) Perimeter() {
	fmt.Println("Circle Perimeter")
}
