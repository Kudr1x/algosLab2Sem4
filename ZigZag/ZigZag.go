package ZigZag

func ZigZagOrder() [][2]int {
	return [][2]int{
		{0, 0}, {0, 1}, {1, 0}, {2, 0}, {1, 1}, {0, 2}, {0, 3}, {1, 2},
		{2, 1}, {3, 0}, {4, 0}, {3, 1}, {2, 2}, {1, 3}, {0, 4}, {0, 5},
		{1, 4}, {2, 3}, {3, 2}, {4, 1}, {5, 0}, {6, 0}, {5, 1}, {4, 2},
		{3, 3}, {2, 4}, {1, 5}, {0, 6}, {0, 7}, {1, 6}, {2, 5}, {3, 4},
		{4, 3}, {5, 2}, {6, 1}, {7, 0}, {7, 1}, {6, 2}, {5, 3}, {4, 4},
		{3, 5}, {2, 6}, {1, 7}, {2, 7}, {3, 6}, {4, 5}, {5, 4}, {6, 3},
		{7, 2}, {7, 3}, {6, 4}, {5, 5}, {4, 6}, {3, 7}, {4, 7}, {5, 6},
		{6, 5}, {7, 4}, {7, 5}, {6, 6}, {5, 7}, {6, 7}, {7, 6}, {7, 7},
	}
}

func ZigZagScan(block [][]int) []int {
	result := make([]int, 64)
	order := ZigZagOrder()
	for i, pos := range order {
		result[i] = block[pos[0]][pos[1]]
	}
	return result
}

func InverseZigZagScan(coefficients []int) [][]int {
	block := make([][]int, 8)
	for i := range block {
		block[i] = make([]int, 8)
	}
	order := ZigZagOrder()
	for i, pos := range order {
		if i < len(coefficients) {
			block[pos[0]][pos[1]] = coefficients[i]
		}
	}
	return block
}
