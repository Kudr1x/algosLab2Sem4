package main

import (
	"AlgosLab2Sem4v4/DCT"
	"AlgosLab2Sem4v4/RLE"
	"AlgosLab2Sem4v4/ZigZag"
	"AlgosLab2Sem4v4/color"
	"AlgosLab2Sem4v4/differential"
	"AlgosLab2Sem4v4/quant"
	"encoding/binary"
	"fmt"
	"image"
	"os"
)

type JpegCompressor struct {
	Quality int
}

func NewJpegCompressor(quality int) *JpegCompressor {
	if quality < 1 {
		quality = 1
	} else if quality > 100 {
		quality = 100
	}
	return &JpegCompressor{Quality: quality}
}

func (jc *JpegCompressor) Compress(img image.Image, outputPath string) error {
	ycbcrImg := color.ConvertToYCbCr(img)

	cbDownsampled := color.Downsample(ycbcrImg.Cb, 2)
	crDownsampled := color.Downsample(ycbcrImg.Cr, 2)

	yBlocks := GetBlocks(ycbcrImg.Y)
	cbBlocks := GetBlocks(cbDownsampled)
	crBlocks := GetBlocks(crDownsampled)

	yDctBlocks := make([][][]float64, len(yBlocks))
	cbDctBlocks := make([][][]float64, len(cbBlocks))
	crDctBlocks := make([][][]float64, len(crBlocks))

	for i, block := range yBlocks {
		yDctBlocks[i] = DCT.DCT2D(block)
	}
	for i, block := range cbBlocks {
		cbDctBlocks[i] = DCT.DCT2D(block)
	}
	for i, block := range crBlocks {
		crDctBlocks[i] = DCT.DCT2D(block)
	}

	yQuantMatrix := quant.QuantCoeff(quant.QuantMatrixY, jc.Quality)
	cbcrQuantMatrix := quant.QuantCoeff(quant.QuantMatrixCbCr, jc.Quality)

	yQuantBlocks := make([][][]int, len(yDctBlocks))
	cbQuantBlocks := make([][][]int, len(cbDctBlocks))
	crQuantBlocks := make([][][]int, len(crDctBlocks))

	for i, block := range yDctBlocks {
		yQuantBlocks[i] = quant.Quantize(block, yQuantMatrix)
	}
	for i, block := range cbDctBlocks {
		cbQuantBlocks[i] = quant.Quantize(block, cbcrQuantMatrix)
	}
	for i, block := range crDctBlocks {
		crQuantBlocks[i] = quant.Quantize(block, cbcrQuantMatrix)
	}

	yDC := make([]int, len(yQuantBlocks))
	cbDC := make([]int, len(cbQuantBlocks))
	crDC := make([]int, len(crQuantBlocks))

	for i, block := range yQuantBlocks {
		yDC[i] = block[0][0]
	}
	for i, block := range cbQuantBlocks {
		cbDC[i] = block[0][0]
	}
	for i, block := range crQuantBlocks {
		crDC[i] = block[0][0]
	}

	yDcDiff := differential.DifferentialEncode(yDC)
	cbDcDiff := differential.DifferentialEncode(cbDC)
	crDcDiff := differential.DifferentialEncode(crDC)

	yZigzag := make([][]int, len(yQuantBlocks))
	cbZigzag := make([][]int, len(cbQuantBlocks))
	crZigzag := make([][]int, len(crQuantBlocks))

	for i, block := range yQuantBlocks {
		yZigzag[i] = ZigZag.ZigZagScan(block)
	}
	for i, block := range cbQuantBlocks {
		cbZigzag[i] = ZigZag.ZigZagScan(block)
	}
	for i, block := range crQuantBlocks {
		crZigzag[i] = ZigZag.ZigZagScan(block)
	}

	yAcRle := make([][][]int, len(yZigzag))
	cbAcRle := make([][][]int, len(cbZigzag))
	crAcRle := make([][][]int, len(crZigzag))

	for i, block := range yZigzag {
		yAcRle[i] = RLE.RunLengthEncode(block[1:])
	}
	for i, block := range cbZigzag {
		cbAcRle[i] = RLE.RunLengthEncode(block[1:])
	}
	for i, block := range crZigzag {
		crAcRle[i] = RLE.RunLengthEncode(block[1:])
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %v", err)
	}
	defer file.Close()

	header := []byte{
		byte(ycbcrImg.Width >> 8), byte(ycbcrImg.Width & 0xFF),
		byte(ycbcrImg.Height >> 8), byte(ycbcrImg.Height & 0xFF),
		byte(jc.Quality),
	}
	if _, err := file.Write(header); err != nil {
		return fmt.Errorf("ошибка записи заголовка: %v", err)
	}

	binary.Write(file, binary.BigEndian, uint32(len(yZigzag)))
	binary.Write(file, binary.BigEndian, uint32(len(cbZigzag)))
	binary.Write(file, binary.BigEndian, uint32(len(crZigzag)))

	for _, val := range yDcDiff {
		binary.Write(file, binary.BigEndian, int16(val))
	}
	for _, val := range cbDcDiff {
		binary.Write(file, binary.BigEndian, int16(val))
	}
	for _, val := range crDcDiff {
		binary.Write(file, binary.BigEndian, int16(val))
	}

	for _, blockRle := range yAcRle {
		binary.Write(file, binary.BigEndian, uint16(len(blockRle)))
		for _, run := range blockRle {
			binary.Write(file, binary.BigEndian, int8(run[0]))
			if len(run) > 1 {
				binary.Write(file, binary.BigEndian, int16(run[1]))
			}
		}
	}
	for _, blockRle := range cbAcRle {
		binary.Write(file, binary.BigEndian, uint16(len(blockRle)))
		for _, run := range blockRle {
			binary.Write(file, binary.BigEndian, int8(run[0]))
			if len(run) > 1 {
				binary.Write(file, binary.BigEndian, int16(run[1]))
			}
		}
	}
	for _, blockRle := range crAcRle {
		binary.Write(file, binary.BigEndian, uint16(len(blockRle)))
		for _, run := range blockRle {
			binary.Write(file, binary.BigEndian, int8(run[0]))
			if len(run) > 1 {
				binary.Write(file, binary.BigEndian, int16(run[1]))
			}
		}
	}

	return nil
}

func GetBlocks(channel [][]uint8) [][][]float64 {
	height := len(channel)
	if height == 0 {
		return nil
	}
	width := len(channel[0])

	paddedHeight := ((height + 7) / 8) * 8
	paddedWidth := ((width + 7) / 8) * 8

	padded := make([][]uint8, paddedHeight)
	for i := range padded {
		padded[i] = make([]uint8, paddedWidth)
		if i < height {
			copy(padded[i], channel[i])
		}
	}

	blocksY := paddedHeight / 8
	blocksX := paddedWidth / 8
	blocks := make([][][]float64, blocksY*blocksX)

	blockIndex := 0
	for by := 0; by < blocksY; by++ {
		for bx := 0; bx < blocksX; bx++ {
			block := make([][]float64, 8)
			for i := range block {
				block[i] = make([]float64, 8)
				for j := range block[i] {
					block[i][j] = float64(padded[by*8+i][bx*8+j]) - 128.0
				}
			}
			blocks[blockIndex] = block
			blockIndex++
		}
	}

	return blocks
}
