package processing

import (
	"image"
	"image/color"
	"math/rand"
)

// ConvertToGrayscale преобразует изображение в оттенки серого
func ConvertToGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Стандартная формула для преобразования RGB в оттенки серого
			gray := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)
			grayImg.SetGray(x, y, color.Gray{Y: gray})
		}
	}

	return grayImg
}

// ConvertToBWWithDithering преобразует изображение в ч/б с дизерингом (метод Флойда-Стейнберга)
func ConvertToBWWithDithering(img image.Image) *image.Gray {
	bounds := img.Bounds()
	bwImg := image.NewGray(bounds)

	// Создаем временное изображение в оттенках серого для работы с дизерингом
	tempImg := make([][]int, bounds.Dy())
	for i := range tempImg {
		tempImg[i] = make([]int, bounds.Dx())
	}

	// Заполняем временное изображение значениями яркости
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			tempImg[y][x] = int((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)
		}
	}

	// Применяем алгоритм дизеринга Флойда-Стейнберга
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			oldPixel := tempImg[y][x]
			newPixel := 0
			if oldPixel > 127 {
				newPixel = 255
			}

			// Устанавливаем новое значение пикселя
			bwImg.SetGray(x+bounds.Min.X, y+bounds.Min.Y, color.Gray{Y: uint8(newPixel)})

			// Распространяем ошибку квантования на соседние пиксели
			quantError := oldPixel - newPixel

			if x+1 < bounds.Dx() {
				tempImg[y][x+1] += quantError * 7 / 16
			}
			if y+1 < bounds.Dy() {
				if x-1 >= 0 {
					tempImg[y+1][x-1] += quantError * 3 / 16
				}
				tempImg[y+1][x] += quantError * 5 / 16
				if x+1 < bounds.Dx() {
					tempImg[y+1][x+1] += quantError * 1 / 16
				}
			}
		}
	}

	return bwImg
}

// ConvertToBWNoDithering преобразует изображение в ч/б без дизеринга (4:2:0)
func ConvertToBWNoDithering(img image.Image) *image.Gray {
	bounds := img.Bounds()
	bwImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)

			// Простое пороговое преобразование
			if gray > 127 {
				bwImg.SetGray(x, y, color.Gray{Y: 255})
			} else {
				bwImg.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	// Применяем субдискретизацию 4:2:0 (уменьшаем разрешение цветовых компонент)
	// В данном случае, так как мы уже в ч/б, просто имитируем эффект
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 2 {
			// Берем значение из верхнего левого пикселя блока 2x2
			val := bwImg.GrayAt(x, y).Y

			// Применяем небольшой шум для имитации артефактов 4:2:0
			noise := uint8(rand.Intn(10) - 5)
			if int(val)+int(noise) > 255 {
				val = 255
			} else if int(val)+int(noise) < 0 {
				val = 0
			} else {
				val += noise
			}

			// Устанавливаем одинаковое значение для блока 2x2
			for dy := 0; dy < 2 && y+dy < bounds.Max.Y; dy++ {
				for dx := 0; dx < 2 && x+dx < bounds.Max.X; dx++ {
					bwImg.SetGray(x+dx, y+dy, color.Gray{Y: val})
				}
			}
		}
	}

	return bwImg
}
