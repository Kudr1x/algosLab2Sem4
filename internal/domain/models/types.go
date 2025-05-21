package models

type Block struct {
	Data    []float64 // Данные блока в виде одномерного массива
	Size    int       // Размер блока (N для блока NxN)
	Channel string    // Канал ("Y", "Cb" или "Cr")
}

type YCbCrData struct {
	Y                []byte
	Cb               []byte
	Cr               []byte
	Width            int
	Height           int
	CbWidth          int
	CbHeight         int
	CrWidth          int
	CrHeight         int
	SubsamplingRatio string // например, "4:2:0", "4:2:2", "4:4:4"
}

type ACSymbol struct {
	RunLength uint8 // Количество предшествующих нулей (0-15)
	Size      uint8 // Категория ненулевого коэффициента (1-10)
}
