package zigzag

func Order() [][2]int {
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

func Scan(block []float64, size int) []float64 {
	if size != 8 {
		return fallbackZigZagScan(block, size)
	}

	blockMatrix := make([][]float64, size)
	for i := range blockMatrix {
		blockMatrix[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			blockMatrix[i][j] = block[i*size+j]
		}
	}

	// Используем порядок зигзаг-обхода
	result := make([]float64, size*size)
	order := Order()
	for i, pos := range order {
		result[i] = blockMatrix[pos[0]][pos[1]]
	}

	return result
}

func IScan(zigzagData []float64, size int) []float64 {
	if size != 8 {
		return fallbackInverseZigZagScan(zigzagData, size)
	}

	blockMatrix := make([][]float64, size)
	for i := range blockMatrix {
		blockMatrix[i] = make([]float64, size)
	}

	order := Order()
	for i, pos := range order {
		if i < len(zigzagData) {
			blockMatrix[pos[0]][pos[1]] = zigzagData[i]
		}
	}

	result := make([]float64, size*size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			result[i*size+j] = blockMatrix[i][j]
		}
	}

	return result
}

func fallbackZigZagScan(block []float64, size int) []float64 {
	result := make([]float64, size*size)
	index := 0

	for sum := 0; sum <= 2*(size-1); sum++ {
		if sum%2 == 0 {
			for i := min(sum, size-1); i >= max(0, sum-(size-1)); i-- {
				j := sum - i
				result[index] = block[i*size+j]
				index++
			}
		} else {
			for i := max(0, sum-(size-1)); i <= min(sum, size-1); i++ {
				j := sum - i
				result[index] = block[i*size+j]
				index++
			}
		}
	}

	return result
}

func fallbackInverseZigZagScan(zigzagData []float64, size int) []float64 {
	result := make([]float64, size*size)
	index := 0

	for sum := 0; sum <= 2*(size-1); sum++ {
		if sum%2 == 0 {
			for i := min(sum, size-1); i >= max(0, sum-(size-1)); i-- {
				j := sum - i
				if index < len(zigzagData) {
					result[i*size+j] = zigzagData[index]
					index++
				}
			}
		} else {
			for i := max(0, sum-(size-1)); i <= min(sum, size-1); i++ {
				j := sum - i
				if index < len(zigzagData) {
					result[i*size+j] = zigzagData[index]
					index++
				}
			}
		}
	}

	return result
}
