package encoding

type BitWriter struct {
	buffer []byte
	bitPos int
}

func NewBitWriter() *BitWriter {
	return &BitWriter{
		buffer: make([]byte, 0, 1024), // Начальная емкость
		bitPos: 0,
	}
}

// WriteBit записывает один бит
func (bw *BitWriter) WriteBit(bit byte) {
	bytePos := bw.bitPos / 8
	bitOffset := 7 - (bw.bitPos % 8) // Биты идут от старшего к младшему

	// Расширяем буфер при необходимости
	for bytePos >= len(bw.buffer) {
		bw.buffer = append(bw.buffer, 0)
	}

	// Устанавливаем бит
	if bit != 0 {
		bw.buffer[bytePos] |= 1 << bitOffset
	}

	bw.bitPos++
}

// WriteBits записывает последовательность битов
func (bw *BitWriter) WriteBits(bits []byte) {
	for _, bit := range bits {
		bw.WriteBit(bit)
	}
}

// WriteValue записывает значение заданной длины
func (bw *BitWriter) WriteValue(value int, length int) {
	for i := length - 1; i >= 0; i-- {
		bit := (value >> i) & 1
		bw.WriteBit(byte(bit))
	}
}

// Bytes возвращает байтовый буфер
func (bw *BitWriter) Bytes() []byte {
	// Если последний байт не полностью заполнен, он все равно включается
	byteLen := (bw.bitPos + 7) / 8
	return bw.buffer[:byteLen]
}

// BitReader представляет собой структуру для чтения битов из байтового потока
type BitReader struct {
	buffer []byte
	bitPos int
}

// NewBitReader создает новый BitReader
func NewBitReader(data []byte) *BitReader {
	return &BitReader{
		buffer: data,
		bitPos: 0,
	}
}

// ReadBit читает один бит
func (br *BitReader) ReadBit() (byte, bool) {
	bytePos := br.bitPos / 8
	if bytePos >= len(br.buffer) {
		return 0, false
	}

	bitOffset := 7 - (br.bitPos % 8) // Биты идут от старшего к младшему
	bit := (br.buffer[bytePos] >> bitOffset) & 1

	br.bitPos++
	return bit, true
}

// ReadBits читает указанное количество битов
func (br *BitReader) ReadBits(count int) ([]byte, bool) {
	bits := make([]byte, count)
	for i := 0; i < count; i++ {
		bit, ok := br.ReadBit()
		if !ok {
			return nil, false
		}
		bits[i] = bit
	}
	return bits, true
}

// ReadValue читает значение заданной длины
func (br *BitReader) ReadValue(length int) (int, bool) {
	value := 0
	for i := 0; i < length; i++ {
		bit, ok := br.ReadBit()
		if !ok {
			return 0, false
		}
		value = (value << 1) | int(bit)
	}
	return value, true
}
