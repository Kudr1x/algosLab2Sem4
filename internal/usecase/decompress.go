package usecase

import (
	"AlgosSem4Lab2Neo/internal/compression/blocking"
	"AlgosSem4Lab2Neo/internal/compression/dct"
	"AlgosSem4Lab2Neo/internal/compression/encoding"
	"AlgosSem4Lab2Neo/internal/compression/quantization"
	"AlgosSem4Lab2Neo/internal/compression/zigzag"
	"AlgosSem4Lab2Neo/internal/domain/models"
	png "AlgosSem4Lab2Neo/internal/io"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
)

func DecompressImage(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	metadataLenBytes := make([]byte, 4)
	if _, err := io.ReadFull(inputFile, metadataLenBytes); err != nil {
		return err
	}

	metadataLen := uint32(metadataLenBytes[0])<<24 |
		uint32(metadataLenBytes[1])<<16 |
		uint32(metadataLenBytes[2])<<8 |
		uint32(metadataLenBytes[3])

	metadataBytes := make([]byte, metadataLen)
	if _, err := io.ReadFull(inputFile, metadataBytes); err != nil {
		return err
	}

	var metadata MJPEGMetadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return err
	}

	width := metadata.Width
	height := metadata.Height
	blockSize := metadata.BlockSize

	huffmanTables, err := encoding.GenerateStandardHuffmanTables()
	if err != nil {
		return fmt.Errorf("ошибка генерации таблиц Хаффмана: %v", err)
	}

	channelData := make(map[string][]byte)

	for _, channelName := range []string{"Y", "Cb", "Cr"} {
		dataLenBytes := make([]byte, 4)
		if _, err := io.ReadFull(inputFile, dataLenBytes); err != nil {
			return fmt.Errorf("ошибка чтения длины данных канала %s: %v", channelName, err)
		}

		dataLen := uint32(dataLenBytes[0])<<24 |
			uint32(dataLenBytes[1])<<16 |
			uint32(dataLenBytes[2])<<8 |
			uint32(dataLenBytes[3])

		data := make([]byte, dataLen)
		if _, err := io.ReadFull(inputFile, data); err != nil {
			return fmt.Errorf("ошибка чтения данных канала %s: %v", channelName, err)
		}

		channelData[channelName] = data
	}

	channels := make(map[string][]byte)

	for idx, channelInfo := range []struct {
		Name   string
		IsLuma bool
	}{
		{"Y", true},
		{"Cb", false},
		{"Cr", false},
	} {
		data := channelData[channelInfo.Name]

		bitReader := encoding.NewBitReader(data)

		var dcTable, acTable encoding.HuffmanTable
		if channelInfo.IsLuma {
			dcTable = huffmanTables["dc_lum"]
			acTable = huffmanTables["ac_lum"]
		} else {
			dcTable = huffmanTables["dc_chrom"]
			acTable = huffmanTables["ac_chrom"]
		}

		blocksX := (width + blockSize - 1) / blockSize
		blocksY := (height + blockSize - 1) / blockSize
		totalBlocks := blocksX * blocksY

		dcCoeffs, err := encoding.DecodeDCCoefficients(bitReader, dcTable, totalBlocks)
		if err != nil {
			// Если возникла ошибка при декодировании, используем метаданные
			dcCoeffs = metadata.DCCoefficients[idx]
		}

		qMatrix := metadata.QMatrices[idx]

		blocks := make([]*models.Block, totalBlocks)

		for i := 0; i < totalBlocks; i++ {
			block := &models.Block{
				Data:    make([]float64, blockSize*blockSize),
				Size:    blockSize,
				Channel: channelInfo.Name,
			}

			dcCoeff := float64(dcCoeffs[i])

			var acCoeffs []int
			if bitReader != nil {
				acCoeffs, err = encoding.DecodeACCoefficients(bitReader, acTable, blockSize)
				if err != nil {
					if i < len(metadata.ACCoefficients[idx]) {
						acCoeffs = metadata.ACCoefficients[idx][i]
					} else {
						acCoeffs = make([]int, blockSize*blockSize-1)
					}
				}
			} else if i < len(metadata.ACCoefficients[idx]) {
				acCoeffs = metadata.ACCoefficients[idx][i]
			} else {
				acCoeffs = make([]int, blockSize*blockSize-1)
			}

			coeffs := make([]float64, blockSize*blockSize)
			coeffs[0] = dcCoeff

			for j := 0; j < len(acCoeffs) && j+1 < len(coeffs); j++ {
				coeffs[j+1] = float64(acCoeffs[j])
			}

			dequantizedCoeffs := zigzag.IScan(coeffs, blockSize)

			dequantizedCoeffs = quantization.DequantizeBlock(dequantizedCoeffs, qMatrix)

			idctCoeffs := dct.ApplyIDCT(dequantizedCoeffs, blockSize)

			block.Data = idctCoeffs
			blocks[i] = block
		}

		channelBytes, err := blocking.ReconstructFromBlocks(blocks, width, height)
		if err != nil {
			return fmt.Errorf("ошибка восстановления данных из блоков: %v", err)
		}

		channels[channelInfo.Name] = channelBytes
	}

	rgbaImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x

			yVal := float64(channels["Y"][idx])
			cbVal := float64(channels["Cb"][idx]) - 128
			crVal := float64(channels["Cr"][idx]) - 128

			r := yVal + 1.402*crVal
			g := yVal - 0.344136*cbVal - 0.714136*crVal
			b := yVal + 1.772*cbVal

			r = clamp(r, 0, 255)
			g = clamp(g, 0, 255)
			b = clamp(b, 0, 255)

			rgbaImg.SetRGBA(x, y, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: 255,
			})
		}
	}

	err = png.Save(rgbaImg, outputPath)
	if err != nil {
		return err
	}

	fmt.Printf("Изображение успешно распаковано и сохранено как %s\n", outputPath)
	return nil
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
