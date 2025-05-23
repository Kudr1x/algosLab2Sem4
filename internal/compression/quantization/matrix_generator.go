package quantization

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

func GenerateQuantizationMatrix(blockSize int, quality int, isLuma bool) []float64 {
	var baseMatrix [][]int
	if isLuma {
		baseMatrix = QuantMatrixY
	} else {
		baseMatrix = QuantMatrixCbCr
	}

	if blockSize != 8 {
		result := make([]float64, blockSize*blockSize)
		for i := 0; i < blockSize*blockSize; i++ {
			result[i] = 16.0
		}
		return result
	}

	scale := 5000.0 / float64(quality)

	result := make([]float64, blockSize*blockSize)
	for i := 0; i < blockSize; i++ {
		for j := 0; j < blockSize; j++ {
			val := math.Round(float64(baseMatrix[i][j]) * scale / 100.0)

			if i == 0 && j == 0 {
				if isLuma {
					val = math.Min(val, 10.0)
				} else {
					val = math.Min(val, 16.0)
				}
			}

			result[i*blockSize+j] = math.Max(1.0, val)
		}
	}

	return result
}
