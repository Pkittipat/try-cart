package sort

func MergeSort(nums []int) []int {
	if len(nums) <= 1 {
		return nums
	}
	mid := len(nums) / 2
	leftHalf := nums[:mid]
	rightHalf := nums[mid:]

	sortedLeft := MergeSort(leftHalf)
	sortedRight := MergeSort(rightHalf)

	return merge(sortedLeft, sortedRight)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
			continue
		}

		result = append(result, right[j])
		j++
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}
