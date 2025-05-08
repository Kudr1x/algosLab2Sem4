package differential

func DifferentialEncode(dc []int) []int {
	if len(dc) == 0 {
		return []int{}
	}
	result := make([]int, len(dc))
	result[0] = dc[0]
	for i := 1; i < len(dc); i++ {
		result[i] = dc[i] - dc[i-1]
	}
	return result
}

func DifferentialDecode(diff []int) []int {
	if len(diff) == 0 {
		return []int{}
	}
	result := make([]int, len(diff))
	result[0] = diff[0]
	for i := 1; i < len(diff); i++ {
		result[i] = result[i-1] + diff[i]
	}
	return result
}
