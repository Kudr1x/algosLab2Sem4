package encoding

import (
	"fmt"
)

func DecodeDCCoefficient(br *BitReader, table HuffmanTable) (int, error) {
	reverseTable := make(map[uint32]byte)
	for symbol, huffCode := range table {
		key := (uint32(huffCode.Code) << 8) | uint32(huffCode.Length)
		reverseTable[key] = symbol
	}

	var code uint16
	var codeLen uint8

	for codeLen < 16 {
		bit, ok := br.ReadBit()
		if !ok {
			return 0, fmt.Errorf("неожиданный конец данных")
		}

		code = (code << 1) | uint16(bit)
		codeLen++

		key := (uint32(code) << 8) | uint32(codeLen)
		if symbol, found := reverseTable[key]; found {
			category := int(symbol)

			if category == 0 {
				return 0, nil
			}

			value, ok := br.ReadValue(int(category))
			if !ok {
				return 0, fmt.Errorf("неожиданный конец данных при чтении дополнительных битов")
			}

			if value < (1 << (category - 1)) {
				value = value - (1 << category) + 1
			}

			return value, nil
		}
	}

	return 0, fmt.Errorf("код не найден в таблице Хаффмана")
}

func DecodeACCoefficient(br *BitReader, table HuffmanTable) (uint8, int, error) {
	reverseTable := make(map[uint32]byte)
	for symbol, huffCode := range table {
		key := (uint32(huffCode.Code) << 8) | uint32(huffCode.Length)
		reverseTable[key] = symbol
	}

	var code uint16
	var codeLen uint8

	for codeLen < 16 {
		bit, ok := br.ReadBit()
		if !ok {
			return 0, 0, fmt.Errorf("неожиданный конец данных")
		}

		code = (code << 1) | uint16(bit)
		codeLen++

		key := (uint32(code) << 8) | uint32(codeLen)
		if symbol, found := reverseTable[key]; found {
			runLength := symbol >> 4
			size := symbol & 0x0F

			if runLength == 0 && size == 0 {
				return runLength, 0, nil
			}

			if runLength == 15 && size == 0 {
				return runLength, 0, nil
			}

			value, ok := br.ReadValue(int(size))
			if !ok {
				return 0, 0, fmt.Errorf("неожиданный конец данных при чтении дополнительных битов")
			}

			if value < (1 << (size - 1)) {
				value = value - (1 << size) + 1
			}

			return runLength, value, nil
		}
	}

	return 0, 0, fmt.Errorf("код не найден в таблице Хаффмана")
}

func DecodeDCCoefficients(br *BitReader, table HuffmanTable, count int) ([]int, error) {
	result := make([]int, count)

	value, err := DecodeDCCoefficient(br, table)
	if err != nil {
		return nil, err
	}
	result[0] = value

	for i := 1; i < count; i++ {
		diff, err := DecodeDCCoefficient(br, table)
		if err != nil {
			return nil, err
		}
		result[i] = result[i-1] + diff
	}

	return result, nil
}

func DecodeACCoefficients(br *BitReader, table HuffmanTable, blockSize int) ([]int, error) {
	result := make([]int, blockSize*blockSize-1)

	var i int
	for i < len(result) {
		runLength, value, err := DecodeACCoefficient(br, table)
		if err != nil {
			return nil, err
		}

		if runLength == 0 && value == 0 {
			break
		}

		i += int(runLength)
		if i >= len(result) {
			return nil, fmt.Errorf("выход за пределы блока")
		}

		result[i] = value
		i++
	}

	return result, nil
}
