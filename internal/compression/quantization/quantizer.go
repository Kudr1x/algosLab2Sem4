package quantization

import (
	"math"
)

func QuantizeBlock(dctCoeffs []float64, qMatrix []float64) []float64 {
	size := int(math.Sqrt(float64(len(dctCoeffs))))
	result := make([]float64, size*size)

	for i := 0; i < size*size; i++ {
		result[i] = math.Round(dctCoeffs[i] / qMatrix[i])
	}

	return result
}

func DequantizeBlock(quantizedCoeffs []float64, qMatrix []float64) []float64 {
	size := int(math.Sqrt(float64(len(quantizedCoeffs))))
	result := make([]float64, size*size)

	for i := 0; i < size*size; i++ {
		result[i] = quantizedCoeffs[i] * qMatrix[i]
	}

	return result
}
