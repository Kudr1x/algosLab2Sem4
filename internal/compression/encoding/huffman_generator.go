package encoding

import (
	"fmt"
)

type HuffmanCode struct {
	Code   uint16 // Код
	Length uint8  // Длина кода в битах
}

type HuffmanTable map[byte]HuffmanCode

func GenerateHuffmanTable(bits []byte, values []byte) (HuffmanTable, error) {
	if len(bits) != 16 {
		return nil, fmt.Errorf("список BITS должен содержать 16 элементов")
	}

	var totalCodes int
	for _, count := range bits {
		totalCodes += int(count)
	}

	if totalCodes > 256 || totalCodes != len(values) {
		return nil, fmt.Errorf("несоответствие количества кодов и значений")
	}

	codes := make([]uint16, totalCodes)
	var code uint16
	var index int

	for i := 0; i < 16; i++ {
		for j := 0; j < int(bits[i]); j++ {
			codes[index] = code
			code++
			index++
		}
		code <<= 1
	}

	table := make(HuffmanTable)
	index = 0
	for i := 0; i < 16; i++ {
		for j := 0; j < int(bits[i]); j++ {
			table[values[index]] = HuffmanCode{
				Code:   codes[index],
				Length: uint8(i + 1),
			}
			index++
		}
	}

	return table, nil
}

func GenerateStandardHuffmanTables() (map[string]HuffmanTable, error) {
	tables := make(map[string]HuffmanTable)

	dcLumTable, err := GenerateHuffmanTable(DCLuminanceBits, DCLuminanceValues)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации DC Luminance таблицы: %v", err)
	}
	tables["dc_lum"] = dcLumTable

	dcChromTable, err := GenerateHuffmanTable(DCChrominanceBits, DCChrominanceValues)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации DC Chrominance таблицы: %v", err)
	}
	tables["dc_chrom"] = dcChromTable

	acLumTable, err := GenerateHuffmanTable(ACLuminanceBits, ACLuminanceValues)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации AC Luminance таблицы: %v", err)
	}
	tables["ac_lum"] = acLumTable

	acChromTable, err := GenerateHuffmanTable(ACChrominanceBits, ACChrominanceValues)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации AC Chrominance таблицы: %v", err)
	}
	tables["ac_chrom"] = acChromTable

	return tables, nil
}
