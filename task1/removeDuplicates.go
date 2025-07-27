package task1

func RemoveDuplicates(nums []int) int {
	// [1,2,1,5,9,2,3]
	mapInt := make(map[int]int)

	for i := 0; i < len(nums); i++ {
		if mapInt[nums[i]] == 0 {
			mapInt[nums[i]] = 1
		} else {
			nums = append(nums[:i], nums[i+1:]...)
			i--
		}
	}

	return len(nums)
}

func RemoveDuplicates1(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	slow := 1
	for fast := 1; fast < n; fast++ {
		if nums[fast] != nums[fast-1] {
			nums[slow] = nums[fast]
			slow++
		}
	}
	return slow
}
