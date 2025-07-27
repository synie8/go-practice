package task1

import "fmt"

func Merge(intervals [][]int) [][]int {

	for i := 1; i < len(intervals); i++ {
		//i-1[0]
		a := intervals[i-1][0]
		//i-1[1]
		b := intervals[i-1][1]
		//i[0]
		c := intervals[i][0]
		//i[1]
		d := intervals[i][1]
		fmt.Println("i ******", i, intervals)

		fmt.Println("f ******", !((a > d) || (c > b)) || (c-1 == b))
		if !((a > d) || (c > b)) || (c-1 == b) {

			if b >= d {
				intervals[i-1][1] = b
			} else {
				intervals[i-1][1] = d
			}
			if a >= c {
				intervals[i-1][0] = c
			} else {
				intervals[i-1][0] = a
			}
			fmt.Println("i ******", i, intervals)

			intervals = append(intervals[:i], intervals[i+1:]...)
			i--
		}
	}
	return intervals
}
