package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
)

func main() {
	inputFile := "/home/kudrix/GolandProjects/AlgosLab2Sem4v4/lenna.png"
	compressedFile := "/home/kudrix/GolandProjects/AlgosLab2Sem4v4/lenna.bin"
	decompressedFile := "/home/kudrix/GolandProjects/AlgosLab2Sem4v4/lenna_decompressed.png"
	quality := 100

	if filepath.Ext(inputFile) != ".png" {
		fmt.Println("Ошибка: входной файл должен быть в формате PNG")
		return
	}

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Ошибка открытия файла: %v\n", err)
		return
	}
	img, err := png.Decode(file)
	file.Close()
	if err != nil {
		fmt.Printf("Ошибка декодирования PNG: %v\n", err)
		return
	}

	compressor := NewJpegCompressor(quality)
	err = compressor.Compress(img, compressedFile)
	if err != nil {
		fmt.Printf("Ошибка при сжатии: %v\n", err)
		return
	}
	fmt.Println("Изображение успешно сжато и сохранено в", compressedFile)

	err = Decompress(compressedFile, decompressedFile)
	if err != nil {
		fmt.Printf("Ошибка при декодировании: %v\n", err)
		return
	}
	fmt.Println("Изображение успешно распаковано и сохранено в", decompressedFile)
}
