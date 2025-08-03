package task2

import "fmt"

func Add(a *int) {
	*a += 10
	fmt.Println("a1:", a)
}

func Mul(a *[]int) {
	for i := 0; i < len(*a); i++ {
		(*a)[i] *= 2
	}
	fmt.Println("a1:", a)
}
