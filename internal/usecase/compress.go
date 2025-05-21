package usecase

import (
	"AlgosSem4Lab2Neo/internal/compression/blocking"
	"AlgosSem4Lab2Neo/internal/compression/colorspace"
	"AlgosSem4Lab2Neo/internal/compression/dct"
	encoding "AlgosSem4Lab2Neo/internal/compression/encoding"
	"AlgosSem4Lab2Neo/internal/compression/quantization"
	"AlgosSem4Lab2Neo/internal/compression/zigzag"
	"AlgosSem4Lab2Neo/internal/io"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type MJPEGMetadata struct {
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	Quality        int         `json:"quality"`
	BlockSize      int         `json:"blockSize"`
	DCCoefficients [][]int     `json:"dcCoefficients"` // DC коэффициенты для каждого канала
	ACCoefficients [][][]int   `json:"acCoefficients"` // AC коэффициенты для каждого канала и блока
	QMatrices      [][]float64 `json:"qMatrices"`      // Матрицы квантования для каждого канала
}

func CompressImage(inputPath, outputPath string, quality int, blockSize int) error {
	img, err := io.Load(inputPath)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	ycbcrData, err := colorspace.RGBToYCbCr(img)
	if err != nil {
		return fmt.Errorf("ошибка преобразования в YCbCr: %v", err)
	}

	metadata := MJPEGMetadata{
		Width:          width,
		Height:         height,
		Quality:        quality,
		BlockSize:      blockSize,
		DCCoefficients: make([][]int, 3),
		ACCoefficients: make([][][]int, 3),
		QMatrices:      make([][]float64, 3),
	}

	huffmanTables, err := encoding.GenerateStandardHuffmanTables()
	if err != nil {
		return fmt.Errorf("ошибка генерации таблиц Хаффмана: %v", err)
	}

	compressedData := make(map[string][]byte)

	for idx, channelInfo := range []struct {
		Name   string
		Data   []byte
		IsLuma bool
	}{
		{"Y", ycbcrData.Y, true},
		{"Cb", ycbcrData.Cb, false},
		{"Cr", ycbcrData.Cr, false},
	} {
		blocks, err := blocking.DivideIntoBlocks(channelInfo.Data, width, height, blockSize, channelInfo.Name)
		if err != nil {
			return fmt.Errorf("ошибка разбиения на блоки: %v", err)
		}

		qMatrix := quantization.GenerateQuantizationMatrix(blockSize, quality, channelInfo.IsLuma)

		metadata.QMatrices[idx] = qMatrix

		dcCoeffs := make([]int, len(blocks))
		allACCoeffs := make([][]int, len(blocks))

		for i, block := range blocks {
			dctCoeffs := dct.ApplyDCT(block.Data, blockSize)

			quantizedCoeffs := quantization.QuantizeBlock(dctCoeffs, qMatrix)

			zigzagCoeffs := zigzag.Scan(quantizedCoeffs, blockSize)

			dcCoeffs[i] = int(zigzagCoeffs[0])

			acCoeffs := make([]int, len(zigzagCoeffs)-1)
			for j := 1; j < len(zigzagCoeffs); j++ {
				acCoeffs[j-1] = int(zigzagCoeffs[j])
			}
			allACCoeffs[i] = acCoeffs
		}

		metadata.DCCoefficients[idx] = dcCoeffs
		metadata.ACCoefficients[idx] = allACCoeffs

		var dcTable, acTable encoding.HuffmanTable
		if channelInfo.IsLuma {
			dcTable = huffmanTables["dc_lum"]
			acTable = huffmanTables["ac_lum"]
		} else {
			dcTable = huffmanTables["dc_chrom"]
			acTable = huffmanTables["ac_chrom"]
		}

		dcBits, err := encoding.EncodeDCCoefficients(dcCoeffs, dcTable)
		if err != nil {
			return fmt.Errorf("ошибка кодирования DC коэффициентов: %v", err)
		}

		bitWriter := encoding.NewBitWriter()

		bitWriter.WriteBits(dcBits)

		for _, acCoeffs := range allACCoeffs {
			acBits, err := encoding.EncodeACCoefficients(acCoeffs, acTable)
			if err != nil {
				return fmt.Errorf("ошибка кодирования AC коэффициентов: %v", err)
			}
			bitWriter.WriteBits(acBits)
		}

		compressedData[channelInfo.Name] = bitWriter.Bytes()
	}

	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("ошибка создания директории: %v", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer outputFile.Close()

	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга метаданных: %v", err)
	}

	metadataLen := uint32(len(metadataBytes))
	outputFile.Write([]byte{
		byte(metadataLen >> 24),
		byte(metadataLen >> 16),
		byte(metadataLen >> 8),
		byte(metadataLen),
	})

	outputFile.Write(metadataBytes)

	for _, channelName := range []string{"Y", "Cb", "Cr"} {
		channelData := compressedData[channelName]

		dataLen := uint32(len(channelData))
		outputFile.Write([]byte{
			byte(dataLen >> 24),
			byte(dataLen >> 16),
			byte(dataLen >> 8),
			byte(dataLen),
		})

		outputFile.Write(channelData)
	}

	fmt.Printf("Изображение успешно сжато и сохранено как %s\n", outputPath)
	return nil
}
