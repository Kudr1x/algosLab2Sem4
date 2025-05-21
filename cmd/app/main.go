package main

import (
	"AlgosSem4Lab2Neo/internal/io"
	"AlgosSem4Lab2Neo/internal/processing"
	"AlgosSem4Lab2Neo/internal/usecase"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {
	forestPath := "/home/kudrix/GolandProjects/AlgosSem4Lab2Neo/assets/forest.png"
	lennaPath := "/home/kudrix/GolandProjects/AlgosSem4Lab2Neo/assets/lenna.png"

	outputDir := "/home/kudrix/GolandProjects/AlgosSem4Lab2Neo/output"
	tempDir := filepath.Join(outputDir, "temp")
	compressedDir := filepath.Join(outputDir, "compressed")
	decompressedDir := filepath.Join(outputDir, "decompressed")

	os.MkdirAll(outputDir, 0755)
	os.MkdirAll(tempDir, 0755)
	os.MkdirAll(compressedDir, 0755)
	os.MkdirAll(decompressedDir, 0755)

	imagePaths := convertImages(forestPath, lennaPath, outputDir)

	buildSizeQualityGraphs(imagePaths, tempDir)

	compressAndDecompress(imagePaths, compressedDir, decompressedDir)
}

func convertImages(forestPath, lennaPath, outputDir string) []string {
	forestImg, _ := io.Load(forestPath)

	lennaImg, _ := io.Load(lennaPath)

	forestGray := processing.ConvertToGrayscale(forestImg)
	forestBWDithered := processing.ConvertToBWWithDithering(forestImg)
	forestBWNoDither := processing.ConvertToBWNoDithering(forestImg)

	lennaGray := processing.ConvertToGrayscale(lennaImg)
	lennaBWDithered := processing.ConvertToBWWithDithering(lennaImg)
	lennaBWNoDither := processing.ConvertToBWNoDithering(lennaImg)

	forestGrayPath := filepath.Join(outputDir, "forest_grayscale.png")
	forestBWDitheredPath := filepath.Join(outputDir, "forest_bw_dithered.png")
	forestBWNoDitherPath := filepath.Join(outputDir, "forest_bw_no_dither_420.png")
	lennaGrayPath := filepath.Join(outputDir, "lenna_grayscale.png")
	lennaBWDitheredPath := filepath.Join(outputDir, "lenna_bw_dithered.png")
	lennaBWNoDitherPath := filepath.Join(outputDir, "lenna_bw_no_dither_420.png")

	io.Save(forestGray, forestGrayPath)
	io.Save(forestBWDithered, forestBWDitheredPath)
	io.Save(forestBWNoDither, forestBWNoDitherPath)
	io.Save(lennaGray, lennaGrayPath)
	io.Save(lennaBWDithered, lennaBWDitheredPath)
	io.Save(lennaBWNoDither, lennaBWNoDitherPath)

	fmt.Println("Преобразование изображений завершено. Результаты сохранены в", outputDir)

	return []string{
		forestGrayPath,
		forestBWDitheredPath,
		forestBWNoDitherPath,
		lennaGrayPath,
		lennaBWDitheredPath,
		lennaBWNoDitherPath,
	}
}

func buildSizeQualityGraphs(imagePaths []string, tempDir string) {
	qualities := []int{1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100}

	fileSizes := make(map[string][]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	maxConcurrency := runtime.NumCPU()
	if maxConcurrency < 1 {
		maxConcurrency = 1
	}

	semaphore := make(chan struct{}, maxConcurrency)

	for _, imgPath := range imagePaths {
		imgName := filepath.Base(imgPath)
		imgName = imgName[:len(imgName)-4]

		mu.Lock()
		fileSizes[imgName] = make([]int, len(qualities))
		mu.Unlock()

		for i, quality := range qualities {
			wg.Add(1)

			semaphore <- struct{}{}

			go func(imgPath, imgName string, i, quality int) {
				defer func() {
					<-semaphore
					wg.Done()
				}()

				outputPath := filepath.Join(tempDir, fmt.Sprintf("%s_q%d.mjpeg", imgName, quality))

				err := usecase.CompressImage(imgPath, outputPath, quality, 8)
				if err != nil {
					fmt.Printf("Ошибка сжатия %s с качеством %d: %v\n", imgName, quality, err)
					return
				}

				fileInfo, err := os.Stat(outputPath)
				if err != nil {
					fmt.Printf("Ошибка получения информации о файле %s: %v\n", outputPath, err)
					return
				}

				mu.Lock()
				fileSizes[imgName][i] = int(fileInfo.Size())
				mu.Unlock()
			}(imgPath, imgName, i, quality)
		}
	}

	wg.Wait()

	jsonData, err := json.Marshal(fileSizes)
	if err != nil {
		fmt.Println("Ошибка маршалинга данных:", err)
		return
	}

	jsonPath := filepath.Join(tempDir, "compression_sizes.json")
	err = os.WriteFile(jsonPath, jsonData, 0644)
	if err != nil {
		fmt.Println("Ошибка сохранения JSON-файла:", err)
		return
	}

	cmd := exec.Command("python3", "/home/kudrix/GolandProjects/AlgosSem4Lab2Neo/scripts/plot_graph.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Ошибка запуска Python-скрипта:", err)
		return
	}

}

func compressAndDecompress(imagePaths []string, compressedDir, decompressedDir string) {
	qualities := []int{1, 20, 40, 60, 80, 100}

	var wg sync.WaitGroup

	maxConcurrency := runtime.NumCPU()
	if maxConcurrency < 1 {
		maxConcurrency = 1
	}

	semaphore := make(chan struct{}, maxConcurrency)

	for _, imgPath := range imagePaths {
		imgName := filepath.Base(imgPath)
		imgName = imgName[:len(imgName)-4]

		for _, quality := range qualities {
			wg.Add(1)

			semaphore <- struct{}{}

			go func(imgPath, imgName string, quality int) {
				defer func() {
					<-semaphore
					wg.Done()
				}()

				compressedPath := filepath.Join(compressedDir, fmt.Sprintf("%s_q%d.mjpeg", imgName, quality))
				decompressedPath := filepath.Join(decompressedDir, fmt.Sprintf("%s_q%d_decompressed.png", imgName, quality))

				err := usecase.CompressImage(imgPath, compressedPath, quality, 8)
				if err != nil {
					fmt.Printf("Ошибка сжатия %s с качеством %d: %v\n", imgName, quality, err)
					return
				}

				err = usecase.DecompressImage(compressedPath, decompressedPath)
				if err != nil {
					fmt.Printf("Ошибка декомпрессии %s с качеством %d: %v\n", imgName, quality, err)
					return
				}

				fmt.Printf("Обработано %s с качеством %d\n", imgName, quality)
			}(imgPath, imgName, quality)
		}
	}

	wg.Wait()
}
