package sort

func SelectSort(nums []int) []int {
	for i := 0; i < len(nums)-1; i++ {
		minIdx := i
		for j := i + 1; j < len(nums); j++ {
			if nums[j] < nums[minIdx] {
				minIdx = j
			}
		}

		if minIdx != i {
			nums[i], nums[minIdx] = nums[minIdx], nums[i]
		}
	}

	return nums
}
