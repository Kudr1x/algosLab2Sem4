package vli

import (
	"fmt"
	"math"
)

func GetVLICategoryAndValue(value int) (int, string) {
	if value == 0 {
		return 0, ""
	}

	absValue := int(math.Abs(float64(value)))

	category := 0
	temp := absValue
	for temp > 0 {
		category++
		temp >>= 1
	}

	var bits string
	if value >= 0 {
		bits = fmt.Sprintf("%0*b", category, value)
	} else {
		invValue := (1 << uint(category)) - 1 + value
		bits = fmt.Sprintf("%0*b", category, invValue)
	}

	return category, bits
}

func DecodeVLI(category int, bits string) int {
	if category == 0 {
		return 0
	}

	bitValue, _ := fmt.Sscanf(bits, "%b", new(int))

	if len(bits) > 0 && bits[0] == '0' {
		lowerBound := -(1 << uint(category-1))
		return bitValue + lowerBound
	}

	return bitValue
}
