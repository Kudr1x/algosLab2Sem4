package encoding

import (
	"fmt"
	"math"
)

func ACCategory(value int) uint8 {
	if value == 0 {
		return 0
	}

	absValue := int(math.Abs(float64(value)))
	category := uint8(math.Floor(math.Log2(float64(absValue))) + 1)

	return category
}

func EncodeACCoefficient(runLength uint8, value int, table HuffmanTable) ([]byte, error) {
	category := ACCategory(value)
	if category > 10 {
		return nil, fmt.Errorf("категория %d превышает максимальное значение 10", category)
	}

	symbol := byte((runLength << 4) | category)

	huffCode, ok := table[symbol]
	if !ok {
		return nil, fmt.Errorf("символ %02x не найден в таблице Хаффмана", symbol)
	}

	result := make([]byte, 0, 32) // Предварительно выделяем память

	for i := int(huffCode.Length) - 1; i >= 0; i-- {
		bit := (huffCode.Code >> uint(i)) & 1
		result = append(result, byte(bit))
	}

	if category > 0 {
		var additionalBits int
		if value < 0 {
			additionalBits = value + (1 << int(category)) - 1
		} else {
			additionalBits = value
		}

		for i := int(category) - 1; i >= 0; i-- {
			bit := (additionalBits >> uint(i)) & 1
			result = append(result, byte(bit))
		}
	}

	return result, nil
}

func EncodeACCoefficients(coeffs []int, table HuffmanTable) ([]byte, error) {
	if len(coeffs) == 0 {
		return []byte{}, nil
	}

	result := make([]byte, 0, len(coeffs)*8)

	var runLength uint8

	for i := 0; i < len(coeffs); i++ {
		if coeffs[i] == 0 {
			runLength++

			if i == len(coeffs)-1 || runLength == 16 {
				if runLength == 16 {
					encodedAC, err := EncodeACCoefficient(15, 0, table)
					if err != nil {
						return nil, err
					}
					result = append(result, encodedAC...)
					runLength = 0
				} else {
					encodedAC, err := EncodeACCoefficient(0, 0, table)
					if err != nil {
						return nil, err
					}
					result = append(result, encodedAC...)
				}
			}
		} else {
			encodedAC, err := EncodeACCoefficient(runLength, coeffs[i], table)
			if err != nil {
				return nil, err
			}
			result = append(result, encodedAC...)
			runLength = 0
		}
	}

	return result, nil
}
