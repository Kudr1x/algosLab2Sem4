package main

import (
	"AlgosLab2Sem4v4/DCT"
	"AlgosLab2Sem4v4/RLE"
	"AlgosLab2Sem4v4/ZigZag"
	c "AlgosLab2Sem4v4/color"
	"AlgosLab2Sem4v4/differential"
	"AlgosLab2Sem4v4/quant"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func Decompress(inputPath string, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	header := make([]byte, 5)
	if _, err := file.Read(header); err != nil {
		return fmt.Errorf("не удалось прочитать заголовок: %v", err)
	}
	width := int(header[0])<<8 | int(header[1])
	height := int(header[2])<<8 | int(header[3])
	quality := int(header[4])

	var yBlockCount, cbBlockCount, crBlockCount uint32
	if err := binary.Read(file, binary.BigEndian, &yBlockCount); err != nil {
		return fmt.Errorf("не удалось прочитать кол-во Y блоков: %v", err)
	}
	if err := binary.Read(file, binary.BigEndian, &cbBlockCount); err != nil {
		return fmt.Errorf("не удалось прочитать кол-во Cb блоков: %v", err)
	}
	if err := binary.Read(file, binary.BigEndian, &crBlockCount); err != nil {
		return fmt.Errorf("не удалось прочитать кол-во Cr блоков: %v", err)
	}

	yDcDiff := make([]int16, yBlockCount)
	cbDcDiff := make([]int16, cbBlockCount)
	crDcDiff := make([]int16, crBlockCount)

	for i := 0; i < int(yBlockCount); i++ {
		if err := binary.Read(file, binary.BigEndian, &yDcDiff[i]); err != nil {
			return fmt.Errorf("не удалось прочитать Y DC коэфф. %d: %v", i, err)
		}
	}
	for i := 0; i < int(cbBlockCount); i++ {
		if err := binary.Read(file, binary.BigEndian, &cbDcDiff[i]); err != nil {
			return fmt.Errorf("не удалось прочитать Cb DC коэфф. %d: %v", i, err)
		}
	}
	for i := 0; i < int(crBlockCount); i++ {
		if err := binary.Read(file, binary.BigEndian, &crDcDiff[i]); err != nil {
			return fmt.Errorf("не удалось прочитать Cr DC коэфф. %d: %v", i, err)
		}
	}

	yDc := differential.DifferentialDecode(toIntSlice(yDcDiff))
	cbDc := differential.DifferentialDecode(toIntSlice(cbDcDiff))
	crDc := differential.DifferentialDecode(toIntSlice(crDcDiff))

	yAcRle := make([][][]int, yBlockCount)
	cbAcRle := make([][][]int, cbBlockCount)
	crAcRle := make([][][]int, crBlockCount)

	if err := readACData(file, yAcRle); err != nil {
		return fmt.Errorf("ошибка чтения Y AC RLE: %v", err)
	}
	if err := readACData(file, cbAcRle); err != nil {
		return fmt.Errorf("ошибка чтения Cb AC RLE: %v", err)
	}
	if err := readACData(file, crAcRle); err != nil {
		return fmt.Errorf("ошибка чтения Cr AC RLE: %v", err)
	}

	yQuantBlocks := reconstructBlocks(yDc, yAcRle)
	cbQuantBlocks := reconstructBlocks(cbDc, cbAcRle)
	crQuantBlocks := reconstructBlocks(crDc, crAcRle)

	yQuantMatrix := quant.QuantCoeff(quant.QuantMatrixY, quality)
	cbcrQuantMatrix := quant.QuantCoeff(quant.QuantMatrixCbCr, quality)

	yDctBlocks := quant.InverseQuantizeBlocks(yQuantBlocks, yQuantMatrix)
	cbDctBlocks := quant.InverseQuantizeBlocks(cbQuantBlocks, cbcrQuantMatrix)
	crDctBlocks := quant.InverseQuantizeBlocks(crQuantBlocks, cbcrQuantMatrix)

	yIdctBlocks := DCT.Idct2DBlocks(yDctBlocks)
	cbIdctBlocks := DCT.Idct2DBlocks(cbDctBlocks)
	crIdctBlocks := DCT.Idct2DBlocks(crDctBlocks)

	yChannel := MergeBlocks(yIdctBlocks, width, height)
	cbChannel := MergeBlocks(cbIdctBlocks, (width+1)/2, (height+1)/2) // Chroma channels are downsampled
	crChannel := MergeBlocks(crIdctBlocks, (width+1)/2, (height+1)/2)

	cbUpsampled := c.Upsample(cbChannel, width, height)
	crUpsampled := c.Upsample(crChannel, width, height)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b := c.YCbCrToRGB(yChannel[y][x], cbUpsampled[y][x], crUpsampled[y][x])
			img.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("не удалось создать выходной файл PNG: %v", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return fmt.Errorf("не удалось закодировать PNG: %v", err)
	}

	return nil
}

func toIntSlice(slice []int16) []int {
	result := make([]int, len(slice))
	for i, v := range slice {
		result[i] = int(v)
	}
	return result
}

func MergeBlocks(blocks [][][]float64, width, height int) [][]uint8 {
	paddedWidth := ((width + 7) / 8) * 8

	channel := make([][]uint8, height)
	for i := range channel {
		channel[i] = make([]uint8, width)
	}

	blockXCount := paddedWidth / 8
	if blockXCount == 0 {
		blockXCount = 1
	}

	for blockIndex, blockData := range blocks {
		blockY := blockIndex / blockXCount
		blockX := blockIndex % blockXCount

		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				y := blockY*8 + i
				x := blockX*8 + j
				if y < height && x < width {
					val := blockData[i][j] + 128.0
					if val < 0 {
						val = 0
					} else if val > 255 {
						val = 255
					}
					channel[y][x] = uint8(math.Round(val))
				}
			}
		}
	}
	return channel
}

func readACData(file *os.File, acRleBlocks [][][]int) error {
	for i := range acRleBlocks {
		var rleListLength uint16
		if err := binary.Read(file, binary.BigEndian, &rleListLength); err != nil {
			return fmt.Errorf("чтение длины списка RLE для блока %d: %v", i, err)
		}

		reconstructedRleList := make([][]int, 0, rleListLength)
		var valueExpected bool = false

		for itemCounter := 0; itemCounter < int(rleListLength); itemCounter++ {
			if valueExpected {
				var acVal int8
				if err := binary.Read(file, binary.BigEndian, &acVal); err != nil {

					return fmt.Errorf("чтение значения AC для блока %d, RLE-элемент %d (ожидалось значение): %v", i, itemCounter, err)
				}
				reconstructedRleList = append(reconstructedRleList, []int{int(acVal)})
				valueExpected = false
			} else {
				var runPart int8
				if err := binary.Read(file, binary.BigEndian, &runPart); err != nil {
					return fmt.Errorf("чтение run-части RLE для блока %d, RLE-элемент %d: %v", i, itemCounter, err)
				}
				var catPart int16
				if err := binary.Read(file, binary.BigEndian, &catPart); err != nil {

					return fmt.Errorf("чтение cat-части RLE для блока %d, RLE-элемент %d (run-часть была %d): %v", i, itemCounter, runPart, err)
				}

				reconstructedRleList = append(reconstructedRleList, []int{int(runPart), int(catPart)})

				isEOB := (runPart == 0 && catPart == 0)
				isZRL := (runPart == 15 && catPart == 0)

				if !isEOB && !isZRL {
					valueExpected = true
				} else {
					valueExpected = false
				}
			}
		}

		if valueExpected {
			return fmt.Errorf("ошибка логики RLE в блоке %d: после обработки %d элементов ожидалось значение AC, но список RLE закончился", i, rleListLength)
		}
		acRleBlocks[i] = reconstructedRleList
	}
	return nil
}

func reconstructBlocks(dc []int, acRle [][][]int) [][][]int {
	blocks := make([][][]int, len(dc))
	for i := range blocks {
		acCoefficients := RLE.RunLengthDecode(acRle[i])
		allCoefficients := make([]int, 0, 64)
		allCoefficients = append(allCoefficients, dc[i])
		allCoefficients = append(allCoefficients, acCoefficients...)
		blocks[i] = ZigZag.InverseZigZagScan(allCoefficients)
	}
	return blocks
}
