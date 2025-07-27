package main

import (
	"fmt"
	"go-practice/task1"
)

func main() {
	digits := [][]int{[]int{2, 3}, []int{4, 5}, []int{6, 7}, []int{8, 9}, []int{1, 10}}
	ok := task1.Merge(digits)
	fmt.Println("****", ok)
}
