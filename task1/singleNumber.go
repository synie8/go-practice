package task1

import "fmt"

func SingleNumber(nums []int) int {
	var a int
	var mapRes = make(map[int]int)
	for _, v := range nums {
		mapRes[v]++
	}
	for k, v := range mapRes {
		if v == 1 {
			a = k
		}
	}
	fmt.Println("map ****", mapRes)
	return a
}
