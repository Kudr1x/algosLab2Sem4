package huffman

import (
	"AlgosLab2Sem4v4/vli"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

var DefaultDCLuminanceBits = []int{0, 1, 5, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0}
var DefaultDCLuminanceHuffval = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

var DefaultDCChrominanceBits = []int{0, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0}
var DefaultDCChrominanceHuffval = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

var DefaultACLuminanceBits = []int{0, 2, 1, 3, 3, 2, 4, 3, 5, 5, 4, 4, 0, 0, 1, 125}
var DefaultACLuminanceHuffval = []int{
	0x01, 0x02, 0x03, 0x00, 0x04, 0x11, 0x05, 0x12, 0x21, 0x31, 0x41, 0x06, 0x13, 0x51, 0x61, 0x07,
	0x22, 0x71, 0x14, 0x32, 0x81, 0x91, 0xA1, 0x08, 0x23, 0x42, 0xB1, 0xC1, 0x15, 0x52, 0xD1, 0xF0,
	0x24, 0x33, 0x62, 0x72, 0x82, 0x09, 0x0A, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x25, 0x26, 0x27, 0x28,
	0x29, 0x2A, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49,
	0x4A, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69,
	0x6A, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89,
	0x8A, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9A, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7,
	0xA8, 0xA9, 0xAA, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7, 0xB8, 0xB9, 0xBA, 0xC2, 0xC3, 0xC4, 0xC5,
	0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9, 0xDA, 0xE1, 0xE2,
	0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8,
	0xF9, 0xFA,
}

var DefaultACChrominanceBits = []int{0, 2, 1, 2, 4, 4, 3, 4, 7, 5, 4, 4, 0, 1, 2, 119}
var DefaultACChrominanceHuffval = []int{
	0x00, 0x01, 0x02, 0x03, 0x11, 0x04, 0x05, 0x21, 0x31, 0x06, 0x12, 0x41, 0x51, 0x07, 0x61, 0x71,
	0x13, 0x22, 0x32, 0x81, 0x08, 0x14, 0x42, 0x91, 0xA1, 0xB1, 0xC1, 0x09, 0x23, 0x33, 0x52, 0xF0,
	0x15, 0x62, 0x72, 0xD1, 0x0A, 0x16, 0x24, 0x34, 0xE1, 0x25, 0xF1, 0x17, 0x18, 0x19, 0x1A, 0x26,
	0x27, 0x28, 0x29, 0x2A, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48,
	0x49, 0x4A, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68,
	0x69, 0x6A, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
	0x88, 0x89, 0x8A, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9A, 0xA2, 0xA3, 0xA4, 0xA5,
	0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7, 0xB8, 0xB9, 0xBA, 0xC2, 0xC3,
	0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9, 0xDA,
	0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8,
	0xF9, 0xFA,
}

type HuffmanTable struct {
	bits        []int
	huffval     []int
	encodeTable map[int]CodeLength
	decodeTable map[string]int
	maxCodeLen  int
}

type CodeLength struct {
	Code   int
	Length int
}

func NewHuffmanTable(bits []int, huffval []int) (*HuffmanTable, error) {
	if len(bits) != 16 {
		return nil, fmt.Errorf("список BITS должен содержать 16 элементов")
	}

	totalCodesInBits := sum(bits)
	if totalCodesInBits != len(huffval) {
		return nil, fmt.Errorf("сумма BITS (%d) не равна длине HUFFVAL (%d)", totalCodesInBits, len(huffval))
	}

	table := &HuffmanTable{
		bits:        make([]int, len(bits)),
		huffval:     make([]int, len(huffval)),
		encodeTable: make(map[int]CodeLength),
		decodeTable: make(map[string]int),
	}

	copy(table.bits, bits)
	copy(table.huffval, huffval)

	err := table.generateHuffmanCodes()
	if err != nil {
		return nil, err
	}

	table.buildDecodeStructure()

	return table, nil
}

func sum(arr []int) int {
	result := 0
	for _, v := range arr {
		result += v
	}
	return result
}

func (w *BitWriter) WriteBit(bit int) error {
	if bit != 0 && bit != 1 {
		return fmt.Errorf("бит должен быть 0 или 1")
	}

	w.buffer = (w.buffer << 1) | bit
	w.bitCount++

	if w.bitCount == 8 {
		return w.flushByte()
	}

	return nil
}

func (w *BitWriter) WriteBits(value int, numBits int) error {
	if numBits <= 0 {
		return nil
	}

	mask := (1 << numBits) - 1
	bitsToWrite := value & mask

	for i := numBits - 1; i >= 0; i-- {
		bit := (bitsToWrite >> i) & 1
		err := w.WriteBit(bit)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *BitWriter) flushByte() error {
	if w.bitCount != 8 {
		return fmt.Errorf("попытка сбросить неполный байт")
	}

	byteToWrite := byte(w.buffer)
	w.byteStream = append(w.byteStream, byteToWrite)

	if byteToWrite == 0xFF {
		w.byteStream = append(w.byteStream, 0x00) // Byte stuffing для 0xFF
	}

	w.buffer = 0
	w.bitCount = 0

	return nil
}

func (w *BitWriter) GetByteString() []byte {
	if w.bitCount > 0 {
		paddingBits := 8 - w.bitCount
		padMask := (1 << paddingBits) - 1
		w.buffer = (w.buffer << paddingBits) | padMask
		w.bitCount = 8
		_ = w.flushByte()
	}

	return w.byteStream
}

func (h *HuffmanTable) GetSpec() ([]int, []int) {
	return h.bits, h.huffval
}

func (h *HuffmanTable) generateHuffmanCodes() error {
	code := 0
	si := 1
	numCodesGenerated := 0
	huffvalIdx := 0

	for i := 0; i < 16; i++ {
		numCodesOfLengthSi := h.bits[i]

		for j := 0; j < numCodesOfLengthSi; j++ {
			if huffvalIdx >= len(h.huffval) {
				return fmt.Errorf("ошибка генерации кодов: индекс HUFFVAL (%d) выходит за пределы "+
					"для длины %d (i=%d, j=%d). Сумма BITS=%d, длина HUFFVAL=%d",
					huffvalIdx, si, i, j, sum(h.bits), len(h.huffval))
			}

			symbol := h.huffval[huffvalIdx]
			h.encodeTable[symbol] = CodeLength{Code: code, Length: si}
			huffvalIdx++
			code++
		}

		numCodesGenerated += numCodesOfLengthSi
		code <<= 1
		si++

		if numCodesOfLengthSi > 0 {
			h.maxCodeLen = i + 1
		}
	}

	if numCodesGenerated != len(h.huffval) {
		fmt.Printf("Предупреждение: Сгенерировано %d кодов, но HUFFVAL содержит %d символов.\n",
			numCodesGenerated, len(h.huffval))
	}

	return nil
}

func (h *HuffmanTable) buildDecodeStructure() {
	for symbol, codeLength := range h.encodeTable {
		codeStr := fmt.Sprintf("%0*b", codeLength.Length, codeLength.Code)
		h.decodeTable[codeStr] = symbol
	}
}

func (h *HuffmanTable) GetCode(symbol int) (CodeLength, bool) {
	codeLength, ok := h.encodeTable[symbol]
	return codeLength, ok
}

func (h *HuffmanTable) DecodeSymbol(bitReader *BitReader) (int, error) {
	currentCodeStr := ""

	for i := 0; i < h.maxCodeLen; i++ {
		bit, err := bitReader.ReadBit()
		if err != nil {
			if len(currentCodeStr) > 0 {
				fmt.Printf("Ошибка декодирования: конец потока после неполного кода '%s'\n", currentCodeStr)
			}
			return 0, err
		}

		currentCodeStr += fmt.Sprintf("%d", bit)

		if symbol, ok := h.decodeTable[currentCodeStr]; ok {
			return symbol, nil
		}
	}

	fmt.Printf("Ошибка декодирования: не найден символ для кода '%s' (макс. длина %d)\n",
		currentCodeStr, h.maxCodeLen)
	return 0, fmt.Errorf("не найден символ для кода")
}

type BitWriter struct {
	buffer     int
	bitCount   int
	byteStream []byte
}

func NewBitWriter() *BitWriter {
	return &BitWriter{
		buffer:     0,
		bitCount:   0,
		byteStream: make([]byte, 0),
	}
}

type BitReader struct {
	byteStream  *bytes.Reader
	currentByte byte
	bitPos      int
	markerFound bool
}

func NewBitReader(byteData []byte) *BitReader {
	return &BitReader{
		byteStream:  bytes.NewReader(byteData),
		currentByte: 0,
		bitPos:      8,
		markerFound: false,
	}
}

func (r *BitReader) loadByte() error {
	if r.markerFound {
		return io.EOF
	}

	byte, err := r.byteStream.ReadByte()
	if err != nil {
		return err
	}

	if byte == 0xFF {
		nextByte, err := r.byteStream.ReadByte()
		if err != nil {
			r.markerFound = true
			return io.EOF
		}

		if nextByte == 0x00 {
			r.currentByte = 0xFF
			r.bitPos = 0
			return nil
		} else {
			if _, err := r.byteStream.Seek(-2, io.SeekCurrent); err != nil {
				return err
			}
			r.markerFound = true
			return io.EOF
		}
	} else {
		r.currentByte = byte
		r.bitPos = 0
		return nil
	}
}

func (r *BitReader) ReadBit() (int, error) {
	if r.bitPos > 7 {
		if err := r.loadByte(); err != nil {
			return 0, err
		}
	}

	bit := (r.currentByte >> (7 - r.bitPos)) & 1
	r.bitPos++

	return int(bit), nil
}

func (r *BitReader) ReadBits(numBits int) (int, error) {
	if numBits < 0 {
		return 0, fmt.Errorf("количество бит не может быть отрицательным")
	}

	if numBits == 0 {
		return 0, nil
	}

	value := 0
	for i := 0; i < numBits; i++ {
		bit, err := r.ReadBit()
		if err != nil {
			return 0, fmt.Errorf("неожиданный конец потока/маркер при попытке чтения %d бит (прочитано %d бит): %w",
				numBits, i, err)
		}

		value = (value << 1) | bit
	}

	return value, nil
}

func GetVLICategoryAndValue(value int) (int, string) {
	if value == 0 {
		return 0, ""
	}

	absValue := value
	if absValue < 0 {
		absValue = -absValue
	}

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

	bitValue, _ := strconv.ParseInt(bits, 2, 64)

	if len(bits) > 0 && bits[0] == '0' {
		lowerBound := -(1 << uint(category-1))
		return int(bitValue) + lowerBound
	}

	return int(bitValue)
}

func HuffmanEncodeData(dcDiffs []int, acRlePairs [][][]int, dcTable, acTable *HuffmanTable) ([]byte, error) {
	bitWriter := NewBitWriter()

	for i, dcDiff := range dcDiffs {
		dcCategory, dcVliBits := vli.GetVLICategoryAndValue(dcDiff)

		dcCodeLength, ok := dcTable.GetCode(dcCategory)
		if !ok {
			return nil, fmt.Errorf("символ DC категории %d не найден в таблице Хаффмана", dcCategory)
		}

		err := bitWriter.WriteBits(dcCodeLength.Code, dcCodeLength.Length)
		if err != nil {
			return nil, err
		}

		// Записываем биты значения DC
		if dcCategory > 0 {
			dcVliVal, _ := fmt.Sscanf(dcVliBits, "%b", new(int))
			err = bitWriter.WriteBits(dcVliVal, dcCategory)
			if err != nil {
				return nil, err
			}
		}

		// Кодируем AC коэффициенты
		for _, pair := range acRlePairs[i] {
			runLength := pair[0]
			acValue := pair[1]

			if runLength == 0 && acValue == 0 {
				// EOB
				acSymbol := 0x00
				acCodeLength, ok := acTable.GetCode(acSymbol)
				if !ok {
					return nil, fmt.Errorf("символ EOB (0x00) не найден в AC таблице Хаффмана")
				}

				err := bitWriter.WriteBits(acCodeLength.Code, acCodeLength.Length)
				if err != nil {
					return nil, err
				}

				continue
			} else if runLength == 15 && acValue == 0 {
				// ZRL
				acSymbol := 0xF0
				acCodeLength, ok := acTable.GetCode(acSymbol)
				if !ok {
					return nil, fmt.Errorf("символ ZRL (0xF0) не найден в AC таблице Хаффмана")
				}

				err := bitWriter.WriteBits(acCodeLength.Code, acCodeLength.Length)
				if err != nil {
					return nil, err
				}
			} else {
				acCategory, acVliBits := vli.GetVLICategoryAndValue(acValue)

				if acCategory == 0 {
					return nil, fmt.Errorf("получена нулевая категория для ненулевого AC: %d", acValue)
				}

				if acCategory > 15 {
					return nil, fmt.Errorf("AC VLI категория %d не может быть > 15", acCategory)
				}

				// Составной символ: старшие 4 бита - runLength, младшие 4 - категория
				acSymbol := (runLength << 4) | acCategory
				acCodeLength, ok := acTable.GetCode(acSymbol)
				if !ok {
					return nil, fmt.Errorf("символ AC (run=%d, size=%d) не найден в таблице Хаффмана", runLength, acCategory)
				}

				err := bitWriter.WriteBits(acCodeLength.Code, acCodeLength.Length)
				if err != nil {
					return nil, err
				}

				// Записываем биты значения AC
				acVliVal, _ := fmt.Sscanf(acVliBits, "%b", new(int))
				err = bitWriter.WriteBits(acVliVal, acCategory)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return bitWriter.GetByteString(), nil
}

// HuffmanDecodeData декодирует Хаффман-закодированные данные для нескольких блоков
func HuffmanDecodeData(byteData []byte, dcTable *HuffmanTable, acTable *HuffmanTable, numBlocks int) ([][]interface{}, error) {
	bitReader := NewBitReader(byteData)
	decodedUnits := make([][]interface{}, 0, numBlocks)

	for blockIndex := 0; blockIndex < numBlocks; blockIndex++ {
		dcCategory, err := dcTable.DecodeSymbol(bitReader)
		if err != nil {
			fmt.Printf("Предупреждение: Ошибка конца потока при декодировании блока %d: %v. Декодировано %d блоков.\n",
				blockIndex+1, err, len(decodedUnits))
			break
		}

		dcVliBitsStr := ""
		if dcCategory > 15 {
			return decodedUnits, fmt.Errorf("декодирована некорректная DC категория %d > 15", dcCategory)
		}

		if dcCategory > 0 {
			dcVliVal, err := bitReader.ReadBits(dcCategory)
			if err != nil {
				return decodedUnits, err
			}

			dcVliBitsStr = fmt.Sprintf("%0*b", dcCategory, dcVliVal)
		}

		acRlePairs := make([][]int, 0)
		acCount := 0
		var acSymbol int

		for acCount < 64 {
			acSymbol, err = acTable.DecodeSymbol(bitReader)
			if err != nil {
				return decodedUnits, fmt.Errorf("не удалось декодировать AC символ в блоке %d после %d пар: %w",
					blockIndex+1, len(acRlePairs), err)
			}

			if acSymbol == 0x00 {
				// EOB
				acRlePairs = append(acRlePairs, []int{0, 0})
				break
			} else if acSymbol == 0xF0 {
				// ZRL
				acRlePairs = append(acRlePairs, []int{15, 0})
				acCount += 16
			} else {
				runLength := (acSymbol >> 4) & 0x0F
				acCategory := acSymbol & 0x0F

				if acCategory == 0 || acCategory > 15 {
					return decodedUnits, fmt.Errorf("некорректный AC символ 0x%02X (run=%d, size=%d)",
						acSymbol, runLength, acCategory)
				}

				acVliVal, err := bitReader.ReadBits(acCategory)
				if err != nil {
					return decodedUnits, err
				}

				acVliBitsStr := fmt.Sprintf("%0*b", acCategory, acVliVal)
				acValue := DecodeVLI(acCategory, acVliBitsStr)

				acRlePairs = append(acRlePairs, []int{runLength, acValue})
				acCount += runLength + 1
			}

			if acCount > 63 {
				fmt.Printf("Предупреждение: Счетчик AC (%d) превысил 63 в блоке %d. Возможно, лишние данные.\n",
					acCount, blockIndex+1)
			}
		}

		if acCount > 63 && acSymbol != 0x00 {
			fmt.Printf("Предупреждение: Цикл декодирования AC завершился с ac_count=%d > 63 и без EOB.\n", acCount)
		}

		decodedUnits = append(decodedUnits, []interface{}{dcCategory, dcVliBitsStr, acRlePairs})
	}

	return decodedUnits, nil
}
