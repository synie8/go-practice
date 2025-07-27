package task1

func twoSum(nums []int, target int) []int {
	/**
	*输入：nums = [2,7,11,15], target = 9
	输出：[0,1]
	解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。
	*/

	result := make([]int, 2)

	for i, v := range nums {
		for j := i + 1; j < len(nums); j++ {
			if v+nums[j] == target {
				result[0] = i
				result[1] = j
			}
		}
	}

	return result
}
