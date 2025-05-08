package quant

import "math"

var QuantMatrixY = [][]int{
	{16, 11, 10, 16, 24, 40, 51, 61},
	{12, 12, 14, 19, 26, 58, 60, 55},
	{14, 13, 16, 24, 40, 57, 69, 56},
	{14, 17, 22, 29, 51, 87, 80, 62},
	{18, 22, 37, 56, 68, 109, 103, 77},
	{24, 35, 55, 64, 81, 104, 113, 92},
	{49, 64, 78, 87, 103, 121, 120, 101},
	{72, 92, 95, 98, 112, 100, 103, 99},
}

var QuantMatrixCbCr = [][]int{
	{17, 18, 24, 47, 99, 99, 99, 99},
	{18, 21, 26, 66, 99, 99, 99, 99},
	{24, 26, 56, 99, 99, 99, 99, 99},
	{47, 66, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
}

func InverseQuantize(block [][]int, quantMatrix [][]float64) [][]float64 {
	result := make([][]float64, 8)
	for i := range result {
		result[i] = make([]float64, 8)
		for j := range result[i] {
			result[i][j] = float64(block[i][j]) * quantMatrix[i][j]
		}
	}
	return result
}

func QuantCoeff(quantMatrix [][]int, quality int) [][]float64 {
	result := make([][]float64, 8)
	var scale = 5000.0 / float64(quality)

	for i := range result {
		result[i] = make([]float64, 8)
		for j := range result[i] {
			val := math.Round(float64(quantMatrix[i][j]) * scale / 100.0)

			if i == 0 && j == 0 {
				if i == 0 && j == 0 {
					val = math.Min(val, 10.0)
				} else if i == 0 && j == 0 {
					val = math.Min(val, 16.0)
				}
			}

			result[i][j] = math.Max(1.0, val)
		}
	}
	return result
}

func Quantize(block [][]float64, quantMatrix [][]float64) [][]int {
	result := make([][]int, 8)
	for i := range result {
		result[i] = make([]int, 8)
		for j := range result[i] {
			result[i][j] = int(math.Round(block[i][j] / quantMatrix[i][j]))
		}
	}
	return result
}

func InverseQuantizeBlocks(quantBlocks [][][]int, quantMatrix [][]float64) [][][]float64 {
	dctBlocks := make([][][]float64, len(quantBlocks))
	for i, block := range quantBlocks {
		dctBlocks[i] = InverseQuantize(block, quantMatrix)
	}
	return dctBlocks
}
