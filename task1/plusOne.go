package task1

func PlusOne(digits []int) []int {
	result := make([]int, len(digits))
	m := 1
	len := len(digits)
	for i := len - 1; i >= 0; i-- {
		m += digits[i]
		result[i] = m % 10
		if m >= 10 && i > 0 {
			m = 1
		} else if m < 10 && i > 0 {
			m = 0
		} else if m >= 10 {
			r := []int{1}
			result = append(r, result...)
		}

	}

	return result
}

func PlusOne1(digits []int) []int {
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		if digits[i] != 9 {
			digits[i]++
			for j := i + 1; j < n; j++ {
				digits[j] = 0
			}
			return digits
		}
	}
	// digits 中所有的元素均为 9

	digits = make([]int, n+1)
	digits[0] = 1
	return digits
}
