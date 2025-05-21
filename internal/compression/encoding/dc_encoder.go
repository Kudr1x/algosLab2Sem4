package encoding

import (
	"fmt"
	"math"
)

func DCCategory(diff int) uint8 {
	if diff == 0 {
		return 0
	}

	absDiff := int(math.Abs(float64(diff)))
	category := uint8(math.Floor(math.Log2(float64(absDiff))) + 1)

	return category
}

func EncodeDCCoefficient(diff int, table HuffmanTable) ([]byte, error) {
	category := DCCategory(diff)

	huffCode, ok := table[category]
	if !ok {
		return nil, fmt.Errorf("категория %d не найдена в таблице Хаффмана", category)
	}

	result := make([]byte, 0, 32)

	for i := int(huffCode.Length) - 1; i >= 0; i-- {
		bit := (huffCode.Code >> uint(i)) & 1
		result = append(result, byte(bit))
	}

	if category > 0 {
		var additionalBits int
		if diff < 0 {
			additionalBits = diff + (1 << int(category)) - 1
		} else {
			additionalBits = diff
		}

		for i := int(category) - 1; i >= 0; i-- {
			bit := (additionalBits >> uint(i)) & 1
			result = append(result, byte(bit))
		}
	}

	return result, nil
}

func EncodeDCCoefficients(coeffs []int, table HuffmanTable) ([]byte, error) {
	if len(coeffs) == 0 {
		return []byte{}, nil
	}

	result := make([]byte, 0, len(coeffs)*16)

	encodedDC, err := EncodeDCCoefficient(coeffs[0], table)
	if err != nil {
		return nil, err
	}
	result = append(result, encodedDC...)

	for i := 1; i < len(coeffs); i++ {
		diff := coeffs[i] - coeffs[i-1]
		encodedDC, err := EncodeDCCoefficient(diff, table)
		if err != nil {
			return nil, err
		}
		result = append(result, encodedDC...)
	}

	return result, nil
}
