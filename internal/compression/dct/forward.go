package dct

import (
	"math"
)

func ApplyDCT(block []float64, size int) []float64 {
	blockMatrix := make([][]float64, size)
	for i := range blockMatrix {
		blockMatrix[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			blockMatrix[i][j] = block[i*size+j]
		}
	}

	resultMatrix := make([][]float64, size)
	for i := range resultMatrix {
		resultMatrix[i] = make([]float64, size)
	}

	for u := 0; u < size; u++ {
		for v := 0; v < size; v++ {
			sum := 0.0

			cu := 1.0
			if u == 0 {
				cu = 1.0 / math.Sqrt(2.0)
			}

			cv := 1.0
			if v == 0 {
				cv = 1.0 / math.Sqrt(2.0)
			}

			for x := 0; x < size; x++ {
				for y := 0; y < size; y++ {
					sum += blockMatrix[x][y] *
						math.Cos((float64(2*x+1)*float64(u)*math.Pi)/(2.0*float64(size))) *
						math.Cos((float64(2*y+1)*float64(v)*math.Pi)/(2.0*float64(size)))
				}
			}

			resultMatrix[u][v] = 0.25 * cu * cv * sum
		}
	}

	result := make([]float64, size*size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			result[i*size+j] = resultMatrix[i][j]
		}
	}

	return result
}
